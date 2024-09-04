package model

import "time"

// OpenaiClient openai service secret
//
// Reference: https://docs.alioth.center/akasha-whisper-database.html#openai-clients
type OpenaiClient struct {
	ID          int64     `gorm:"column:id;type:integer;autoIncrement:true;primaryKey;index:idx_ids"`
	Description string    `gorm:"column:description;type:varchar(64);not null;comment:openai_service_description;uniqueIndex:idx_desc"`
	ApiKey      string    `gorm:"column:api_key;type:varchar(64);not null;comment:openai_service_api_key;index:idx_api_key"`
	Endpoint    string    `gorm:"column:endpoint;type:varchar(64);not null;comment:openai_service_endpoint;index:idx_endpoint"`
	Weight      int       `gorm:"column:weight;type:integer;not null;comment:openai_service_weight;index:idx_weight"`
	CreatedAt   time.Time `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP"`
}

func (c OpenaiClient) TableName() string {
	return TableNameOpenaiClients
}
