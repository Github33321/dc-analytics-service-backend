package handler

import (
	"dc-analytics-service-backend/internal/middleware"
	"dc-analytics-service-backend/internal/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	UserHandler               *UserHandler
	DeviceHandler             *DeviceHandler
	DeviceCloudWebhookHandler *DeviceCloudWebhookHandler
}

func NewHandler(userService service.UserService, deviceService service.DeviceService, clickhouseService service.ClickhouseService) *Handler {
	return &Handler{
		UserHandler:               NewUserHandler(userService),
		DeviceHandler:             NewDeviceHandler(deviceService),
		DeviceCloudWebhookHandler: NewDeviceCloudWebhookHandler(clickhouseService),
	}
}

func (h *Handler) InitRoutes(router *gin.Engine, jwtSecret string) {
	router.Use(middleware.DynamicCORSMiddleware())
	router.POST("/login", LoginHandler)

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

		secure.GET("/deviceCloudWebhooks", h.DeviceCloudWebhookHandler.GetDeviceCloudWebhooks)
	}
}
