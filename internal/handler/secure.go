package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func SecureHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Доступ разрешён по JWT",
	})
}
