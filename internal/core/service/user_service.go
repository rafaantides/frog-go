package service

import (
	"context"

	"frog-go/internal/core/domain"
	"frog-go/internal/core/ports/inbound"
	"frog-go/internal/core/ports/outbound/repository"
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
