package model

import "time"

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
