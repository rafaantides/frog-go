package dto

type ChartFilters struct {
	Period    string `form:"period"`
	StartDate string `form:"start_date"`
	EndDate   string `form:"end_date"`
	DateField string `form:"date_field"`
}

type CategorySummary struct {
	Category            string  `json:"category"`
	Income              float64 `json:"income"`
	Expense             float64 `json:"expense"`
	Tax                 float64 `json:"tax"`
	IncomeTransactions  int     `json:"income_transactions"`
	ExpenseTransactions int     `json:"expense_transactions"`
}

type SummaryByDate struct {
	Date       string            `json:"date"`
	Income     float64           `json:"income"`
	Tax        float64           `json:"tax"`
	Expense    float64           `json:"expense"`
	Categories []CategorySummary `json:"categories"`
}

type TransactionStatsSummary struct {
	Income              float64 `json:"income"`
	Expense             float64 `json:"expense"`
	Tax                 float64 `json:"tax"`
	Balance             float64 `json:"balance"`
	IncomeTransactions  int     `json:"income_transactions"`
	ExpenseTransactions int     `json:"expense_transactions"`
}
