package service

import (
	"context"
	"dc-analytics-service-backend/internal/models"
	"dc-analytics-service-backend/internal/repository"
	"fmt"
	"time"
)

type DeviceStatsService interface {
	GetTaskStats(ctx context.Context, date string) ([]models.TaskStat, error)
	GetCallStats(ctx context.Context, deviceID, date string) (models.DeviceCallStatsResponse, error)
	GetDeviceScreenshots(ctx context.Context, deviceID string, page, limit int) ([]models.DeviceScreenshot, error)
	GetDeviceCarrierStats(ctx context.Context) ([]models.DeviceCarrierStats, error)
	GetOriginatingCarrierStats(ctx context.Context, fromDate string) (models.CarrierStatsResponse, error)
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

func (s *deviceStatsService) GetDeviceCarrierStats(ctx context.Context) ([]models.DeviceCarrierStats, error) {
	return s.repo.GetDeviceCarrierStats(ctx)
}

func (s *deviceStatsService) GetOriginatingCarrierStats(ctx context.Context, fromDate string) (models.CarrierStatsResponse, error) {
	parsedDate, err := time.Parse("2006-01-02", fromDate)
	if err != nil {
		return models.CarrierStatsResponse{}, fmt.Errorf("некорректный формат даты: %w", err)
	}
	return s.repo.GetOriginatingCarrierStats(ctx, parsedDate)
}
