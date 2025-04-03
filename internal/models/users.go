package models

import "time"

type User struct {
	ID         int64      `json:"id" db:"id"`
	Email      string     `json:"email" db:"email"`
	Username   string     `json:"username" db:"username"`
	Password   string     `json:"password" db:"password"`
	Role       string     `json:"role" db:"role"`
	IsActive   bool       `json:"is_active" db:"is_active"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at" db:"updated_at"`
	VerifiedAt *time.Time `json:"verified_at,omitempty" db:"verified_at"`
	LastLogin  *time.Time `json:"last_login,omitempty" db:"last_login"`
}
