package dto

import "time"

type DataRecord struct {
	Date     time.Time
	Category string
	Amount   float64
}

type ChartFilters struct {
	CategoryID *[]string `form:"category_id"`
	StatusID   *[]string `form:"status_id"`
	Period     string    `form:"period"`
	StartDate  string    `form:"start_date"`
	EndDate    string    `form:"end_date"`
}

type CategorySummary struct {
	Category     string  `json:"category"`
	Total        float64 `json:"total"`
	Transactions int     `json:"transactions"`
}

type SummaryByDate struct {
	Date       string            `json:"date"`
	Total      float64           `json:"total"`
	Categories []CategorySummary `json:"categories"`
}

type DebtStatsSummary struct {
	TotalAmount           float64 `json:"total_amount"`
	TotalTransactions     int     `json:"total_transactions"`
	UniqueEstablishments  int     `json:"unique_establishments"`
	AveragePerTransaction float64 `json:"average_per_transaction"`
}