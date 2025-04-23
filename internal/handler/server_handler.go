package handler

import (
	"dc-analytics-service-backend/internal/models"
	"go.uber.org/zap"
	"math"
	"net/http"
	"strconv"
	"strings"

	"dc-analytics-service-backend/internal/service"

	"github.com/gin-gonic/gin"
)

type ServerHandler struct {
	Logger        *zap.Logger
	ServerService service.ServerService
}

func NewServerHandler(logger *zap.Logger, s service.ServerService) *ServerHandler {
	return &ServerHandler{
		Logger:        logger,
		ServerService: s,
	}
}

// GetServers godoc
// @Summary GetServers
// @Description Возвращает все сервера или постранично с total_pages
// @Tags servers
// @Accept json
// @Produce json
// @Param page query int false "Страница"
// @Param limit query int false "Количество на странице"
// @Success 200 {object} models.ServersResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Security BearerAuth
// @Router /v1/analytics/servers [get]
func (h *ServerHandler) GetServers(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "")
	pageStr := c.DefaultQuery("page", "")

	ctx := c.Request.Context()

	if limitStr == "" || pageStr == "" {
		servers, _, err := h.ServerService.GetAllServers(ctx, 1000, 0)
		if err != nil {
			handleClientError(c, h.Logger, http.StatusInternalServerError, "StatusInternalServerError", err)
			return
		}
		c.JSON(http.StatusOK, models.ServersResponse{Servers: servers})
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		handleClientError(c, h.Logger, http.StatusBadRequest, "BadRequest", err)
		return
	}
	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		handleClientError(c, h.Logger, http.StatusBadRequest, "BadRequest", err)
		return
	}
	offset := (page - 1) * limit

	servers, total, err := h.ServerService.GetAllServers(ctx, limit, offset)
	if err != nil {
		handleClientError(c, h.Logger, http.StatusInternalServerError, "StatusInternalServerError", err)
		return
	}
	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	c.JSON(http.StatusOK, models.ServersResponse{
		Servers:    servers,
		TotalPages: totalPages,
	})
}

// GetServerByID godoc
// @Summary GetServerByID
// @Description Возвращает сервер по заданному ID
// @Tags servers
// @Param id path int true "ID сервера"
// @Produce json
// @Success 200 {object} models.Server
// @Failure 400 {object} models.ErrorResponse  "Неверный формат ID"
// @Failure 500 {object} models.ErrorResponse  "Внутренняя ошибка сервера"
// @Security BearerAuth
// @Router /v1/analytics/servers/{id} [get]
func (h *ServerHandler) GetServerByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		handleClientError(c, h.Logger, http.StatusBadRequest, "BadRequest", err)
		return
	}
	server, err := h.ServerService.GetServerByID(c.Request.Context(), id)
	if err != nil {
		handleClientError(c, h.Logger, http.StatusInternalServerError, "StatusInternalServerError", err)
		c.Error(err)
		return
	}
	if server == nil {
		handleClientError(c, h.Logger, http.StatusBadRequest, "BadRequest", err)
		return
	}
	c.JSON(http.StatusOK, server)
}

// UpdateServer godoc
// @Summary UpdateServer
// @Description Обновляет запись сервера с заданным ID. Передается JSON с полями для обновления.
// @Tags servers
// @Accept json
// @Produce json
// @Param id path int true "ID сервера"
// @Param server body models.UpdateServerRequest true "Данные для обновления сервера"
// @Success 200 {object} models.Server
// @Failure 400 {object} models.ErrorResponse  "Неверный формат запроса или ID"
// @Failure 404 {object} models.ErrorResponse  "Сервер не найден"
// @Failure 500 {object} models.ErrorResponse  "Внутренняя ошибка сервера"
// @Security BearerAuth
// @Router /v1/analytics/servers/{id} [put]
func (h *ServerHandler) UpdateServer(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		handleClientError(c, h.Logger, http.StatusBadRequest, "BadRequest", err)
		return
	}

	var req models.UpdateServerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleClientError(c, h.Logger, http.StatusBadRequest, "BadRequest", err)
		return
	}

	updated, err := h.ServerService.UpdateServer(c.Request.Context(), id, req)
	if err != nil {
		if strings.Contains(err.Error(), "не найден") {
			handleClientError(c, h.Logger, http.StatusNotFound, "NotFound", err)
			return
		}
		handleClientError(c, h.Logger, http.StatusInternalServerError, "StatusInternalServerError", err)
		return
	}

	c.JSON(http.StatusOK, updated)
}

// GetDevices godoc
// @Summary     GetDevicesByServerID
// @Description Возвращает устройства сервера по его ID с пагинацией.
// @Tags        servers
// @Accept      json
// @Produce     json
// @Param       id     path      int  true   "ID сервера"
// @Param       limit  query     int  false  "Размер страницы"   default(10)
// @Param       page   query     int  false  "Номер страницы"    default(1)
// @Success     200    {array}   models.Device
// @Failure     400    {object}  models.ErrorResponse  "Неверный формат параметров"
// @Failure     500    {object}  models.ErrorResponse  "Внутренняя ошибка сервера"
// @Security    BearerAuth
// @Router      /v1/analytics/servers/{id}/devices [get]
func (h *ServerHandler) GetDevices(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		handleClientError(c, h.Logger, http.StatusBadRequest, "BadRequest", err)
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit < 1 {
		handleClientError(c, h.Logger, http.StatusBadRequest, "BadRequest", err)
		return
	}
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		handleClientError(c, h.Logger, http.StatusBadRequest, "BadRequest", err)
		return
	}
	offset := (page - 1) * limit

	devices, err := h.ServerService.GetDevicesByServerID(c.Request.Context(), id, limit, offset)
	if err != nil {
		handleClientError(c, h.Logger, http.StatusInternalServerError, "StatusInternalServerError", err)
		return
	}
	if devices == nil {
		devices = []models.Device{}
	}

	c.JSON(http.StatusOK, devices)
}
