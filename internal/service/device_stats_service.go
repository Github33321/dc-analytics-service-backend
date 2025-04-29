package service

import (
	"context"
	"dc-analytics-service-backend/internal/config"
	"dc-analytics-service-backend/internal/models"
	"dc-analytics-service-backend/internal/repository"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const sourcesURL = "https://ms1.calleridreputation.com/v1/sources/list"
const authKey = "4f9c5e8f-2c8b-4f7d-91d8-1a2b3c4d5e6f"

type DeviceStatsService interface {
	GetTaskStats(ctx context.Context, date string) ([]models.TaskStat, error)
	GetDeviceCallStats(ctx context.Context, deviceID, date string) (*models.DeviceCallStatsResponse, error)
	GetDeviceScreenshots(ctx context.Context, deviceID string, page, limit int) ([]models.DeviceScreenshot, error)
	GetDeviceCarrierStats(ctx context.Context) ([]*models.DeviceCarrierStats, error)
	GetOriginatingCarrierStats(ctx context.Context, fromDate string) (*models.CarrierStatsResponse, error)
	GetDeviceGroupStats(ctx context.Context) (*[]models.DeviceGroupStats, error)
	GetTasksReadyCounts(ctx context.Context) (*[]models.TasksReadyCount, error)
	GetByUserID(ctx context.Context, userID uint64) (*[]models.DedicatedDevice, error)
	GetCountedUsers(ctx context.Context) (*[]models.DCUser, error)
	TodayStats(ctx context.Context) (*[]models.UserGroupStats, error)
	GetDistinctUsers(ctx context.Context, date string) (*[]models.DCUser2, error)
	GetSourcesStats(ctx context.Context) (*[]models.SourceStatResponse, error)
}

type deviceStatsService struct {
	repo   repository.DeviceStatsRepository
	config *config.Config
}

func NewDeviceStatsService(repo repository.DeviceStatsRepository, cfg *config.Config) DeviceStatsService {
	return &deviceStatsService{
		repo:   repo,
		config: cfg,
	}
}

func (s *deviceStatsService) GetDeviceCallStats(ctx context.Context, deviceID, date string) (*models.DeviceCallStatsResponse, error) {
	return s.repo.GetDeviceCallStats(ctx, deviceID, date)
}

func (s *deviceStatsService) GetTaskStats(ctx context.Context, date string) ([]models.TaskStat, error) {
	return s.repo.GetTaskStats(ctx, date)
}

func (s *deviceStatsService) GetDeviceScreenshots(ctx context.Context, deviceID string, page, limit int) ([]models.DeviceScreenshot, error) {
	offset := (page - 1) * limit
	return s.repo.GetDeviceScreenshots(ctx, deviceID, limit, offset)
}

func (s *deviceStatsService) GetDeviceCarrierStats(ctx context.Context) ([]*models.DeviceCarrierStats, error) {
	return s.repo.GetDeviceCarrierStats(ctx)
}

func (s *deviceStatsService) GetOriginatingCarrierStats(ctx context.Context, fromDate string) (*models.CarrierStatsResponse, error) {
	parsedDate, err := time.Parse("2006-01-02", fromDate)
	if err != nil {
		return nil, fmt.Errorf("некорректный формат даты: %w", err)
	}
	return s.repo.GetOriginatingCarrierStats(ctx, parsedDate)
}
func (s *deviceStatsService) GetDeviceGroupStats(ctx context.Context) (*[]models.DeviceGroupStats, error) {
	return s.repo.GetDeviceGroupStats(ctx)
}

func (s *deviceStatsService) GetTasksReadyCounts(ctx context.Context) (*[]models.TasksReadyCount, error) {
	return s.repo.GetTasksReadyCounts(ctx)
}

func (s *deviceStatsService) GetByUserID(ctx context.Context, userID uint64) (*[]models.DedicatedDevice, error) {
	return s.repo.GetByUserID(ctx, userID)
}

func (s *deviceStatsService) GetCountedUsers(ctx context.Context) (*[]models.DCUser, error) {
	return s.repo.GetCountedUsers(ctx)
}

func (s *deviceStatsService) TodayStats(ctx context.Context) (*[]models.UserGroupStats, error) {
	return s.repo.GetTodayStats(ctx)
}

func (s *deviceStatsService) GetDistinctUsers(ctx context.Context, date string) (*[]models.DCUser2, error) {
	return s.repo.GetDistinctUsers(ctx, date)
}

func (s *deviceStatsService) GetSourcesStats(ctx context.Context) (*[]models.SourceStatResponse, error) {
	stats, err := s.repo.GetSourceTypeStats(ctx)
	if err != nil {
		return nil, err
	}

	sources, err := s.GetSources(ctx)
	if err != nil {
		return nil, err
	}

	sourceNameMap := make(map[int]string)
	for _, src := range sources {
		sourceNameMap[src.ID] = src.Name
	}

	var response []models.SourceStatResponse
	for _, stat := range *stats {
		name, ok := sourceNameMap[int(stat.SourceTypeID)]
		if !ok {
			name = "Unknown"
		}
		response = append(response, models.SourceStatResponse{
			SourceTypeID: int(stat.SourceTypeID),
			SourceName:   name,
			Count:        int(stat.Count),
		})
	}

	return &response, nil
}

func (s *deviceStatsService) GetSources(ctx context.Context) ([]models.Source, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.config.ExternalAPIs.SourcesAPIURL, nil)
	if err != nil {
		return nil, fmt.Errorf("создание запроса GetSources: %w", err)
	}
	req.Header.Set("X-Auth-Key", s.config.ExternalAPIs.SourcesAPIKey)

	client := &http.Client{Timeout: 5 * time.Second}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса GetSources: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("неожиданный статус код GetSources: %d", resp.StatusCode)
	}

	var sourcesResp models.SourcesResponse
	if err := json.NewDecoder(resp.Body).Decode(&sourcesResp); err != nil {
		return nil, fmt.Errorf("ошибка декодирования ответа GetSources: %w", err)
	}

	return sourcesResp.Content, nil
}
