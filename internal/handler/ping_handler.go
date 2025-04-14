package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// PingHandler проверяет доступность сервера.
// @Summary PingHandler
// @Description Возвращает "pong" если сервер работает
// @Tags ping
// @Produce plain
// @Success 200 {string} string "pong"
// @Failure 500 {string} string "Internal Server Error"
// @Security BearerAuth
// @Router /ping [get]
func PingHandler(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}
