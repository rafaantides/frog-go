package dto

import (
	"frog-go/internal/core/domain"
	"frog-go/internal/core/errors"
	"frog-go/internal/utils"

	"github.com/google/uuid"
)

type TransactionRequest struct {
	Title      string  `json:"title"`
	Amount     float64 `json:"amount"`
	RecordDate string  `json:"record_date"`
	CategoryID *string `json:"category_id"`
	InvoiceID  *string `json:"invoice_id"`
	Status     string  `json:"status" validate:"required,oneof=pending paid canceled"`
	RecordType string  `json:"record_type" validate:"required,oneof=income expense"`
}

// TODO: fazer um bind que funcione com uuid.UUID o ShouldBindQuery n esta reconhecendo o *[]uuid.UUID
type TransactionFilters struct {
	InvoiceIDs  *[]string `json:"invoice_ids"`
	CategoryIDs *[]string `json:"category_ids"`
	Statuses    *[]string `form:"statuses"`
	RecordTypes *[]string `form:"record_types"`
	MinAmount   *float64  `form:"min_amount"`
	MaxAmount   *float64  `form:"max_amount"`
	StartDate   *string   `form:"start_date"`
	EndDate     *string   `form:"end_date"`
}

type TransactionResponse struct {
	ID         uuid.UUID                    `json:"id"`
	Title      string                       `json:"title"`
	Amount     float64                      `json:"amount"`
	RecordDate string                       `json:"record_date"`
	Category   *TransactionCategoryResponse `json:"category"`
	Invoice    *TransactionInvoiceResponse  `json:"invoice"`
	RecordType string                       `json:"record_type"`
	Status     string                       `json:"status"`
	CreatedAt  string                       `json:"created_at"`
	UpdatedAt  string                       `json:"updated_at"`
}

type TransactionInvoiceResponse struct {
	ID    uuid.UUID `json:"id"`
	Title string    `json:"title"`
}

type TransactionCategoryResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

func (r *TransactionRequest) ToDomain() (*domain.Transaction, error) {
	RecordDate, err := utils.ToDateTime(r.RecordDate)
	if err != nil {
		return nil, errors.InvalidParam("record_date", err)
	}

	var invoiceID *uuid.UUID
	if r.InvoiceID != nil {
		invoiceID, err = utils.ToNillableUUID(*r.InvoiceID)
		if err != nil {
			return nil, errors.InvalidParam("invoice_id", err)
		}
	}

	var categoryID *uuid.UUID
	if r.CategoryID != nil {
		categoryID, err = utils.ToNillableUUID(*r.CategoryID)
		if err != nil {
			return nil, errors.InvalidParam("category_id", err)
		}
	}

	status := domain.TxnStatus(r.Status)
	recordType := domain.RecordType(r.RecordType)

	return domain.NewTransaction(
		r.Title,
		r.Amount,
		RecordDate,
		invoiceID,
		categoryID,
		&status,
		&recordType,
	)
}
