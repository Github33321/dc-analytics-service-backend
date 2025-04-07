package service

import (
	"context"
	"dc-analytics-service-backend/internal/models"
	"dc-analytics-service-backend/internal/repository"
)

type ClickhouseService interface {
	GetResults(ctx context.Context) ([]models.DeviceCloudResult, error)
}

type clickhouseService struct {
	repo repository.IClickhouse
}

func NewClickhouseService(repo repository.IClickhouse) ClickhouseService {
	return &clickhouseService{repo: repo}
}

func (s *clickhouseService) GetResults(ctx context.Context) ([]models.DeviceCloudResult, error) {
	query := `
		SELECT 
			_id,
			created_at,
			updated_at,
			toString(MessageId) as MessageId,
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
			toInt8(incoming_number_match) as incoming_number_match,
			log_cnam,
			ocr_cloud_id,
			screenshot,
			toInt8(spam) as spam,
			text,
			toInt8(text_recognized) as text_recognized,
			to_num
		FROM device_cloud_webhooks
		LIMIT 10
	`
	rows, err := s.repo.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.DeviceCloudResult
	for rows.Next() {
		var result models.DeviceCloudResult
		var incomingNumberMatch, spam, textRecognized uint8
		err := rows.Scan(
			&result.ID,
			&result.CreatedAt,
			&result.UpdatedAt,
			&result.MessageID,
			&result.FromNum,
			&result.OriginatingCarrier,
			&result.CreatedAtStr,
			&result.DeviceOS,
			&result.DeviceCarrier,
			&result.Status,
			&result.DeviceID,
			&result.DeviceConfigID,
			&result.UserID,
			&result.GroupID,
			&result.HiyaEnabled,
			&result.SourceTypeID,
			&result.NotificationURL,
			&result.CallDuration,
			&result.CallEnd,
			&result.CallStart,
			&result.Cnam,
			&result.DeviceModel,
			&result.DisplayCnam,
			&result.Hiya,
			&result.IncomingNumber,
			&incomingNumberMatch,
			&result.LogCnam,
			&result.OcrCloudID,
			&result.Screenshot,
			&spam,
			&result.Text,
			&textRecognized,
			&result.ToNum,
		)
		if err != nil {
			return nil, err
		}
		result.IncomingNumberMatch = incomingNumberMatch == 1
		result.Spam = spam == 1
		result.TextRecognized = textRecognized == 1

		results = append(results, result)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}
