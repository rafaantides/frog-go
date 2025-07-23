package dto

type ChartFilters struct {
	Period    string `form:"period"`
	StartDate string `form:"start_date"`
	EndDate   string `form:"end_date"`
}

type CategorySummary struct {
	Category            string  `json:"category"`
	Income              float64 `json:"income"`
	Expense             float64 `json:"expense"`
	IncomeTransactions  int     `json:"income_transactions"`
	ExpenseTransactions int     `json:"expense_transactions"`
}

type SummaryByDate struct {
	Date       string            `json:"date"`
	Income     float64           `json:"income"`
	Expense    float64           `json:"expense"`
	Categories []CategorySummary `json:"categories"`
}

type TransactionStatsSummary struct {
	Income              float64 `json:"income"`
	Expense             float64 `json:"expense"`
	Balance             float64 `json:"balance"`
	IncomeTransactions  int     `json:"income_transactions"`
	ExpenseTransactions int     `json:"expense_transactions"`
}
