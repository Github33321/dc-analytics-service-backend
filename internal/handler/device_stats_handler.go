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
// @Description  Возвращает статистику звонков для указанного устройства. Если параметр date не указан, агрегируются данные по всем датам; если указан, то только для указанной даты.
// @Tags         devices
// @Accept       json
// @Produce      json
// @Param        id    path      int     true  "ID устройства"
// @Param        date  query     string  false "Дата для фильтрации (формат YYYY-MM-DD)"
// @Success      200   {object}  models.DeviceCallStatsResponse  "Агрегированная статистика звонков устройства"
// @Failure      400   {object}  models.ErrorResponse            "Неверный формат запроса"
// @Failure      500   {object}  models.ErrorResponse            "Внутренняя ошибка сервера"
// @Security     BearerAuth
// @Router       /v1/analytics/devices/{id}/call-stats [get]
func (h *DeviceStatsHandler) GetCallStats(c *gin.Context) {
	idStr := c.Param("id")
	deviceID := idStr

	date := c.Query("date")

	stats, err := h.statsService.GetCallStats(c.Request.Context(), deviceID, date)
	if err != nil {
		//c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения данных: " + err.Error()})
		c.Error(err)
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
// @Success      200   {array}   models.TaskStat  "Массив статистических данных"
// @Failure      500   {object}  models.ErrorResponse   "Внутренняя ошибка сервера"
// @Security BearerAuth
// @Router       /v1/analytics/tasks/stats [get]
func (h *DeviceStatsHandler) GetTaskStats(c *gin.Context) {
	ctx := c.Request.Context()

	date := c.Query("date")

	stats, err := h.statsService.GetTaskStats(ctx, date)
	if err != nil {
		//c.JSON(http.StatusInternalServerError, gin.H{
		//	"error": "Ошибка получения данных: " + err.Error(),
		//})
		c.Error(err)
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
// @Failure      400   {object}  models.ErrorResponse  "Неверный формат параметров"
// @Failure      500   {object}  models.ErrorResponse  "Внутренняя ошибка сервера"
// @Security BearerAuth
// @Router       /v1/analytics/devices/{id}/screenshots [get]
func (h *DeviceStatsHandler) GetDeviceScreenshots(c *gin.Context) {
	deviceID := c.Param("id")
	page := 1
	limit := 10

	if p := c.Query("page"); p != "" {
		pInt, err := strconv.Atoi(p)
		if err != nil || pInt < 1 {
			//c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат параметра page"})
			c.Error(err)
			return
		}
		page = pInt
	}

	if l := c.Query("limit"); l != "" {
		lInt, err := strconv.Atoi(l)
		if err != nil || lInt < 1 {
			//c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат параметра limit"})
			c.Error(err)
			return
		}
		limit = lInt
	}

	screenshots, err := h.statsService.GetDeviceScreenshots(c.Request.Context(), deviceID, page, limit)
	if err != nil {
		//c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения данных: " + err.Error()})
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, screenshots)
}
