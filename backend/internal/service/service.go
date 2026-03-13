package service

import (
	"github.com/MaiWittawat/openclaw-empire/internal/model"
	"github.com/MaiWittawat/openclaw-empire/internal/repository"
)

type Service struct {
	repo *repository.Repository
}

func NewService(repo *repository.Repository) *Service {
	s := &Service{repo: repo}
	s.repo.SeedData()
	return s
}

func (s *Service) GetAgents() ([]model.Agent, error) {
	return s.repo.GetAgents()
}

func (s *Service) GetTasks() ([]model.Task, error) {
	return s.repo.GetTasks()
}

func (s *Service) CreateTask(title, agentID string) (*model.Task, error) {
	return s.repo.CreateTask(title, agentID)
}

func (s *Service) GetStats() (*model.Stats, error) {
	return s.repo.GetStats()
}
