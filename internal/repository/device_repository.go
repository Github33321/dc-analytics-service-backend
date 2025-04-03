package repository

import (
	"context"
	"dc-analytics-service-backend/internal/models"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type DeviceRepository interface {
	GetDevices(ctx context.Context) ([]models.Device, error)
	GetDeviceByID(ctx context.Context, id int64) (*models.Device, error)
	UpdateDevice(ctx context.Context, device *models.Device) (*models.Device, error)
	DeleteDevice(ctx context.Context, id int64) error
}

type deviceRepository struct {
	db *sqlx.DB
}

func NewDeviceRepository(db *sqlx.DB) DeviceRepository {
	return &deviceRepository{db: db}
}

func (r *deviceRepository) GetDevices(ctx context.Context) ([]models.Device, error) {
	query := `SELECT * FROM devices`
	var devices []models.Device
	err := r.db.SelectContext(ctx, &devices, query)
	if err != nil {
		return nil, err
	}
	return devices, nil
}

func (r *deviceRepository) GetDeviceByID(ctx context.Context, id int64) (*models.Device, error) {
	query := `SELECT * FROM devices WHERE id = $1`
	var device models.Device
	err := r.db.GetContext(ctx, &device, query, id)
	if err != nil {
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
	err := r.db.GetContext(ctx, device, query,
		device.SmartCallHiya, device.Platform, device.Serial, device.Imei, device.Number,
		device.Carrier, device.Priority, device.Model, device.OSVersion, device.Server,
		device.Cloud, device.Config, device.Hub, device.Port, device.Active,
		device.UIVersion, device.BuildNumber, device.BasebandVersion, device.SPSoftwareVersion,
		device.ID,
	)
	if err != nil {
		return nil, err
	}
	return device, nil
}

func (r *deviceRepository) DeleteDevice(ctx context.Context, id int64) error {
	query := `DELETE FROM devices WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("устройство с id %d не найдено", id)
	}

	return nil
}
