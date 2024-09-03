package model

import "time"

type CheckDTO struct {
	ID   int    `gorm:"column:id"`
	Name string `gorm:"column:description"`
}

type AvailableClientDTO struct {
	ClientID             int     `gorm:"column:client_id"`
	ClientBalance        float64 `gorm:"column:client_balance"`
	ClientWeight         int     `gorm:"column:client_weight"`
	UserID               int     `gorm:"column:user_id"`
	UserBalance          float64 `gorm:"column:user_balance"`
	UserAllowIPs         string  `gorm:"column:user_allow_ips"`
	UserRole             string  `gorm:"column:user_role"`
	ModelName            string  `gorm:"column:model_name"`
	ModelID              int     `gorm:"column:model_id"`
	ModelMaxToken        int     `gorm:"column:model_max_tokens"`
	ModelPromptPrice     float64 `gorm:"column:model_prompt_price"`
	ModelCompletionPrice float64 `gorm:"column:model_completion_price"`
}

type ClientSecretDTO struct {
	ClientID       int     `gorm:"column:client_id"`
	ClientKey      string  `gorm:"column:client_key"`
	ClientEndpoint string  `gorm:"column:client_endpoint"`
	ClientWeight   int     `gorm:"column:client_weight"`
	ClientBalance  float64 `gorm:"column:client_balance"`
}

type RelatedModelDTO struct {
	ModelID         int       `gorm:"column:model_id"`
	ModelName       string    `gorm:"column:model_name"`
	MaxTokens       int       `gorm:"column:model_max_tokens"`
	PromptPrice     float64   `gorm:"column:model_prompt_price"`
	CompletionPrice float64   `gorm:"column:model_completion_price"`
	ModelRpmLimit   int       `gorm:"column:model_rpm_limit"`
	ModelTpmLimit   int       `gorm:"column:model_tpm_limit"`
	LastUpdatedAt   time.Time `gorm:"column:last_updated_at"`
}

type ClientModelDTO struct {
	Name            string
	MaxTokens       int
	PromptPrice     float64
	CompletionPrice float64
	RpmLimit        int
	TpmLimit        int
}

type ClientDTO struct {
	Description string
	ApiKey      string
	Endpoint    string
	Balance     float64
	Weight      int
}

type RequestRecordDTO struct {
	ClientID             int
	ModelID              int
	UserID               int
	PromptTokenUsage     int
	CompletionTokenUsage int
	RequestIP            string
	BalanceCost          float64
}

type WhisperUserDTO struct {
	Email    string   `gorm:"column:email"`
	ApiKey   string   `gorm:"column:api_key"`
	Balance  float64  `gorm:"column:balance"`
	Role     string   `gorm:"column:role"`
	AllowIPs []string `gorm:"column:allow_ips"`
}

type QueryWhisperUserDTO struct {
	Email    string  `gorm:"column:email"`
	ApiKey   string  `gorm:"column:api_key"`
	Balance  float64 `gorm:"column:balance"`
	Role     string  `gorm:"column:role"`
	AllowIPs string  `gorm:"column:allow_ips"`
}

type UserPermissionItem struct {
	Desc   string
	Models []string
}
