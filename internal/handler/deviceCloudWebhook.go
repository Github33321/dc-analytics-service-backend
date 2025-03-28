package handler

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type DeviceCloudWebhook struct {
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

var db *sql.DB

func SetDB(database *sql.DB) {
	db = database
}

func SelectHandler(c *gin.Context) {
	query := `
		SELECT 
			_id,
			created_at,
			updated_at,
			MessageId,
			from_num,
			originating_carrier,
			created_at_str,
			device_os,
			device_carrier,
			status,
			device_id,
			device_config_id,
			user_id,
			group_id,
			hiya_enabled,
			source_type_id,
			notification_url,
			call_duration,
			call_end,
			call_start,
			cnam,
			device_model,
			display_cnam,
			hiya,
			incoming_number,
			incoming_number_match,
			log_cnam,
			ocr_cloud_id,
			screenshot,
			spam,
			text,
			text_recognized,
			to_num
		FROM device_cloud_webhooks
		LIMIT 10
	`

	rows, err := db.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var results []DeviceCloudWebhook

	for rows.Next() {
		var webhook DeviceCloudWebhook
		err := rows.Scan(
			&webhook.ID,
			&webhook.CreatedAt,
			&webhook.UpdatedAt,
			&webhook.MessageID,
			&webhook.FromNum,
			&webhook.OriginatingCarrier,
			&webhook.CreatedAtStr,
			&webhook.DeviceOS,
			&webhook.DeviceCarrier,
			&webhook.Status,
			&webhook.DeviceID,
			&webhook.DeviceConfigID,
			&webhook.UserID,
			&webhook.GroupID,
			&webhook.HiyaEnabled,
			&webhook.SourceTypeID,
			&webhook.NotificationURL,
			&webhook.CallDuration,
			&webhook.CallEnd,
			&webhook.CallStart,
			&webhook.Cnam,
			&webhook.DeviceModel,
			&webhook.DisplayCnam,
			&webhook.Hiya,
			&webhook.IncomingNumber,
			&webhook.IncomingNumberMatch,
			&webhook.LogCnam,
			&webhook.OcrCloudID,
			&webhook.Screenshot,
			&webhook.Spam,
			&webhook.Text,
			&webhook.TextRecognized,
			&webhook.ToNum,
		)
		if err != nil {
			log.Printf("Ошибка при чтении строки: %v", err)
			continue
		}
		results = append(results, webhook)
	}

	if err = rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}
