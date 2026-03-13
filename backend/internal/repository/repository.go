package repository

import (
	"time"

	"gorm.io/gorm"

	"github.com/MaiWittawat/openclaw-empire/internal/model"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetAgents() ([]model.Agent, error) {
	var agents []model.Agent
	err := r.db.Order("name").Find(&agents).Error
	return agents, err
}

func (r *Repository) GetTasks() ([]model.Task, error) {
	var tasks []model.Task
	err := r.db.Preload("Agent").Order("created_at DESC").Limit(50).Find(&tasks).Error
	return tasks, err
}

func (r *Repository) CreateTask(title, agentID string) (*model.Task, error) {
	task := &model.Task{
		Title:   title,
		AgentID: agentID,
		Status:  "pending",
		Tokens:  0,
	}
	err := r.db.Create(task).Error
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (r *Repository) GetStats() (*model.Stats, error) {
	var stats model.Stats
	var totalAgents, activeAgents, tasksToday, pendingTasks, totalTasks, doneTasks int64

	r.db.Model(&model.Agent{}).Count(&totalAgents)
	r.db.Model(&model.Agent{}).Where("status = ?", "working").Count(&activeAgents)
	r.db.Model(&model.Task{}).Where("DATE(created_at) = ?", time.Now().Format("2006-01-02")).Count(&tasksToday)
	r.db.Model(&model.Task{}).Where("status = ?", "pending").Count(&pendingTasks)
	r.db.Model(&model.Task{}).Select("COALESCE(SUM(tokens), 0)").Row().Scan(&stats.TokensUsed)
	r.db.Model(&model.Task{}).Count(&totalTasks)
	r.db.Model(&model.Task{}).Where("status = ?", "done").Count(&doneTasks)

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

func (r *Repository) SeedData() error {
	err := r.db.AutoMigrate(&model.Agent{}, &model.Task{})
	if err != nil {
		return err
	}

	var count int64
	r.db.Model(&model.Agent{}).Count(&count)
	if count == 0 {
		agents := []model.Agent{
			{ID: "parae", Name: "แพร", Role: "Tech Lead", Avatar: "👩‍💻", Status: "working", Color: "#a78bfa", Tokens: 12450, MaxTokens: 20000, UpdatedAt: time.Now()},
			{ID: "nova", Name: "NOVA", Role: "Researcher", Avatar: "🔍", Status: "working", Color: "#34d399", Tokens: 8820, MaxTokens: 20000, UpdatedAt: time.Now()},
			{ID: "forge", Name: "FORGE", Role: "DevOps", Avatar: "🛠️", Status: "idle", Color: "#60a5fa", Tokens: 3200, MaxTokens: 20000, UpdatedAt: time.Now()},
			{ID: "lyra", Name: "LYRA", Role: "Writer", Avatar: "✍️", Status: "done", Color: "#f472b6", Tokens: 5400, MaxTokens: 20000, UpdatedAt: time.Now()},
		}
		r.db.Create(&agents)

		tasks := []model.Task{
			{Title: "เขียน FastAPI endpoint สำหรับ Telegram webhook", AgentID: "parae", Status: "working", Tokens: 12450, CreatedAt: time.Now().Add(-2 * time.Minute), UpdatedAt: time.Now()},
			{Title: "ค้นหา best practice สำหรับ multi-agent routing", AgentID: "nova", Status: "working", Tokens: 8820, CreatedAt: time.Now().Add(-5 * time.Minute), UpdatedAt: time.Now()},
			{Title: "เขียน README สำหรับ project", AgentID: "lyra", Status: "done", Tokens: 5400, CreatedAt: time.Now().Add(-18 * time.Minute), UpdatedAt: time.Now()},
			{Title: "สร้าง Docker Compose config", AgentID: "forge", Status: "done", Tokens: 3200, CreatedAt: time.Now().Add(-1 * time.Hour), UpdatedAt: time.Now()},
			{Title: "Deploy to VPS — timeout error", AgentID: "forge", Status: "error", Tokens: 1100, CreatedAt: time.Now().Add(-2 * time.Hour), UpdatedAt: time.Now()},
		}
		r.db.Create(&tasks)
	}
	return nil
}
