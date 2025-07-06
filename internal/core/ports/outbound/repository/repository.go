package repository

import (
	"context"
	"frog-go/internal/core/domain"
	"frog-go/internal/core/dto"
	"frog-go/internal/utils/pagination"

	"github.com/google/uuid"
)

type Repository interface {
	Close()

	GetCategoryByID(ctx context.Context, id uuid.UUID) (*dto.CategoryResponse, error)
	GetCategoryIDByName(ctx context.Context, name *string) (*uuid.UUID, error)
	CreateCategory(ctx context.Context, input domain.Category) (*dto.CategoryResponse, error)
	UpdateCategory(ctx context.Context, id uuid.UUID, input domain.Category) (*dto.CategoryResponse, error)
	DeleteCategoryByID(ctx context.Context, id uuid.UUID) error
	ListCategories(ctx context.Context, pgn *pagination.Pagination) ([]dto.CategoryResponse, error)
	CountCategories(ctx context.Context, pgn *pagination.Pagination) (int, error)

	GetDebtByID(ctx context.Context, id uuid.UUID) (*dto.DebtResponse, error)
	CreateDebt(ctx context.Context, input domain.Debt) (*dto.DebtResponse, error)
	UpdateDebt(ctx context.Context, id uuid.UUID, input domain.Debt) (*dto.DebtResponse, error)
	DeleteDebtByID(ctx context.Context, id uuid.UUID) error
	ListDebts(ctx context.Context, flt dto.DebtFilters, pgn *pagination.Pagination) ([]dto.DebtResponse, error)
	CountDebts(ctx context.Context, flt dto.DebtFilters, pgn *pagination.Pagination) (int, error)
	DebtsSummary(ctx context.Context, flt dto.ChartFilters) ([]dto.SummaryByDate, error)
	DebtsGeneralStats(ctx context.Context, flt dto.ChartFilters) (*dto.DebtStatsSummary, error)
}
