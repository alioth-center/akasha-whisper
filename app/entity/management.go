package entity

import (
	"github.com/alioth-center/infrastructure/network/http"
	"github.com/shopspring/decimal"
)

type OverviewRequest = http.NoBody

type OverviewResponse = http.BaseResponse[*OverviewResult]

type OverviewResult struct {
	Clients           []ClientItem               `json:"clients"`
	ClientBalanceLogs []OverviewClientBalanceLog `json:"client_balance_logs"`
}

type OverviewClientBalanceLog struct {
	ClientName   string          `json:"client_name"`
	TotalRequest int             `json:"total_request"`
	TotalCost    decimal.Decimal `json:"total_cost"`
	Date         string          `json:"date"`
}
