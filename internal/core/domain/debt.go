package domain

import (
	"fmt"
	"frog-go/internal/core/errors"
	"slices"
	"time"

	"github.com/google/uuid"
)

type DebtStatus string

const (
	DebtStatusPending  DebtStatus = "pending"
	DebtStatusPaid     DebtStatus = "paid"
	DebtStatusCanceled DebtStatus = "canceled"
)

func ValidDebtStatus() []string {
	return []string{
		string(DebtStatusPending),
		string(DebtStatusPaid),
		string(DebtStatusCanceled),
	}
}

func (a DebtStatus) IsValid() bool {
	return slices.Contains(ValidDebtStatus(), string(a))
}

type Debt struct {
	ID           uuid.UUID  `json:"id"`
	Title        string     `json:"title"`
	Amount       float64    `json:"amount"`
	PurchaseDate time.Time  `json:"purchase_date"`
	DueDate      *time.Time `json:"due_date"`
	CategoryID   *uuid.UUID `json:"category_id"`
	Status       DebtStatus `json:"status"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

func NewDebt(
	title string,
	amount float64,
	purchaseDate time.Time,
	dueDate *time.Time,
	categoryID *uuid.UUID,
	status *DebtStatus,
) (*Debt, error) {
	if title == "" {
		return nil, errors.EmptyField("name")
	}

	if amount == 0 {
		return nil, errors.EmptyField("amount")
	}

	statusValue := DebtStatusPending
	if status != nil {
		statusValue = *status
	}

	if !statusValue.IsValid() {
		return nil, errors.InvalidParam("status", fmt.Errorf("invalid value"))
	}

	return &Debt{
		Title:        title,
		Amount:       amount,
		PurchaseDate: purchaseDate,
		DueDate:      dueDate,
		Status:       statusValue,
		CategoryID:   categoryID,
	}, nil
}
