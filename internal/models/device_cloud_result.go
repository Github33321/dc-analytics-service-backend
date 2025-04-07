package models

type DeviceCloudResult struct {
	ID                  string `json:"_id"`
	CreatedAt           string `json:"created_at"`
	UpdatedAt           string `json:"updated_at"`
	MessageID           string `json:"MessageId"`
	FromNum             string `json:"from_num"`
	OriginatingCarrier  string `json:"originating_carrier"`
	CreatedAtStr        string `json:"created_at_str"`
	DeviceOS            string `json:"device_os"`
	DeviceCarrier       string `json:"device_carrier"`
	Status              string `json:"status"`
	DeviceID            uint64 `json:"device_id"`
	DeviceConfigID      uint64 `json:"device_config_id"`
	UserID              uint64 `json:"user_id"`
	GroupID             uint64 `json:"group_id"`
	HiyaEnabled         uint8  `json:"hiya_enabled"`
	SourceTypeID        uint8  `json:"source_type_id"`
	NotificationURL     string `json:"notification_url"`
	CallDuration        uint64 `json:"call_duration"`
	CallEnd             string `json:"call_end"`
	CallStart           string `json:"call_start"`
	Cnam                string `json:"cnam"`
	DeviceModel         string `json:"device_model"`
	DisplayCnam         string `json:"display_cnam"`
	Hiya                uint8  `json:"hiya"`
	IncomingNumber      string `json:"incoming_number"`
	IncomingNumberMatch bool   `json:"incoming_number_match"`
	LogCnam             string `json:"log_cnam"`
	OcrCloudID          uint64 `json:"ocr_cloud_id"`
	Screenshot          string `json:"screenshot"`
	Spam                bool   `json:"spam"`
	Text                string `json:"text"`
	TextRecognized      bool   `json:"text_recognized"`
	ToNum               string `json:"to_num"`
}
