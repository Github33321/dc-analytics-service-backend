package main

import (
	"context"
	_ "dc-analytics-service-backend/docs"
	"dc-analytics-service-backend/internal/config"
	"dc-analytics-service-backend/internal/handler"
	"dc-analytics-service-backend/internal/repository"
	"dc-analytics-service-backend/internal/service"
	"dc-analytics-service-backend/pkg/clickhouse"
	"dc-analytics-service-backend/pkg/logger"
	"dc-analytics-service-backend/pkg/postgres"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx := context.Background()

	appConfig, err := config.NewConfig()
	if err != nil {
		log.Fatal("Ошибка загрузки конфигурации: " + err.Error())
	}

	logg, err := logger.NewLogger(appConfig.LogLevel)
	if err != nil {
		log.Fatal("Ошибка инициализации логгера: " + err.Error())
	}
	defer logg.Sync()
	logg.Sugar().Infof("Логгер инициализирован с уровнем %s", appConfig.LogLevel)

	router := gin.Default()

	// Подключение к PostgreSQL
	pgDB, err := postgres.OpenDB(ctx, appConfig.PostgresDSN)
	if err != nil {
		logg.Sugar().Fatalf("Ошибка подключения к PostgreSQL: %v", err)
	}
	//if err := postgres.PingDB(ctx, pgDB); err != nil {
	//	logg.Sugar().Fatalf("Не удалось проверить соединение с PostgreSQL: %v", err)
	//}
	logg.Sugar().Info("Соединение с PostgreSQL установлено")

	userRepo := repository.NewUserRepository(pgDB)
	userService := service.NewUserService(userRepo)

	deviceRepo := repository.NewDeviceRepository(pgDB)
	deviceService := service.NewDeviceService(deviceRepo)

	// Подключение к ClickHouse
	chClient, err := clickhouse.NewClient(ctx, appConfig.ClickhouseConfig)
	if err != nil {
		logg.Sugar().Fatalf("Ошибка подключения к ClickHouse: %v", err)
	}
	logg.Sugar().Info("Соединение с ClickHouse установлено")

	clickhouseRepo := repository.NewClickhouseRepo(chClient.Conn, appConfig)
	clickhouseService := service.NewClickhouseService(clickhouseRepo)

	deviceStatsRepo := repository.NewDeviceStatsRepository(clickhouseRepo)
	deviceStatsService := service.NewDeviceStatsService(deviceStatsRepo)

	h := handler.NewHandler(logg, userService, deviceService, clickhouseService, deviceStatsService)
	h.InitRoutes(router, appConfig.JWTSecret)

	srv := &http.Server{
		Addr:    ":" + appConfig.Port,
		Handler: router,
	}

	go func() {
		logg.Sugar().Infof("Запуск сервера на порту %s", appConfig.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logg.Sugar().Fatalf("Ошибка сервера: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	logg.Sugar().Info("Получен сигнал завершения, останавливаем сервер...")

	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctxShutDown); err != nil {
		logg.Sugar().Fatal("Ошибка при завершении работы сервера: ", err)
	}

	pgDB.Close()

	if err := chClient.Conn.Close(); err != nil {
		logg.Sugar().Errorf("Ошибка при закрытии ClickHouse: %v", err)
	}

	logg.Sugar().Info("Сервер и соединения остановлены корректно")
}
