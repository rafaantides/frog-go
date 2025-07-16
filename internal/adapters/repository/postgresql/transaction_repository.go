package postgresql

import (
	"fmt"
	"frog-go/internal/config"
	"frog-go/internal/core/domain"
	"frog-go/internal/core/dto"
	"frog-go/internal/core/errors"
	"frog-go/internal/ent"
	entCategory "frog-go/internal/ent/category"
	entInvoice "frog-go/internal/ent/invoice"
	"frog-go/internal/ent/transaction"
	"frog-go/internal/utils"
	"sort"
	"time"

	"context"
	"frog-go/internal/utils/pagination"

	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
)

const transactionEntity = "transactions"

func (p *PostgreSQL) GetTransactionByID(ctx context.Context, id uuid.UUID) (*dto.TransactionResponse, error) {
	row, err := p.Client.Transaction.Query().
		Where(transaction.IDEQ(id)).
		WithCategory().
		Only(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.ErrNotFound
		}
		return nil, err
	}
	return newTransactionResponse(row)
}

func (p *PostgreSQL) DeleteTransactionByID(ctx context.Context, id uuid.UUID) error {
	err := p.Client.Transaction.DeleteOneID(id).Exec(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return errors.ErrNotFound
		}
		return err
	}
	return nil
}

func (p *PostgreSQL) CreateTransaction(ctx context.Context, input domain.Transaction) (*dto.TransactionResponse, error) {
	created, err := p.Client.Transaction.
		Create().
		SetTitle(input.Title).
		SetAmount(input.Amount).
		SetRecordType(string(input.RecordType)).
		SetStatus(string(input.Status)).
		SetRecordDate(input.RecordDate).
		SetNillableCategoryID(input.CategoryID).
		SetNillableInvoiceID(input.InvoiceID).
		Save(ctx)

	if err != nil {
		return nil, errors.FailedToSave(transactionEntity, err)
	}

	row, err := p.Client.Transaction.Query().
		Where(transaction.ID(created.ID)).
		WithCategory().
		WithInvoice().
		Only(ctx)

	if err != nil {
		return nil, errors.FailedToFind(transactionEntity, err)
	}

	return newTransactionResponse(row)
}

func (p *PostgreSQL) UpdateTransaction(ctx context.Context, id uuid.UUID, input domain.Transaction) (*dto.TransactionResponse, error) {
	updated, err := p.Client.Transaction.
		UpdateOneID(id).
		SetTitle(input.Title).
		SetAmount(input.Amount).
		SetRecordType(string(input.RecordType)).
		SetStatus(string(input.Status)).
		SetRecordDate(input.RecordDate).
		SetNillableCategoryID(input.CategoryID).
		SetNillableInvoiceID(input.InvoiceID).
		Save(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.ErrNotFound
		}
		return nil, errors.FailedToSave(transactionEntity, err)
	}

	row, err := p.Client.Transaction.Query().
		Where(transaction.ID(updated.ID)).
		WithCategory().
		WithInvoice().
		Only(ctx)

	if err != nil {
		return nil, errors.FailedToFind(transactionEntity, err)
	}

	return newTransactionResponse(row)
}

func (p *PostgreSQL) ListTransactions(ctx context.Context, flt dto.TransactionFilters, pgn *pagination.Pagination) ([]dto.TransactionResponse, error) {
	query := p.Client.Transaction.Query().
		WithCategory().
		WithInvoice()

	query = applyTransactionFilters(query, flt, pgn)
	query = applyTransactionOrderBy(query, pgn)
	query = query.Limit(pgn.PageSize).Offset(pgn.Offset())

	data, err := query.All(ctx)
	if err != nil {
		return nil, err
	}

	return newTransactionResponseList(data)
}

func (p *PostgreSQL) CountTransactions(ctx context.Context, flt dto.TransactionFilters, pgn *pagination.Pagination) (int, error) {
	query := p.Client.Transaction.Query()
	query = applyTransactionFilters(query, flt, pgn)

	total, err := query.Count(ctx)
	if err != nil {
		return 0, err
	}
	return total, nil
}

func (p *PostgreSQL) TransactionsGeneralStats(ctx context.Context, flt dto.ChartFilters) (*dto.TransactionStatsSummary, error) {
	startDate, err := utils.ToDateTime(flt.StartDate)
	if err != nil {
		return nil, err
	}
	endDate, err := utils.ToDateTime(flt.EndDate)
	if err != nil {
		return nil, err
	}

	query := `
		SELECT
			COALESCE(SUM(t.amount), 0) AS total_amount,
			COUNT(*) AS total_transactions,
			COUNT(DISTINCT t.title) AS unique_establishments
		FROM transactions t
		WHERE t.record_date BETWEEN $1 AND $2 AND t.record_type = 'expense'
	`

	var totalAmount float64
	var totalTransactions int
	var uniqueEstablishments int

	err = p.db.QueryRowContext(ctx, query, startDate, endDate).Scan(&totalAmount, &totalTransactions, &uniqueEstablishments)
	if err != nil {
		return nil, err
	}

	var averagePerTransaction float64
	if totalTransactions > 0 {
		averagePerTransaction = totalAmount / float64(totalTransactions)
	}

	return &dto.TransactionStatsSummary{
		TotalAmount:           totalAmount,
		TotalTransactions:     totalTransactions,
		UniqueEstablishments:  uniqueEstablishments,
		AveragePerTransaction: averagePerTransaction,
	}, nil
}

func (p *PostgreSQL) TransactionsSummary(ctx context.Context, flt dto.ChartFilters) ([]dto.SummaryByDate, error) {
	var periodTrunc string

	switch flt.Period {
	case "daily":
		periodTrunc = "day"
	case "weekly":
		periodTrunc = "week"
	case "monthly":
		periodTrunc = "month"
	case "year", "yearly":
		periodTrunc = "year"
	default:
		return nil, fmt.Errorf("invalid period: %s", flt.Period)
	}

	startDate, err := utils.ToDateTime(flt.StartDate)
	if err != nil {
		return nil, err
	}
	endDate, err := utils.ToDateTime(flt.EndDate)
	if err != nil {
		return nil, err
	}

	// Buscar todas as categorias
	allCategories := []string{"Sem categoria"}
	catQuery := "SELECT name FROM categories"
	catRows, err := p.db.QueryContext(ctx, catQuery)
	if err != nil {
		return nil, err
	}
	defer catRows.Close()

	for catRows.Next() {
		var name string
		if err := catRows.Scan(&name); err != nil {
			return nil, err
		}
		allCategories = append(allCategories, name)
	}

	// Buscar os dados de débitos agrupados por data e categoria
	query := `
		SELECT 
			DATE_TRUNC($1, t.record_date) AS period,
			COALESCE(c.name, 'Sem categoria') AS category,
			SUM(t.amount) AS total,
			COUNT(*) AS transactions
		FROM transactions t
		LEFT JOIN categories c ON t.category_id = c.id
		WHERE t.record_date BETWEEN $2 AND $3
		AND t.record_type = 'expense'
		GROUP BY period, category
		ORDER BY period
	`

	rows, err := p.db.QueryContext(ctx, query, periodTrunc, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	type rawData struct {
		date         string
		category     string
		total        float64
		transactions int
	}

	dataByDate := map[string][]rawData{}

	for rows.Next() {
		var date time.Time
		var category string
		var total float64
		var transactions int

		if err := rows.Scan(&date, &category, &total, &transactions); err != nil {
			return nil, err
		}

		// TODO: usar utils para formatar a data
		key := date.Format("2006-01-02")
		dataByDate[key] = append(dataByDate[key], rawData{
			date:         key,
			category:     category,
			total:        total,
			transactions: transactions,
		})
	}

	// Montar resposta final
	var result []dto.SummaryByDate
	for date, entries := range dataByDate {
		summary := dto.SummaryByDate{
			Date:       date,
			Total:      0,
			Categories: []dto.CategorySummary{},
		}

		catMap := map[string]rawData{}
		for _, entry := range entries {
			catMap[entry.category] = entry
			summary.Total += entry.total
		}

		for _, category := range allCategories {
			data, exists := catMap[category]
			total := 0.0
			transactions := 0
			if exists {
				total = data.total
				transactions = data.transactions
			}
			summary.Categories = append(summary.Categories, dto.CategorySummary{
				Category:     category,
				Total:        total,
				Transactions: transactions,
			})
		}

		result = append(result, summary)
	}

	// Ordenar por data
	sort.Slice(result, func(i, j int) bool {
		return result[i].Date < result[j].Date
	})

	return result, nil
}

func mapTransactionToResponse(row *ent.Transaction) dto.TransactionResponse {
	response := dto.TransactionResponse{
		ID:         row.ID,
		Title:      row.Title,
		Amount:     row.Amount,
		Status:     row.Status,
		RecordType: row.RecordType,
		RecordDate: utils.ToDateTimeString(row.RecordDate),
		CreatedAt:  utils.ToDateTimeString(row.CreatedAt),
		UpdatedAt:  utils.ToDateTimeString(row.UpdatedAt),
	}

	if row.Edges.Invoice != nil {
		response.Invoice = &dto.TransactionInvoiceResponse{
			ID:    row.Edges.Invoice.ID,
			Title: row.Edges.Invoice.Title,
		}
	}

	if row.Edges.Category != nil {
		response.Category = &dto.TransactionCategoryResponse{
			ID:   row.Edges.Category.ID,
			Name: row.Edges.Category.Name,
		}
	}

	return response
}

func newTransactionResponse(row *ent.Transaction) (*dto.TransactionResponse, error) {
	if row == nil {
		return nil, nil
	}
	response := mapTransactionToResponse(row)
	return &response, nil
}

func newTransactionResponseList(rows []*ent.Transaction) ([]dto.TransactionResponse, error) {
	if rows == nil {
		return nil, nil
	}
	response := make([]dto.TransactionResponse, 0, len(rows))
	for _, row := range rows {
		response = append(response, mapTransactionToResponse(row))
	}
	return response, nil
}

func applyTransactionOrderBy(query *ent.TransactionQuery, pgn *pagination.Pagination) *ent.TransactionQuery {

	var orderDirection sql.OrderTermOption
	if pgn.OrderDirection == config.OrderAsc {
		orderDirection = sql.OrderAsc()
	} else {
		orderDirection = sql.OrderDesc()
	}

	switch pgn.OrderBy {
	case "invoice":
		query.Order(
			transaction.ByInvoiceField(entInvoice.FieldTitle, orderDirection),
			transaction.ByID(sql.OrderAsc()),
		)
	case "category":
		query.Order(
			transaction.ByCategoryField(entCategory.FieldName, orderDirection),
			transaction.ByID(sql.OrderAsc()),
		)
	default:
		if pgn.OrderDirection == config.OrderAsc {
			query = query.Order(
				ent.Asc(pgn.OrderBy),
				ent.Asc(transaction.FieldID),
			)
		} else {
			query = query.Order(
				ent.Desc(pgn.OrderBy),
				ent.Asc(transaction.FieldID),
			)
		}
	}

	return query
}

func applyTransactionFilters(query *ent.TransactionQuery, flt dto.TransactionFilters, pgn *pagination.Pagination) *ent.TransactionQuery {
	if pgn.Search != "" {
		query = query.Where(
			transaction.Or(
				transaction.TitleContainsFold(pgn.Search),
				transaction.StatusContainsFold(pgn.Search),
				transaction.HasCategoryWith(
					entCategory.NameContainsFold(pgn.Search),
				),
			),
		)
	}

	if flt.Statuses != nil && len(*flt.Statuses) > 0 {
		query = query.Where(
			transaction.StatusIn(*flt.Statuses...),
		)
	}

	if flt.RecordTypes != nil && len(*flt.RecordTypes) > 0 {
		query = query.Where(
			transaction.RecordTypeIn(*flt.RecordTypes...),
		)
	}

	if flt.InvoiceIDs != nil {
		invoiceIDs := utils.ToUUIDSlice(*flt.InvoiceIDs)
		if len(invoiceIDs) > 0 {
			query = query.Where(
				transaction.HasInvoiceWith(entInvoice.IDIn(invoiceIDs...)),
			)
		}
	}

	if flt.CategoryIDs != nil {
		categoryIDs := utils.ToUUIDSlice(*flt.CategoryIDs)
		if len(categoryIDs) > 0 {
			query = query.Where(
				transaction.HasCategoryWith(entCategory.IDIn(categoryIDs...)),
			)
		}
	}

	if flt.MinAmount != nil {
		query = query.Where(
			transaction.AmountGTE(*flt.MinAmount),
		)
	}
	if flt.MaxAmount != nil {
		query = query.Where(
			transaction.AmountLTE(*flt.MaxAmount),
		)
	}
	if t := utils.ToDateTimeUnsafe(flt.StartDate); t != nil {
		query = query.Where(transaction.RecordDateGTE(*t))
	}

	if t := utils.ToDateTimeUnsafe(flt.EndDate); t != nil {
		query = query.Where(transaction.RecordDateLTE(*t))
	}

	return query
}
