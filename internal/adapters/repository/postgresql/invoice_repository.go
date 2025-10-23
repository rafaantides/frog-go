package postgresql

import (
	"context"
	"frog-go/internal/config"
	"frog-go/internal/core/domain"
	"frog-go/internal/core/dto"
	appError "frog-go/internal/core/errors"
	"frog-go/internal/ent"
	"frog-go/internal/ent/invoice"
	"frog-go/internal/ent/transaction"
	"frog-go/internal/ent/user"
	"frog-go/internal/utils"
	"frog-go/internal/utils/pagination"
	"frog-go/internal/utils/utilsctx"

	"github.com/google/uuid"
)

func (d *PostgreSQL) GetInvoiceByID(ctx context.Context, id uuid.UUID) (*dto.InvoiceResponse, error) {
	userID, err := utilsctx.GetUserID(ctx)

	if err != nil {
		return nil, err
	}

	row, err := d.Client.Invoice.Query().
		Where(invoice.IDEQ(id)).
		Where(invoice.HasUserWith(user.IDEQ(userID))).
		Only(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, appError.ErrNotFound
		}
		return nil, err
	}
	return newInvoiceResponse(row)
}

func (d *PostgreSQL) DeleteInvoiceByID(ctx context.Context, id uuid.UUID) error {
	userID, err := utilsctx.GetUserID(ctx)
	if err != nil {
		return err
	}

	err = d.Client.Invoice.DeleteOneID(id).
		Where(invoice.HasUserWith(user.IDEQ(userID))).
		Exec(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return appError.ErrNotFound
		}
		return err
	}
	return nil
}

func (d *PostgreSQL) CreateInvoice(ctx context.Context, input domain.Invoice) (*dto.InvoiceResponse, error) {
	userID, err := utilsctx.GetUserID(ctx)
	if err != nil {
		return nil, err
	}

	created, err := d.Client.Invoice.
		Create().
		SetUserID(userID).
		SetTitle(input.Title).
		SetDueDate(input.DueDate).
		SetStatus(string(input.Status)).
		Save(ctx)

	if err != nil {
		return nil, appError.FailedToSave("invoices", err)
	}

	row, err := d.Client.Invoice.
		Query().
		Where(invoice.ID(created.ID)).
		Only(ctx)

	if err != nil {
		return nil, appError.FailedToFind("invoice", err)
	}

	return newInvoiceResponse(row)
}

func (d *PostgreSQL) UpdateInvoice(ctx context.Context, id uuid.UUID, input domain.Invoice) (*dto.InvoiceResponse, error) {
	userID, err := utilsctx.GetUserID(ctx)
	if err != nil {
		return nil, err
	}

	updated, err := d.Client.Invoice.
		UpdateOneID(id).
		Where(invoice.HasUserWith(user.IDEQ(userID))).
		SetTitle(input.Title).
		SetDueDate(input.DueDate).
		SetStatus(string(input.Status)).
		Save(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, appError.ErrNotFound
		}
		return nil, appError.FailedToSave("invoices", err)
	}

	_, err = d.Client.Transaction.
		Update().
		Where(transaction.HasInvoiceWith(invoice.ID(id))).
		SetStatus(string(input.Status)).
		Save(ctx)
	if err != nil {
		return nil, appError.FailedToSave("transactions", err)
	}

	row, err := d.Client.Invoice.
		Query().
		Where(invoice.ID(updated.ID)).
		Only(ctx)

	if err != nil {
		return nil, appError.FailedToFind("invoice", err)
	}

	return newInvoiceResponse(row)
}

func (d *PostgreSQL) ListInvoices(ctx context.Context, flt dto.InvoiceFilters, pgn *pagination.Pagination) ([]dto.InvoiceResponse, error) {
	userID, err := utilsctx.GetUserID(ctx)

	if err != nil {
		return nil, err
	}

	query := d.Client.Invoice.Query().
		Where(invoice.HasUserWith(user.IDEQ(userID)))

	query = applyInvoiceFilters(query, flt, pgn)
	query = apllyInvoiceOrderBy(query, pgn)
	query = query.Limit(pgn.PageSize).Offset(pgn.Offset())

	data, err := query.All(ctx)
	if err != nil {
		return nil, err
	}

	return newInvoiceResponseList(data)
}

func (d *PostgreSQL) CountInvoices(ctx context.Context, flt dto.InvoiceFilters, pgn *pagination.Pagination) (int, error) {
	userID, err := utilsctx.GetUserID(ctx)
	if err != nil {
		return 0, err
	}

	query := d.Client.Invoice.Query().
		Where(invoice.HasUserWith(user.IDEQ(userID)))

	query = applyInvoiceFilters(query, flt, pgn)

	total, err := query.Count(ctx)
	if err != nil {
		return 0, err
	}
	return total, nil
}

func mapInvoiceToResponse(row *ent.Invoice) dto.InvoiceResponse {
	response := dto.InvoiceResponse{
		ID:        row.ID,
		Title:     row.Title,
		Amount:    row.Amount,
		Status:    row.Status,
		DueDate:   utils.ToDateTimeString(row.DueDate),
		CreatedAt: utils.ToDateTimeString(row.CreatedAt),
		UpdatedAt: utils.ToDateTimeString(row.UpdatedAt),
	}

	return response
}

func newInvoiceResponse(row *ent.Invoice) (*dto.InvoiceResponse, error) {
	if row == nil {
		return nil, nil
	}
	response := mapInvoiceToResponse(row)
	return &response, nil
}

func newInvoiceResponseList(rows []*ent.Invoice) ([]dto.InvoiceResponse, error) {
	if rows == nil {
		return nil, nil
	}
	response := make([]dto.InvoiceResponse, 0, len(rows))
	for _, row := range rows {
		response = append(response, mapInvoiceToResponse(row))
	}
	return response, nil
}

func apllyInvoiceOrderBy(query *ent.InvoiceQuery, pgn *pagination.Pagination) *ent.InvoiceQuery {

	if pgn.OrderDirection == config.OrderAsc {
		query = query.Order(
			ent.Asc(pgn.OrderBy),
			ent.Asc(invoice.FieldID),
		)
	} else {
		query = query.Order(
			ent.Desc(pgn.OrderBy),
			ent.Asc(invoice.FieldID),
		)
	}

	return query
}

func applyInvoiceFilters(query *ent.InvoiceQuery, flt dto.InvoiceFilters, pgn *pagination.Pagination) *ent.InvoiceQuery {
	if pgn.Search != "" {
		query = query.Where(
			invoice.Or(
				invoice.TitleContainsFold(pgn.Search),
				invoice.StatusContainsFold(pgn.Search),
			),
		)
	}

	if flt.Statuses != nil && len(*flt.Statuses) > 0 {
		query = query.Where(
			invoice.StatusIn(*flt.Statuses...),
		)
	}

	if flt.MinAmount != nil {
		query = query.Where(
			invoice.AmountGTE(*flt.MinAmount),
		)
	}

	if flt.MaxAmount != nil {
		query = query.Where(
			invoice.AmountLTE(*flt.MaxAmount),
		)
	}

	if t := utils.ToDateTimeUnsafe(flt.StartDate); t != nil {
		query = query.Where(invoice.DueDateGTE(*t))
	}

	if t := utils.ToDateTimeUnsafe(flt.EndDate); t != nil {
		query = query.Where(invoice.DueDateLTE(*t))
	}

	return query
}
