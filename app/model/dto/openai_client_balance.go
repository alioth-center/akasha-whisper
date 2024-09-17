package dto

import (
	"time"

	"github.com/shopspring/decimal"
)

type OpenaiClientBalanceStatisticsDTO struct {
	ClientID     int             `gorm:"column:client_id" json:"client_id"`
	ClientName   string          `gorm:"column:client_name" json:"client_name"`
	DateDay      time.Time       `gorm:"column:date_day" json:"date_day"`
	TotalCost    decimal.Decimal `gorm:"column:total_cost" json:"total_cost"`
	RequestCount int             `gorm:"column:request_count" json:"request_count"`
}
