package repository

import (
	"context"
	"fmt"

	"dc-analytics-service-backend/internal/models"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ServerRepository interface {
	GetAllServers(ctx context.Context) ([]models.Server, error)
	GetServerByID(ctx context.Context, id int) (*models.Server, error)
}

type serverRepository struct {
	db *pgxpool.Pool
}

func NewServerRepository(db *pgxpool.Pool) ServerRepository {
	return &serverRepository{db: db}
}

func (r *serverRepository) GetAllServers(ctx context.Context) ([]models.Server, error) {
	query := `SELECT server_id, ip, cloud_name, cloud_type, cloud_device_type, cloud_status, cloud_state, created_at, updated_at, deleted_at
			  FROM servers`
	var servers []models.Server
	if err := pgxscan.Select(ctx, r.db, &servers, query); err != nil {
		return nil, fmt.Errorf("ошибка получения серверов: %w", err)
	}
	return servers, nil
}

func (r *serverRepository) GetServerByID(ctx context.Context, id int) (*models.Server, error) {
	query := `SELECT server_id, ip, cloud_name, cloud_type, cloud_device_type, cloud_status, cloud_state, created_at, updated_at, deleted_at
			  FROM servers
			  WHERE server_id = $1`
	var server models.Server
	if err := pgxscan.Get(ctx, r.db, &server, query, id); err != nil {
		return nil, fmt.Errorf("сервер с id %d не найден: %w", id, err)
	}
	return &server, nil
}
