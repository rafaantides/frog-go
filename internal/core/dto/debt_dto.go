package dto

import (
	"frog-go/internal/core/domain"
	"frog-go/internal/core/errors"
	"frog-go/internal/utils"
	"time"

	"github.com/google/uuid"
)

type DebtRequest struct {
	Title        string  `json:"title"`
	Amount       float64 `json:"amount"`
	PurchaseDate string  `json:"purchase_date"`
	DueDate      *string `json:"due_date"`
	CategoryID   *string `json:"category_id"`
	Status       string  `json:"status" validate:"required,oneof=pending paid canceled"`
}

// TODO: fazer um bind que funcione com uuid.UUID o ShouldBindQuery n esta reconhecendo o *[]uuid.UUID
type DebtFilters struct {
	CategoryID *[]string `json:"category_id"`
	Status     *[]string `form:"status"`
	MinAmount  *float64  `form:"min_amount"`
	MaxAmount  *float64  `form:"max_amount"`
	StartDate  *string   `form:"start_date"`
	EndDate    *string   `form:"end_date"`
}

type DebtResponse struct {
	ID           uuid.UUID             `json:"id"`
	Title        string                `json:"title"`
	Amount       float64               `json:"amount"`
	PurchaseDate string                `json:"purchase_date"`
	DueDate      *string               `json:"due_date"`
	Category     *DebtCategoryResponse `json:"category"`
	Status       string                `json:"status"`
	CreatedAt    string                `json:"created_at"`
	UpdatedAt    string                `json:"updated_at"`
}

type DebtCategoryResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

func (r *DebtRequest) ToDomain() (*domain.Debt, error) {
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

	status := domain.DebtStatus(r.Status)

	return domain.NewDebt(
		r.Title,
		r.Amount,
		purchaseDate,
		dueDate,
		categoryID,
		&status,
	)
}
