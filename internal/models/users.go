package models

import "time"

type User struct {
	ID         int64      `json:"id"`
	Email      string     `json:"email"`
	Username   string     `json:"username"`
	Password   string     `json:"password"`
	Role       string     `json:"role"`
	IsActive   bool       `json:"is_active"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	VerifiedAt *time.Time `json:"verified_at,omitempty"`
	LastLogin  *time.Time `json:"last_login,omitempty"`
}
