package service

import (
	"context"

	"dc-analytics-service-backend/internal/models"
	"dc-analytics-service-backend/internal/repository"
)

type ServerService interface {
	GetAllServers(ctx context.Context) ([]models.Server, error)
	GetServerByID(ctx context.Context, id int) (*models.Server, error)
}

type serverService struct {
	repo repository.ServerRepository
}

func NewServerService(repo repository.ServerRepository) ServerService {
	return &serverService{repo: repo}
}

func (s *serverService) GetAllServers(ctx context.Context) ([]models.Server, error) {
	return s.repo.GetAllServers(ctx)
}

func (s *serverService) GetServerByID(ctx context.Context, id int) (*models.Server, error) {
	return s.repo.GetServerByID(ctx, id)
}
