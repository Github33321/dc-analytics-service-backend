package handler

import (
	"dc-analytics-service-backend/internal/models"
	"net/http"
	"strconv"

	"dc-analytics-service-backend/internal/service"

	"github.com/gin-gonic/gin"
)

type ServerHandler struct {
	ServerService service.ServerService
}

func NewServerHandler(s service.ServerService) *ServerHandler {
	return &ServerHandler{
		ServerService: s,
	}
}

// GetServers godoc
// @Summary GetServers
// @Description Возвращает список всех серверов
// @Tags servers
// @Produce json
// @Success 200 {array} models.Server
// @Failure 500 {object} models.ErrorResponse
// @Security BearerAuth
// @Router /v1/analytics/servers [get]
func (h *ServerHandler) GetServers(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	pageStr := c.DefaultQuery("page", "1")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		//c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit"})
		c.Error(err)
		return
	}
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		//c.JSON(http.StatusBadRequest, gin.H{"error": "invalid page"})
		c.Error(err)
		return
	}
	offset := (page - 1) * limit
	servers, err := h.ServerService.GetAllServers(c.Request.Context(), limit, offset)
	if err != nil {
		//c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, servers)
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
		//c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат ID"})
		c.Error(err)
		return
	}
	server, err := h.ServerService.GetServerByID(c.Request.Context(), id)
	if err != nil {
		//c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}
	if server == nil {
		//c.JSON(http.StatusNotFound, gin.H{"error": "Сервер не найден"})
		c.Error(&customError{Msg: "Неверные учетные данные", Status: http.StatusNotFound})
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
		//c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат ID"})
		c.Error(err)
		return
	}

	var req models.UpdateServerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		//c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат запроса"})
		c.Error(err)
		return
	}

	updatedServer, err := h.ServerService.UpdateServer(c.Request.Context(), id, req)
	if err != nil {
		//c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, updatedServer)
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
		c.Error(err)
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit < 1 {
		c.Error(err)
		return
	}
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		c.Error(err)
		return
	}
	offset := (page - 1) * limit

	devices, err := h.ServerService.GetDevicesByServerID(c.Request.Context(), id, limit, offset)
	if err != nil {
		c.Error(err)
		return
	}
	if devices == nil {
		devices = []models.Device{}
	}

	c.JSON(http.StatusOK, devices)
}
