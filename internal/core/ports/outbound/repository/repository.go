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

	GetCategoryByID(ctx context.Context, id uuid.UUID) (*domain.Category, error)
	GetCategoryIDByName(ctx context.Context, name *string) (*uuid.UUID, error)
	CreateCategory(ctx context.Context, input domain.Category) (*domain.Category, error)
	UpdateCategory(ctx context.Context, id uuid.UUID, input domain.Category) (*domain.Category, error)
	DeleteCategoryByID(ctx context.Context, id uuid.UUID) error
	ListCategories(ctx context.Context, pgn *pagination.Pagination) ([]domain.Category, error)
	CountCategories(ctx context.Context, pgn *pagination.Pagination) (int, error)

	GetTransactionByID(ctx context.Context, userID uuid.UUID, id uuid.UUID) (*dto.TransactionResponse, error)
	CreateTransaction(ctx context.Context, userID uuid.UUID, input domain.Transaction) (*dto.TransactionResponse, error)
	UpdateTransaction(ctx context.Context, userID uuid.UUID, id uuid.UUID, input domain.Transaction) (*dto.TransactionResponse, error)
	DeleteTransactionByID(ctx context.Context, userID uuid.UUID, id uuid.UUID) error
	ListTransactions(ctx context.Context, userID uuid.UUID, flt dto.TransactionFilters, pgn *pagination.Pagination) ([]dto.TransactionResponse, error)
	CountTransactions(ctx context.Context, userID uuid.UUID, flt dto.TransactionFilters, pgn *pagination.Pagination) (int, error)
	TransactionsSummary(ctx context.Context, userID uuid.UUID, flt dto.ChartFilters) ([]dto.SummaryByDate, error)
	TransactionsGeneralStats(ctx context.Context, userID uuid.UUID, flt dto.ChartFilters) (*dto.TransactionStatsSummary, error)

	GetInvoiceByID(ctx context.Context, userID uuid.UUID, id uuid.UUID) (*dto.InvoiceResponse, error)
	CreateInvoice(ctx context.Context, userID uuid.UUID, input domain.Invoice) (*dto.InvoiceResponse, error)
	UpdateInvoice(ctx context.Context, userID uuid.UUID, id uuid.UUID, input domain.Invoice) (*dto.InvoiceResponse, error)
	DeleteInvoiceByID(ctx context.Context, userID uuid.UUID, id uuid.UUID) error
	ListInvoices(ctx context.Context, userID uuid.UUID, flt dto.InvoiceFilters, pgn *pagination.Pagination) ([]dto.InvoiceResponse, error)
	CountInvoices(ctx context.Context, userID uuid.UUID, flt dto.InvoiceFilters, pgn *pagination.Pagination) (int, error)

	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	GetUserByUsername(ctx context.Context, username string) (*domain.User, error)

	GetUserSessionByToken(ctx context.Context, token string) (*domain.UserSession, error)
	CreateUserSession(ctx context.Context, session domain.UserSession) error
	DeleteUserSession(ctx context.Context, session domain.UserSession) error
}
