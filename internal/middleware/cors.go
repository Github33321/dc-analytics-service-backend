package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strings"
)

var allowedOrigins []string

func init() {
	allowedOrigins = GetAllowedOrigins()
}

func GetAllowedOrigins() []string {
	origins := os.Getenv("ALLOWED_ORIGINS")
	if origins == "" {
		return []string{"http://11.111.111.111"}
	}
	parts := strings.Split(origins, ",")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}

func isAllowedOrigin(origin string) bool {
	for _, o := range allowedOrigins {
		if origin == o {
			return true
		}
	}
	return false
}
func DynamicCORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
		}

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusOK)
			return
		}
		c.Next()
	}
}

//func DynamicCORSMiddleware() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		origin := c.Request.Header.Get("Origin")
//		if origin != "" && isAllowedOrigin(origin) {
//			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
//			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
//			c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
//			c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
//			c.Writer.Header().Set("Access-Control-Max-Age", "86400")
//			c.Writer.Header().Set("Vary", "Origin")
//		}
//
//		if c.Request.Method == http.MethodOptions {
//			c.AbortWithStatus(http.StatusOK)
//			return
//		}
//
//		c.Next()
//	}
//}
