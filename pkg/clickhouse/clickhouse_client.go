package clickhouse

import (
	"context"
	"fmt"
	"time"

	"dc-analytics-service-backend/internal/config"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/sirupsen/logrus"
)

type Client struct {
	Conn     clickhouse.Conn
	Database string
}

func NewClient(ctx context.Context, chConfig *config.ClickhouseConfig) (*Client, error) {
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{chConfig.Host}, // например, "localhost:9000"
		Auth: clickhouse.Auth{
			Database: chConfig.Database,
			Username: chConfig.Username,
			Password: chConfig.Password,
		},
		Debug:       chConfig.Debug,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to ClickHouse: %w", err)
	}

	if err := conn.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping ClickHouse: %w", err)
	}

	client := &Client{
		Conn:     conn,
		Database: chConfig.Database,
	}

	go client.checkPing()

	return client, nil
}

func (c *Client) checkPing() {
	for {
		time.Sleep(time.Minute)
		if err := c.Conn.Ping(context.Background()); err != nil {
			logrus.Error(fmt.Errorf("clickhouse ping error: %v", err))
		}
	}
}
