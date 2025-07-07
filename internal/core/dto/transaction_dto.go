package dto

import (
	"frog-go/internal/core/domain"
	"frog-go/internal/core/errors"
	"frog-go/internal/utils"
	"time"

	"github.com/google/uuid"
)

type TransactionRequest struct {
	Title        string  `json:"title"`
	Amount       float64 `json:"amount"`
	PurchaseDate string  `json:"purchase_date"`
	DueDate      *string `json:"due_date"`
	CategoryID   *string `json:"category_id"`
	Status       string  `json:"status" validate:"required,oneof=pending paid canceled"`
	Kind         string  `json:"kind" validate:"required,oneof=income expense"`
}

// TODO: fazer um bind que funcione com uuid.UUID o ShouldBindQuery n esta reconhecendo o *[]uuid.UUID
type TransactionFilters struct {
	CategoryIDs *[]string `json:"category_ids"`
	Statuses    *[]string `form:"statuses"`
	Kinds       *[]string `form:"kinds"`
	MinAmount   *float64  `form:"min_amount"`
	MaxAmount   *float64  `form:"max_amount"`
	StartDate   *string   `form:"start_date"`
	EndDate     *string   `form:"end_date"`
}

type TransactionResponse struct {
	ID           uuid.UUID                    `json:"id"`
	Title        string                       `json:"title"`
	Amount       float64                      `json:"amount"`
	PurchaseDate string                       `json:"purchase_date"`
	DueDate      *string                      `json:"due_date"`
	Category     *TransactionCategoryResponse `json:"category"`
	Kind         string                       `json:"kind"`
	Status       string                       `json:"status"`
	CreatedAt    string                       `json:"created_at"`
	UpdatedAt    string                       `json:"updated_at"`
}

type TransactionCategoryResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

func (r *TransactionRequest) ToDomain() (*domain.Transaction, error) {
	purchaseDate, err := utils.ToDateTime(r.PurchaseDate)
	if err != nil {
		return nil, errors.InvalidParam("purchase_date", err)
	}

	var dueDate *time.Time
	if r.DueDate != nil {
		dueDate, err = utils.ToNillableDateTime(*r.DueDate)
		if err != nil {
			return nil, errors.InvalidParam("due_date", err)
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
	kind := domain.TxnKind(r.Kind)

	return domain.NewTransaction(
		r.Title,
		r.Amount,
		purchaseDate,
		dueDate,
		categoryID,
		&status,
		&kind,
	)
}
