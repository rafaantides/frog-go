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

type debtService struct {
	repo repository.Repository
}

func NewDebtService(repo repository.Repository) inbound.DebtService {
	return &debtService{repo: repo}
}
func (s *debtService) GetDebtByID(ctx context.Context, id uuid.UUID) (*dto.DebtResponse, error) {
	return s.repo.GetDebtByID(ctx, id)
}

func (s *debtService) CreateDebt(ctx context.Context, input domain.Debt) (*dto.DebtResponse, error) {
	return s.repo.CreateDebt(ctx, input)
}

func (s *debtService) UpdateDebt(ctx context.Context, id uuid.UUID, input domain.Debt) (*dto.DebtResponse, error) {
	return s.repo.UpdateDebt(ctx, id, input)
}

func (s *debtService) DeleteDebtByID(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteDebtByID(ctx, id)
}

func (s *debtService) ListDebts(ctx context.Context, flt dto.DebtFilters, pgn *pagination.Pagination) ([]dto.DebtResponse, int, error) {
	data, err := s.repo.ListDebts(ctx, flt, pgn)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.repo.CountDebts(ctx, flt, pgn)
	if err != nil {
		return nil, 0, err
	}

	return data, total, nil
}

func (s *debtService) DebtsSummary(ctx context.Context, flt dto.ChartFilters) ([]dto.SummaryByDate, error) {
	return s.repo.DebtsSummary(ctx, flt)
}

func (s *debtService) DebtsGeneralStats(ctx context.Context, flt dto.ChartFilters) (*dto.DebtStatsSummary, error) {
	return s.repo.DebtsGeneralStats(ctx, flt)
}
