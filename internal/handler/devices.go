package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type SpamStats struct {
	SpamPercent float64 `json:"spam_percent"`
}

func SpamStatsHandler(c *gin.Context) {
	queryType := c.Query("type")
	if queryType != "spam" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unknown type"})
		return
	}

	query := `
        SELECT 
            (SUM(toInt8(spam)) * 100.0 / COUNT(*)) AS spam_percent
        FROM device_cloud_webhooks
    `

	var stats SpamStats
	err := db.QueryRow(query).Scan(&stats.SpamPercent)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stats)
}
