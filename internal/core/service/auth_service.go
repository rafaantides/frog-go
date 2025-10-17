package service

import (
	"context"
	"time"

	"frog-go/internal/config"
	"frog-go/internal/core/domain"
	appError "frog-go/internal/core/errors"
	"frog-go/internal/core/ports/inbound"
	"frog-go/internal/core/ports/outbound/repository"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type authService struct {
	repo repository.Repository
}

func NewAuthService(repo repository.Repository) inbound.AuthService {
	return &authService{repo: repo}
}

func (s *authService) GenerateToken(ctx context.Context, userID uuid.UUID, duration time.Duration) (string, error) {
	claims := domain.Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)), // expiração
			IssuedAt:  jwt.NewNumericDate(time.Now()),               // data de emissão
		},
	}

	// Cria o token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Assina com a secret
	return token.SignedString(config.JwtSecret)
}

func (s *authService) ValidateToken(tokenString string) (*domain.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &domain.Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Garante que está usando HS256
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, appError.ErrUnexpectedSigningMethod
		}
		return config.JwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	// Valida claims
	if claims, ok := token.Claims.(*domain.Claims); ok && token.Valid {
		if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
			return nil, appError.ErrTokenExpired
		}
		return claims, nil
	}

	return nil, appError.ErrInvalidToken
}
