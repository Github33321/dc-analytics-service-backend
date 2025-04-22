package handler

import (
	"dc-analytics-service-backend/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type DeviceStatsHandler struct {
	statsService service.DeviceStatsService
}

func NewDeviceStatsHandler(s service.DeviceStatsService) *DeviceStatsHandler {
	return &DeviceStatsHandler{statsService: s}
}

// GetDeviceCallStats godoc
// @Summary     GetDeviceCallStats
// @Description Возвращает:
//   - today_calls: число звонков за сегодня
//   - calls_by_day: разбивку по последним 31 дню (или по конкретной дате)
//   - status_counts: количество звонков по статусам за сегодня
//
// @Tags        stats
// @Accept      json
// @Produce     json
// @Param       id    path     string true  "ID устройства"
// @Param       date  query    string false "Дата в формате YYYY-MM-DD (если не указана — последние 31 день)"
// @Success     200   {object} models.DeviceCallStatsResponse
// @Failure     400   {object} models.ErrorResponse "Неверный формат параметров"
// @Failure     404   {object} models.ErrorResponse "Устройство не найдено"
// @Failure     500   {object} models.ErrorResponse "Внутренняя ошибка сервера"
// @Security    BearerAuth
// @Router      /v1/analytics/devices/{id}/call-stats [get]
func (h *DeviceStatsHandler) GetDeviceCallStats(c *gin.Context) {
	deviceID := c.Param("id")
	date := c.DefaultQuery("date", "")

	resp, err := h.statsService.GetDeviceCallStats(c.Request.Context(), deviceID, date)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, resp)
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

// GetDeviceCarrierStats godoc
// @Summary     GetDeviceCarrierStats
// @Description Возвращает для каждого оператора общее число устройств.
// @Tags        stats
// @Accept      json
// @Produce     json
// @Success     200 {object} map[string][]models.DeviceCarrierStats "device_carriers"
// @Failure     500 {object} models.ErrorResponse            "Внутренняя ошибка сервера"
// @Security    BearerAuth
// @Router      /v1/analytics/device-carriers/operator [get]
func (h *DeviceStatsHandler) GetDeviceCarrierStats(c *gin.Context) {
	stats, err := h.statsService.GetDeviceCarrierStats(c)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"device_carriers": stats,
	})
}

// GetOriginatingCarrierStats godoc
// @Summary GetOriginatingCarrierStats
// @Description Получает статистику проверок по операторам за выбранный период (если дата не указана, возвращает данные за сегодня).
// @Tags stats
// @Accept json
// @Produce json
// @Param fromDate query string false "Дата начала периода (формат YYYY-MM-DD, по умолчанию — сегодня)"
// @Success 200 {object} models.CarrierStatsResponse
// @Failure 400 {object} models.ErrorResponse "Некорректный запрос"
// @Failure 500 {object} models.ErrorResponse "Внутренняя ошибка сервера"
// @Security BearerAuth
// @Router /v1/analytics/device-carriers/originating [get]
func (h *DeviceStatsHandler) GetOriginatingCarrierStats(c *gin.Context) {
	fromDate := c.Query("fromDate")
	if fromDate == "" {
		fromDate = time.Now().UTC().Format("2006-01-02")
	}

	resp, err := h.statsService.GetOriginatingCarrierStats(c.Request.Context(), fromDate)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetDeviceGroupStats godoc
// @Summary     GetDeviceGroupStats
// @Description Возвращает количество результатов по каждому group_id
// @Tags        stats
// @Produce     json
// @Success     200 {array} models.DeviceGroupStats
// @Failure     500 {object} models.ErrorResponse
// @Security    BearerAuth
// @Router      /v1/analytics/device-carriers/group [get]
func (h *DeviceStatsHandler) GetDeviceGroupStats(c *gin.Context) {
	stats, err := h.statsService.GetDeviceGroupStats(c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, stats)
}

// GetSources godoc
// @Summary    	GetSources
// @Description Возвращает список источников по source_type_id
// @Tags        stats
// @Produce     json
// @Success     200 {array} models.Source
// @Failure     500 {object} models.ErrorResponse
// @Security    BearerAuth
// @Router      /v1/analytics/device-carriers/source [get]
func (h *DeviceStatsHandler) GetSources(c *gin.Context) {
	sources, err := h.statsService.GetSources(c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, sources)
}

// GetTasksReadyCounts godoc
// @Summary     GetTasksReadyCounts
// @Description Возвращает количество записей с group_id 1 и 2
// @Tags        tasks
// @Produce     json
// @Success     200 {array} models.TasksReadyCount
// @Failure     500 {object} models.ErrorResponse
// @Security    BearerAuth
// @Router      /v1/analytics/tasks/group [get]
func (h *DeviceStatsHandler) GetTasksReadyCounts(c *gin.Context) {
	stats, err := h.statsService.GetTasksReadyCounts(c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, stats)
}

// GetByUserID godoc
// @Summary     GetByUserID
// @Description Возвращает все записи по user_id
// @Tags       	tasks
// @Produce     json
// @Param       id   path      int  true  "User ID"
// @Success     200  {array}   models.DedicatedDevice
// @Failure     400  {object}  models.ErrorResponse "Неверный формат user_id"
// @Failure     500  {object}  models.ErrorResponse "Внутренняя ошибка сервера"
// @Security    BearerAuth
// @Router      /v1/analytics/tasks/{id}/devices [get]
func (h *DeviceStatsHandler) GetByUserID(c *gin.Context) {
	idStr := c.Param("id")
	uid, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.Error(err)
		return
	}

	list, err := h.statsService.GetByUserID(c.Request.Context(), uid)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, list)
}

// GetCountedUsers godoc
// @Summary     GetCountedUsers
// @Description Возвращает всех user и число их проверенных номеров
// @Tags        tasks
// @Produce     json
// @Success     200 {array} models.DCUser
// @Failure     500 {object} models.ErrorResponse "Внутренняя ошибка сервера"
// @Security    BearerAuth
// @Router      /v1/analytics/tasks/users [get]
func (h *DeviceStatsHandler) GetCountedUsers(c *gin.Context) {
	list, err := h.statsService.GetCountedUsers(c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, list)
}

// GetTodayStats godoc
// @Summary     GetTodayStats
// @Description Для каждого user и group выдаёт сколько раз их проверили сегодня
// @Tags        tasks
// @Produce     json
// @Success     200 {array} models.UserGroupStats
// @Failure     500 {object} models.ErrorResponse "Внутренняя ошибка сервера"
// @Security    BearerAuth
// @Router      /v1/analytics/tasks/users-today [get]
func (h *DeviceStatsHandler) GetTodayStats(c *gin.Context) {
	stats, err := h.statsService.TodayStats(c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, stats)
}

// GetDistinctUsers godoc
// @Summary     GetDistinctUsers
// @Description Возвращает список всех users
//
//	Если передан query-параметр date=YYYY-MM-DD — только за эту дату.
//
// @Tags        tasks
// @Accept      json
// @Produce     json
// @Param       date  query    string false "Дата в формате YYYY-MM-DD"
// @Success     200   {array}   models.DCUser2
// @Failure     500   {object}  models.ErrorResponse "Внутренняя ошибка сервера"
// @Security    BearerAuth
// @Router      /v1/analytics/tasks/user-list [get]
func (h *DeviceStatsHandler) GetDistinctUsers(c *gin.Context) {
	date := c.Query("date")
	users, err := h.statsService.GetDistinctUsers(c.Request.Context(), date)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, users)
}
