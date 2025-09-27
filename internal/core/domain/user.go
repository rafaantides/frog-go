package domain

import (
	"frog-go/internal/core/errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"password_hash"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	jwt.RegisteredClaims
}

func NewUser(name, username, email, passwordHash string) (*User, error) {
	if username == "" {
		return nil, errors.EmptyField("username")
	}
	if passwordHash == "" {
		return nil, errors.EmptyField("password_hash")
	}

	return &User{
		ID:           uuid.New(),
		Name:         name,
		Username:     username,
		Email:        email,
		PasswordHash: passwordHash,
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}, nil
}
