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
	secure.Use(middleware.JWTMiddleware(h.JWTSecret))
	{
		secure.GET("/ping", PingHandler)

		secure.GET("/users/:id", h.UserHandler.GetUserByID)
		secure.GET("/users", h.UserHandler.GetUsers)
		secure.POST("/users", h.UserHandler.CreateUser)
		secure.DELETE("/users/:id", h.UserHandler.DeleteUser)

		secure.GET("/devices/:id", h.DeviceHandler.GetDeviceByID)
		secure.GET("/devices", h.DeviceHandler.GetDevices)
		secure.PATCH("/devices/:id", h.DeviceHandler.UpdateDevice)
		secure.DELETE("/devices/:id", h.DeviceHandler.DeleteDevice)
		secure.GET("/devices/stats", h.DeviceHandler.GetDeviceStats)
		secure.GET("/devices/:id/screenshots", h.DeviceStatsHandler.GetDeviceScreenshots)

		secure.GET("/servers", h.ServerHandler.GetServers)
		secure.GET("/servers/:id", h.ServerHandler.GetServerByID)
		secure.PATCH("/servers/:id", h.ServerHandler.UpdateServer)
		secure.GET("/servers/:id/devices", h.ServerHandler.GetDevices)

		secure.GET("/calls/stats", h.DeviceStatsHandler.GetTaskStats)
		secure.GET("/calls/:id/stats", deviceExists, h.DeviceStatsHandler.GetDeviceCallStats)
		secure.GET("/calls/carriers", h.DeviceStatsHandler.GetDeviceCarrierStats)
		secure.GET("/calls/originating-carriers", h.DeviceStatsHandler.GetOriginatingCarrierStats)
		secure.GET("/calls/groups", h.DeviceStatsHandler.GetDeviceGroupStats)
		secure.GET("/calls/sources", h.DeviceStatsHandler.GetSourcesStats)
		secure.GET("/devices/dedicated/users/:id", h.DeviceStatsHandler.GetByUserID)
		secure.GET("/calls/users", h.DeviceStatsHandler.GetCountedUsers)
		secure.GET("/calls/users-today", h.DeviceStatsHandler.GetTodayStats)
		secure.GET("/calls/user-list", h.UserHandler.GetDCUsers)
		secure.GET("/tasks/ready", h.DeviceStatsHandler.GetTasksReadyCounts)
	}
}
