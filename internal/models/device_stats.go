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
