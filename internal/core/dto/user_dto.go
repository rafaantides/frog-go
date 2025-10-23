package dto

import (
	"frog-go/internal/core/domain"
	"frog-go/internal/utils"

	"github.com/google/uuid"
)

type UserRequest struct {
	Name     string `json:"name" binding:"required"`
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	IsActive bool   `json:"is_active"`
}

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	IsActive  bool      `json:"is_active"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
}

func (r *UserRequest) ToDomain() (*domain.User, error) {

	passwordHash, err := utils.HashPassword(r.Password)

	if err != nil {
		return nil, err
	}

	return domain.NewUser(
		r.Name,
		r.Username,
		r.Email,
		string(passwordHash),
		r.IsActive,
	)
}
