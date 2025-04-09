package service

import (
	"context"
	"dc-analytics-service-backend/internal/models"
	"dc-analytics-service-backend/internal/repository"
)

type deviceService struct {
	deviceRepo repository.DeviceRepository
}

func NewDeviceService(deviceRepo repository.DeviceRepository) DeviceService {
	return &deviceService{deviceRepo: deviceRepo}
}

func (s *deviceService) GetDevices(ctx context.Context, page, limit int) ([]models.Device, error) {
	offset := (page - 1) * limit
	return s.deviceRepo.GetDevices(ctx, limit, offset)
}

func (s *deviceService) GetDeviceByID(ctx context.Context, id int64) (*models.Device, error) {
	return s.deviceRepo.GetDeviceByID(ctx, id)
}

func (s *deviceService) UpdateDevice(ctx context.Context, id int64, req UpdateDeviceRequest) (*models.Device, error) {
	device, err := s.deviceRepo.GetDeviceByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if device == nil {
		return nil, nil
	}

	if req.SmartCallHiya != nil {
		device.SmartCallHiya = *req.SmartCallHiya
	}
	if req.Platform != nil {
		device.Platform = *req.Platform
	}
	if req.Serial != nil {
		device.Serial = *req.Serial
	}
	if req.Imei != nil {
		device.Imei = req.Imei
	}
	if req.Number != nil {
		device.Number = *req.Number
	}
	if req.Carrier != nil {
		device.Carrier = *req.Carrier
	}
	if req.Priority != nil {
		device.Priority = *req.Priority
	}
	if req.Model != nil {
		device.Model = *req.Model
	}
	if req.OSVersion != nil {
		device.OSVersion = req.OSVersion
	}
	if req.Server != nil {
		device.Server = *req.Server
	}
	if req.Cloud != nil {
		device.Cloud = *req.Cloud
	}
	if req.Config != nil {
		device.Config = *req.Config
	}
	if req.Hub != nil {
		device.Hub = *req.Hub
	}
	if req.Port != nil {
		device.Port = *req.Port
	}
	if req.Active != nil {
		device.Active = *req.Active
	}
	if req.UIVersion != nil {
		device.UIVersion = req.UIVersion
	}
	if req.BuildNumber != nil {
		device.BuildNumber = req.BuildNumber
	}
	if req.BasebandVersion != nil {
		device.BasebandVersion = req.BasebandVersion
	}
	if req.SPSoftwareVersion != nil {
		device.SPSoftwareVersion = req.SPSoftwareVersion
	}

	return s.deviceRepo.UpdateDevice(ctx, device)
}

func (s *deviceService) DeleteDevice(ctx context.Context, id int64) error {
	return s.deviceRepo.DeleteDevice(ctx, id)
}

func (s *deviceService) GetDeviceStats(ctx context.Context) (models.DeviceStatsResponse, error) {
	return s.deviceRepo.GetDeviceStats(ctx)
}
