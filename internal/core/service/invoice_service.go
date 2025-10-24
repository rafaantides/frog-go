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

type invoiceService struct {
	repo repository.Repository
}

func NewInvoiceService(repo repository.Repository) inbound.InvoiceService {
	return &invoiceService{repo: repo}
}

func (s *invoiceService) GetInvoiceByID(ctx context.Context, userID uuid.UUID, id uuid.UUID) (*dto.InvoiceResponse, error) {
	return s.repo.GetInvoiceByID(ctx, userID, id)
}

func (s *invoiceService) CreateInvoice(ctx context.Context, userID uuid.UUID, input domain.Invoice) (*dto.InvoiceResponse, error) {
	return s.repo.CreateInvoice(ctx, userID, input)
}

func (s *invoiceService) UpdateInvoice(ctx context.Context, userID uuid.UUID, id uuid.UUID, input domain.Invoice) (*dto.InvoiceResponse, error) {
	return s.repo.UpdateInvoice(ctx, userID, id, input)
}

func (s *invoiceService) DeleteInvoiceByID(ctx context.Context, userID uuid.UUID, id uuid.UUID) error {
	return s.repo.DeleteInvoiceByID(ctx, userID, id)
}

func (s *invoiceService) ListInvoices(ctx context.Context, userID uuid.UUID, flt dto.InvoiceFilters, pgn *pagination.Pagination) ([]dto.InvoiceResponse, int, error) {
	data, err := s.repo.ListInvoices(ctx, userID, flt, pgn)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.repo.CountInvoices(ctx, userID, flt, pgn)
	if err != nil {
		return nil, 0, err
	}

	return data, total, nil
}

func (s *invoiceService) ListInvoiceDebts(ctx context.Context, userID uuid.UUID, id uuid.UUID, flt dto.TransactionFilters, pgn *pagination.Pagination) ([]dto.TransactionResponse, int, error) {
	flt.InvoiceIDs = &[]string{id.String()}

	data, err := s.repo.ListTransactions(ctx, userID, flt, pgn)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.repo.CountTransactions(ctx, userID, flt, pgn)
	if err != nil {
		return nil, 0, err
	}

	return data, total, nil
}
