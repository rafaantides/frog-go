package postgresql

import (
	"fmt"
	"frog-go/internal/config"
	"frog-go/internal/core/domain"
	"frog-go/internal/core/dto"
	appError "frog-go/internal/core/errors"
	"frog-go/internal/ent"
	entCategory "frog-go/internal/ent/category"
	entInvoice "frog-go/internal/ent/invoice"
	"frog-go/internal/ent/transaction"
	"frog-go/internal/ent/user"
	"frog-go/internal/utils"
	"sort"
	"time"

	"context"
	"frog-go/internal/utils/authctx"
	"frog-go/internal/utils/pagination"

	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
)

const transactionEntity = "transactions"

func (p *PostgreSQL) GetTransactionByID(ctx context.Context, id uuid.UUID) (*dto.TransactionResponse, error) {
	userID, err := authctx.GetUserID(ctx)
	if err != nil {
		return nil, err
	}

	row, err := p.Client.Transaction.Query().
		Where(transaction.IDEQ(id)).
		Where(transaction.HasUserWith(user.IDEQ(userID))).
		WithCategory().
		Only(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, appError.ErrNotFound
		}
		return nil, err
	}
	return newTransactionResponse(row)
}

func (p *PostgreSQL) DeleteTransactionByID(ctx context.Context, id uuid.UUID) error {
	userID, err := authctx.GetUserID(ctx)
	if err != nil {
		return err
	}

	err = p.Client.Transaction.DeleteOneID(id).
		Where(transaction.HasUserWith(user.IDEQ(userID))).
		Exec(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return appError.ErrNotFound
		}
		return err
	}
	return nil
}

func (p *PostgreSQL) CreateTransaction(ctx context.Context, input domain.Transaction) (*dto.TransactionResponse, error) {
	created, err := p.Client.Transaction.
		Create().
		SetUserID(input.UserID).
		SetTitle(input.Title).
		SetAmount(input.Amount).
		SetRecordType(string(input.RecordType)).
		SetStatus(string(input.Status)).
		SetRecordDate(input.RecordDate).
		SetNillableCategoryID(input.CategoryID).
		SetNillableInvoiceID(input.InvoiceID).
		Save(ctx)

	if err != nil {
		return nil, appError.FailedToSave(transactionEntity, err)
	}

	row, err := p.Client.Transaction.Query().
		Where(transaction.ID(created.ID)).
		WithCategory().
		WithInvoice().
		Only(ctx)

	if err != nil {
		return nil, appError.FailedToFind(transactionEntity, err)
	}

	return newTransactionResponse(row)
}

func (p *PostgreSQL) UpdateTransaction(ctx context.Context, id uuid.UUID, input domain.Transaction) (*dto.TransactionResponse, error) {
	userID, err := authctx.GetUserID(ctx)
	if err != nil {
		return nil, err
	}

	updated, err := p.Client.Transaction.
		UpdateOneID(id).
		Where(transaction.HasUserWith(user.IDEQ(userID))).
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
			return nil, appError.ErrNotFound
		}
		return nil, appError.FailedToSave(transactionEntity, err)
	}

	row, err := p.Client.Transaction.Query().
		Where(transaction.ID(updated.ID)).
		WithCategory().
		WithInvoice().
		Only(ctx)

	if err != nil {
		return nil, appError.FailedToFind(transactionEntity, err)
	}

	return newTransactionResponse(row)
}

func (p *PostgreSQL) ListTransactions(ctx context.Context, flt dto.TransactionFilters, pgn *pagination.Pagination) ([]dto.TransactionResponse, error) {
	userID, err := authctx.GetUserID(ctx)
	if err != nil {
		return nil, err
	}

	query := p.Client.Transaction.Query().
		Where(transaction.HasUserWith(user.IDEQ(userID))).
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
	userID, err := authctx.GetUserID(ctx)
	if err != nil {
		return 0, err
	}

	query := p.Client.Transaction.Query().
		Where(transaction.HasUserWith(user.IDEQ(userID)))

	query = applyTransactionFilters(query, flt, pgn)

	total, err := query.Count(ctx)
	if err != nil {
		return 0, err
	}
	return total, nil
}

func (p *PostgreSQL) TransactionsGeneralStats(ctx context.Context, flt dto.ChartFilters) (*dto.TransactionStatsSummary, error) {
	userID, err := authctx.GetUserID(ctx)
	if err != nil {
		return nil, err
	}

	startDate, err := utils.ToDateTime(flt.StartDate)
	if err != nil {
		return nil, err
	}
	endDate, err := utils.ToDateTime(flt.EndDate)
	if err != nil {
		return nil, err
	}

	var dateExpr string
	switch flt.DateField {
	case "due_date":
		dateExpr = "COALESCE(i.due_date, t.record_date)"
	case "record_date":
		dateExpr = "t.record_date"
	default:
		return nil, fmt.Errorf("invalid dateField: %s", flt.DateField)
	}

	query := fmt.Sprintf(`
		SELECT
			SUM(CASE WHEN t.record_type = 'income' THEN t.amount ELSE 0 END) AS income,
			SUM(CASE WHEN t.record_type = 'expense' THEN t.amount ELSE 0 END) AS expense,
			SUM(CASE WHEN t.record_type = 'tax' THEN t.amount ELSE 0 END) AS tax,
			COUNT(CASE WHEN t.record_type = 'income' THEN 1 END) AS incomeTransactions,
	        COUNT(CASE WHEN t.record_type = 'expense' THEN 1 END) AS expenseTransactions
		FROM transactions AS t
			LEFT JOIN invoices AS i ON t.invoice_id = i.id
		WHERE t.user_id = $1
		AND %s BETWEEN $2 AND $3
	`, dateExpr)

	var income float64
	var expense float64
	var tax float64
	var incomeTransactions int
	var expenseTransactions int

	err = p.db.QueryRowContext(ctx, query, userID, startDate, endDate).
		Scan(&income, &expense, &tax, &incomeTransactions, &expenseTransactions)

	if err != nil {
		return nil, err
	}

	return &dto.TransactionStatsSummary{
		Income:              income - tax,
		Expense:             expense,
		Tax:                 tax,
		Balance:             income - expense - tax,
		IncomeTransactions:  incomeTransactions,
		ExpenseTransactions: expenseTransactions,
	}, nil
}

func (p *PostgreSQL) TransactionsSummary(ctx context.Context, flt dto.ChartFilters) ([]dto.SummaryByDate, error) {
	userID, err := authctx.GetUserID(ctx)
	if err != nil {
		return nil, err
	}

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
	allCategories := []string{}
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

	var dateExpr string
	switch flt.DateField {
	case "due_date":
		dateExpr = "COALESCE(i.due_date, t.record_date)"
	case "record_date":
		dateExpr = "t.record_date"
	default:
		return nil, fmt.Errorf("invalid dateField: %s", flt.DateField)
	}

	query := fmt.Sprintf(`
		SELECT DATE_TRUNC($1, %s) AS period,
			c.name AS category,
			SUM(CASE WHEN t.record_type = 'income' THEN t.amount ELSE 0 END) AS income,
			SUM(CASE WHEN t.record_type = 'expense' THEN t.amount ELSE 0 END) AS expense,
			SUM(CASE WHEN t.record_type = 'tax' THEN t.amount ELSE 0 END) AS tax,
			COUNT(CASE WHEN t.record_type = 'income' THEN 1 END) AS incomeTransactions,
			COUNT(CASE WHEN t.record_type = 'expense' THEN 1 END) AS expenseTransactions
		FROM transactions t
			LEFT JOIN invoices AS i ON t.invoice_id = i.id
			LEFT JOIN categories AS c ON t.category_id = c.id
		WHERE t.user_id = $2
		AND %s BETWEEN $3 AND $4
		GROUP BY period, c.name
		ORDER BY period
	`, dateExpr, dateExpr)

	rows, err := p.db.QueryContext(ctx, query, periodTrunc, userID, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	type rawData struct {
		period              string
		category            string
		income              float64
		expense             float64
		tax                 float64
		incomeTransactions  int
		expenseTransactions int
	}

	periodByDate := map[string][]rawData{}

	for rows.Next() {
		var period time.Time
		var category string
		var income, expense, tax float64
		var incomeTransactions, expenseTransactions int

		if err := rows.Scan(&period, &category, &income, &expense, &tax, &incomeTransactions, &expenseTransactions); err != nil {
			return nil, err
		}

		key := period.Format("2006-01-02")
		periodByDate[key] = append(periodByDate[key], rawData{
			period:              key,
			category:            category,
			income:              income,
			expense:             expense,
			tax:                 tax,
			incomeTransactions:  incomeTransactions,
			expenseTransactions: expenseTransactions,
		})
	}

	var result []dto.SummaryByDate
	for date, entries := range periodByDate {
		summary := dto.SummaryByDate{
			Date:       date,
			Income:     0,
			Expense:    0,
			Categories: []dto.CategorySummary{},
		}

		catMap := map[string]rawData{}
		for _, entry := range entries {
			catMap[entry.category] = entry
			summary.Income += entry.income
			summary.Expense += entry.expense
			summary.Tax += entry.tax
		}

		for _, category := range allCategories {
			data, exists := catMap[category]
			income := 0.0
			expense := 0.0
			tax := 0.0
			incomeTransactions := 0
			expenseTransactions := 0

			if exists {
				income = data.income
				expense = data.expense
				tax = data.tax
				incomeTransactions = data.incomeTransactions
				expenseTransactions = data.expenseTransactions
			}

			summary.Categories = append(summary.Categories, dto.CategorySummary{
				Category:            category,
				Income:              income,
				Expense:             expense,
				Tax:                 tax,
				IncomeTransactions:  incomeTransactions,
				ExpenseTransactions: expenseTransactions,
			})
		}

		result = append(result, summary)
	}

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
