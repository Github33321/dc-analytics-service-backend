package repository

import (
	"context"
	"dc-analytics-service-backend/internal/models"
	"errors"
	"fmt"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	GetUserByID(ctx context.Context, id int64) (*models.User, error)
	CreateUser(ctx context.Context, user *models.User) (*models.User, error)
	GetUsers(ctx context.Context) ([]models.User, error)
	DeleteUser(ctx context.Context, id int64) error
	//users stat
	GetDCUsers(ctx context.Context) (*[]models.DCUserProfile, error)
}

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetUserByID(ctx context.Context, id int64) (*models.User, error) {
	query := `
		SELECT id, email, username, password, role, is_active, 
		       created_at, updated_at, verified_at, last_login 
		FROM users 
		WHERE id = $1`
	var user models.User
	if err := pgxscan.Get(ctx, r.db, &user, query, id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	query := `
		INSERT INTO users (email, username, password, role, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING id, email, username, password, role, is_active, created_at, updated_at, verified_at, last_login`
	if err := pgxscan.Get(ctx, r.db, user, query,
		user.Email, user.Username, user.Password, user.Role, user.IsActive,
	); err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) GetUsers(ctx context.Context) ([]models.User, error) {
	query := `
		SELECT id, email, username, password, role, is_active, 
		       created_at, updated_at, verified_at, last_login 
		FROM users`
	var users []models.User
	if err := pgxscan.Select(ctx, r.db, &users, query); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			empty := []models.User{}
			return empty, nil
		}
		return nil, err
	}
	return users, nil
}

func (r *userRepository) DeleteUser(ctx context.Context, id int64) error {
	query := `DELETE FROM users WHERE id = $1`
	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("Пользователь с id %d не найден", id)
	}
	return nil
}

func (r *userRepository) GetDCUsers(ctx context.Context) (*[]models.DCUserProfile, error) {
	query := `
		SELECT user_id, username, email
		FROM users_dc
		ORDER BY user_id
	`

	var users []models.DCUserProfile
	if err := pgxscan.Select(ctx, r.db, &users, query); err != nil {
		return nil, err
	}

	return &users, nil
}
