package middleware

import (
	"dc-analytics-service-backend/internal/repository"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func DeviceExistsMiddleware(repo repository.DeviceRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid device ID"})
			c.Abort()
			return
		}

		exists, err := repo.ExistsByID(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			c.Abort()
			return
		}
		if !exists {
			c.JSON(http.StatusNotFound, gin.H{"error": "Device not found"})
			c.Abort()
			return
		}

		c.Next()
	}
}
