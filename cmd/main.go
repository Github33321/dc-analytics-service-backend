package main

import (
	"context"
	"log"

	"dc-analytics-service-backend/internal/config"
	"dc-analytics-service-backend/internal/handler"
	"dc-analytics-service-backend/internal/repository"
	"dc-analytics-service-backend/internal/service"
	"dc-analytics-service-backend/pkg/clickhouse"
	"dc-analytics-service-backend/pkg/logger"
	"dc-analytics-service-backend/pkg/postgres"

	"github.com/gin-gonic/gin"
)

func main() {
	ctx := context.Background()

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal("Ошибка загрузки конфигурации: " + err.Error())
	}

	logg, err := logger.NewLogger(cfg.LogLevel)
	if err != nil {
		log.Fatal("Ошибка инициализации логгера: " + err.Error())
	}
	defer logg.Sync()
	logg.Sugar().Infof("Логгер инициализирован с уровнем %s", cfg.LogLevel)

	router := gin.Default()

	pgDB, err := postgres.OpenDB(ctx, cfg.PostgresDSN)
	if err != nil {
		logg.Sugar().Fatalf("Ошибка подключения к PostgreSQL: %v", err)
	}
	defer pgDB.Close()

	if err := postgres.PingDB(ctx, pgDB); err != nil {
		logg.Sugar().Fatalf("Не удалось проверить соединение с PostgreSQL: %v", err)
	}
	logg.Sugar().Info("Соединение с PostgreSQL установлено")

	userRepo := repository.NewUserRepository(pgDB)
	userService := service.NewUserService(userRepo)

	deviceRepo := repository.NewDeviceRepository(pgDB)
	deviceService := service.NewDeviceService(deviceRepo)

	chClient, err := clickhouse.NewClickhouseClient(ctx, cfg.ClickhouseConfig)
	if err != nil {
		logg.Sugar().Fatalf("Ошибка подключения к ClickHouse: %v", err)
	}

	clickhouseRepo := repository.NewClickhouseRepo(chClient.DB, cfg)
	clickhouseService := service.NewClickhouseService(clickhouseRepo)

	h := handler.NewHandler(userService, deviceService, clickhouseService)
	h.InitRoutes(router, cfg.JWTSecret)

	logg.Sugar().Infof("Запуск сервера на порту %s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		logg.Sugar().Fatalf("Ошибка сервера: %v", err)
	}
}
