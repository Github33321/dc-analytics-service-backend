package main

import (
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

	chConn, err := clickhouse.WaitForClickHouse(cfg.ClickHouseDSN, 10, 5*time.Second)
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
