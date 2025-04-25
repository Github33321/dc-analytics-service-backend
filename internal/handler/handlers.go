package handler

import (
	"dc-analytics-service-backend/internal/middleware"
	"dc-analytics-service-backend/internal/repository"
	"dc-analytics-service-backend/internal/service"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Logger                    *zap.Logger
	JWTSecret                 string
	UserHandler               *UserHandler
	DeviceHandler             *DeviceHandler
	DeviceCloudWebhookHandler *DeviceCloudWebhookHandler
	DeviceStatsHandler        *DeviceStatsHandler
	ServerHandler             *ServerHandler
}

func handleClientError(c *gin.Context, logger *zap.Logger, httpStatus int, clientMsg string, err error) {
	logger.Error(clientMsg, zap.Error(err))
	c.JSON(httpStatus, gin.H{"error": clientMsg})
	c.Abort()
}
func NewHandler(
	logger *zap.Logger,
	jwtSecret string,
	userService service.UserService,
	deviceService service.DeviceService,
	clickhouseService service.ClickhouseService,
	deviceStatsService service.DeviceStatsService,
	serverService service.ServerService,
) *Handler {
	return &Handler{
		Logger:                    logger,
		JWTSecret:                 jwtSecret,
		UserHandler:               NewUserHandler(logger, userService),
		DeviceHandler:             NewDeviceHandler(logger, deviceService),
		DeviceCloudWebhookHandler: NewDeviceCloudWebhookHandler(clickhouseService),
		DeviceStatsHandler:        NewDeviceStatsHandler(logger, deviceStatsService),
		ServerHandler:             NewServerHandler(logger, serverService),
	}
}

func (h *Handler) InitRoutes(router *gin.Engine, deviceRepo repository.DeviceRepository) {
	router.Use(middleware.DynamicCORSMiddleware())

	router.Use(GlobalErrorHandler(h.Logger))
	router.POST("/login", LoginHandler)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	deviceExists := middleware.DeviceExistsMiddleware(deviceRepo)
	secure := router.Group("/v1/analytics")
	//secure.Use(middleware.JWTMiddleware(jwtSecret))
	secure.Use(middleware.JWTMiddleware(h.JWTSecret))
	{
		//time curl -X GET "http://localhost:7002/v1/analytics/device-carriers/source" -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJmcm9udGVuZC11c2VyIiwiZXhwIjoxNzQ1NjYyNjcwLCJpYXQiOjE3NDU1NzYyNzB9.Ssv6_u8PyFgMWzvCn6j9cNVGn2HRiGIKe2zNmYnt6W0" -H "Accept: application/json"
		//curl -X GET 'https://ms1.calleridreputation.com/v1/sources/list' -H "X-Auth-Key: 4f9c5e8f-2c8b-4f7d-91d8-1a2b3c4d5e6f"
		//postgres
		secure.GET("/ping", PingHandler)

		secure.GET("/users/:id", h.UserHandler.GetUserByID)
		secure.GET("/users", h.UserHandler.GetUsers)
		secure.POST("/users", h.UserHandler.CreateUser)
		secure.DELETE("/users/:id", h.UserHandler.DeleteUser)

		secure.GET("/devices/:id", h.DeviceHandler.GetDeviceByID)
		secure.GET("/devices", h.DeviceHandler.GetDevices)
		secure.PATCH("/devices/:id", h.DeviceHandler.UpdateDevice)
		secure.DELETE("/devices/:id", h.DeviceHandler.DeleteDevice)

		secure.GET("/servers", h.ServerHandler.GetServers)
		secure.GET("/servers/:id", h.ServerHandler.GetServerByID)
		secure.PATCH("/servers/:id", h.ServerHandler.UpdateServer)
		secure.GET("/servers/:id/devices", h.ServerHandler.GetDevices)
		//secure.GET("/deviceCloudWebhooks", h.DeviceCloudWebhookHandler.GetDeviceCloudWebhooks)
		//clickhouse
		secure.GET("/tasks/stats", h.DeviceStatsHandler.GetTaskStats)
		secure.GET("/devices/:id/call-stats", deviceExists, h.DeviceStatsHandler.GetDeviceCallStats)
		secure.GET("/devices/stats", h.DeviceHandler.GetDeviceStats)
		secure.GET("/devices/:id/screenshots", h.DeviceStatsHandler.GetDeviceScreenshots)

		secure.GET("/device-carriers/operator", h.DeviceStatsHandler.GetDeviceCarrierStats)
		secure.GET("/device-carriers/originating", h.DeviceStatsHandler.GetOriginatingCarrierStats)
		secure.GET("/device-carriers/group", h.DeviceStatsHandler.GetDeviceGroupStats)
		//Иммитация
		secure.GET("/device-carriers/source", h.DeviceStatsHandler.GetSourcesStats)

		secure.GET("/tasks/group", h.DeviceStatsHandler.GetTasksReadyCounts)
		secure.GET("/tasks/:id/devices", h.DeviceStatsHandler.GetByUserID)
		secure.GET("/tasks/users", h.DeviceStatsHandler.GetCountedUsers)
		secure.GET("/tasks/users-today", h.DeviceStatsHandler.GetTodayStats)
		secure.GET("/tasks/user-list", h.DeviceStatsHandler.GetDistinctUsers)
	}
}
