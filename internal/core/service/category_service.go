package service

import (
	"context"

	"frog-go/internal/core/domain"
	"frog-go/internal/core/ports/inbound"
	"frog-go/internal/core/ports/outbound/repository"
	"frog-go/internal/utils/pagination"

	"github.com/google/uuid"
)

type categoryService struct {
	repo repository.Repository
}

func NewCategoryService(repo repository.Repository) inbound.CategoryService {
	return &categoryService{repo: repo}
}
func (s *categoryService) GetCategoryByID(ctx context.Context, id uuid.UUID) (*domain.Category, error) {
	return s.repo.GetCategoryByID(ctx, id)
}

func (s *categoryService) CreateCategory(ctx context.Context, input domain.Category) (*domain.Category, error) {
	return s.repo.CreateCategory(ctx, input)
}

func (s *categoryService) UpdateCategory(ctx context.Context, id uuid.UUID, input domain.Category) (*domain.Category, error) {
	return s.repo.UpdateCategory(ctx, id, input)
}

func (s *categoryService) DeleteCategoryByID(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteCategoryByID(ctx, id)
}

func (s *categoryService) ListCategories(ctx context.Context, pgn *pagination.Pagination) ([]domain.Category, int, error) {
	data, err := s.repo.ListCategories(ctx, pgn)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.repo.CountCategories(ctx, pgn)
	if err != nil {
		return nil, 0, err
	}

	return data, total, nil
}
