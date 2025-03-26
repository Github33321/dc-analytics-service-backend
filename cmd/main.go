package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"dc-analytics-service-backend/internal/config"
	"dc-analytics-service-backend/internal/handler"
	"dc-analytics-service-backend/pkg/clickhouse"
	"dc-analytics-service-backend/pkg/logger"
)

func main() {

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	logg := logger.InitLogger(cfg.LogLevel)
	logg.Infof("Логгер инициализирован с уровнем %s", cfg.LogLevel)

	chConn, err := waitForClickHouse(cfg.ClickHouseDSN, 10, 5*time.Second)
	if err != nil {
		logg.Fatalf("Ошибка подключения к ClickHouse: %v", err)
	}
	defer chConn.Close()
	logg.Infof("Подключение к ClickHouse установлено")

	http.HandleFunc("/ping", handler.PingHandler)

	logg.Infof("Запуск сервера на порту %s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, nil); err != nil {
		logg.Fatalf("Ошибка сервера: %v", err)
	}
}

func waitForClickHouse(dsn string, maxRetries int, delay time.Duration) (*sql.DB, error) {
	var db *sql.DB
	var err error
	for i := 0; i < maxRetries; i++ {
		db, err = clickhouse.Connect(dsn)
		if err == nil {
			return db, nil
		}
		log.Printf("Не удалось подключиться к ClickHouse (попытка %d/%d): %v", i+1, maxRetries, err)
		time.Sleep(delay)
	}
	return nil, fmt.Errorf("не удалось подключиться к ClickHouse после %d попыток: %w", maxRetries, err)
}
