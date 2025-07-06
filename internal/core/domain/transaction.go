package domain

import (
	"fmt"
	"frog-go/internal/core/errors"
	"slices"
	"time"

	"github.com/google/uuid"
)

type TxnKind string

const (
	TxnKindPending TxnStatus = "pending"
	TxnKindExpense TxnKind   = "expense"
	TxnKindIncome  TxnKind   = "income"
	TxnKindSavings TxnKind   = "savings"
	TxnKindLoan    TxnKind   = "loan"
)

func ValidTxnKind() []string {
	return []string{
		string(TxnKindPending),
		string(TxnKindExpense),
		string(TxnKindIncome),
		string(TxnKindSavings),
		string(TxnKindLoan),
	}
}

func (a TxnKind) IsValid() bool {
	return slices.Contains(ValidTxnKind(), string(a))
}

type TxnStatus string

const (
	TxnStatusPending  TxnStatus = "pending"
	TxnStatusPaid     TxnStatus = "paid"
	TxnStatusFailed   TxnStatus = "failed"
	TxnStatusCanceled TxnStatus = "canceled"
)

func ValidTxnStatus() []string {
	return []string{
		string(TxnStatusPending),
		string(TxnStatusPaid),
		string(TxnStatusFailed),
		string(TxnStatusCanceled),
	}
}

func (a TxnStatus) IsValid() bool {
	return slices.Contains(ValidTxnStatus(), string(a))
}

type Transaction struct {
	ID           uuid.UUID  `json:"id"`
	Title        string     `json:"title"`
	Amount       float64    `json:"amount"`
	PurchaseDate time.Time  `json:"purchase_date"`
	DueDate      *time.Time `json:"due_date"`
	CategoryID   *uuid.UUID `json:"category_id"`
	Kind         TxnKind    `json:"kind"`
	Status       TxnStatus  `json:"status"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

func NewTransaction(
	title string,
	amount float64,
	accountID uuid.UUID,
	purchaseDate time.Time,
	dueDate *time.Time,
	invoiceID *uuid.UUID,
	categoryID *uuid.UUID,
	kind TxnKind,
	status *TxnStatus,
) (*Transaction, error) {
	if title == "" {
		return nil, errors.EmptyField("name")
	}

	if !kind.IsValid() {
		return nil, errors.InvalidParam("kind", fmt.Errorf("invalid value"))
	}

	statusValue := TxnStatusPending
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
		CategoryID:   categoryID,
		Status:       statusValue,
	}, nil
}
