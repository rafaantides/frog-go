package domain

import (
	"fmt"
	appError "frog-go/internal/core/errors"
	"time"

	"github.com/google/uuid"
)

type Invoice struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Title     string    `json:"title"`
	Amount    float64   `json:"amount"`
	DueDate   time.Time `json:"due_date"`
	Status    TxnStatus `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewInvoice(
	userID uuid.UUID,
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
		UserID:  userID,
		Title:   title,
		DueDate: dueDate,
		Status:  statusValue,
	}, nil
}
