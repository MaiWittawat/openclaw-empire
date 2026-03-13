package model

import "time"

type Agent struct {
	ID        string    `gorm:"primaryKey" json:"id"`
	Name      string    `json:"name"`
	Role      string    `json:"role"`
	Avatar    string    `json:"avatar"`
	Status    string    `json:"status"`
	Color     string    `json:"color"`
	Tokens    int       `json:"tokens"`
	MaxTokens int       `json:"max_tokens"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Task struct {
	ID        string    `gorm:"primaryKey" json:"id"`
	Title     string    `json:"title"`
	AgentID   string    `json:"agent_id"`
	Agent     *Agent    `gorm:"foreignKey:AgentID" json:"agent,omitempty"`
	AgentName string    `gorm:"-" json:"agent_name"`
	Status    string    `json:"status"`
	Tokens    int       `json:"tokens"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Task) TableName() string {
	return "tasks"
}

type Stats struct {
	ActiveAgents int     `json:"active_agents"`
	TotalAgents  int     `json:"total_agents"`
	TasksToday   int     `json:"tasks_today"`
	PendingTasks int     `json:"pending_tasks"`
	TokensUsed   int     `json:"tokens_used"`
	SuccessRate  float64 `json:"success_rate"`
	TotalTasks   int     `json:"total_tasks"`
	DoneTasks    int     `json:"done_tasks"`
}

type WebSocketMessage struct {
	Event   string      `json:"event"`
	Payload interface{} `json:"payload"`
}
