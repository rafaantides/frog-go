package domain

import (
	"fmt"
	"frog-go/internal/core/errors"
	"slices"
	"time"

	"github.com/google/uuid"
)

type TransactionStatus string

const (
	TransactionStatusPending  TransactionStatus = "pending"
	TransactionStatusPaid     TransactionStatus = "paid"
	TransactionStatusCanceled TransactionStatus = "canceled"
)

func ValidTransactionStatus() []string {
	return []string{
		string(TransactionStatusPending),
		string(TransactionStatusPaid),
		string(TransactionStatusCanceled),
	}
}

func (a TransactionStatus) IsValid() bool {
	return slices.Contains(ValidTransactionStatus(), string(a))
}

type Transaction struct {
	ID           uuid.UUID  `json:"id"`
	Title        string     `json:"title"`
	Amount       float64    `json:"amount"`
	PurchaseDate time.Time  `json:"purchase_date"`
	DueDate      *time.Time `json:"due_date"`
	CategoryID   *uuid.UUID `json:"category_id"`
	Status       TransactionStatus `json:"status"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

func NewTransaction(
	title string,
	amount float64,
	purchaseDate time.Time,
	dueDate *time.Time,
	categoryID *uuid.UUID,
	status *TransactionStatus,
) (*Transaction, error) {
	if title == "" {
		return nil, errors.EmptyField("name")
	}

	if amount == 0 {
		return nil, errors.EmptyField("amount")
	}

	statusValue := TransactionStatusPending
	if status != nil {
		statusValue = *status
	}

	if !statusValue.IsValid() {
		return nil, errors.InvalidParam("status", fmt.Errorf("invalid value"))
	}

	return &Transaction{
		Title:        title,
		Amount:       amount,
		PurchaseDate: purchaseDate,
		DueDate:      dueDate,
		Status:       statusValue,
		CategoryID:   categoryID,
	}, nil
}
