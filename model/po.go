package model

import (
	"time"
)

const (
	TableNameOpenaiClients          = "openai_clients"
	TableNameOpenaiModels           = "openai_models"
	TableNameOpenaiRequests         = "openai_requests"
	TableNameWhisperUsers           = "whisper_users"
	TableNameWhisperUserPermissions = "whisper_user_permissions"
)

// OpenaiClient openai service secret
//
// Reference: https://docs.alioth.center/akasha-whisper-database.html#openai-clients
type OpenaiClient struct {
	ID          int64     `gorm:"column:id;type:integer;autoIncrement:true;primaryKey;index:idx_ids"`
	Description string    `gorm:"column:description;type:varchar(64);not null;comment:openai_service_description;uniqueIndex:idx_desc"`
	ApiKey      string    `gorm:"column:api_key;type:varchar(64);not null;comment:openai_service_api_key;index:idx_api_key"`
	Endpoint    string    `gorm:"column:endpoint;type:varchar(64);not null;comment:openai_service_endpoint;index:idx_endpoint"`
	Weight      int       `gorm:"column:weight;type:integer;not null;comment:openai_service_weight;index:idx_weight"`
	Balance     float64   `gorm:"column:balance;type:decimal(16,8);not null;comment:openai_service_balance;index:idx_balance"`
	CreatedAt   time.Time `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP"`
}

func (c OpenaiClient) TableName() string {
	return TableNameOpenaiClients
}

// OpenaiModel openai model
//
// Reference: https://docs.alioth.center/akasha-whisper-database.html#openai-models
type OpenaiModel struct {
	ID              int64     `gorm:"column:id;type:integer;autoIncrement:true;primaryKey;uniqueIndex:idx_ids"`
	ClientID        int64     `gorm:"column:client_id;type:integer;not null;comment:openai_client_id;uniqueIndex:idx_ids;uniqueIndex:idx_names;index:idx_client_ids"`
	Model           string    `gorm:"column:model;type:varchar(32);not null;comment:openai_model_name;index:idx_name;uniqueIndex:idx_names"`
	MaxTokens       int       `gorm:"column:max_tokens;type:integer;not null;comment:openai_max_tokens"`
	PromptPrice     float64   `gorm:"column:prompt_price;type:decimal(16,8);not null;comment:openai_prompt_price"`
	CompletionPrice float64   `gorm:"column:completion_price;type:decimal(16,8);not null;comment:openai_completion_price"`
	RpmLimit        int       `gorm:"column:rpm_limit;type:integer;not null;default:-1;comment:openai_rpm_limit"`
	TpmLimit        int       `gorm:"column:tpm_limit;type:integer;not null;default:-1;comment:openai_tpm_limit"`
	CreatedAt       time.Time `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt       time.Time `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP"`
}

func (m OpenaiModel) TableName() string {
	return TableNameOpenaiModels
}

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

// WhisperUser whisper user
//
// Reference: https://docs.alioth.center/akasha-whisper-database.html#whisper-users
type WhisperUser struct {
	ID        int64     `gorm:"column:id;type:integer;autoIncrement:true;primaryKey;index:idx_ids"`
	Email     string    `gorm:"column:email;type:varchar(64);not null;comment:whisper_user_email;uniqueIndex:idx_emails"`
	ApiKey    string    `gorm:"column:api_key;type:varchar(64);not null;comment:whisper_user_api_key;uniqueIndex:idx_api_keys"`
	Balance   float64   `gorm:"column:balance;type:decimal(16,8);not null;comment:whisper_user_balance;index:idx_balances"`
	Role      string    `gorm:"column:role;type:varchar(10);not null;comment:whisper_user_role;index:idx_roles"`
	Language  string    `gorm:"column:language;type:varchar(2);not null;default:en;comment:whisper_user_language"`
	AllowIps  string    `gorm:"column:allow_ips;type:varchar(128);not null;comment:whisper_user_allow_ips;index:idx_allow_ips"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP"`
}

func (u WhisperUser) TableName() string {
	return TableNameWhisperUsers
}

// WhisperUserPermission whisper user permission
//
// Reference: https://docs.alioth.center/akasha-whisper-database.html#whisper-user-permissions
type WhisperUserPermission struct {
	ID        int64     `gorm:"column:id;type:integer;autoIncrement:true;primaryKey;index:idx_ids"`
	UserID    int64     `gorm:"column:user_id;type:integer;not null;comment:whisper_user_id;index:idx_user_ids;uniqueIndex:idx_perm"`
	ModelID   int64     `gorm:"column:model_id;type:integer;not null;comment:whisper_model_id;index:idx_model_ids;uniqueIndex:idx_perm"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP"`
}

func (p WhisperUserPermission) TableName() string {
	return TableNameWhisperUserPermissions
}
