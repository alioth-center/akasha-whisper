package model

import "time"

// OpenaiRequest openai request
//
// Reference: https://docs.alioth.center/akasha-whisper-database.html#openai-requests
type OpenaiRequest struct {
	ID                   int64     `gorm:"column:id;type:integer;autoIncrement:true;primaryKey;index:idx_ids"`
	ClientID             int64     `gorm:"column:client_id;type:integer;not null;comment:openai_client_id;index:idx_client_ids"`
	ModelID              int64     `gorm:"column:model_id;type:integer;not null;comment:openai_model_id;index:idx_model_ids"`
	UserID               int64     `gorm:"column:user_id;type:integer;not null;comment:openai_user_id;index:idx_user_ids"`
	RequestIP            string    `gorm:"column:request_ip;type:varchar(40);not null;comment:openai_request_ip;index:idx_request_ips"`
	PromptTokenUsage     int       `gorm:"column:prompt_token_usage;type:integer;not null;comment:openai_prompt_token_usage"`
	CompletionTokenUsage int       `gorm:"column:completion_token_usage;type:integer;not null;comment:openai_completion_token_usage"`
	BalanceCost          float64   `gorm:"column:balance_cost;type:decimal(16,8);not null;comment:openai_balance_cost;index:idx_balance_costs"`
	CreatedAt            time.Time `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP"`
}

func (r OpenaiRequest) TableName() string {
	return TableNameOpenaiRequests
}
