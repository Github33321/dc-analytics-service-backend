package service

import (
	"context"
	"dc-analytics-service-backend/internal/models"
	"dc-analytics-service-backend/internal/repository"
)

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) GetUserByID(ctx context.Context, id int64) (*models.User, error) {
	return s.userRepo.GetUserByID(ctx, id)
}

func (s *userService) CreateUser(ctx context.Context, req CreateUserRequest) (*models.User, error) {
	user := &models.User{
		Email:    req.Email,
		Username: req.Username,
		Password: req.Password,
		Role:     req.Role,
		IsActive: req.IsActive,
	}

	return s.userRepo.CreateUser(ctx, user)
}

func (s *userService) GetUsers(ctx context.Context) ([]models.User, error) {
	return s.userRepo.GetUsers(ctx)
}
func (s *userService) DeleteUser(ctx context.Context, id int64) error {
	return s.userRepo.DeleteUser(ctx, id)
}
