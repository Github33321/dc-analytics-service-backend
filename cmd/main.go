package main

import (
	"dc-analytics-service-backend/internal/config"
	"dc-analytics-service-backend/internal/handler"
	"dc-analytics-service-backend/internal/middleware"
	"dc-analytics-service-backend/internal/repository"
	"dc-analytics-service-backend/internal/service"
	"dc-analytics-service-backend/pkg/clickhouse"
	"dc-analytics-service-backend/pkg/logger"
	"dc-analytics-service-backend/pkg/postgres"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"path/filepath"
	"time"
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

	router := gin.Default()

	router.GET("/deviceCloudWebhooks", func(c *gin.Context) {

		chConn, err := clickhouse.WaitForClickHouse(cfg.ClickHouseDSN, 10, 5*time.Second)
		if err != nil {
			logg.Sugar().Fatalf("Ошибка подключения к ClickHouse: %v", err)
		}
		defer chConn.Close()
		logg.Sugar().Info("Подключение к ClickHouse установлено")

		handler.SetDB(chConn)

		handler.SelectHandler(c)
	})

	db, err := postgres.OpenDB()
	if err != nil {
		logg.Sugar().Fatalf("Ошибка подключения к базе данных: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			logg.Sugar().Errorf("Ошибка при закрытии соединения с БД: %v", err)
		}
	}()
	if err := db.Ping(); err != nil {
		logg.Sugar().Fatalf("Не удалось проверить соединение с БД: %v", err)
	}

	router.POST("/login", handler.LoginHandler)

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	deviceRepo := repository.NewDeviceRepository(db)
	deviceService := service.NewDeviceService(deviceRepo)
	deviceHandler := handler.NewDeviceHandler(deviceService)

	secure := router.Group("/v1/analytics")
	secure.Use(middleware.JWTMiddleware(cfg.JWTSecret))
	{
		secure.GET("/ping", handler.PingHandler)

		secure.GET("/users/:id", userHandler.GetUserByID)
		secure.GET("/users", userHandler.GetUsers)
		secure.POST("/users", userHandler.CreateUser)
		secure.DELETE("/users/:id", userHandler.DeleteUser)

		secure.GET("/devices/:id", deviceHandler.GetDeviceByID)
		secure.GET("/devices", deviceHandler.GetDevices)
		secure.PATCH("/devices/:id", deviceHandler.UpdateDevice)
		secure.DELETE("/devices/:id", deviceHandler.DeleteDevice)

	}

	logg.Sugar().Infof("Запуск сервера на порту %s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		logg.Sugar().Fatalf("Ошибка сервера: %v", err)
	}
}
