package repository

import (
	"context"
	"dc-analytics-service-backend/internal/models"
	"fmt"
)

type DeviceStatsRepository interface {
	GetTaskStats(ctx context.Context, date string) ([]models.TaskStat, error)
	GetDeviceCallStats(ctx context.Context, deviceID string, date string) (models.DeviceCallStatsResponse, error)
	GetDeviceScreenshots(ctx context.Context, deviceID string, limit, offset int) ([]models.DeviceScreenshot, error)
}

type deviceStatsRepo struct {
	ch IClickhouse
}

func NewDeviceStatsRepository(ch IClickhouse) DeviceStatsRepository {
	return &deviceStatsRepo{ch: ch}
}

//	func (r *deviceStatsRepo) GetDeviceCallStats(ctx context.Context, deviceID string, date string) (models.DeviceCallStatsResponse, error) {
//		var resp models.DeviceCallStatsResponse
//
//		todayQuery := fmt.Sprintf(`
//			SELECT count(*)
//			FROM device_cloud_webhooks
//			WHERE toString(device_id) = '%s'
//			  AND toDate(parseDateTimeBestEffort(created_at)) = today()
//		`, deviceID)
//		rows, err := r.ch.Query(ctx, todayQuery)
//		if err != nil {
//			return resp, err
//		}
//		if rows.Next() {
//			if err := rows.Scan(&resp.TodayCalls); err != nil {
//				rows.Close()
//				return resp, err
//			}
//		}
//		rows.Close()
//
//		var dayQuery string
//		if date == "" {
//			dayQuery = fmt.Sprintf(`
//				SELECT
//					formatDateTime(parseDateTimeBestEffort(created_at), '%%Y-%%m-%%d') AS created_at_str,
//					count(*) AS count
//				FROM device_cloud_webhooks
//				WHERE toString(device_id) = '%s'
//				GROUP BY formatDateTime(parseDateTimeBestEffort(created_at), '%%Y-%%m-%%d')
//				ORDER BY formatDateTime(parseDateTimeBestEffort(created_at), '%%Y-%%m-%%d') ASC
//			`, deviceID)
//		} else {
//			dayQuery = fmt.Sprintf(`
//				SELECT
//					formatDateTime(parseDateTimeBestEffort(created_at), '%%Y-%%m-%%d') AS created_at_str,
//					count(*) AS count
//				FROM device_cloud_webhooks
//				WHERE toString(device_id) = '%s'
//				  AND formatDateTime(parseDateTimeBestEffort(created_at), '%%Y-%%m-%%d') = '%s'
//				GROUP BY formatDateTime(parseDateTimeBestEffort(created_at), '%%Y-%%m-%%d')
//				ORDER BY formatDateTime(parseDateTimeBestEffort(created_at), '%%Y-%%m-%%d') ASC
//			`, deviceID, date)
//		}
//		rows, err = r.ch.Query(ctx, dayQuery)
//		if err != nil {
//			return resp, err
//		}
//		var callsByDay []models.TaskStat
//		for rows.Next() {
//			var stat models.TaskStat
//			if err := rows.Scan(&stat.CreatedAtStr, &stat.Count); err != nil {
//				rows.Close()
//				return resp, err
//			}
//			callsByDay = append(callsByDay, stat)
//		}
//		if err := rows.Err(); err != nil {
//			rows.Close()
//			return resp, err
//		}
//		rows.Close()
//		resp.CallsByDay = callsByDay
//
//		statusQuery := fmt.Sprintf(`
//			SELECT status, count(*) AS count
//			FROM device_cloud_webhooks
//			WHERE toString(device_id) = '%s'
//			GROUP BY status
//			ORDER BY status
//		`, deviceID)
//		rows, err = r.ch.Query(ctx, statusQuery)
//		if err != nil {
//			return resp, err
//		}
//		var statusCounts []models.StatusCount
//		for rows.Next() {
//			var sc models.StatusCount
//			if err := rows.Scan(&sc.Status, &sc.Count); err != nil {
//				rows.Close()
//				return resp, err
//			}
//			statusCounts = append(statusCounts, sc)
//		}
//		if err := rows.Err(); err != nil {
//			rows.Close()
//			return resp, err
//		}
//		rows.Close()
//		resp.StatusCounts = statusCounts
//
//		return resp, nil
//	}
func (r *deviceStatsRepo) GetDeviceCallStats(ctx context.Context, deviceID string, date string) (models.DeviceCallStatsResponse, error) {
	var resp models.DeviceCallStatsResponse

	todayQuery := fmt.Sprintf(`
		SELECT count(*) 
		FROM device_cloud_webhooks
		WHERE toString(device_id) = '%s'
		  AND toDate(parseDateTimeBestEffort(created_at)) = today()
	`, deviceID)

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
			SELECT 
				toString(toDate(parseDateTimeBestEffort(created_at))) AS created_at_str, 
				count(*) AS count
			FROM device_cloud_webhooks
			WHERE toString(device_id) = '%s'
			GROUP BY toDate(parseDateTimeBestEffort(created_at))
			ORDER BY toDate(parseDateTimeBestEffort(created_at)) ASC
		`, deviceID)
	} else {
		dayQuery = fmt.Sprintf(`
			SELECT 
				toString(toDate(parseDateTimeBestEffort(created_at))) AS created_at_str, 
				count(*) AS count
			FROM device_cloud_webhooks
			WHERE toString(device_id) = '%s'
			  AND toString(toDate(parseDateTimeBestEffort(created_at))) = '%s'
			GROUP BY toDate(parseDateTimeBestEffort(created_at))
			ORDER BY toDate(parseDateTimeBestEffort(created_at)) ASC
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
		SELECT status, sum(count) AS count 
		FROM (
			SELECT 'success' AS status, sum(multiIf(lower(status) = 'success', 1, 0)) AS count
			FROM device_cloud_webhooks
			WHERE toString(device_id) = '%s'
			UNION ALL
			SELECT 'call_mismatch', sum(multiIf(lower(status) = 'call_mismatch', 1, 0))
			FROM device_cloud_webhooks
			WHERE toString(device_id) = '%s'
			UNION ALL
			SELECT 'wait', sum(multiIf(lower(status) = 'wait', 1, 0))
			FROM device_cloud_webhooks
			WHERE toString(device_id) = '%s'
			UNION ALL
			SELECT 'no_result', sum(multiIf(lower(status) = 'no_result', 1, 0))
			FROM device_cloud_webhooks
			WHERE toString(device_id) = '%s'
			UNION ALL
			SELECT 'call_failed', sum(multiIf(lower(status) = 'call_failed', 1, 0))
			FROM device_cloud_webhooks
			WHERE toString(device_id) = '%s'
		) AS t
		GROUP BY status
		ORDER BY status ASC
	`, deviceID, deviceID, deviceID, deviceID, deviceID)

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

	return resp, nil
}

func (r *deviceStatsRepo) GetTaskStats(ctx context.Context, date string) ([]models.TaskStat, error) {
	var query string
	if date == "" {
		query = `
            SELECT 
                toString(toDate(parseDateTimeBestEffort(created_at))) AS created_at_str,
                count(*) AS count
            FROM device_cloud_webhooks
            GROUP BY toDate(parseDateTimeBestEffort(created_at))
            ORDER BY toDate(parseDateTimeBestEffort(created_at)) ASC
        `
	} else {
		query = `
            SELECT 
                toString(toDate(parseDateTimeBestEffort(created_at))) AS created_at_str,
                count(*) AS count
            FROM device_cloud_webhooks
            WHERE toString(toDate(parseDateTimeBestEffort(created_at))) = '` + date + `'
            GROUP BY toDate(parseDateTimeBestEffort(created_at))
            ORDER BY toDate(parseDateTimeBestEffort(created_at)) ASC
        `
	}

	rows, err := r.ch.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []models.TaskStat
	for rows.Next() {
		var st models.TaskStat
		if err := rows.Scan(&st.CreatedAtStr, &st.Count); err != nil {
			return nil, err
		}
		stats = append(stats, st)
	}
	if err := rows.Err(); err != nil {
		return nil, err
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
