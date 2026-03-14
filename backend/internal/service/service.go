package service

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/MaiWittawat/openclaw-empire/internal/model"
	"github.com/MaiWittawat/openclaw-empire/internal/openclaw"
	"github.com/MaiWittawat/openclaw-empire/internal/repository"
)

type Service struct {
	repo        *repository.Repository
	openclaw    *openclaw.Client
	broadcaster func(model.WebSocketMessage)
}

func NewService(repo *repository.Repository, openclawClient *openclaw.Client) *Service {
	return &Service{
		repo:     repo,
		openclaw: openclawClient,
	}
}

func (s *Service) SetBroadcaster(broadcaster func(model.WebSocketMessage)) {
	s.broadcaster = broadcaster
}

func (s *Service) Start(ctx context.Context) {
	s.openclaw.Start(ctx, s.handleGatewayEvent)
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				_ = s.syncAgents(ctx)
			}
		}
	}()
}

func (s *Service) GetAgents() ([]model.Agent, error) {
	if err := s.syncAgents(context.Background()); err == nil {
		return s.repo.GetAgents()
	}
	return s.repo.GetAgents()
}

func (s *Service) GetTasks() ([]model.Task, error) {
	return s.repo.GetTasks(100)
}

func (s *Service) GetTaskEvents(taskID string) ([]model.TaskEvent, error) {
	return s.repo.GetTaskEvents(taskID, 200)
}

func (s *Service) CreateTask(title, agentID string) (*model.Task, error) {
	task := &model.Task{
		ID:      repository.NewTaskID(),
		Title:   title,
		AgentID: agentID,
		Status:  "pending",
	}

	if err := s.repo.CreateTask(task); err != nil {
		return nil, err
	}

	if err := s.repo.AddTaskEvent(&model.TaskEvent{
		ID:         repository.NewTaskEventID(),
		TaskID:     task.ID,
		Event:      "task.created",
		Status:     "pending",
		SessionKey: task.SessionKey,
		Payload:    `{"message":"task created"}`,
	}); err != nil {
		return nil, err
	}

	runID, sessionKey, err := s.openclaw.DispatchTask(context.Background(), agentID, title)
	if err != nil {
		_ = s.repo.UpdateTaskProgress(task.ID, "error", 0, err.Error(), true)
		_ = s.repo.AddTaskEvent(&model.TaskEvent{
			ID:         repository.NewTaskEventID(),
			TaskID:     task.ID,
			Event:      "task.dispatch_failed",
			Status:     "error",
			SessionKey: task.SessionKey,
			Payload:    marshalPayload(map[string]any{"error": err.Error()}),
		})
		return nil, err
	}

	task.RunID = runID
	task.SessionKey = sessionKey
	task.Status = "working"
	now := time.Now()
	task.StartedAt = &now
	task.UpdatedAt = now

	if err := s.repo.UpdateTaskDispatch(task.ID, "working", runID, sessionKey); err != nil {
		return nil, err
	}

	_ = s.repo.AddTaskEvent(&model.TaskEvent{
		ID:         repository.NewTaskEventID(),
		TaskID:     task.ID,
		Event:      "task.dispatched",
		Status:     "working",
		RunID:      runID,
		SessionKey: sessionKey,
		Payload:    marshalPayload(map[string]any{"run_id": runID, "session_key": sessionKey}),
	})

	s.broadcast(model.WebSocketMessage{
		Event:   "task_created",
		Payload: task,
	})

	return task, nil
}

func (s *Service) GetStats() (*model.Stats, error) {
	if err := s.syncAgents(context.Background()); err != nil {
		return s.repo.GetStats()
	}
	return s.repo.GetStats()
}

func (s *Service) syncAgents(ctx context.Context) error {
	if err := s.openclaw.RefreshSessions(ctx); err != nil {
		return err
	}

	return s.repo.UpsertAgents(s.openclaw.ListAgents())
}

func (s *Service) handleGatewayEvent(event openclaw.Event) {
	agentID := agentIDFromPayload(event.Payload)
	sessionKey := sessionKeyFromPayload(event.Payload)
	runID := stringFromPayload(event.Payload, "runId", "run_id", "messageId", "message_id", "id")
	status := normalizeTaskStatus(event.Name, event.Payload)
	tokens := tokensFromPayload(event.Payload)
	payload := marshalPayload(event.Payload)

	if agentID != "" {
		_ = s.repo.UpsertAgents(s.openclaw.ListAgents())
	}

	task, err := s.repo.FindLatestTask(agentID, sessionKey, runID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return
	}

	if task != nil && task.ID != "" {
		completed := status == "done" || status == "error"
		_ = s.repo.UpdateTaskProgress(task.ID, status, max(tokens, task.Tokens), errorMessageFromPayload(event.Payload), completed)
		_ = s.repo.AddTaskEvent(&model.TaskEvent{
			ID:         repository.NewTaskEventID(),
			TaskID:     task.ID,
			Event:      event.Name,
			Status:     status,
			RunID:      firstNonEmpty(runID, task.RunID),
			SessionKey: firstNonEmpty(sessionKey, task.SessionKey),
			Payload:    payload,
		})
	}

	s.broadcast(model.WebSocketMessage{
		Event:   event.Name,
		Payload: event.Payload,
	})
}

func (s *Service) broadcast(message model.WebSocketMessage) {
	if s.broadcaster != nil {
		s.broadcaster(message)
	}
}

func marshalPayload(payload any) string {
	raw, err := json.Marshal(payload)
	if err != nil {
		return `{}`
	}
	return string(raw)
}

func normalizeTaskStatus(eventName string, payload map[string]any) string {
	if status := stringFromPayload(payload, "status", "state", "phase"); status != "" {
		switch status {
		case "running", "working", "active", "busy", "in_progress":
			return "working"
		case "completed", "done", "success", "succeeded", "finished":
			return "done"
		case "failed", "error":
			return "error"
		case "pending", "queued", "waiting":
			return "pending"
		}
	}

	switch eventName {
	case "agent", "chat", "presence":
		return "working"
	default:
		return "pending"
	}
}

func agentIDFromPayload(payload map[string]any) string {
	if value := stringFromPayload(payload, "agentId", "agent_id", "agent"); value != "" {
		return value
	}
	sessionKey := sessionKeyFromPayload(payload)
	parts := []rune(sessionKey)
	if len(parts) == 0 || len(parts) < 7 || sessionKey[:6] != "agent:" {
		return ""
	}
	chunks := splitSessionKey(sessionKey)
	if len(chunks) < 2 {
		return ""
	}
	return chunks[1]
}

func sessionKeyFromPayload(payload map[string]any) string {
	return stringFromPayload(payload, "sessionKey", "session_key", "key")
}

func tokensFromPayload(payload map[string]any) int {
	for _, key := range []string{"tokens", "totalTokens", "outputTokens", "inputTokens", "contextTokens"} {
		if value, ok := payload[key].(float64); ok {
			return int(value)
		}
	}
	return 0
}

func errorMessageFromPayload(payload map[string]any) string {
	return stringFromPayload(payload, "error", "message", "detail")
}

func stringFromPayload(payload map[string]any, keys ...string) string {
	for _, key := range keys {
		if value, ok := payload[key].(string); ok && value != "" {
			return value
		}
	}
	return ""
}

func splitSessionKey(sessionKey string) []string {
	parts := make([]string, 0)
	current := ""
	for _, char := range sessionKey {
		if char == ':' {
			parts = append(parts, current)
			current = ""
			continue
		}
		current += string(char)
	}
	if current != "" {
		parts = append(parts, current)
	}
	return parts
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}

func max(left int, right int) int {
	if left > right {
		return left
	}
	return right
}
