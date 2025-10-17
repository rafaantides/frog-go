package inbound

import (
	"context"
	"frog-go/internal/core/domain"
	"frog-go/internal/core/dto"
	"frog-go/internal/utils/pagination"
	"mime/multipart"
	"time"

	"github.com/google/uuid"
)

type CategoryService interface {
	GetCategoryByID(ctx context.Context, id uuid.UUID) (*dto.CategoryResponse, error)
	CreateCategory(ctx context.Context, input domain.Category) (*dto.CategoryResponse, error)
	UpdateCategory(ctx context.Context, id uuid.UUID, input domain.Category) (*dto.CategoryResponse, error)
	DeleteCategoryByID(ctx context.Context, id uuid.UUID) error
	ListCategories(ctx context.Context, pgn *pagination.Pagination) ([]dto.CategoryResponse, int, error)
}

type TransactionService interface {
	GetTransactionByID(ctx context.Context, id uuid.UUID) (*dto.TransactionResponse, error)
	CreateTransaction(ctx context.Context, input domain.Transaction) (*dto.TransactionResponse, error)
	UpdateTransaction(ctx context.Context, id uuid.UUID, input domain.Transaction) (*dto.TransactionResponse, error)
	DeleteTransactionByID(ctx context.Context, id uuid.UUID) error
	ListTransactions(ctx context.Context, flt dto.TransactionFilters, pgn *pagination.Pagination) ([]dto.TransactionResponse, int, error)
	TransactionsSummary(ctx context.Context, flt dto.ChartFilters) ([]dto.SummaryByDate, error)
	TransactionsGeneralStats(ctx context.Context, flt dto.ChartFilters) (*dto.TransactionStatsSummary, error)
}

type InvoiceService interface {
	GetInvoiceByID(ctx context.Context, id uuid.UUID) (*dto.InvoiceResponse, error)
	CreateInvoice(ctx context.Context, input domain.Invoice) (*dto.InvoiceResponse, error)
	UpdateInvoice(ctx context.Context, id uuid.UUID, input domain.Invoice) (*dto.InvoiceResponse, error)
	DeleteInvoiceByID(ctx context.Context, id uuid.UUID) error
	ListInvoices(ctx context.Context, flt dto.InvoiceFilters, pgn *pagination.Pagination) ([]dto.InvoiceResponse, int, error)
	ListInvoiceDebts(ctx context.Context, id uuid.UUID, flt dto.TransactionFilters, pgn *pagination.Pagination) ([]dto.TransactionResponse, int, error)
}
type UploadService interface {
	ImportFile(userID uuid.UUID, model, action string, invoiceID *uuid.UUID, file multipart.File, fileHeader *multipart.FileHeader) error
}

type AuthService interface {
	GenerateToken(ctx context.Context, userID uuid.UUID, duration time.Duration) (string, error)
	ValidateToken(tokenString string) (*domain.Claims, error)
}

type UserService interface {
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	GetUserByUsername(ctx context.Context, username string) (*domain.User, error)
}
