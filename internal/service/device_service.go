package service

import (
	"context"
	"dc-analytics-service-backend/internal/models"
)

type DeviceService interface {
	GetDevices(ctx context.Context, page, limit int) (models.PaginatedDevices, error)
	GetDeviceByID(ctx context.Context, id int64) (*models.Device, error)
	UpdateDevice(ctx context.Context, id int64, req UpdateDeviceRequest) (*models.Device, error)
	DeleteDevice(ctx context.Context, id int64) error
	GetDeviceStats(ctx context.Context) (models.DeviceStatsResponse, error)
}

type UpdateDeviceRequest struct {
	SmartCallHiya     *int     `json:"smart_call_hiya,omitempty"`
	Platform          *string  `json:"platform,omitempty"`
	Serial            *string  `json:"serial,omitempty"`
	Imei              *string  `json:"imei,omitempty"`
	Number            *string  `json:"number,omitempty"`
	Carrier           *string  `json:"carrier,omitempty"`
	Priority          *string  `json:"priority,omitempty"`
	Model             *string  `json:"model,omitempty"`
	OSVersion         *string  `json:"os_version,omitempty"`
	Server            *string  `json:"server,omitempty"`
	Cloud             *int     `json:"cloud,omitempty"`
	Config            *int     `json:"config,omitempty"`
	Hub               *string  `json:"hub,omitempty"`
	Port              *int     `json:"port,omitempty"`
	Active            *int     `json:"active,omitempty"`
	UIVersion         *float64 `json:"ui_version,omitempty"`
	BuildNumber       *string  `json:"build_number,omitempty"`
	BasebandVersion   *string  `json:"baseband_version,omitempty"`
	SPSoftwareVersion *string  `json:"sp_software_version,omitempty"`
}
