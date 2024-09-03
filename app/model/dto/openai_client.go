package dto

import "github.com/shopspring/decimal"

type ClientCheckDTO struct {
	ID   int    `gorm:"column:id"`
	Name string `gorm:"column:description"`
}

type AvailableClientDTO struct {
	ClientID             string          `gorm:"column:client_id"`
	ClientWeight         int             `gorm:"column:client_weight"`
	ClientBalance        decimal.Decimal `gorm:"column:client_balance"`
	UserID               int             `gorm:"column:user_id"`
	UserBalance          decimal.Decimal `gorm:"column:user_balance"`
	UserRole             string          `gorm:"column:user_role"`
	ModelID              int             `gorm:"column:model_id"`
	ModelName            string          `gorm:"column:model_name"`
	ModelMaxToken        int             `gorm:"column:model_max_token"`
	ModelPromptPrice     decimal.Decimal `gorm:"column:model_prompt_price"`
	ModelCompletionPrice decimal.Decimal `gorm:"column:model_completion_price"`
}

type ClientSecretDTO struct {
	ClientID       int             `gorm:"column:id"`
	ClientKey      string          `gorm:"column:api_key"`
	ClientEndpoint string          `gorm:"column:endpoint"`
	ClientWeight   int             `gorm:"column:weight"`
	ClientBalance  decimal.Decimal `gorm:"column:balance"`
}
