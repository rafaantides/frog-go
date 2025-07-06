package domain

import (
	"frog-go/internal/core/errors"
	"time"

	"github.com/google/uuid"
)

type Debt struct {
	ID           uuid.UUID  `json:"id"`
	Title        string     `json:"title"`
	Amount       float64    `json:"amount"`
	PurchaseDate time.Time  `json:"purchase_date"`
	DueDate      *time.Time `json:"due_date"`
	CategoryID   *uuid.UUID `json:"category_id"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

func NewDebt(
	title string,
	amount float64,
	purchaseDate time.Time,
	dueDate *time.Time,
	categoryID *uuid.UUID,
) (*Debt, error) {
	if title == "" {
		return nil, errors.EmptyField("name")
	}

	if amount == 0 {
		return nil, errors.EmptyField("amount")
	}

	return &Debt{
		Title:        title,
		Amount:       amount,
		PurchaseDate: purchaseDate,
		DueDate:      dueDate,
		CategoryID:   categoryID,
	}, nil
}
