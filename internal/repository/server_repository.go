package repository

import (
	"context"
	"fmt"

	"dc-analytics-service-backend/internal/models"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ServerRepository interface {
	GetAllServers(ctx context.Context, limit, offset int) ([]models.Server, error)
	GetServerByID(ctx context.Context, id int) (*models.Server, error)
	UpdateServer(ctx context.Context, id int, req models.UpdateServerRequest) error
	GetDevicesByServerID(ctx context.Context, serverID int) ([]models.Device, error)
}

type serverRepository struct {
	db *pgxpool.Pool
}

func NewServerRepository(db *pgxpool.Pool) ServerRepository {
	return &serverRepository{db: db}
}

func (r *serverRepository) GetAllServers(ctx context.Context, limit, offset int) ([]models.Server, error) {
	const query = `
        SELECT
            server_id,
            ip,
            cloud_name,
            cloud_type,
            cloud_device_type,
            cloud_status,
            cloud_state,
            created_at,
            updated_at,
            deleted_at
        FROM servers
        ORDER BY server_id
        LIMIT $1 OFFSET $2
    `

	var servers []models.Server
	if err := pgxscan.Select(ctx, r.db, &servers, query, limit, offset); err != nil {
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

func (r *serverRepository) UpdateServer(ctx context.Context, id int, req models.UpdateServerRequest) error {
	query := `
		UPDATE servers
		SET 
		  ip = COALESCE($1, ip),
		  cloud_name = COALESCE($2, cloud_name),
		  cloud_type = COALESCE($3, cloud_type),
		  cloud_status = COALESCE($4, cloud_status),
		  cloud_state = COALESCE($5, cloud_state),
		  updated_at = CURRENT_TIMESTAMP
		WHERE server_id = $6
	`
	result, err := r.db.Exec(ctx, query,
		req.IP,
		req.CloudName,
		req.CloudType,
		req.CloudStatus,
		req.CloudState,
		id,
	)
	if err != nil {
		return fmt.Errorf("ошибка обновления сервера: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("сервер с id %d не найден", id)
	}
	return nil
}

func (r *serverRepository) GetDevicesByServerID(
	ctx context.Context,
	serverID int,
) ([]models.Device, error) {
	const query = `
        SELECT
            id, smart_call_hiya, platform, serial, imei,
            number, carrier, priority, model, os_version,
            server, cloud, config, hub, port,
            created_at, updated_at, active, deleted_at,
            ui_version, build_number, baseband_version,
            sp_software_version, model_image_url
        FROM devices
        WHERE cloud = $1
        ORDER BY id
    `
	var devices []models.Device
	if err := pgxscan.Select(ctx, r.db, &devices, query, serverID); err != nil {
		return nil, fmt.Errorf("GetDevicesByServerID: %w", err)
	}
	return devices, nil
}
