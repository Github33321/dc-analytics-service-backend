package service

import (
	"context"
	"dc-analytics-service-backend/internal/models"
	"dc-analytics-service-backend/internal/repository"
)

type TaskStatService interface {
	GetTaskStats(ctx context.Context, date string) ([]models.TaskStat, error)
}

type taskStatService struct {
	repo repository.TaskStatRepository
}

func NewTaskStatService(repo repository.TaskStatRepository) TaskStatService {
	return &taskStatService{repo: repo}
}

func (s *taskStatService) GetTaskStats(ctx context.Context, date string) ([]models.TaskStat, error) {
	return s.repo.GetTaskStats(ctx, date)
}
