package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"os"
	"time"
)

type LoginRequest struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}
type AuthResponse struct {
	Token   string `json:"token" example:"eyJhbGciOiJIUz..."`
	Message string `json:"message" example:"Успешная авторизация"`
}

// LoginHandler godoc
// @Summary      LoginHandler
// @Description  Принимает логин и пароль, возвращает JWT-токен при успехе
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        loginData  body    LoginRequest  true  "Логин и пароль"
// @Success      200  {object}  AuthResponse  "Возвращает поле token и message"
// @Failure      400  {object}  models.ErrorResponse   "Неверный формат запроса"
// @Failure      401  {object}  models.ErrorResponse   "Неверные учетные данные"
// @Failure      500  {object}  models.ErrorResponse   "Ошибка генерации токена"
// @Router       /login [post]
func LoginHandler(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		//c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат запроса"})
		c.Error(err)
		return
	}

	expectedLogin := os.Getenv("LOGIN")
	if expectedLogin == "" {
		expectedLogin = "login"
	}
	expectedPassword := os.Getenv("PASSWORD")
	if expectedPassword == "" {
		expectedPassword = "password"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "default_super_secret"
	}

	if req.Login != expectedLogin || req.Password != expectedPassword {
		//c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверные учетные данные"})
		c.Error(&customError{Msg: "Неверные учетные данные", Status: http.StatusUnauthorized})
		return
	}

	claims := jwt.MapClaims{
		"sub": req.Login,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		//c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка генерации токена"})
		c.Error(&customError{Msg: "Ошибка генерации токена", Status: http.StatusInternalServerError})
		return
	}

	c.SetCookie("token", tokenString, 24*3600, "/", "", false, false)

	c.JSON(http.StatusOK, gin.H{
		"token":   tokenString,
		"message": "Авторизация прошла успешно",
	})
}
