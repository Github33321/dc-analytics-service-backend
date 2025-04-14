package handler

import (
	"net/http"
	"strconv"

	"dc-analytics-service-backend/internal/service"
	"github.com/gin-gonic/gin"
)

type DeviceStatsHandler struct {
	statsService service.DeviceStatsService
}

func NewDeviceStatsHandler(s service.DeviceStatsService) *DeviceStatsHandler {
	return &DeviceStatsHandler{statsService: s}
}

// GetCallStats godoc
// @Summary      GetCallStats
// @Description  Возвращает статистику звонков для указанного устройства. Ответ включает:
//   - today_calls: количество звонков за сегодня,
//   - calls_by_day: массив объектов, где каждый объект содержит дату (created_at_str в формате "YYYY-MM-DD") и количество звонков за этот день,
//   - status_counts: массив объектов, где для каждого из ожидаемых статусов (call_failed, call_mismatch, wait, no_result, success) возвращается агрегированное количество звонков.
//
// Если параметр date указан, статистика будет возвращена только для этой даты, иначе агрегируются данные по всем датам.
// @Tags         devices
// @Accept       json
// @Produce      json
// @Param        id    path      int     true  "ID устройства"
// @Param        date  query     string  false "Дата для фильтрации (формат YYYY-MM-DD)"
// @Success      200   {object}  models.DeviceCallStatsResponse  "Агрегированная статистика звонков устройства"
// @Failure      500   {object}  map[string]string               "Внутренняя ошибка сервера"
// @Security BearerAuth
// @Router       /v1/analytics/devices/{id}/call-stats [get]
func (h *DeviceStatsHandler) GetCallStats(c *gin.Context) {
	idStr := c.Param("id")
	deviceID := idStr

	date := c.Query("date")

	stats, err := h.statsService.GetCallStats(c.Request.Context(), deviceID, date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения данных: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// GetTaskStats godoc
// @Summary      GetTaskStats
// @Description  Возвращает статистику звонков, сгруппированную по датам.
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        date  query     string  false  "Дата для фильтрации (YYYY-MM-DD). Если не указан, возвращаются данные по всем датам."
// @Success      200   {array}   models.TaskStat  "Массив агрегированных статистических данных"
// @Failure      500   {object}  map[string]string  "Внутренняя ошибка сервера"
// @Security BearerAuth
// @Router       /v1/analytics/tasks/stats [get]
func (h *DeviceStatsHandler) GetTaskStats(c *gin.Context) {
	ctx := c.Request.Context()

	date := c.Query("date")

	stats, err := h.statsService.GetTaskStats(ctx, date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Ошибка получения данных: " + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, stats)
}

// GetDeviceScreenshots godoc
// @Summary      GetDeviceScreenshots
// @Description  Возвращает последние скриншоты устройства с пагинацией.
// @Tags         devices
// @Accept       json
// @Produce      json
// @Param        id    path      int    true  "ID устройства"
// @Param        page  query     int    false "Номер страницы" default(1)
// @Param        limit query     int    false "Количество элементов на странице" default(10)
// @Success      200   {array}   models.DeviceScreenshot
// @Failure      400   {object}  map[string]string "Неверный формат параметров"
// @Failure      500   {object}  map[string]string "Внутренняя ошибка сервера"
// @Security BearerAuth
// @Router       /devices/{id}/screenshots [get]
func (h *DeviceStatsHandler) GetDeviceScreenshots(c *gin.Context) {
	deviceID := c.Param("id")
	page := 1
	limit := 10

	if p := c.Query("page"); p != "" {
		pInt, err := strconv.Atoi(p)
		if err != nil || pInt < 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат параметра page"})
			return
		}
		page = pInt
	}

	if l := c.Query("limit"); l != "" {
		lInt, err := strconv.Atoi(l)
		if err != nil || lInt < 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат параметра limit"})
			return
		}
		limit = lInt
	}

	screenshots, err := h.statsService.GetDeviceScreenshots(c.Request.Context(), deviceID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения данных: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, screenshots)
}
