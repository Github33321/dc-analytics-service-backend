package handler

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

type customError struct {
	Msg    string
	Status int
}

func (e *customError) Error() string {
	return e.Msg
}

func GlobalErrorHandler(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			var errs []error
			var customStatus int
			var customMsg string
			customStatus = http.StatusInternalServerError
			customMsg = "StatusInternalServerError"

			for _, e := range c.Errors {
				errs = append(errs, e.Err)
				if cerr, ok := e.Err.(*customError); ok {
					customStatus = cerr.Status
					customMsg = cerr.Msg
					break
				}
				lowerMsg := strings.ToLower(e.Error())
				if strings.Contains(lowerMsg, "invalid character") ||
					strings.Contains(lowerMsg, "cannot unmarshal") ||
					strings.Contains(lowerMsg, "unexpected end of json") {
					customStatus = http.StatusBadRequest
					customMsg = "StatusBadRequest"
					break
				}
			}

			logger.Error("Возникли ошибки при обработке запроса", zap.Errors("errors", errs))
			c.JSON(customStatus, gin.H{"error": customMsg})
			c.Abort()
		}
	}
}
