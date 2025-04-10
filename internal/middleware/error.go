package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

func GlobalErrorHandler(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			status := http.StatusInternalServerError
			message := "Внутренняя ошибка сервера"

			for _, e := range c.Errors {
				lowerMsg := strings.ToLower(e.Error())
				if strings.Contains(lowerMsg, "invalid character") ||
					strings.Contains(lowerMsg, "cannot unmarshal") ||
					strings.Contains(lowerMsg, "unexpected end of json") {
					status = http.StatusBadRequest
					message = "StatusBadRequest"
					break
				}
			}

			var errorsSlice []error
			for _, e := range c.Errors {
				errorsSlice = append(errorsSlice, e.Err)
			}
			logger.Error("Error(s) occurred", zap.Errors("errors", errorsSlice))

			c.JSON(status, gin.H{"error": message})
			c.Abort()
		}
	}
}
