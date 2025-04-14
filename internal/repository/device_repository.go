package repository

import (
	"context"
	"dc-analytics-service-backend/internal/models"
	"fmt"
	"github.com/georgysavva/scany/v2/pgxscan"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DeviceRepository interface {
	GetDevices(ctx context.Context, limit, offset int) ([]models.Device, error)
	GetDeviceByID(ctx context.Context, id int64) (*models.Device, error)
	UpdateDevice(ctx context.Context, device *models.Device) (*models.Device, error)
	DeleteDevice(ctx context.Context, id int64) error
	GetDeviceStats(ctx context.Context) (models.DeviceStatsResponse, error)
	GetDevicesCount(ctx context.Context) (int64, error)
}

type deviceRepository struct {
	db *pgxpool.Pool
}

func NewDeviceRepository(db *pgxpool.Pool) DeviceRepository {
	return &deviceRepository{db: db}
}

func (r *deviceRepository) GetDevices(ctx context.Context, limit, offset int) ([]models.Device, error) {
	query := `
		SELECT * FROM devices 
		ORDER BY id 
		LIMIT $1 OFFSET $2
	`
	var devices []models.Device
	if err := pgxscan.Select(ctx, r.db, &devices, query, limit, offset); err != nil {
		return nil, err
	}
	return devices, nil
}

func (r *deviceRepository) GetDevicesCount(ctx context.Context) (int64, error) {
	query := `SELECT count(*) FROM devices`
	var count int64
	if err := pgxscan.Get(ctx, r.db, &count, query); err != nil {
		return 0, err
	}
	return count, nil
}

func (r *deviceRepository) GetDeviceByID(ctx context.Context, id int64) (*models.Device, error) {
	query := `SELECT * FROM devices WHERE id = $1`
	var device models.Device
	if err := pgxscan.Get(ctx, r.db, &device, query, id); err != nil {
		return nil, err
	}
	return &device, nil
}

func (r *deviceRepository) UpdateDevice(ctx context.Context, device *models.Device) (*models.Device, error) {
	query := `
		UPDATE devices
		SET smart_call_hiya = $1, platform = $2, serial = $3, imei = $4, number = $5,
			carrier = $6, priority = $7, model = $8, os_version = $9, server = $10,
			cloud = $11, config = $12, hub = $13, port = $14, active = $15,
			ui_version = $16, build_number = $17, baseband_version = $18, sp_software_version = $19,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $20
		RETURNING *
	`
	if err := pgxscan.Get(ctx, r.db, device, query,
		device.SmartCallHiya, device.Platform, device.Serial, device.Imei, device.Number,
		device.Carrier, device.Priority, device.Model, device.OSVersion, device.Server,
		device.Cloud, device.Config, device.Hub, device.Port, device.Active,
		device.UIVersion, device.BuildNumber, device.BasebandVersion, device.SPSoftwareVersion,
		device.ID,
	); err != nil {
		return nil, err
	}
	return device, nil
}

func (r *deviceRepository) DeleteDevice(ctx context.Context, id int64) error {
	query := `UPDATE devices SET deleted_at = CURRENT_TIMESTAMP WHERE id = $1`
	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("устройство с id %d не найдено", id)
	}
	return nil
}

func (r *deviceRepository) GetDeviceStats(ctx context.Context) (models.DeviceStatsResponse, error) {
	var stats models.DeviceStatsResponse

	query := `
		SELECT 
			count(*) AS total_count,
			sum(CASE WHEN lower(platform) = 'android' THEN 1 ELSE 0 END) AS android_count,
			sum(CASE WHEN lower(platform) = 'ios' THEN 1 ELSE 0 END) AS ios_count,
			sum(CASE WHEN lower(model) LIKE 'pixel%' THEN 1 ELSE 0 END) AS pixel_count,
			sum(CASE WHEN smart_call_hiya = 1 THEN 1 ELSE 0 END) AS smart_call_hiya_count
		FROM devices;
	`

	err := pgxscan.Get(ctx, r.db, &stats, query)
	if err != nil {
		return stats, err
	}
	return stats, nil
}
