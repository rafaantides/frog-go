package dto

import (
	"frog-go/internal/core/domain"
	appError "frog-go/internal/core/errors"
	"frog-go/internal/utils"

	"github.com/google/uuid"
)

type InvoiceRequest struct {
	Title   string `json:"title"`
	DueDate string `json:"due_date"`
	Status  string `json:"status"`
}

type InvoiceFilters struct {
	MinAmount *float64  `form:"min_amount"`
	MaxAmount *float64  `form:"max_amount"`
	StartDate *string   `form:"start_date"`
	EndDate   *string   `form:"end_date"`
	Statuses  *[]string `form:"statuses"`
}

type InvoiceResponse struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	Amount    float64   `json:"amount"`
	DueDate   string    `json:"due_date"`
	Status    string    `json:"status"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
}

func (r *InvoiceRequest) ToDomain(userID uuid.UUID) (*domain.Invoice, error) {

	dueDate, err := utils.ToDateTime(r.DueDate)
	if err != nil {
		return nil, appError.InvalidParam("due_date", err)
	}

	status := domain.TxnStatus(r.Status)

	return domain.NewInvoice(
		userID,
		r.Title,
		dueDate,
		&status,
	)
}
