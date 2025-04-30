package handler

import (
	"dc-analytics-service-backend/internal/service"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	Logger      *zap.Logger
	UserService service.UserService
}

func NewUserHandler(logger *zap.Logger, userService service.UserService) *UserHandler {
	return &UserHandler{
		Logger:      logger,
		UserService: userService}
}

// GetUserByID godoc
// @Summary      GetUserByID
// @Description  Возвращает данные пользователя, если он существует
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID пользователя"
// @Success      200  {object}  models.User
// @Failure      400  {object}  models.ErrorResponse  "Неверный формат ID"
// @Failure      404  {object}  models.ErrorResponse  "Пользователь не найден"
// @Failure      500  {object}  models.ErrorResponse  "Внутренняя ошибка сервера"
// @Security BearerAuth
// @Router       /v1/analytics/users/{id} [get]
func (h *UserHandler) GetUserByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		handleClientError(c, h.Logger, http.StatusBadRequest, "StatusBadRequest", err)
		return
	}
	user, err := h.UserService.GetUserByID(c, id)
	if err != nil {
		handleClientError(c, h.Logger, http.StatusInternalServerError, "StatusInternalServerError", err)
		return
	}

	if user == nil {
		handleClientError(c, h.Logger, http.StatusNotFound, "NotFound", err)
		return
	}

	c.JSON(http.StatusOK, user)
}

// CreateUser godoc
// @Summary      CreateUser
// @Description  Принимает данные и создает пользователя в системе
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        userData  body   service.CreateUserRequest  true  "Данные пользователя"
// @Success      201  {object}  models.User
// @Failure      400  {object}  models.ErrorResponse  "Неверные данные для создания пользователя"
// @Failure      500  {object}  models.ErrorResponse  "Внутренняя ошибка сервера"
// @Security BearerAuth
// @Router       /v1/analytics/users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var newUser service.CreateUserRequest
	if err := c.ShouldBindJSON(&newUser); err != nil {
		handleClientError(c, h.Logger, http.StatusBadRequest, "StatusBadRequest", err)
		return
	}

	user, err := h.UserService.CreateUser(c, newUser)
	if err != nil {
		handleClientError(c, h.Logger, http.StatusInternalServerError, "StatusInternalServerError", err)
		return
	}

	c.JSON(http.StatusCreated, user)
}

// GetUsers godoc
// @Summary      GetUsers
// @Description  Возвращает массив пользователей
// @Tags         users
// @Accept       json
// @Produce      json
// @Success      200  {array}   models.User
// @Failure      500  {object}  models.ErrorResponse  "Внутренняя ошибка сервера"
// @Security BearerAuth
// @Router       /v1/analytics/users [get]
func (h *UserHandler) GetUsers(c *gin.Context) {
	users, err := h.UserService.GetUsers(c)
	if err != nil {
		handleClientError(c, h.Logger, http.StatusInternalServerError, "StatusInternalServerError", err)
		return
	}

	c.JSON(http.StatusOK, users)
}

// DeleteUser godoc
// @Summary      DeleteUser
// @Description  Удаляет пользователя по его ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID пользователя"
// @Success      200  {object}  map[string]string "Пользователь удален"
// @Failure      400  {object}  models.ErrorResponse  "Неверный формат ID"
// @Failure      404  {object}  models.ErrorResponse  "Пользователь не найден"
// @Failure      500  {object}  models.ErrorResponse  "Внутренняя ошибка сервера"
// @Security BearerAuth
// @Router       /v1/analytics/users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		handleClientError(c, h.Logger, http.StatusBadRequest, "StatusBadRequest", err)
		return
	}

	err = h.UserService.DeleteUser(c.Request.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "не найден") {
			handleClientError(c, h.Logger, http.StatusNotFound, "NotFound", err)
			return
		}
		handleClientError(c, h.Logger, http.StatusInternalServerError, "StatusInternalServerError", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Пользователь удален"})
}

// GetDCUsers godoc
// @Summary     GetDCUsers
// @Description Возвращает список всех пользователей.
// @Tags        calls
// @Accept      json
// @Produce     json
// @Success     200 {array} models.DCUserProfile
// @Failure     500 {object} models.ErrorResponse "Внутренняя ошибка сервера"
// @Security    BearerAuth
// @Router      /v1/analytics/calls/user-list [get]
func (h *UserHandler) GetDCUsers(c *gin.Context) {
	ctx := c.Request.Context()

	users, err := h.UserService.GetDCUsers(ctx)
	if err != nil {
		handleClientError(c, h.Logger, http.StatusInternalServerError, "StatusInternalServerError", err)
		return
	}

	c.JSON(http.StatusOK, users)
}
