package postgresql

import (
	"fmt"
	"frog-go/internal/config"
	"frog-go/internal/core/domain"
	"frog-go/internal/core/dto"
	"frog-go/internal/core/errors"
	"frog-go/internal/ent"
	"frog-go/internal/ent/category"
	"frog-go/internal/ent/debt"
	"frog-go/internal/utils"
	"sort"
	"time"

	"context"
	"frog-go/internal/utils/pagination"

	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
)

const debtEntity = "debts"

func (d *PostgreSQL) GetDebtByID(ctx context.Context, id uuid.UUID) (*dto.DebtResponse, error) {
	row, err := d.Client.Debt.Query().
		Where(debt.IDEQ(id)).
		WithCategory().
		Only(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.ErrNotFound
		}
		return nil, err
	}
	return newDebtResponse(row)
}

func (d *PostgreSQL) DeleteDebtByID(ctx context.Context, id uuid.UUID) error {
	err := d.Client.Debt.DeleteOneID(id).Exec(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return errors.ErrNotFound
		}
		return err
	}
	return nil
}

func (d *PostgreSQL) CreateDebt(ctx context.Context, input domain.Debt) (*dto.DebtResponse, error) {
	created, err := d.Client.Debt.
		Create().
		SetTitle(input.Title).
		SetAmount(input.Amount).
		SetStatus(string(input.Status)).
		SetNillableDueDate(input.DueDate).
		SetPurchaseDate(input.PurchaseDate).
		SetNillableCategoryID(input.CategoryID).
		Save(ctx)

	if err != nil {
		return nil, errors.FailedToSave(debtEntity, err)
	}

	row, err := d.Client.Debt.Query().
		Where(debt.ID(created.ID)).
		WithCategory().
		Only(ctx)

	if err != nil {
		return nil, errors.FailedToFind(debtEntity, err)
	}

	return newDebtResponse(row)
}

func (d *PostgreSQL) UpdateDebt(ctx context.Context, id uuid.UUID, input domain.Debt) (*dto.DebtResponse, error) {
	updated, err := d.Client.Debt.
		UpdateOneID(id).
		SetTitle(input.Title).
		SetAmount(input.Amount).
		SetStatus(string(input.Status)).
		SetNillableDueDate(input.DueDate).
		SetPurchaseDate(input.PurchaseDate).
		SetNillableCategoryID(input.CategoryID).
		Save(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.ErrNotFound
		}
		return nil, errors.FailedToSave(debtEntity, err)
	}

	row, err := d.Client.Debt.Query().
		Where(debt.ID(updated.ID)).
		WithCategory().
		Only(ctx)

	if err != nil {
		return nil, errors.FailedToFind(debtEntity, err)
	}

	return newDebtResponse(row)
}

func (d *PostgreSQL) ListDebts(ctx context.Context, flt dto.DebtFilters, pgn *pagination.Pagination) ([]dto.DebtResponse, error) {
	query := d.Client.Debt.Query().
		WithCategory()

	query = applyDebtFilters(query, flt, pgn)
	query = apllyDebtOrderBy(query, pgn)
	query = query.Limit(pgn.PageSize).Offset(pgn.Offset())

	data, err := query.All(ctx)
	if err != nil {
		return nil, err
	}

	return newDebtResponseList(data)
}

func (d *PostgreSQL) CountDebts(ctx context.Context, flt dto.DebtFilters, pgn *pagination.Pagination) (int, error) {
	query := d.Client.Debt.Query()
	query = applyDebtFilters(query, flt, pgn)

	total, err := query.Count(ctx)
	if err != nil {
		return 0, err
	}
	return total, nil
}

func (d *PostgreSQL) DebtsGeneralStats(ctx context.Context, flt dto.ChartFilters) (*dto.DebtStatsSummary, error) {
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
			COALESCE(SUM(d.amount), 0) AS total_amount,
			COUNT(*) AS total_transactions,
			COUNT(DISTINCT d.title) AS unique_establishments
		FROM debts d
		WHERE d.purchase_date BETWEEN $1 AND $2
	`

	var totalAmount float64
	var totalTransactions int
	var uniqueEstablishments int

	err = d.db.QueryRowContext(ctx, query, startDate, endDate).Scan(&totalAmount, &totalTransactions, &uniqueEstablishments)
	if err != nil {
		return nil, err
	}

	var averagePerTransaction float64
	if totalTransactions > 0 {
		averagePerTransaction = totalAmount / float64(totalTransactions)
	}

	return &dto.DebtStatsSummary{
		TotalAmount:           totalAmount,
		TotalTransactions:     totalTransactions,
		UniqueEstablishments:  uniqueEstablishments,
		AveragePerTransaction: averagePerTransaction,
	}, nil
}

func (d *PostgreSQL) DebtsSummary(ctx context.Context, flt dto.ChartFilters) ([]dto.SummaryByDate, error) {
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
	catRows, err := d.db.QueryContext(ctx, catQuery)
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
			DATE_TRUNC($1, d.purchase_date) AS period,
			COALESCE(c.name, 'Sem categoria') AS category,
			SUM(d.amount) AS total,
			COUNT(*) AS transactions
		FROM debts d
		LEFT JOIN categories c ON d.category_id = c.id
		WHERE d.purchase_date BETWEEN $2 AND $3
		GROUP BY period, category
		ORDER BY period
	`

	rows, err := d.db.QueryContext(ctx, query, periodTrunc, startDate, endDate)
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

func mapDebtToResponse(row *ent.Debt) dto.DebtResponse {
	response := dto.DebtResponse{
		ID:           row.ID,
		Title:        row.Title,
		Amount:       row.Amount,
		Status:       row.Status,
		PurchaseDate: utils.ToDateTimeString(row.PurchaseDate),
		DueDate:      utils.ToNillableDateTimeString(row.DueDate),
		CreatedAt:    utils.ToDateTimeString(row.CreatedAt),
		UpdatedAt:    utils.ToDateTimeString(row.UpdatedAt),
	}

	if row.Edges.Category != nil {
		response.Category = &dto.DebtCategoryResponse{
			ID:   row.Edges.Category.ID,
			Name: row.Edges.Category.Name,
		}
	}

	return response
}

func newDebtResponse(row *ent.Debt) (*dto.DebtResponse, error) {
	if row == nil {
		return nil, nil
	}
	response := mapDebtToResponse(row)
	return &response, nil
}

func newDebtResponseList(rows []*ent.Debt) ([]dto.DebtResponse, error) {
	if rows == nil {
		return nil, nil
	}
	response := make([]dto.DebtResponse, 0, len(rows))
	for _, row := range rows {
		response = append(response, mapDebtToResponse(row))
	}
	return response, nil
}

func apllyDebtOrderBy(query *ent.DebtQuery, pgn *pagination.Pagination) *ent.DebtQuery {

	var orderDirection sql.OrderTermOption
	if pgn.OrderDirection == config.OrderAsc {
		orderDirection = sql.OrderAsc()
	} else {
		orderDirection = sql.OrderDesc()
	}

	switch pgn.OrderBy {
	case "category":
		query.Order(
			debt.ByCategoryField(category.FieldName, orderDirection),
			debt.ByID(sql.OrderAsc()),
		)
	case "status":
		query.Order(
			debt.ByStatus(orderDirection),
			debt.ByID(sql.OrderAsc()),
		)
	default:
		if pgn.OrderDirection == config.OrderAsc {
			query = query.Order(
				ent.Asc(pgn.OrderBy),
				ent.Asc(debt.FieldID),
			)
		} else {
			query = query.Order(
				ent.Desc(pgn.OrderBy),
				ent.Asc(debt.FieldID),
			)
		}
	}

	return query
}

func applyDebtFilters(query *ent.DebtQuery, flt dto.DebtFilters, pgn *pagination.Pagination) *ent.DebtQuery {
	if pgn.Search != "" {
		query = query.Where(
			debt.Or(
				debt.TitleContainsFold(pgn.Search),
				debt.StatusContainsFold(pgn.Search),
				debt.HasCategoryWith(
					category.NameContainsFold(pgn.Search),
				),
			),
		)
	}

	if flt.Status != nil && len(*flt.Status) > 0 {
		query = query.Where(
			debt.StatusIn((*flt.Status)...),
		)
	}

	if flt.CategoryID != nil {
		categoryIds := utils.ToUUIDSlice(*flt.CategoryID)
		if len(categoryIds) > 0 {
			query = query.Where(
				debt.HasCategoryWith(category.IDIn(categoryIds...)),
			)
		}
	}

	if flt.MinAmount != nil {
		query = query.Where(
			debt.AmountGTE(*flt.MinAmount),
		)
	}
	if flt.MaxAmount != nil {
		query = query.Where(
			debt.AmountLTE(*flt.MaxAmount),
		)
	}
	if t := utils.ToDateTimeUnsafe(flt.StartDate); t != nil {
		query = query.Where(debt.PurchaseDateGTE(*t))
	}

	if t := utils.ToDateTimeUnsafe(flt.EndDate); t != nil {
		query = query.Where(debt.PurchaseDateLTE(*t))
	}

	return query
}
