package model

import "time"

// WhisperUser whisper user
//
// Reference: https://docs.alioth.center/akasha-whisper-database.html#whisper-users
type WhisperUser struct {
	ID        int64     `gorm:"column:id;type:integer;autoIncrement:true;primaryKey;index:idx_ids"`
	Email     string    `gorm:"column:email;type:varchar(64);not null;comment:whisper_user_email;uniqueIndex:idx_emails"`
	ApiKey    string    `gorm:"column:api_key;type:varchar(64);not null;comment:whisper_user_api_key;uniqueIndex:idx_api_keys"`
	Role      string    `gorm:"column:role;type:varchar(10);not null;comment:whisper_user_role;index:idx_roles"`
	Language  string    `gorm:"column:language;type:varchar(2);not null;default:en;comment:whisper_user_language"`
	AllowIps  string    `gorm:"column:allow_ips;type:varchar(128);not null;comment:whisper_user_allow_ips;index:idx_allow_ips"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP"`
}

func (u WhisperUser) TableName() string {
	return TableNameWhisperUsers
}
