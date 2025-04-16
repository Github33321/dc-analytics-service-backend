package service

import (
	"context"

	"dc-analytics-service-backend/internal/models"
	"dc-analytics-service-backend/internal/repository"
)

type ServerService interface {
	GetAllServers(ctx context.Context, limit, offset int) ([]models.Server, error)
	GetServerByID(ctx context.Context, id int) (*models.Server, error)
	UpdateServer(ctx context.Context, id int, req models.UpdateServerRequest) (*models.Server, error)
	GetDevicesByServerID(ctx context.Context, serverID int) ([]models.Device, error)
}

type serverService struct {
	repo repository.ServerRepository
}

func NewServerService(repo repository.ServerRepository) ServerService {
	return &serverService{repo: repo}
}

func (s *serverService) GetAllServers(
	ctx context.Context,
	limit, offset int,
) ([]models.Server, error) {
	return s.repo.GetAllServers(ctx, limit, offset)
}

func (s *serverService) GetServerByID(ctx context.Context, id int) (*models.Server, error) {
	return s.repo.GetServerByID(ctx, id)
}
func (s *serverService) UpdateServer(ctx context.Context, id int, req models.UpdateServerRequest) (*models.Server, error) {
	if err := s.repo.UpdateServer(ctx, id, req); err != nil {
		return nil, err
	}
	return s.repo.GetServerByID(ctx, id)
}

func (s *serverService) GetDevicesByServerID(
	ctx context.Context,
	serverID int,
) ([]models.Device, error) {
	return s.repo.GetDevicesByServerID(ctx, serverID)
}
