package repository

import (
	"context"
	"dc-analytics-service-backend/internal/models"
	"fmt"
	"golang.org/x/sync/errgroup"
	"time"
)

type DeviceStatsRepository interface {
	GetTaskStats(ctx context.Context, date string) ([]models.TaskStat, error)
	GetDeviceCallStats(ctx context.Context, deviceID, date string) (*models.DeviceCallStatsResponse, error)
	GetDeviceScreenshots(ctx context.Context, deviceID string, limit, offset int) ([]models.DeviceScreenshot, error)
	GetDeviceCarrierStats(ctx context.Context) ([]*models.DeviceCarrierStats, error)
	GetOriginatingCarrierStats(ctx context.Context, fromDate time.Time) (*models.CarrierStatsResponse, error)
	GetDeviceGroupStats(ctx context.Context) (*[]models.DeviceGroupStats, error)
	GetTasksReadyCounts(ctx context.Context) (*[]models.TasksReadyCount, error)
	GetByUserID(ctx context.Context, userID uint64) (*[]models.DedicatedDevice, error)
	GetCountedUsers(ctx context.Context) (*[]models.DCUser, error)
	GetTodayStats(ctx context.Context) (*[]models.UserGroupStats, error)
	GetDistinctUsers(ctx context.Context, date string) (*[]models.DCUser2, error)
	GetSourceTypeStats(ctx context.Context) (*[]models.SourceTypeStat, error)
}

type deviceStatsRepo struct {
	ch IClickhouse
}

func NewDeviceStatsRepository(ch IClickhouse) DeviceStatsRepository {
	return &deviceStatsRepo{ch: ch}
}

func (r *deviceStatsRepo) GetDeviceCallStats(ctx context.Context, deviceID, date string) (*models.DeviceCallStatsResponse, error) {
	resp := &models.DeviceCallStatsResponse{}

	today := time.Now().UTC().Format("2006-01-02")
	queryDate := today
	if date != "" {
		queryDate = date
	}

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		q := fmt.Sprintf(`
            SELECT count(*)
            FROM device_cloud_webhooks
            WHERE device_id = %s
              AND created_at_str = '%s'
        `, deviceID, today)

		rows, err := r.ch.Query(ctx, q)
		if err != nil {
			return fmt.Errorf("today query: %w", err)
		}
		defer rows.Close()

		if rows.Next() {
			if err := rows.Scan(&resp.TodayCalls); err != nil {
				return fmt.Errorf("today scan: %w", err)
			}
		}
		return nil
	})

	g.Go(func() error {
		var q string
		if date == "" {
			q = fmt.Sprintf(`
                SELECT created_at_str, count() AS count
                FROM device_cloud_webhooks
                WHERE device_id = %s
                GROUP BY created_at_str
                ORDER BY created_at_str DESC
                LIMIT 31
            `, deviceID)
		} else {
			q = fmt.Sprintf(`
                SELECT created_at_str, count() AS count
                FROM device_cloud_webhooks
                WHERE device_id = %s
                  AND created_at_str = '%s'
                GROUP BY created_at_str
                ORDER BY created_at_str DESC
            `, deviceID, queryDate)
		}

		rows, err := r.ch.Query(ctx, q)
		if err != nil {
			return fmt.Errorf("day query: %w", err)
		}
		defer rows.Close()

		var list []models.TaskStat
		for rows.Next() {
			var ts models.TaskStat
			if err := rows.Scan(&ts.CreatedAtStr, &ts.Count); err != nil {
				return fmt.Errorf("day scan: %w", err)
			}
			list = append(list, ts)
		}
		resp.CallsByDay = &list
		return nil
	})

	g.Go(func() error {
		q := fmt.Sprintf(`
            SELECT status, count() AS count
            FROM device_cloud_webhooks
            WHERE device_id = %s
              AND created_at_str = '%s'
            GROUP BY status
            ORDER BY status
        `, deviceID, today)

		rows, err := r.ch.Query(ctx, q)
		if err != nil {
			return fmt.Errorf("status query: %w", err)
		}
		defer rows.Close()

		var stats []models.StatusCount
		for rows.Next() {
			var sc models.StatusCount
			if err := rows.Scan(&sc.Status, &sc.Count); err != nil {
				return fmt.Errorf("status scan: %w", err)
			}
			stats = append(stats, sc)
		}
		resp.StatusCounts = &stats
		return nil
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}

	if resp.CallsByDay == nil {
		empty := []models.TaskStat{}
		resp.CallsByDay = &empty
	}
	if resp.StatusCounts == nil {
		empty := []models.StatusCount{}
		resp.StatusCounts = &empty
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

func (r *deviceStatsRepo) GetDeviceCarrierStats(ctx context.Context) ([]*models.DeviceCarrierStats, error) {
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

	var stats []*models.DeviceCarrierStats
	for rows.Next() {
		stat := new(models.DeviceCarrierStats)
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

func (r *deviceStatsRepo) GetOriginatingCarrierStats(ctx context.Context, fromDate time.Time) (*models.CarrierStatsResponse, error) {
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
		return nil, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var stat models.CarrierStat
		if err := rows.Scan(&stat.Date, &stat.OriginatingCarrier, &stat.Count); err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}
		stats = append(stats, stat)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return &models.CarrierStatsResponse{
		Stats: &stats,
	}, nil
}
func (r *deviceStatsRepo) GetDeviceGroupStats(ctx context.Context) (*[]models.DeviceGroupStats, error) {
	query := `
		SELECT group_id, count() AS count 
		FROM device_cloud_webhooks 
		GROUP BY group_id 
		ORDER BY count DESC
	`

	rows, err := r.ch.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()

	var result []models.DeviceGroupStats
	for rows.Next() {
		var stat models.DeviceGroupStats
		if err := rows.Scan(&stat.GroupID, &stat.Count); err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}
		result = append(result, stat)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return &result, nil
}

func (r *deviceStatsRepo) GetTasksReadyCounts(ctx context.Context) (*[]models.TasksReadyCount, error) {
	const q = `
        SELECT
            server_group_id,
            count() AS count
        FROM device_cloud_tasks
        WHERE job_id IS NULL
          AND server_group_id IN (1,2)
        GROUP BY server_group_id
        ORDER BY server_group_id
    `
	rows, err := r.ch.Query(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("query tasks_ready: %w", err)
	}
	defer rows.Close()

	var res []models.TasksReadyCount
	for rows.Next() {
		var t models.TasksReadyCount
		if err := rows.Scan(&t.ServerGroupID, &t.Count); err != nil {
			return nil, fmt.Errorf("scan tasks_ready: %w", err)
		}
		res = append(res, t)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error tasks_ready: %w", err)
	}

	return &res, nil
}

func (r *deviceStatsRepo) GetByUserID(ctx context.Context, userID uint64) (*[]models.DedicatedDevice, error) {
	const q = `
        SELECT user_id, device_id, created_at, updated_at
        FROM dedicated_devices
        WHERE user_id = ? 
        ORDER BY created_at DESC
    `
	rows, err := r.ch.Query(ctx, q, userID)
	if err != nil {
		return nil, fmt.Errorf("query dedicated_devices: %w", err)
	}
	defer rows.Close()

	var out []models.DedicatedDevice
	for rows.Next() {
		var dd models.DedicatedDevice
		if err := rows.Scan(&dd.UserID, &dd.DeviceID, &dd.CreatedAt, &dd.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan dedicated_devices: %w", err)
		}
		out = append(out, dd)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error dedicated_devices: %w", err)
	}

	return &out, nil
}

func (r *deviceStatsRepo) GetCountedUsers(ctx context.Context) (*[]models.DCUser, error) {
	const q = `
        SELECT
            user_id,
            count() AS count
        FROM device_cloud_webhooks
        GROUP BY user_id
        ORDER BY count DESC
    `
	rows, err := r.ch.Query(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("query device_cloud_webhooks users: %w", err)
	}
	defer rows.Close()

	var out []models.DCUser
	for rows.Next() {
		var u models.DCUser
		if err := rows.Scan(&u.UserID, &u.Count); err != nil {
			return nil, fmt.Errorf("scan dc user: %w", err)
		}
		out = append(out, u)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows err dc user: %w", err)
	}

	return &out, nil
}

func (r *deviceStatsRepo) GetTodayStats(ctx context.Context) (*[]models.UserGroupStats, error) {
	today := time.Now().UTC().Format("2006-01-02")

	q := fmt.Sprintf(`
        SELECT
            user_id,
            group_id,
            count() AS checks_count
        FROM device_cloud_webhooks
        WHERE created_at_str = '%s'
        GROUP BY user_id, group_id
        ORDER BY user_id, group_id
    `, today)

	rows, err := r.ch.Query(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("query today stats: %w", err)
	}
	defer rows.Close()

	var out []models.UserGroupStats
	for rows.Next() {
		var s models.UserGroupStats
		if err := rows.Scan(&s.UserID, &s.GroupID, &s.ChecksCount); err != nil {
			return nil, fmt.Errorf("scan today stats: %w", err)
		}
		out = append(out, s)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows err today stats: %w", err)
	}

	return &out, nil
}

func (r *deviceStatsRepo) GetDistinctUsers(ctx context.Context, date string) (*[]models.DCUser2, error) {
	var q string
	if date == "" {
		q = `
            SELECT DISTINCT user_id
            FROM device_cloud_webhooks
            ORDER BY user_id
        `
	} else {
		q = fmt.Sprintf(`
            SELECT DISTINCT user_id
            FROM device_cloud_webhooks
            WHERE created_at_str = '%s'
            ORDER BY user_id
        `, date)
	}

	rows, err := r.ch.Query(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("query distinct users: %w", err)
	}
	defer rows.Close()

	var out []models.DCUser2
	for rows.Next() {
		var u models.DCUser2
		if err := rows.Scan(&u.UserID); err != nil {
			return nil, fmt.Errorf("scan distinct user: %w", err)
		}
		out = append(out, u)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows err distinct users: %w", err)
	}

	return &out, nil
}

func (r *deviceStatsRepo) GetSourceTypeStats(ctx context.Context) (*[]models.SourceTypeStat, error) {
	const query = `
		SELECT 
			source_type_id, 
			count() AS count
		FROM device_cloud_webhooks
		GROUP BY source_type_id
		ORDER BY source_type_id
	`

	rows, err := r.ch.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения статистики по источникам: %w", err)
	}
	defer rows.Close()

	var stats []models.SourceTypeStat
	for rows.Next() {
		var stat models.SourceTypeStat
		if err := rows.Scan(&stat.SourceTypeID, &stat.Count); err != nil {
			return nil, fmt.Errorf("ошибка сканирования строки: %w", err)
		}
		stats = append(stats, stat)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("ошибка в rows: %w", rows.Err())
	}

	return &stats, nil
}
