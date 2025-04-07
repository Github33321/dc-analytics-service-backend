package repository

import (
	"context"
	"dc-analytics-service-backend/internal/config"

	"github.com/jackc/pgx/v5"
	pgxpgconn "github.com/jackc/pgx/v5/pgconn"
)

type IClickhouse interface {
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	Exec(ctx context.Context, sql string, args ...interface{}) (pgxpgconn.CommandTag, error)
}

type Clickhouse struct {
	db     *pgx.Conn
	config *config.Config
}

func NewClickhouseRepo(conn *pgx.Conn, cfg *config.Config) IClickhouse {
	return &Clickhouse{
		db:     conn,
		config: cfg,
	}
}

func (c *Clickhouse) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	return c.db.Query(ctx, sql, args...)
}

func (c *Clickhouse) Exec(ctx context.Context, sql string, args ...interface{}) (pgxpgconn.CommandTag, error) {
	return c.db.Exec(ctx, sql, args...)
}
