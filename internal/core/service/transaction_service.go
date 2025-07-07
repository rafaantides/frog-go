package service

import (
	"context"

	"frog-go/internal/core/domain"
	"frog-go/internal/core/dto"
	"frog-go/internal/core/ports/inbound"
	"frog-go/internal/core/ports/outbound/repository"
	"frog-go/internal/utils/pagination"

	"github.com/google/uuid"
)

type transactionService struct {
	repo repository.Repository
}

func NewTransactionService(repo repository.Repository) inbound.TransactionService {
	return &transactionService{repo: repo}
}
func (s *transactionService) GetTransactionByID(ctx context.Context, id uuid.UUID) (*dto.TransactionResponse, error) {
	return s.repo.GetTransactionByID(ctx, id)
}

func (s *transactionService) CreateTransaction(ctx context.Context, input domain.Transaction) (*dto.TransactionResponse, error) {
	return s.repo.CreateTransaction(ctx, input)
}

func (s *transactionService) UpdateTransaction(ctx context.Context, id uuid.UUID, input domain.Transaction) (*dto.TransactionResponse, error) {
	return s.repo.UpdateTransaction(ctx, id, input)
}

func (s *transactionService) DeleteTransactionByID(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteTransactionByID(ctx, id)
}

func (s *transactionService) ListTransactions(ctx context.Context, flt dto.TransactionFilters, pgn *pagination.Pagination) ([]dto.TransactionResponse, int, error) {
	data, err := s.repo.ListTransactions(ctx, flt, pgn)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.repo.CountTransactions(ctx, flt, pgn)
	if err != nil {
		return nil, 0, err
	}

	return data, total, nil
}

func (s *transactionService) TransactionsSummary(ctx context.Context, flt dto.ChartFilters) ([]dto.SummaryByDate, error) {
	return s.repo.TransactionsSummary(ctx, flt)
}

func (s *transactionService) TransactionsGeneralStats(ctx context.Context, flt dto.ChartFilters) (*dto.TransactionStatsSummary, error) {
	return s.repo.TransactionsGeneralStats(ctx, flt)
}
