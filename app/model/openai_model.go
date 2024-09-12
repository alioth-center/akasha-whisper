package model

import (
	"time"

	"github.com/shopspring/decimal"
)

// OpenaiModel openai model
//
// Reference: https://docs.alioth.center/akasha-whisper-database.html#openai-models
type OpenaiModel struct {
	ID              int64           `gorm:"column:id;type:integer;autoIncrement:true;primaryKey;uniqueIndex:idx_ids"`
	ClientID        int64           `gorm:"column:client_id;type:integer;not null;comment:openai_client_id;uniqueIndex:idx_ids;uniqueIndex:idx_names;index:idx_client_ids"`
	Model           string          `gorm:"column:model;type:varchar(32);not null;comment:openai_model_name;index:idx_name;uniqueIndex:idx_names"`
	MaxTokens       int             `gorm:"column:max_tokens;type:integer;not null;comment:openai_max_tokens"`
	PromptPrice     decimal.Decimal `gorm:"column:prompt_price;type:decimal(16,8);not null;comment:openai_prompt_price"`
	CompletionPrice decimal.Decimal `gorm:"column:completion_price;type:decimal(16,8);not null;comment:openai_completion_price"`
	RpmLimit        int             `gorm:"column:rpm_limit;type:integer;not null;default:-1;comment:openai_rpm_limit"`
	TpmLimit        int             `gorm:"column:tpm_limit;type:integer;not null;default:-1;comment:openai_tpm_limit"`
	CreatedAt       time.Time       `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt       time.Time       `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP"`
}

func (m OpenaiModel) TableName() string {
	return TableNameOpenaiModels
}
