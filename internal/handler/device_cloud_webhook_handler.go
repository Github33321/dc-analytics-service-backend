package handler

import (
	"dc-analytics-service-backend/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DeviceCloudWebhookHandler struct {
	clickhouseService service.ClickhouseService
}

func NewDeviceCloudWebhookHandler(s service.ClickhouseService) *DeviceCloudWebhookHandler {
	return &DeviceCloudWebhookHandler{
		clickhouseService: s,
	}
}

func (h *DeviceCloudWebhookHandler) GetDeviceCloudWebhooks(c *gin.Context) {
	ctx := c.Request.Context()
	results, err := h.clickhouseService.GetResults(ctx)
	if err != nil {
		//c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения данных: " + err.Error()})
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, results)
}
