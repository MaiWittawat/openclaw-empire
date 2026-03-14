package repository

import (
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/MaiWittawat/openclaw-empire/internal/model"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) AutoMigrate() error {
	return r.db.AutoMigrate(&model.Agent{}, &model.Task{}, &model.TaskEvent{})
}

func (r *Repository) GetAgents() ([]model.Agent, error) {
	var agents []model.Agent
	err := r.db.Order("name ASC").Find(&agents).Error
	return agents, err
}

func (r *Repository) UpsertAgents(agents []model.Agent) error {
	if len(agents) == 0 {
		return nil
	}

	return r.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"name",
			"role",
			"avatar",
			"status",
			"color",
			"tokens",
			"max_tokens",
			"session_key",
			"last_seen_at",
			"updated_at",
		}),
	}).Create(&agents).Error
}

func (r *Repository) GetTasks(limit int) ([]model.Task, error) {
	if limit <= 0 {
		limit = 50
	}

	var tasks []model.Task
	err := r.db.Preload("Agent").Order("created_at DESC").Limit(limit).Find(&tasks).Error
	return tasks, err
}

func (r *Repository) CreateTask(task *model.Task) error {
	return r.db.Create(task).Error
}

func (r *Repository) UpdateTaskDispatch(taskID string, status string, runID string, sessionKey string) error {
	updates := map[string]any{
		"status":      status,
		"run_id":      runID,
		"session_key": sessionKey,
		"updated_at":  time.Now(),
	}

	return r.db.Model(&model.Task{}).Where("id = ?", taskID).Updates(updates).Error
}

func (r *Repository) UpdateTaskProgress(taskID string, status string, tokens int, errMessage string, completed bool) error {
	updates := map[string]any{
		"status":        status,
		"tokens":        tokens,
		"error_message": errMessage,
		"updated_at":    time.Now(),
	}

	now := time.Now()
	if status == "working" {
		updates["started_at"] = &now
	}
	if completed {
		updates["completed_at"] = &now
	}

	return r.db.Model(&model.Task{}).Where("id = ?", taskID).Updates(updates).Error
}

func (r *Repository) FindLatestTask(agentID string, sessionKey string, runID string) (*model.Task, error) {
	var task model.Task

	if runID != "" {
		if err := r.db.Where("run_id = ?", runID).Order("created_at DESC").First(&task).Error; err == nil {
			return &task, nil
		}
	}

	if sessionKey != "" {
		if err := r.db.Where("session_key = ?", sessionKey).Order("created_at DESC").First(&task).Error; err == nil {
			return &task, nil
		}
	}

	if agentID != "" {
		if err := r.db.Where("agent_id = ?", agentID).Order("created_at DESC").First(&task).Error; err == nil {
			return &task, nil
		}
	}

	return nil, gorm.ErrRecordNotFound
}

func (r *Repository) AddTaskEvent(event *model.TaskEvent) error {
	return r.db.Create(event).Error
}

func (r *Repository) GetTaskEvents(taskID string, limit int) ([]model.TaskEvent, error) {
	if limit <= 0 {
		limit = 100
	}

	var events []model.TaskEvent
	err := r.db.Where("task_id = ?", taskID).Order("created_at DESC").Limit(limit).Find(&events).Error
	return events, err
}

func (r *Repository) GetStats() (*model.Stats, error) {
	var stats model.Stats
	var totalAgents, activeAgents, tasksToday, pendingTasks, totalTasks, doneTasks int64

	if err := r.db.Model(&model.Agent{}).Count(&totalAgents).Error; err != nil {
		return nil, err
	}
	if err := r.db.Model(&model.Agent{}).Where("status = ?", "working").Count(&activeAgents).Error; err != nil {
		return nil, err
	}
	if err := r.db.Model(&model.Task{}).Where("DATE(created_at) = ?", time.Now().Format("2006-01-02")).Count(&tasksToday).Error; err != nil {
		return nil, err
	}
	if err := r.db.Model(&model.Task{}).Where("status = ?", "pending").Count(&pendingTasks).Error; err != nil {
		return nil, err
	}
	if err := r.db.Model(&model.Task{}).Select("COALESCE(SUM(tokens), 0)").Row().Scan(&stats.TokensUsed); err != nil {
		return nil, err
	}
	if err := r.db.Model(&model.Task{}).Count(&totalTasks).Error; err != nil {
		return nil, err
	}
	if err := r.db.Model(&model.Task{}).Where("status = ?", "done").Count(&doneTasks).Error; err != nil {
		return nil, err
	}

	stats.TotalAgents = int(totalAgents)
	stats.ActiveAgents = int(activeAgents)
	stats.TasksToday = int(tasksToday)
	stats.PendingTasks = int(pendingTasks)
	stats.TotalTasks = int(totalTasks)
	stats.DoneTasks = int(doneTasks)

	if stats.TotalTasks > 0 {
		stats.SuccessRate = float64(stats.DoneTasks) / float64(stats.TotalTasks) * 100
	}

	return &stats, nil
}

func NewTaskID() string {
	return fmt.Sprintf("task_%d", time.Now().UnixNano())
}

func NewTaskEventID() string {
	return fmt.Sprintf("te_%d", time.Now().UnixNano())
}
