package main

import (
	"dc-analytics-service-backend/internal/config"
	"dc-analytics-service-backend/internal/handler"
	"dc-analytics-service-backend/internal/middleware"
	"dc-analytics-service-backend/pkg/clickhouse"
	"dc-analytics-service-backend/pkg/logger"
	"github.com/joho/godotenv"
	"log"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {

	if err := godotenv.Load(filepath.Join("..", ".env")); err != nil {
		log.Println(".env не найден, используются переменные окружения по умолчанию")
	}
	cfg, err := config.NewConfig()

	if err != nil {
		panic("Ошибка загрузки конфигурации: " + err.Error())
	}

	logg, err := logger.NewLogger(cfg.LogLevel)
	if err != nil {
		panic("Ошибка инициализации логгера: " + err.Error())
	}
	defer logg.Sync()
	logg.Sugar().Infof("Логгер инициализирован с уровнем %s", cfg.LogLevel)

	chConn, err := clickhouse.WaitForClickHouse(cfg.ClickHouseDSN, 10, 5*time.Second)
	if err != nil {
		logg.Sugar().Fatalf("Ошибка подключения к ClickHouse: %v", err)
	}
	defer chConn.Close()
	logg.Sugar().Info("Подключение к ClickHouse установлено")

	handler.SetDB(chConn)

	router := gin.Default()
	router.GET("/ping", handler.PingHandler)

	secure := router.Group("/secure")
	secure.Use(middleware.JWTMiddleware(cfg.JWTSecret))
	secure.GET("", handler.SecureHandler)

	deviceCloudWebhooks := router.GET("/deviceCloudWebhooks", handler.SelectHandler)
	deviceCloudWebhooks.Use(middleware.JWTMiddleware(cfg.JWTSecret))
	deviceCloudWebhooks.GET("", handler.SecureHandler)

	logg.Sugar().Infof("Запуск сервера на порту %s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		logg.Sugar().Fatalf("Ошибка сервера: %v", err)
	}
}
