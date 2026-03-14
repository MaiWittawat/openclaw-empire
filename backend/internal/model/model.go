package model

import "time"

type Agent struct {
	ID         string    `gorm:"primaryKey" json:"id"`
	Name       string    `json:"name"`
	Role       string    `json:"role"`
	Avatar     string    `json:"avatar"`
	Status     string    `json:"status"`
	Color      string    `json:"color"`
	Tokens     int       `json:"tokens"`
	MaxTokens  int       `json:"max_tokens"`
	SessionKey string    `json:"session_key"`
	LastSeenAt time.Time `json:"last_seen_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	CreatedAt  time.Time `json:"created_at"`
}

type Task struct {
	ID           string     `gorm:"primaryKey" json:"id"`
	Title        string     `json:"title"`
	AgentID      string     `gorm:"index" json:"agent_id"`
	Agent        *Agent     `gorm:"foreignKey:AgentID" json:"agent,omitempty"`
	AgentName    string     `gorm:"-" json:"agent_name"`
	SessionKey   string     `gorm:"index" json:"session_key"`
	RunID        string     `gorm:"index" json:"run_id"`
	Status       string     `gorm:"index" json:"status"`
	Tokens       int        `json:"tokens"`
	ErrorMessage string     `json:"error_message"`
	StartedAt    *time.Time `json:"started_at,omitempty"`
	CompletedAt  *time.Time `json:"completed_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

func (Task) TableName() string {
	return "tasks"
}

type TaskEvent struct {
	ID         string    `gorm:"primaryKey" json:"id"`
	TaskID     string    `gorm:"index" json:"task_id"`
	Task       *Task     `gorm:"foreignKey:TaskID" json:"task,omitempty"`
	Event      string    `gorm:"index" json:"event"`
	Status     string    `gorm:"index" json:"status"`
	RunID      string    `gorm:"index" json:"run_id"`
	SessionKey string    `gorm:"index" json:"session_key"`
	Payload    string    `gorm:"type:text" json:"payload"`
	CreatedAt  time.Time `json:"created_at"`
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
