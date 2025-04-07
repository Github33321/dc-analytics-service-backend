package clickhouse

import (
	"context"
	"fmt"
	"strings"
	"time"

	"dc-analytics-service-backend/internal/config"

	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
)

type Clickhouse struct {
	DB       *pgx.Conn
	Database string
}

func NewClickhouseClient(ctx context.Context, chConfig *config.ClickhouseConfig) (*Clickhouse, error) {
	dsn := strings.TrimSpace(chConfig.ConnectionURL)
	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to ClickHouse using pgx: %w", err)
	}

	if err := conn.Ping(ctx); err != nil {
		conn.Close(ctx)
		return nil, fmt.Errorf("failed to ping ClickHouse: %w", err)
	}

	ch := &Clickhouse{
		DB:       conn,
		Database: chConfig.Database,
	}

	go ch.checkClickhousePing()

	return ch, nil
}

func (c *Clickhouse) checkClickhousePing() {
	for {
		time.Sleep(time.Minute)
		if err := c.DB.Ping(context.Background()); err != nil {
			logrus.Error(fmt.Errorf("clickhouse ping error: %v", err))
		}
	}
}
