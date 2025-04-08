package repository

import (
	"context"
	"dc-analytics-service-backend/internal/models"
)

type TaskStatRepository interface {
	GetTaskStats(ctx context.Context, date string) ([]models.TaskStat, error)
}

type taskStatRepository struct {
	ch IClickhouse
}

func NewTaskStatRepository(ch IClickhouse) TaskStatRepository {
	return &taskStatRepository{ch: ch}
}

func (r *taskStatRepository) GetTaskStats(ctx context.Context, date string) ([]models.TaskStat, error) {
	var query string
	if date == "" {
		query = `
            SELECT 
                toString(created_at) AS created_at_str,
                count(*) AS count
            FROM device_cloud_webhooks
            GROUP BY created_at_str
            ORDER BY created_at_str
        `
	} else {
		query = `
            SELECT
                toString(created_at) AS created_at_str,
                count(*) AS count
            FROM device_cloud_webhooks
            WHERE toString(created_at) = '` + date + `'
            GROUP BY created_at_str
            ORDER BY created_at_str
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
