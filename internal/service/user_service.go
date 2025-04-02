package service

import (
	"context"
	"dc-analytics-service-backend/internal/models"
)

type UserService interface {
	GetUserByID(ctx context.Context, id int64) (*models.User, error)
	CreateUser(ctx context.Context, req CreateUserRequest) (*models.User, error)
	GetUsers(ctx context.Context) ([]models.User, error)
}

type CreateUserRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
	IsActive bool   `json:"is_active"`
}
