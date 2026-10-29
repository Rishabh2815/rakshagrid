package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           string    `db:"id"`
	Email        string    `db:"email"`
	PasswordHash string    `db:"password_hash"`
	Role         string    `db:"role"` // user, responder, admin
	Name         string    `db:"name"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type LoginResponse struct {
	Token string `json:"token"`
	Role  string `json:"role"`
	Name  string `json:"name"`
}

type SignupRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Role     string `json:"role"`
}

type SignupResponse struct {
	Token string `json:"token"`
	Role  string `json:"role"`
	Name  string `json:"name"`
}

// NewUser creates a new user with a generated ID
func NewUser(email, passwordHash, role, name string) *User {
	now := time.Now()
	return &User{
		ID:           uuid.New().String(),
		Email:        email,
		PasswordHash: passwordHash,
		Role:         role,
		Name:         name,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}
