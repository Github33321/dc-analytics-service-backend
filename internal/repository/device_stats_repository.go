package repository

import (
	"context"
	"dc-analytics-service-backend/internal/models"
	"fmt"
	"time"
)

type DeviceStatsRepository interface {
	GetTaskStats(ctx context.Context, date string) ([]models.TaskStat, error)
	GetDeviceCallStats(ctx context.Context, deviceID string, date string) (models.DeviceCallStatsResponse, error)
	GetDeviceScreenshots(ctx context.Context, deviceID string, limit, offset int) ([]models.DeviceScreenshot, error)
	GetDeviceCarrierStats(ctx context.Context) ([]models.DeviceCarrierStats, error)
	GetOriginatingCarrierStats(ctx context.Context, fromDate time.Time) (models.CarrierStatsResponse, error)
}

type deviceStatsRepo struct {
	ch IClickhouse
}

func NewDeviceStatsRepository(ch IClickhouse) DeviceStatsRepository {
	return &deviceStatsRepo{ch: ch}
}

func (r *deviceStatsRepo) GetDeviceCallStats(ctx context.Context, deviceID string, date string) (models.DeviceCallStatsResponse, error) {
	var resp models.DeviceCallStatsResponse

	today := time.Now().UTC().Format("2006-01-02")

	todayQuery := fmt.Sprintf(`
		SELECT count(*)
		FROM device_cloud_webhooks
		WHERE toString(device_id) = '%s'
		  AND created_at_str = '%s'
	`, deviceID, today)

	rows, err := r.ch.Query(ctx, todayQuery)
	if err != nil {
		return resp, err
	}
	if rows.Next() {
		if err := rows.Scan(&resp.TodayCalls); err != nil {
			rows.Close()
			return resp, err
		}
	}
	rows.Close()

	var dayQuery string
	if date == "" {
		dayQuery = fmt.Sprintf(`
			SELECT created_at_str, count(*) AS count
			FROM device_cloud_webhooks
			WHERE toString(device_id) = '%s'
			GROUP BY created_at_str
			ORDER BY created_at_str DESC
			LIMIT 31
		`, deviceID)
	} else {
		dayQuery = fmt.Sprintf(`
			SELECT created_at_str, count(*) AS count
			FROM device_cloud_webhooks
			WHERE toString(device_id) = '%s'
			  AND created_at_str = '%s'
			GROUP BY created_at_str
			ORDER BY created_at_str DESC
		`, deviceID, date)
	}

	rows, err = r.ch.Query(ctx, dayQuery)
	if err != nil {
		return resp, err
	}
	var callsByDay []models.TaskStat
	for rows.Next() {
		var stat models.TaskStat
		if err := rows.Scan(&stat.CreatedAtStr, &stat.Count); err != nil {
			rows.Close()
			return resp, err
		}
		callsByDay = append(callsByDay, stat)
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return resp, err
	}
	rows.Close()
	resp.CallsByDay = callsByDay

	statusQuery := fmt.Sprintf(`
SELECT
    status,
    count() AS count
FROM device_cloud_webhooks
WHERE device_id      = '%s'
  AND created_at_str = '%s'
GROUP BY status
ORDER BY status
	`, deviceID, today)

	rows, err = r.ch.Query(ctx, statusQuery)
	if err != nil {
		return resp, err
	}
	var statusCounts []models.StatusCount
	for rows.Next() {
		var sc models.StatusCount
		if err := rows.Scan(&sc.Status, &sc.Count); err != nil {
			rows.Close()
			return resp, err
		}
		statusCounts = append(statusCounts, sc)
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return resp, err
	}
	rows.Close()
	resp.StatusCounts = statusCounts
	if resp.CallsByDay == nil {
		resp.CallsByDay = []models.TaskStat{}
	}
	if resp.StatusCounts == nil {
		resp.StatusCounts = []models.StatusCount{}
	}
	return resp, nil
}

func (r *deviceStatsRepo) GetTaskStats(ctx context.Context, date string) ([]models.TaskStat, error) {
	var query string
	var args []interface{}

	if date == "" {
		query = `
            SELECT 
                created_at_str,
                count(*) AS count
            FROM device_cloud_webhooks
            GROUP BY created_at_str
            ORDER BY created_at_str DESC
            LIMIT 31
        `
	} else {
		query = `
            SELECT 
                created_at_str,
                count(*) AS count
            FROM device_cloud_webhooks
            WHERE created_at_str = ?
            GROUP BY created_at_str
            ORDER BY created_at_str ASC
        `
		args = append(args, date)
	}

	rows, err := r.ch.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var stats []models.TaskStat
	for rows.Next() {
		var st models.TaskStat
		if err := rows.Scan(&st.CreatedAtStr, &st.Count); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		stats = append(stats, st)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return stats, nil
}

func (r *deviceStatsRepo) GetDeviceScreenshots(ctx context.Context, deviceID string, limit, offset int) ([]models.DeviceScreenshot, error) {
	query := fmt.Sprintf(`
		SELECT 
			toString(toDate(parseDateTimeBestEffort(created_at))) AS created_at_str,
			screenshot
		FROM device_cloud_webhooks
		WHERE toString(device_id) = '%s' 
		  AND screenshot <> ''
		ORDER BY parseDateTimeBestEffort(created_at) DESC
		LIMIT %d OFFSET %d
	`, deviceID, limit, offset)

	rows, err := r.ch.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var screenshots []models.DeviceScreenshot
	for rows.Next() {
		var ds models.DeviceScreenshot
		if err := rows.Scan(&ds.CreatedAt, &ds.Screenshot); err != nil {
			return nil, err
		}
		screenshots = append(screenshots, ds)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return screenshots, nil
}

func (r *deviceStatsRepo) GetDeviceCarrierStats(ctx context.Context) ([]models.DeviceCarrierStats, error) {
	query := `
		SELECT device_carrier, COUNT(*) as device_count
		FROM device_cloud_webhooks
		GROUP BY device_carrier
		ORDER BY device_count DESC
	`

	rows, err := r.ch.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %w", err)
	}
	defer rows.Close()

	var stats []models.DeviceCarrierStats
	for rows.Next() {
		var stat models.DeviceCarrierStats
		if err := rows.Scan(&stat.DeviceCarrier, &stat.DeviceCount); err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		stats = append(stats, stat)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return stats, nil
}

func (r *deviceStatsRepo) GetOriginatingCarrierStats(ctx context.Context, fromDate time.Time) (models.CarrierStatsResponse, error) {
	var stats []models.CarrierStat

	query := `
	SELECT
		created_at_str,
		originating_carrier,
		count() AS count
	FROM device_cloud_webhooks
	WHERE created_at_str >= ?
	GROUP BY created_at_str, originating_carrier
	ORDER BY created_at_str DESC, count DESC
	LIMIT 31`

	rows, err := r.ch.Query(ctx, query, fromDate.Format("2006-01-02"))
	if err != nil {
		return models.CarrierStatsResponse{}, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var stat models.CarrierStat
		if err := rows.Scan(&stat.Date, &stat.OriginatingCarrier, &stat.Count); err != nil {
			return models.CarrierStatsResponse{}, fmt.Errorf("scan error: %w", err)
		}
		stats = append(stats, stat)
	}

	if err := rows.Err(); err != nil {
		return models.CarrierStatsResponse{}, fmt.Errorf("rows error: %w", err)
	}

	return models.CarrierStatsResponse{
		Stats: stats,
	}, nil
}
