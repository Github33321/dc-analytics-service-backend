package models

import "time"

// postgres
type Device struct {
	ID                int64      `json:"id" db:"id"`
	SmartCallHiya     int        `json:"smart_call_hiya" db:"smart_call_hiya"`
	Platform          string     `json:"platform" db:"platform"`
	Serial            string     `json:"serial" db:"serial"`
	Imei              *string    `json:"imei,omitempty" db:"imei"`
	Number            string     `json:"number" db:"number"`
	Carrier           string     `json:"carrier" db:"carrier"`
	Priority          string     `json:"priority" db:"priority"`
	Model             string     `json:"model" db:"model"`
	OSVersion         *string    `json:"os_version,omitempty" db:"os_version"`
	Server            string     `json:"server" db:"server"`
	Cloud             int        `json:"cloud" db:"cloud"`
	Config            int        `json:"config" db:"config"`
	Hub               string     `json:"hub" db:"hub"`
	Port              int        `json:"port" db:"port"`
	CreatedAt         *time.Time `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt         *time.Time `json:"updated_at,omitempty" db:"updated_at"`
	Active            int        `json:"active" db:"active"`
	DeletedAt         *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
	UIVersion         *float64   `json:"ui_version" db:"ui_version"`
	BuildNumber       *string    `json:"build_number" db:"build_number"`
	BasebandVersion   *string    `json:"baseband_version" db:"baseband_version"`
	SPSoftwareVersion *string    `json:"sp_software_version" db:"sp_software_version"`
	ModelImageURL     string     `json:"model_image_url" db:"model_image_url"`
}
