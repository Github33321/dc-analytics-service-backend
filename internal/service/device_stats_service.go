package service

import (
	"context"
	"dc-analytics-service-backend/internal/models"
	"dc-analytics-service-backend/internal/repository"
)

type DeviceStatsService interface {
	GetTaskStats(ctx context.Context, date string) ([]models.TaskStat, error)
	GetCallStats(ctx context.Context, deviceID, date string) (models.DeviceCallStatsResponse, error)
	GetDeviceScreenshots(ctx context.Context, deviceID string, page, limit int) ([]models.DeviceScreenshot, error)
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

func (s *deviceStatsService) GetDeviceScreenshots(ctx context.Context, deviceID string, page, limit int) ([]models.DeviceScreenshot, error) {
	offset := (page - 1) * limit
	return s.repo.GetDeviceScreenshots(ctx, deviceID, limit, offset)
}
