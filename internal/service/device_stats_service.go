package service

import (
	"context"
	"dc-analytics-service-backend/internal/models"
	"dc-analytics-service-backend/internal/repository"
)

type DeviceStatsService interface {
	GetTaskStats(ctx context.Context, date string) ([]models.TaskStat, error)
	GetCallStats(ctx context.Context, deviceID, date string) (models.DeviceCallStatsResponse, error)
}

type deviceStatsService struct {
	repo repository.DeviceStatsRepository
}

func NewDeviceStatsService(repo repository.DeviceStatsRepository) DeviceStatsService {
	return &deviceStatsService{repo: repo}
}

func (s *deviceStatsService) GetCallStats(ctx context.Context, deviceID, date string) (models.DeviceCallStatsResponse, error) {
	return s.repo.GetDeviceCallStats(ctx, deviceID, date)
}

func (s *deviceStatsService) GetTaskStats(ctx context.Context, date string) ([]models.TaskStat, error) {
	return s.repo.GetTaskStats(ctx, date)
}
