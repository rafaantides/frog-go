package domain

import (
	"fmt"
	"frog-go/internal/core/errors"
	"slices"
	"time"

	"github.com/google/uuid"
)

type TxnStatus string

const (
	TxnStatusPending  TxnStatus = "pending"
	TxnStatusPaid     TxnStatus = "paid"
	TxnStatusCanceled TxnStatus = "canceled"
)

type TxnKind string

const (
	TxnKindIncome  TxnKind = "income"
	TxnKindExpense TxnKind = "expense"
)

func ValidTxnStatus() []string {
	return []string{
		string(TxnStatusPending),
		string(TxnStatusPaid),
		string(TxnStatusCanceled),
	}
}

func (a TxnStatus) IsValid() bool {
	return slices.Contains(ValidTxnStatus(), string(a))
}

func ValidTxnKind() []string {
	return []string{
		string(TxnKindIncome),
		string(TxnKindExpense),
	}
}

func (a TxnKind) IsValid() bool {
	return slices.Contains(ValidTxnKind(), string(a))
}

type Transaction struct {
	ID           uuid.UUID  `json:"id"`
	Title        string     `json:"title"`
	Amount       float64    `json:"amount"`
	PurchaseDate time.Time  `json:"purchase_date"`
	DueDate      *time.Time `json:"due_date"`
	CategoryID   *uuid.UUID `json:"category_id"`
	Status       TxnStatus  `json:"status"`
	Kind         TxnKind    `json:"kind"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

func NewTransaction(
	title string,
	amount float64,
	purchaseDate time.Time,
	dueDate *time.Time,
	categoryID *uuid.UUID,
	status *TxnStatus,
	kind *TxnKind,
) (*Transaction, error) {
	if title == "" {
		return nil, errors.EmptyField("name")
	}

	if amount == 0 {
		return nil, errors.EmptyField("amount")
	}

	statusValue := TxnStatusPending
	if status != nil {
		statusValue = *status
	}

	kindValue := TxnKindExpense
	if kind != nil {
		kindValue = *kind
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
		Kind:         kindValue,
		CategoryID:   categoryID,
	}, nil
}
