package repository

import (
	"context"
	"dc-analytics-service-backend/internal/models"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	GetUserByID(ctx context.Context, id int64) (*models.User, error)
	CreateUser(ctx context.Context, user *models.User) (*models.User, error)
	GetUsers(ctx context.Context) ([]models.User, error)
	DeleteUser(ctx context.Context, id int64) error
}

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetUserByID(ctx context.Context, id int64) (*models.User, error) {
	query := `SELECT id, email, username, password, role, is_active, created_at, updated_at, verified_at, last_login FROM users WHERE id = $1`
	var user models.User
	err := r.db.GetContext(ctx, &user, query, id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	query := `INSERT INTO users (email, username, password, role, is_active, created_at, updated_at) 
			  VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP) 
			  RETURNING id, email, username, password, role, is_active, created_at, updated_at, verified_at, last_login`

	err := r.db.GetContext(ctx, user, query, user.Email, user.Username, user.Password, user.Role, user.IsActive)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) GetUsers(ctx context.Context) ([]models.User, error) {
	query := `SELECT id, email, username, password, role, is_active, created_at, updated_at, verified_at, last_login FROM users`
	var users []models.User
	err := r.db.SelectContext(ctx, &users, query)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userRepository) DeleteUser(ctx context.Context, id int64) error {
	query := `DELETE FROM users WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("Пользователь с id %d не найден", id)
	}
	return nil
}
