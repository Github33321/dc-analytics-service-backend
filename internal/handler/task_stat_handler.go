package handler

import (
	"dc-analytics-service-backend/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type TaskStatHandler struct {
	service service.TaskStatService
}

func NewTaskStatHandler(s service.TaskStatService) *TaskStatHandler {
	return &TaskStatHandler{service: s}
}

// GetTaskStats godoc
// @Summary      GetTaskStats
// @Description  Если параметр date не задан, возвращает статистику за все даты. Если задан, то только за указанный день.
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        date   query    string  false  "Дата в формате YYYY-MM-DDTHH:MM:SSZ"
// @Success      200    {array}  models.TaskStat
// @Failure      500    {object}  map[string]string  "Internal Server Error"
// @Router       /v1/analytics/tasks/stats [get]
func (h *TaskStatHandler) GetTaskStats(c *gin.Context) {
	ctx := c.Request.Context()

	date := c.Query("date")

	stats, err := h.service.GetTaskStats(ctx, date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Ошибка получения данных: " + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, stats)
}
