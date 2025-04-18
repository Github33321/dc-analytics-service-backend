package handler

import (
	"dc-analytics-service-backend/internal/middleware"
	"dc-analytics-service-backend/internal/service"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Logger                    *zap.Logger
	UserHandler               *UserHandler
	DeviceHandler             *DeviceHandler
	DeviceCloudWebhookHandler *DeviceCloudWebhookHandler
	DeviceStatsHandler        *DeviceStatsHandler
	ServerHandler             *ServerHandler
}

func NewHandler(
	logger *zap.Logger,
	userService service.UserService,
	deviceService service.DeviceService,
	clickhouseService service.ClickhouseService,
	deviceStatsService service.DeviceStatsService,
	serverService service.ServerService,

) *Handler {
	return &Handler{
		Logger:                    logger,
		UserHandler:               NewUserHandler(userService),
		DeviceHandler:             NewDeviceHandler(deviceService),
		DeviceCloudWebhookHandler: NewDeviceCloudWebhookHandler(clickhouseService),
		DeviceStatsHandler:        NewDeviceStatsHandler(deviceStatsService),
		ServerHandler:             NewServerHandler(serverService),
	}
}

func (h *Handler) InitRoutes(router *gin.Engine, jwtSecret string) {
	router.Use(middleware.DynamicCORSMiddleware())
	router.Use(GlobalErrorHandler(h.Logger))
	router.POST("/login", LoginHandler)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	secure := router.Group("/v1/analytics")
	secure.Use(middleware.JWTMiddleware(jwtSecret))
	{
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
		router.PATCH("/servers/:id", h.ServerHandler.UpdateServer)
		secure.GET("/servers/:id/devices", h.ServerHandler.GetDevices)
		//secure.GET("/deviceCloudWebhooks", h.DeviceCloudWebhookHandler.GetDeviceCloudWebhooks)
		//clickhouse
		secure.GET("/tasks/stats", h.DeviceStatsHandler.GetTaskStats)
		secure.GET("/devices/:id/call-stats", h.DeviceStatsHandler.GetCallStats)
		secure.GET("/devices/stats", h.DeviceHandler.GetDeviceStats)
		secure.GET("/devices/:id/screenshots", h.DeviceStatsHandler.GetDeviceScreenshots)

	}
}
