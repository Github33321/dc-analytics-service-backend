package repository

import (
	"context"
	"dc-analytics-service-backend/internal/config"

	"github.com/ClickHouse/clickhouse-go/v2"
	chdriver "github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

type IClickhouse interface {
	Query(ctx context.Context, query string, args ...interface{}) (chdriver.Rows, error)
	Exec(ctx context.Context, query string, args ...interface{}) error
}

type ClickhouseRepo struct {
	conn   clickhouse.Conn
	config *config.Config
}

func NewClickhouseRepo(conn clickhouse.Conn, cfg *config.Config) IClickhouse {
	return &ClickhouseRepo{
		conn:   conn,
		config: cfg,
	}
}

func (r *ClickhouseRepo) Query(ctx context.Context, query string, args ...interface{}) (chdriver.Rows, error) {
	return r.conn.Query(ctx, query, args...)
}

func (r *ClickhouseRepo) Exec(ctx context.Context, query string, args ...interface{}) error {
	return r.conn.Exec(ctx, query, args...)
}
