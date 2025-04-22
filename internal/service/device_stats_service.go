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
	GetDeviceCallStats(ctx context.Context, deviceID, date string) (models.DeviceCallStatsResponse, error)
	GetDeviceScreenshots(ctx context.Context, deviceID string, page, limit int) ([]models.DeviceScreenshot, error)
	GetDeviceCarrierStats(ctx context.Context) ([]models.DeviceCarrierStats, error)
	GetOriginatingCarrierStats(ctx context.Context, fromDate string) (models.CarrierStatsResponse, error)
	GetDeviceGroupStats(ctx context.Context) ([]models.DeviceGroupStats, error)
	// Имитация ответа от внешнего API
	GetSources(ctx context.Context) ([]*models.Source, error)
	GetTasksReadyCounts(ctx context.Context) ([]models.TasksReadyCount, error)
	GetByUserID(ctx context.Context, userID uint64) ([]models.DedicatedDevice, error)
	GetCountedUsers(ctx context.Context) ([]models.DCUser, error)
	TodayStats(ctx context.Context) ([]models.UserGroupStats, error)
	GetDistinctUsers(ctx context.Context, date string) ([]models.DCUser2, error)
}

type deviceStatsService struct {
	repo repository.DeviceStatsRepository
}

func NewDeviceStatsService(repo repository.DeviceStatsRepository) DeviceStatsService {
	return &deviceStatsService{repo: repo}
}

func (s *deviceStatsService) GetDeviceCallStats(ctx context.Context, deviceID, date string) (models.DeviceCallStatsResponse, error) {
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
func (s *deviceStatsService) GetDeviceGroupStats(ctx context.Context) ([]models.DeviceGroupStats, error) {
	return s.repo.GetDeviceGroupStats(ctx)
}

// Имитация ответа от внешнего API
func (s *deviceStatsService) GetSources(ctx context.Context) ([]*models.Source, error) {
	return []*models.Source{
		{ID: 1, Name: "daily-platform", CreatedAt: "2024-12-26 10:09:05", UpdatedAt: "2024-12-26 10:09:05"},
		{ID: 2, Name: "demand-platform", CreatedAt: "2024-12-26 10:09:05", UpdatedAt: "2024-12-26 10:09:05"},
		{ID: 3, Name: "partner-platform", CreatedAt: "2024-12-26 10:09:05", UpdatedAt: "2024-12-26 10:09:05"},
		{ID: 4, Name: "admin-platform", CreatedAt: "2024-12-26 10:09:05", UpdatedAt: "2024-12-26 10:09:05"},
		{ID: 5, Name: "user-api", CreatedAt: "2024-12-26 10:09:05", UpdatedAt: "2024-12-26 10:09:05"},
		{ID: 6, Name: "daily-hiya-platform", CreatedAt: "2025-01-20 13:35:01", UpdatedAt: "2025-01-28 09:07:23"},
		{ID: 7, Name: "daily-pixel-platform", CreatedAt: "2025-04-08 10:08:32", UpdatedAt: "2025-04-08 10:08:32"},
	}, nil
}

func (s *deviceStatsService) GetTasksReadyCounts(ctx context.Context) ([]models.TasksReadyCount, error) {
	return s.repo.GetTasksReadyCounts(ctx)
}

func (s *deviceStatsService) GetByUserID(ctx context.Context, userID uint64) ([]models.DedicatedDevice, error) {
	return s.repo.GetByUserID(ctx, userID)
}

func (s *deviceStatsService) GetCountedUsers(ctx context.Context) ([]models.DCUser, error) {
	return s.repo.GetCountedUsers(ctx)
}

func (s *deviceStatsService) TodayStats(ctx context.Context) ([]models.UserGroupStats, error) {
	return s.repo.GetTodayStats(ctx)
}

func (s *deviceStatsService) GetDistinctUsers(ctx context.Context, date string) ([]models.DCUser2, error) {
	return s.repo.GetDistinctUsers(ctx, date)
}
