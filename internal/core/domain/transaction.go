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
	StatusPending  TxnStatus = "pending"
	StatusPaid     TxnStatus = "paid"
	StatusCanceled TxnStatus = "canceled"
)

type RecordType string

const (
	TypeIncome  RecordType = "income"
	TypeExpense RecordType = "expense"
)

func ValidTxnStatus() []string {
	return []string{
		string(StatusPending),
		string(StatusPaid),
		string(StatusCanceled),
	}
}

func (a TxnStatus) IsValid() bool {
	return slices.Contains(ValidTxnStatus(), string(a))
}

func ValidRecordType() []string {
	return []string{
		string(TypeIncome),
		string(TypeExpense),
	}
}

func (a RecordType) IsValid() bool {
	return slices.Contains(ValidRecordType(), string(a))
}

type Transaction struct {
	ID         uuid.UUID  `json:"id"`
	Title      string     `json:"title"`
	Amount     float64    `json:"amount"`
	RecordDate time.Time  `json:"record_date"`
	CategoryID *uuid.UUID `json:"category_id"`
	InvoiceID  *uuid.UUID `json:"invoice_id"`
	Status     TxnStatus  `json:"status"`
	RecordType RecordType `json:"record_type"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

func NewTransaction(
	title string,
	amount float64,
	RecordDate time.Time,
	invoiceID *uuid.UUID,
	categoryID *uuid.UUID,
	status *TxnStatus,
	recordType *RecordType,
) (*Transaction, error) {
	if title == "" {
		return nil, errors.EmptyField("name")
	}

	if amount == 0 {
		return nil, errors.EmptyField("amount")
	}

	statusValue := StatusPending
	if status != nil && *status != "" {
		statusValue = *status
	}

	recordTypeValue := TypeExpense
	if recordType != nil && *recordType != "" {
		recordTypeValue = *recordType
	}

	if !statusValue.IsValid() {
		return nil, errors.InvalidParam("status", fmt.Errorf("invalid value"))
	}

	return &Transaction{
		Title:      title,
		Amount:     amount,
		RecordDate: RecordDate,
		Status:     statusValue,
		RecordType: recordTypeValue,
		InvoiceID:  invoiceID,
		CategoryID: categoryID,
	}, nil
}
