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
	ListCategories(ctx context.Context, kinds []string, pgn *pagination.Pagination) ([]dto.CategoryResponse, error)
	CountCategories(ctx context.Context, kinds []string, pgn *pagination.Pagination) (int, error)

	GetTransactionByID(ctx context.Context, id uuid.UUID) (*dto.TransactionResponse, error)
	CreateTransaction(ctx context.Context, input domain.Transaction) (*dto.TransactionResponse, error)
	UpdateTransaction(ctx context.Context, id uuid.UUID, input domain.Transaction) (*dto.TransactionResponse, error)
	DeleteTransactionByID(ctx context.Context, id uuid.UUID) error
	ListTransactions(ctx context.Context, flt dto.TransactionFilters, pgn *pagination.Pagination) ([]dto.TransactionResponse, error)
	CountTransactions(ctx context.Context, flt dto.TransactionFilters, pgn *pagination.Pagination) (int, error)
	TransactionsSummary(ctx context.Context, flt dto.ChartFilters) ([]dto.SummaryByDate, error)
	TransactionsGeneralStats(ctx context.Context, flt dto.ChartFilters) (*dto.TransactionStatsSummary, error)
}
