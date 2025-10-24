package domain

import (
	appError "frog-go/internal/core/errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	jwt.RegisteredClaims
}

type UserSession struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Token     string    `json:"token"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

func NewUserSession(userID uuid.UUID, token string) (*UserSession, error) {
	if userID == uuid.Nil {
		return nil, appError.EmptyField("user_id")
	}
	if token == "" {
		return nil, appError.EmptyField("token")
	}

	return &UserSession{
		UserID:    userID,
		Token:     token,
		IsActive:  true,
	}, nil
}
