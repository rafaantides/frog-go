package domain

import (
	"fmt"
	appError "frog-go/internal/core/errors"
	"time"

	"github.com/google/uuid"
)

type Invoice struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Title     string
	Amount    float64
	DueDate   time.Time
	Status    TxnStatus
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewInvoice(
	title string,
	dueDate time.Time,
	status *TxnStatus,
) (*Invoice, error) {
	if title == "" {
		return nil, appError.EmptyField("name")
	}

	statusValue := StatusPending
	if status != nil {
		statusValue = *status
	}

	if !statusValue.IsValid() {
		return nil, appError.InvalidParam("status", fmt.Errorf("invalid value"))
	}

	return &Invoice{
		Title:   title,
		DueDate: dueDate,
		Status:  statusValue,
	}, nil
}
