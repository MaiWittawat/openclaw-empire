package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"github.com/MaiWittawat/openclaw-empire/internal/model"
	"github.com/MaiWittawat/openclaw-empire/internal/service"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // allow all origins, restrict in production
	},
}

type Handler struct {
	svc *service.Service
	hub *Hub
}

type CreateTaskRequest struct {
	Title   string `json:"title"`
	AgentID string `json:"agent_id"`
}

type Hub struct {
	clients map[*websocket.Conn]bool
}

func NewHub() *Hub {
	return &Hub{clients: make(map[*websocket.Conn]bool)}
}

func (h *Hub) Register(c *websocket.Conn) {
	h.clients[c] = true
}

func (h *Hub) Unregister(c *websocket.Conn) {
	delete(h.clients, c)
}

func (h *Hub) Broadcast(msg model.WebSocketMessage) {
	data, _ := json.Marshal(msg)
	for client := range h.clients {
		client.WriteMessage(websocket.TextMessage, data)
	}
}

func NewHandler(svc *service.Service) *Handler {
	h := &Handler{svc: svc, hub: NewHub()}
	go h.simulateUpdates()
	return h
}

func (h *Handler) simulateUpdates() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		h.hub.Broadcast(model.WebSocketMessage{
			Event:   "agent_update",
			Payload: map[string]string{"status": "pulse"},
		})
	}
}

func (h *Handler) GetAgents(c *gin.Context) {
	agents, err := h.svc.GetAgents()
	if err != nil {
		c.JSON(500, gin.H{"data": nil, "error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": agents, "error": nil})
}

func (h *Handler) GetTasks(c *gin.Context) {
	tasks, err := h.svc.GetTasks()
	if err != nil {
		c.JSON(500, gin.H{"data": nil, "error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": tasks, "error": nil})
}

func (h *Handler) CreateTask(c *gin.Context) {
	var req CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"data": nil, "error": "Invalid request"})
		return
	}
	if req.Title == "" {
		c.JSON(400, gin.H{"data": nil, "error": "Title is required"})
		return
	}

	task, err := h.svc.CreateTask(req.Title, req.AgentID)
	if err != nil {
		c.JSON(500, gin.H{"data": nil, "error": err.Error()})
		return
	}

	h.hub.Broadcast(model.WebSocketMessage{
		Event:   "task_created",
		Payload: task,
	})

	c.JSON(200, gin.H{"data": task, "error": nil})
}

func (h *Handler) GetStats(c *gin.Context) {
	stats, err := h.svc.GetStats()
	if err != nil {
		c.JSON(500, gin.H{"data": nil, "error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": stats, "error": nil})
}

func (h *Handler) HandleWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	h.hub.Register(conn)
	defer h.hub.Unregister(conn)

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}
}
