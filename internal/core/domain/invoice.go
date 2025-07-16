package domain

import (
	"fmt"
	"frog-go/internal/core/errors"
	"time"

	"github.com/google/uuid"
)

type Invoice struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	Amount    float64   `json:"amount"`
	DueDate   time.Time `json:"due_date"`
	Status    TxnStatus `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewInvoice(
	title string,
	dueDate time.Time,
	status *TxnStatus,
) (*Invoice, error) {
	if title == "" {
		return nil, errors.EmptyField("name")
	}

	statusValue := StatusPending
	if status != nil {
		statusValue = *status
	}

	if !statusValue.IsValid() {
		return nil, errors.InvalidParam("status", fmt.Errorf("invalid value"))
	}

	return &Invoice{
		Title:   title,
		DueDate: dueDate,
		Status:  statusValue,
	}, nil
}
