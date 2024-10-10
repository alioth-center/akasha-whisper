package dto

import (
	"time"

	"github.com/shopspring/decimal"
)

type RelatedModelDTO struct {
	ModelID         int             `gorm:"column:model_id"`
	ModelName       string          `gorm:"column:model_name"`
	ModelType       string          `gorm:"column:model_type"`
	MaxTokens       int             `gorm:"column:model_max_tokens"`
	ModelRpmLimit   int             `gorm:"column:model_rpm_limit"`
	ModelTpmLimit   int             `gorm:"column:model_tpm_limit"`
	LastUpdatedAt   time.Time       `gorm:"column:last_updated_at"`
	PromptPrice     decimal.Decimal `gorm:"column:model_prompt_price"`
	CompletionPrice decimal.Decimal `gorm:"column:model_completion_price"`
}

type ClientModelDTO struct {
	ClientID   int    `gorm:"column:client_id"`
	ClientName string `gorm:"column:client_name"`
	ModelID    int    `gorm:"column:model_id"`
	ModelName  string `gorm:"column:model_name"`
}
