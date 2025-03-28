package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func JWTMiddleware(secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		fmt.Println("Authorization Header:", authHeader)

		if !strings.HasPrefix(authHeader, "Bearer ") {
			fmt.Println("⚠Header не начинается с Bearer")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid Authorization header"})
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		fmt.Println("Извлечён токен:", tokenString)

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			fmt.Println("Метод подписи:", token.Method.Alg())
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				fmt.Println("Метод подписи не HMAC")
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(secretKey), nil
		})
		if err != nil {
			fmt.Println("Ошибка при разборе токена:", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}
		if !token.Valid {
			fmt.Println("Токен недействителен")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}
		fmt.Println("Токен валидный, доступ разрешён")
		c.Next()
	}
}
