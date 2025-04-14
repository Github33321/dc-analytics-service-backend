package handler

import (
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
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /v1/analytics/servers [get]
func (h *ServerHandler) GetServers(c *gin.Context) {
	servers, err := h.ServerService.GetAllServers(c.Request.Context())
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
// @Failure 400 {object} map[string]string "Неверный формат ID"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
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
