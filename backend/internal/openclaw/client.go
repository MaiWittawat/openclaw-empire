package openclaw

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"

	"github.com/MaiWittawat/openclaw-empire/internal/model"
)

type Config struct {
	GatewayWSURL       string
	GatewayToken       string
	CommandTimeout     time.Duration
	ActiveWindow       time.Duration
	DefaultSessionKey  string
	DeviceStorePath    string
	ClientID           string
	ClientVersion      string
	ClientPlatform     string
	ClientDeviceFamily string
	ReconnectDelay     time.Duration
}

type Event struct {
	Name    string
	Payload map[string]any
}

type DeviceIdentity struct {
	DeviceID    string `json:"device_id"`
	PublicKey   string `json:"public_key"`
	PrivateKey  string `json:"private_key"`
	DeviceToken string `json:"device_token"`
}

type rpcEnvelope struct {
	Type    string         `json:"type"`
	ID      string         `json:"id,omitempty"`
	Method  string         `json:"method,omitempty"`
	Event   string         `json:"event,omitempty"`
	OK      bool           `json:"ok,omitempty"`
	Params  map[string]any `json:"params,omitempty"`
	Payload map[string]any `json:"payload,omitempty"`
	Result  map[string]any `json:"result,omitempty"`
	Error   map[string]any `json:"error,omitempty"`
}

type pendingCall struct {
	result chan rpcEnvelope
}

type Client struct {
	cfg           Config
	device        *DeviceIdentity
	conn          *websocket.Conn
	writeMu       sync.Mutex
	pendingMu     sync.Mutex
	pending       map[string]*pendingCall
	agentsMu      sync.RWMutex
	agents        map[string]model.Agent
	sessionSubsMu sync.Mutex
	sessionSubs   map[string]bool
	eventHandler  func(Event)
}

func NewClient(cfg Config) (*Client, error) {
	device, err := loadOrCreateDevice(cfg.DeviceStorePath)
	if err != nil {
		return nil, err
	}

	return &Client{
		cfg:         cfg,
		device:      device,
		pending:     make(map[string]*pendingCall),
		agents:      make(map[string]model.Agent),
		sessionSubs: make(map[string]bool),
	}, nil
}

func (c *Client) Start(ctx context.Context, eventHandler func(Event)) {
	c.eventHandler = eventHandler

	go func() {
		for {
			if ctx.Err() != nil {
				return
			}

			if err := c.connectAndServe(ctx); err != nil && ctx.Err() == nil {
				time.Sleep(c.cfg.ReconnectDelay)
				continue
			}

			return
		}
	}()
}

func (c *Client) ListAgents() []model.Agent {
	c.agentsMu.RLock()
	defer c.agentsMu.RUnlock()

	agents := make([]model.Agent, 0, len(c.agents))
	for _, agent := range c.agents {
		agents = append(agents, agent)
	}
	slices.SortFunc(agents, func(left model.Agent, right model.Agent) int {
		return strings.Compare(left.Name, right.Name)
	})
	return agents
}

func (c *Client) DispatchTask(ctx context.Context, agentID string, title string) (string, string, error) {
	sessionKey := buildAgentSessionKey(agentID, c.cfg.DefaultSessionKey)
	params := map[string]any{
		"sessionKey":     sessionKey,
		"content":        title,
		"idempotencyKey": newID("idem"),
		"role":           "user",
	}

	result, err := c.call(ctx, "chat.send", params)
	if err != nil {
		return "", sessionKey, err
	}

	runID := firstString(
		result["runId"],
		result["run_id"],
		result["messageId"],
		result["id"],
	)

	return runID, sessionKey, nil
}

func (c *Client) RefreshSessions(ctx context.Context) error {
	result, err := c.call(ctx, "sessions.list", map[string]any{
		"limit":        100,
		"messageLimit": 0,
	})
	if err != nil {
		return err
	}

	items, err := parseCollectionAny(result)
	if err != nil {
		return err
	}

	c.upsertAgentsFromSessions(items)
	for _, item := range items {
		sessionKey := firstString(item["key"], item["sessionKey"])
		if sessionKey == "" {
			continue
		}
		_ = c.subscribeSession(ctx, sessionKey)
	}

	return nil
}

func (c *Client) connectAndServe(ctx context.Context) error {
	conn, _, err := websocket.DefaultDialer.DialContext(ctx, c.cfg.GatewayWSURL, nil)
	if err != nil {
		return err
	}
	c.conn = conn
	c.sessionSubsMu.Lock()
	c.sessionSubs = make(map[string]bool)
	c.sessionSubsMu.Unlock()
	defer func() {
		_ = conn.Close()
		c.conn = nil
	}()

	challenge, err := c.readEnvelope()
	if err != nil {
		return err
	}

	if challenge.Event != "connect.challenge" {
		return fmt.Errorf("unexpected gateway event: %s", challenge.Event)
	}

	challengeNonce := firstString(challenge.Params["nonce"], challenge.Payload["nonce"], challenge.Result["nonce"])
	callID := newID("connect")

	if err := c.writeEnvelope(rpcEnvelope{
		Type:   "req",
		ID:     callID,
		Method: "connect",
		Params: c.buildConnectParams(challengeNonce),
	}); err != nil {
		return err
	}

	response, err := c.readEnvelope()
	if err != nil {
		return err
	}
	if response.Type != "res" || response.ID != callID {
		return fmt.Errorf("unexpected connect response")
	}
	if !response.OK {
		return fmt.Errorf("gateway connect failed: %s", gatewayErrorMessage(response.Error))
	}

	if token := firstString(response.Result["deviceToken"], nestedMapString(response.Result, "auth", "deviceToken")); token != "" && token != c.device.DeviceToken {
		c.device.DeviceToken = token
		if err := saveDevice(c.cfg.DeviceStorePath, c.device); err != nil {
			return err
		}
	}

	errCh := make(chan error, 1)
	go func() {
		errCh <- c.readLoop()
	}()

	if err := c.RefreshSessions(ctx); err != nil {
		return err
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errCh:
		return err
	}
}

func (c *Client) readLoop() error {
	for {
		envelope, err := c.readEnvelope()
		if err != nil {
			return err
		}

		switch envelope.Type {
		case "res":
			c.resolvePending(envelope)
		case "event":
			c.handleEvent(envelope)
		}
	}
}

func (c *Client) handleEvent(envelope rpcEnvelope) {
	payload := envelope.Params
	if len(payload) == 0 {
		payload = envelope.Payload
	}

	if envelope.Event == "chat" || envelope.Event == "agent" || envelope.Event == "presence" {
		c.updateAgentCache(envelope.Event, payload)
	}

	if envelope.Event == "chat" {
		sessionKey := extractSessionKey(payload)
		if sessionKey != "" {
			go func() {
				_ = c.subscribeSession(context.Background(), sessionKey)
			}()
		}
	}

	if c.eventHandler != nil {
		c.eventHandler(Event{
			Name:    envelope.Event,
			Payload: payload,
		})
	}
}

func (c *Client) updateAgentCache(eventName string, payload map[string]any) {
	agentID := extractAgentID(payload)
	if agentID == "" {
		return
	}

	c.agentsMu.Lock()
	defer c.agentsMu.Unlock()

	current := c.agents[agentID]
	if current.ID == "" {
		current = model.Agent{
			ID:        agentID,
			Name:      agentID,
			Role:      "OpenClaw Agent",
			Avatar:    defaultAvatar(agentID),
			Color:     defaultColor(agentID),
			CreatedAt: time.Now(),
		}
	}

	current.Status = eventStatus(eventName, payload, current.Status)
	current.SessionKey = firstNonEmptyString(extractSessionKey(payload), current.SessionKey)
	current.Tokens = maxInt(current.Tokens, extractTokens(payload))
	current.LastSeenAt = time.Now()
	current.UpdatedAt = time.Now()
	c.agents[agentID] = current
}

func (c *Client) upsertAgentsFromSessions(items []map[string]any) {
	c.agentsMu.Lock()
	defer c.agentsMu.Unlock()

	now := time.Now()
	for _, item := range items {
		sessionKey := firstString(item["key"], item["sessionKey"])
		agentID := extractAgentID(item)
		if agentID == "" {
			continue
		}

		current := c.agents[agentID]
		if current.ID == "" {
			current = model.Agent{
				ID:        agentID,
				Name:      agentID,
				Role:      "OpenClaw Agent",
				Avatar:    defaultAvatar(agentID),
				Color:     defaultColor(agentID),
				CreatedAt: now,
			}
		}

		current.Status = normalizeStatus(firstString(item["status"], item["state"]))
		if current.Status == "idle" && recentlyActive(item, c.cfg.ActiveWindow) {
			current.Status = "working"
		}
		current.SessionKey = sessionKey
		current.Tokens = firstNonZero(
			intValue(item, "tokens", "token_usage", "totalTokens"),
			intValue(item, "contextTokens")+intValue(item, "inputTokens")+intValue(item, "outputTokens"),
		)
		current.LastSeenAt = now
		current.UpdatedAt = now
		c.agents[agentID] = current
	}
}

func (c *Client) subscribeSession(ctx context.Context, sessionKey string) error {
	c.sessionSubsMu.Lock()
	if c.sessionSubs[sessionKey] {
		c.sessionSubsMu.Unlock()
		return nil
	}
	c.sessionSubs[sessionKey] = true
	c.sessionSubsMu.Unlock()

	_, err := c.call(ctx, "chat.subscribe", map[string]any{
		"sessionKey": sessionKey,
	})
	return err
}

func (c *Client) call(ctx context.Context, method string, params map[string]any) (map[string]any, error) {
	if c.conn == nil {
		return nil, fmt.Errorf("gateway websocket is not connected")
	}

	if c.cfg.CommandTimeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, c.cfg.CommandTimeout)
		defer cancel()
	}

	callID := newID("rpc")
	pending := &pendingCall{result: make(chan rpcEnvelope, 1)}
	c.pendingMu.Lock()
	c.pending[callID] = pending
	c.pendingMu.Unlock()

	if err := c.writeEnvelope(rpcEnvelope{
		Type:   "req",
		ID:     callID,
		Method: method,
		Params: params,
	}); err != nil {
		c.pendingMu.Lock()
		delete(c.pending, callID)
		c.pendingMu.Unlock()
		return nil, err
	}

	response, err := c.awaitResponse(ctx, callID)
	if err != nil {
		return nil, err
	}
	if !response.OK {
		return nil, fmt.Errorf("%s failed: %s", method, gatewayErrorMessage(response.Error))
	}

	return response.Result, nil
}

func (c *Client) awaitResponse(ctx context.Context, callID string) (rpcEnvelope, error) {
	c.pendingMu.Lock()
	pending := c.pending[callID]
	c.pendingMu.Unlock()

	if pending == nil {
		return rpcEnvelope{}, fmt.Errorf("pending call %s was not found", callID)
	}

	select {
	case response := <-pending.result:
		c.pendingMu.Lock()
		delete(c.pending, callID)
		c.pendingMu.Unlock()
		return response, nil
	case <-ctx.Done():
		c.pendingMu.Lock()
		delete(c.pending, callID)
		c.pendingMu.Unlock()
		return rpcEnvelope{}, ctx.Err()
	}
}

func (c *Client) resolvePending(envelope rpcEnvelope) {
	c.pendingMu.Lock()
	defer c.pendingMu.Unlock()

	pending := c.pending[envelope.ID]
	if pending == nil {
		return
	}

	pending.result <- envelope
}

func (c *Client) writeEnvelope(envelope rpcEnvelope) error {
	c.writeMu.Lock()
	defer c.writeMu.Unlock()

	return c.conn.WriteJSON(envelope)
}

func (c *Client) readEnvelope() (rpcEnvelope, error) {
	var envelope rpcEnvelope
	if err := c.conn.ReadJSON(&envelope); err != nil {
		return rpcEnvelope{}, err
	}
	return envelope, nil
}

func (c *Client) buildConnectParams(nonce string) map[string]any {
	signedAt := time.Now().UTC().Format(time.RFC3339)
	signaturePayload := strings.Join([]string{
		"v2",
		c.device.DeviceID,
		c.cfg.ClientID,
		"server",
		"operator",
		"*",
		nonce,
		signedAt,
	}, "|")

	privateKey := decodePrivateKey(c.device.PrivateKey)
	signature := ed25519.Sign(privateKey, []byte(signaturePayload))

	auth := map[string]any{}
	if c.device.DeviceToken != "" {
		auth["deviceToken"] = c.device.DeviceToken
	} else {
		auth["token"] = c.cfg.GatewayToken
	}

	return map[string]any{
		"client": map[string]any{
			"id":           c.cfg.ClientID,
			"version":      c.cfg.ClientVersion,
			"mode":         "server",
			"platform":     c.cfg.ClientPlatform,
			"deviceFamily": c.cfg.ClientDeviceFamily,
		},
		"role":   "operator",
		"scopes": []string{"*"},
		"auth":   auth,
		"device": map[string]any{
			"id":        c.device.DeviceID,
			"publicKey": c.device.PublicKey,
			"signature": base64.RawURLEncoding.EncodeToString(signature),
			"signedAt":  signedAt,
		},
	}
}

func loadOrCreateDevice(storePath string) (*DeviceIdentity, error) {
	if strings.TrimSpace(storePath) == "" {
		return nil, fmt.Errorf("OPENCLAW_DEVICE_STORE_PATH is required")
	}

	if raw, err := os.ReadFile(filepath.Clean(storePath)); err == nil {
		var device DeviceIdentity
		if err := json.Unmarshal(raw, &device); err == nil {
			return &device, nil
		}
	}

	_, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}

	publicKey := privateKey.Public().(ed25519.PublicKey)
	hash := sha256.Sum256(publicKey)
	device := &DeviceIdentity{
		DeviceID:   "sha256:" + base64.RawURLEncoding.EncodeToString(hash[:]),
		PublicKey:  base64.RawURLEncoding.EncodeToString(publicKey),
		PrivateKey: base64.RawURLEncoding.EncodeToString(privateKey),
	}

	if err := saveDevice(storePath, device); err != nil {
		return nil, err
	}

	return device, nil
}

func saveDevice(storePath string, device *DeviceIdentity) error {
	if err := os.MkdirAll(filepath.Dir(storePath), 0o755); err != nil {
		return err
	}

	payload, err := json.Marshal(device)
	if err != nil {
		return err
	}

	return os.WriteFile(filepath.Clean(storePath), payload, 0o600)
}

func decodePrivateKey(encoded string) ed25519.PrivateKey {
	raw, err := base64.RawURLEncoding.DecodeString(encoded)
	if err != nil {
		return nil
	}
	return ed25519.PrivateKey(raw)
}

func gatewayErrorMessage(errorBody map[string]any) string {
	if len(errorBody) == 0 {
		return "unknown gateway error"
	}

	return firstString(errorBody["message"], errorBody["code"], errorBody["type"])
}

func parseCollectionAny(payload any) ([]map[string]any, error) {
	switch value := payload.(type) {
	case []map[string]any:
		return value, nil
	case []any:
		return toMapSlice(value), nil
	case map[string]any:
		for _, key := range []string{"sessions", "items", "data"} {
			if nested, ok := value[key]; ok {
				return parseCollectionAny(nested)
			}
		}
	}
	return nil, fmt.Errorf("unexpected collection payload")
}

func toMapSlice(items []any) []map[string]any {
	result := make([]map[string]any, 0, len(items))
	for _, item := range items {
		if mapped, ok := item.(map[string]any); ok {
			result = append(result, mapped)
		}
	}
	return result
}

func intValue(item map[string]any, keys ...string) int {
	for _, key := range keys {
		switch typed := item[key].(type) {
		case float64:
			return int(typed)
		case int:
			return typed
		}
	}
	return 0
}

func recentlyActive(item map[string]any, activeWindow time.Duration) bool {
	updatedAt := timeValue(item, "updatedAt", "updated_at", "lastActivityAt", "last_activity_at")
	if updatedAt.IsZero() {
		return false
	}
	return time.Since(updatedAt) <= activeWindow
}

func timeValue(item map[string]any, keys ...string) time.Time {
	for _, key := range keys {
		switch value := item[key].(type) {
		case string:
			for _, layout := range []string{time.RFC3339Nano, time.RFC3339, "2006-01-02 15:04:05"} {
				if parsed, err := time.Parse(layout, value); err == nil {
					return parsed
				}
			}
		}
	}
	return time.Time{}
}

func normalizeStatus(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "running", "active", "busy", "working", "in_progress":
		return "working"
	case "completed", "done", "finished", "success", "succeeded":
		return "done"
	case "error", "failed", "failure":
		return "error"
	case "queued", "pending", "waiting":
		return "pending"
	default:
		return "idle"
	}
}

func defaultAvatar(agentID string) string {
	avatars := []string{"◈", "◆", "⬢", "⬡", "✦", "✧"}
	return avatars[hashIndex(agentID, len(avatars))]
}

func defaultColor(agentID string) string {
	colors := []string{"#22c55e", "#0ea5e9", "#f97316", "#f43f5e", "#eab308", "#14b8a6"}
	return colors[hashIndex(agentID, len(colors))]
}

func hashIndex(value string, length int) int {
	sum := 0
	for _, char := range value {
		sum += int(char)
	}
	if length == 0 {
		return 0
	}
	return sum % length
}

func buildAgentSessionKey(agentID string, defaultSessionKey string) string {
	return fmt.Sprintf("agent:%s:%s", agentID, fallbackString(defaultSessionKey, "main"))
}

func extractAgentID(payload map[string]any) string {
	if agentID := firstString(payload["agentId"], payload["agent_id"], payload["agent"]); agentID != "" {
		return agentID
	}
	sessionKey := extractSessionKey(payload)
	if strings.HasPrefix(sessionKey, "agent:") {
		parts := strings.Split(sessionKey, ":")
		if len(parts) >= 2 {
			return parts[1]
		}
	}
	return ""
}

func extractSessionKey(payload map[string]any) string {
	return firstString(payload["sessionKey"], payload["session_key"], payload["key"])
}

func extractTokens(payload map[string]any) int {
	return firstNonZero(
		intValue(payload, "tokens", "totalTokens"),
		intValue(payload, "contextTokens")+intValue(payload, "inputTokens")+intValue(payload, "outputTokens"),
	)
}

func eventStatus(eventName string, payload map[string]any, fallback string) string {
	if status := normalizeStatus(firstString(payload["status"], payload["state"], payload["phase"])); status != "idle" {
		return status
	}

	switch eventName {
	case "presence":
		return "working"
	case "chat", "agent":
		return "working"
	default:
		return fallback
	}
}

func nestedMapString(payload map[string]any, key string, nestedKey string) any {
	nested, ok := payload[key].(map[string]any)
	if !ok {
		return ""
	}
	return nested[nestedKey]
}

func firstString(values ...any) string {
	for _, value := range values {
		switch typed := value.(type) {
		case string:
			if strings.TrimSpace(typed) != "" {
				return typed
			}
		}
	}
	return ""
}

func firstNonEmptyString(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}

func fallbackString(value string, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return value
}

func firstNonZero(values ...int) int {
	for _, value := range values {
		if value != 0 {
			return value
		}
	}
	return 0
}

func maxInt(left int, right int) int {
	if left > right {
		return left
	}
	return right
}

func newID(prefix string) string {
	return fmt.Sprintf("%s_%d", prefix, time.Now().UnixNano())
}
