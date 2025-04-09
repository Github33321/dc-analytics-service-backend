package handler

import (
	"dc-analytics-service-backend/internal/middleware"
	"dc-analytics-service-backend/internal/service"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	UserHandler               *UserHandler
	DeviceHandler             *DeviceHandler
	DeviceCloudWebhookHandler *DeviceCloudWebhookHandler
	DeviceStatsHandler        *DeviceStatsHandler
}

func NewHandler(
	userService service.UserService,
	deviceService service.DeviceService,
	clickhouseService service.ClickhouseService,
	deviceStatsService service.DeviceStatsService,

) *Handler {
	return &Handler{
		UserHandler:               NewUserHandler(userService),
		DeviceHandler:             NewDeviceHandler(deviceService),
		DeviceCloudWebhookHandler: NewDeviceCloudWebhookHandler(clickhouseService),
		DeviceStatsHandler:        NewDeviceStatsHandler(deviceStatsService),
	}
}

func (h *Handler) InitRoutes(router *gin.Engine, jwtSecret string) {
	router.Use(middleware.DynamicCORSMiddleware())
	router.POST("/login", LoginHandler)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	secure := router.Group("/v1/analytics")
	secure.Use(middleware.JWTMiddleware(jwtSecret))
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

		//secure.GET("/deviceCloudWebhooks", h.DeviceCloudWebhookHandler.GetDeviceCloudWebhooks)

		secure.GET("/tasks/stats", h.DeviceStatsHandler.GetTaskStats)

		secure.GET("/devices/:id/call-stats", h.DeviceStatsHandler.GetCallStats)
		secure.GET("/devices/stats", h.DeviceHandler.GetDeviceStats)
		secure.GET("/devices/:id/screenshots", h.DeviceStatsHandler.GetDeviceScreenshots)
	}
}
