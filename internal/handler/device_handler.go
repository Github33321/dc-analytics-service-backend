package handler

import (
	"net/http"
	"strconv"
	"strings"

	"dc-analytics-service-backend/internal/service"

	"github.com/gin-gonic/gin"
)

type DeviceHandler struct {
	DeviceService service.DeviceService
}

func NewDeviceHandler(deviceService service.DeviceService) *DeviceHandler {
	return &DeviceHandler{
		DeviceService: deviceService,
	}
}

// GetDevices godoc
// @Summary      GetDevices
// @Description  Возвращает список устройств и общее количество устройств (size). Используйте параметры page и limit.
// @Tags         devices
// @Accept       json
// @Produce      json
// @Param        page  query     int  false  "Номер страницы" default(1)
// @Param        limit query     int  false  "Количество элементов на страницу" default(10)
// @Success      200   {object}  models.PaginatedDevices
// @Failure      400   {object}  models.ErrorResponse  "Неверный формат параметров"
// @Failure      500   {object}  models.ErrorResponse  "Внутренняя ошибка сервера"
// @Security BearerAuth
// @Router       /v1/analytics/devices [get]
func (h *DeviceHandler) GetDevices(c *gin.Context) {
	pageStr := c.Query("page")
	limitStr := c.Query("limit")

	page := 1
	limit := 10
	var err error

	if pageStr != "" {
		page, err = strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			//c.JSON(http.StatusBadRequest, gin.H{"error": "BadRequest"})
			c.Error(err)
			return
		}
	}

	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil || limit < 1 {
			//c.JSON(http.StatusBadRequest, gin.H{"error": "BadRequest"})
			c.Error(err)
			return
		}
	}

	response, err := h.DeviceService.GetDevices(c.Request.Context(), page, limit)
	if err != nil {
		//c.JSON(http.StatusInternalServerError, gin.H{"error": "InternalServerError"})
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetDeviceByID godoc
// @Summary      GetDeviceByID
// @Description  Ищет устройство в базе данных и возвращает, если найдено
// @Tags         devices
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID устройства"
// @Success      200  {object}  models.Device
// @Failure      400  {object}  models.ErrorResponse  "Неверный формат ID"
// @Failure      404  {object}  models.ErrorResponse  "Устройство не найдено"
// @Failure      500  {object}  models.ErrorResponse  "Internal Server Error"
// @Security BearerAuth
// @Router       /v1/analytics/devices/{id} [get]
func (h *DeviceHandler) GetDeviceByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		//c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат ID"})
		c.Error(err)
		return
	}
	device, err := h.DeviceService.GetDeviceByID(c.Request.Context(), id)
	if err != nil {
		//c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}
	if device == nil {
		//c.JSON(http.StatusNotFound, gin.H{"message": "Устройство не найдено"})
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, device)
}

// UpdateDevice godoc
// @Summary      UpdateDevice
// @Description  Обновляет поля устройства, переданные в теле запроса (PATCH)
// @Tags         devices
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID устройства"
// @Param        device  body    service.UpdateDeviceRequest  true  "Данные для обновления"
// @Success      200  {object}  models.Device
// @Failure      400  {object}  models.ErrorResponse  "Неверный формат ID или некорректные данные обновления"
// @Failure      500  {object}  models.ErrorResponse  "Internal Server Error"
// @Security BearerAuth
// @Router       /v1/analytics/devices/{id} [patch]
func (h *DeviceHandler) UpdateDevice(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		//c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат ID"})
		c.Error(err)
		return
	}

	var updateReq service.UpdateDeviceRequest
	if err := c.ShouldBindJSON(&updateReq); err != nil {
		//c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные данные для обновления устройства"})
		c.Error(err)
		return
	}

	device, err := h.DeviceService.UpdateDevice(c.Request.Context(), id, updateReq)
	if err != nil {
		//c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, device)
}

// DeleteDevice godoc
// @Summary      DeleteDevice
// @Description  Удаляет устройство, если оно существует
// @Tags         devices
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID устройства"
// @Success      200  {object}  models.MessageResponse  "Устройство удалено"
// @Failure      400  {object}  models.ErrorResponse   "Неверный формат ID"
// @Failure      404  {object}  models.ErrorResponse   "Устройство не найдено"
// @Failure      500  {object}  models.ErrorResponse   "Internal Server Error"
// @Security BearerAuth
// @Router       /v1/analytics/devices/{id} [delete]
func (h *DeviceHandler) DeleteDevice(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		//c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат ID"})
		c.Error(err)
		return
	}

	err = h.DeviceService.DeleteDevice(c.Request.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "не найдено") {
			//c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			c.Error(err)
			return
		}
		//c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Устройство удалено"})
}

// GetDeviceStats godoc
// @Summary      GetDeviceStats
// @Description  Возвращает общее количество устройств, количество устройств с платформой android, ios, Pixel и устройств с smart_call_hiya == 1
// @Tags         devices
// @Accept       json
// @Produce      json
// @Success      200  {object}  models.DeviceStatsResponse
// @Failure      500  {object}  models.ErrorResponse  "Внутренняя ошибка сервера"
// @Security BearerAuth
// @Router       /v1/analytics/devices/stats [get]
func (h *DeviceHandler) GetDeviceStats(c *gin.Context) {
	stats, err := h.DeviceService.GetDeviceStats(c.Request.Context())
	if err != nil {
		//c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, stats)
}
