package repository

import (
	"context"
	"dc-analytics-service-backend/internal/config"

	// Импортируем сам драйвер и дополнительно его внутренние типы через алиас chdriver
	"github.com/ClickHouse/clickhouse-go/v2"
	chdriver "github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

// IClickhouse определяет методы, которые мы хотим поддерживать
// и возвращает тип из самого драйвера clickhouse-go/v2/lib/driver, а не из "database/sql/driver"
type IClickhouse interface {
	Query(ctx context.Context, query string, args ...interface{}) (chdriver.Rows, error)
	Exec(ctx context.Context, query string, args ...interface{}) error
}

// ClickhouseRepo — конкретная реализация интерфейса IClickhouse
type ClickhouseRepo struct {
	conn   clickhouse.Conn
	config *config.Config
}

// NewClickhouseRepo — конструктор, принимающий готовое подключение к ClickHouse
func NewClickhouseRepo(conn clickhouse.Conn, cfg *config.Config) IClickhouse {
	return &ClickhouseRepo{
		conn:   conn,
		config: cfg,
	}
}

// Query возвращает chdriver.Rows (тип, экспортируемый драйвером ClickHouse)
func (r *ClickhouseRepo) Query(ctx context.Context, query string, args ...interface{}) (chdriver.Rows, error) {
	return r.conn.Query(ctx, query, args...)
}

// Exec в clickhouse-go/v2 возвращает error
func (r *ClickhouseRepo) Exec(ctx context.Context, query string, args ...interface{}) error {
	return r.conn.Exec(ctx, query, args...)
}
