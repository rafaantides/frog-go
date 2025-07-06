package inbound

import (
	"context"
	"frog-go/internal/core/domain"
	"frog-go/internal/core/dto"
	"frog-go/internal/utils/pagination"

	"github.com/google/uuid"
)

type CategoryService interface {
	GetCategoryByID(ctx context.Context, id uuid.UUID) (*dto.CategoryResponse, error)
	CreateCategory(ctx context.Context, input domain.Category) (*dto.CategoryResponse, error)
	UpdateCategory(ctx context.Context, id uuid.UUID, input domain.Category) (*dto.CategoryResponse, error)
	DeleteCategoryByID(ctx context.Context, id uuid.UUID) error
	ListCategories(ctx context.Context, pgn *pagination.Pagination) ([]dto.CategoryResponse, int, error)
}

type DebtService interface {
	GetDebtByID(ctx context.Context, id uuid.UUID) (*dto.DebtResponse, error)
	CreateDebt(ctx context.Context, input domain.Debt) (*dto.DebtResponse, error)
	UpdateDebt(ctx context.Context, id uuid.UUID, input domain.Debt) (*dto.DebtResponse, error)
	DeleteDebtByID(ctx context.Context, id uuid.UUID) error
	ListDebts(ctx context.Context, flt dto.DebtFilters, pgn *pagination.Pagination) ([]dto.DebtResponse, int, error)
}
