package models

import "time"

type ErrorResponse struct {
	Error string `json:"error" example:"error"`
}

type MessageResponse struct {
	Message string `json:"message" example:"message"`
}

type TaskStat struct {
	CreatedAtStr time.Time `json:"created_at_str"`
	Count        uint64    `json:"count"`
}

type StatusCount struct {
	Status string `json:"status"`
	Count  uint64 `json:"count"`
}

type DeviceCallStatsResponse struct {
	TodayCalls   uint64        `json:"today_calls"`
	CallsByDay   []TaskStat    `json:"calls_by_day"`
	StatusCounts []StatusCount `json:"status_counts"`
}

type DeviceStatsResponse struct {
	TotalCount         int64 `json:"total_count"`
	AndroidCount       int64 `json:"android_count"`
	IOSCount           int64 `json:"ios_count"`
	PixelCount         int64 `json:"pixel_count"`
	SmartCallHiyaCount int64 `json:"smart_call_hiya_count"`
}

type DeviceScreenshot struct {
	CreatedAt  string `json:"created_at"`
	Screenshot string `json:"screenshot"`
}
type PaginatedDevices struct {
	Devices    []Device `json:"devices"`
	TotalPages int      `json:"total_pages"`
}
type DeviceCarrierStats struct {
	DeviceCarrier string `json:"device_carrier"`
	DeviceCount   uint64 `json:"device_count"`
}

type CarrierStatsResponse struct {
	Stats []CarrierStat `json:"stats"`
}

type CarrierStat struct {
	Date               time.Time `json:"date"`
	OriginatingCarrier string    `json:"originating_carrier"`
	Count              uint64    `json:"count"`
}

type DeviceGroupStats struct {
	GroupID uint64 `json:"group_id"`
	Count   uint64 `json:"count"`
}

// Врушка заглушка
type Source struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type TasksReadyCount struct {
	ServerGroupID uint64 `json:"server_group_id"`
	Count         uint64 `json:"count"`
}

type DedicatedDevice struct {
	UserID    uint64    `json:"user_id"`
	DeviceID  uint64    `json:"device_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type DCUser struct {
	UserID uint64 `json:"user_id"`
	Count  uint64 `json:"count"`
}

type UserGroupStats struct {
	UserID      uint64 `json:"user_id" db:"user_id"`
	GroupID     uint64 `json:"group_id" db:"group_id"`
	ChecksCount uint64 `json:"checks_count" db:"checks_count"`
}

type UserCheckStats struct {
	UserID uint64 `json:"user_id" db:"user_id"`
	Count  uint64 `json:"count"    db:"count"`
}

type DCUser2 struct {
	UserID uint64 `json:"user_id" db:"user_id"`
}
