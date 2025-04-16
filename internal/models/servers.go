package models

import "time"

// postgres
type Server struct {
	ServerID        int        `json:"server_id" db:"server_id"`
	IP              string     `json:"ip" db:"ip"`
	CloudName       string     `json:"cloud_name" db:"cloud_name"`
	CloudType       string     `json:"cloud_type" db:"cloud_type"`
	CloudDeviceType *string    `json:"cloud_device_type" db:"cloud_device_type"`
	CloudStatus     string     `json:"cloud_status" db:"cloud_status"`
	CloudState      string     `json:"cloud_state" db:"cloud_state"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt       *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

type UpdateServerRequest struct {
	IP          *string `json:"ip,omitempty"`
	CloudName   *string `json:"cloud_name,omitempty"`
	CloudType   *string `json:"cloud_type,omitempty"`
	CloudStatus *string `json:"cloud_status,omitempty"`
	CloudState  *string `json:"cloud_state,omitempty"`
}
