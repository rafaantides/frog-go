package service

import (
	"context"

	"frog-go/internal/core/domain"
	"frog-go/internal/core/dto"
	"frog-go/internal/core/ports/inbound"
	"frog-go/internal/core/ports/outbound/repository"

	"github.com/google/uuid"
)

type userService struct {
	repo repository.Repository
}

func NewUserService(repo repository.Repository) inbound.UserService {
	return &userService{repo: repo}
}

func (s *userService) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	return s.repo.GetUserByEmail(ctx, email)
}

func (s *userService) GetUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	return s.repo.GetUserByUsername(ctx, username)
}

func (s *userService) CreateUser(ctx context.Context, input domain.User) (*dto.UserResponse, error) {
	return nil, nil
}

func (s *userService) GetUser(ctx context.Context, userID uuid.UUID) (*dto.UserResponse, error) {
	return nil, nil
}

func (s *userService) UpdateUserPassword(ctx context.Context, userID uuid.UUID, oldPassword string, newPassword string) error {
	return nil
}
