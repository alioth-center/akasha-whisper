package dao

import (
	"context"
	"github.com/alioth-center/akasha-whisper/app/model"
	"github.com/alioth-center/infrastructure/database"
)

type WhisperUserDatabaseAccessor struct {
	db database.DatabaseV2
}

func NewWhisperUserDatabaseAccessor(db database.DatabaseV2) *WhisperUserDatabaseAccessor {
	return &WhisperUserDatabaseAccessor{db: db}
}

func (ac *WhisperUserDatabaseAccessor) CheckWhisperUserApiKey(ctx context.Context, apiKey string) (bool, error) {
	var count int64
	if queryErr := ac.db.GetGormCore(ctx).
		Model(&model.WhisperUser{}).
		Where(model.WhisperUserCols.ApiKey, apiKey).
		Count(&count).
		Error; queryErr != nil {
		return false, queryErr
	}

	return count > 0, nil
}

func (ac *WhisperUserDatabaseAccessor) ListWhisperUserApiKeys(ctx context.Context) ([]string, error) {
	result := make([]string, 0)
	if queryErr := ac.db.GetGormCore(ctx).
		Model(&model.WhisperUser{}).
		Select(model.WhisperUserCols.ApiKey).
		Scan(&result).
		Error; queryErr != nil {
		return nil, queryErr
	}

	return result, nil
}

func (ac *WhisperUserDatabaseAccessor) CreateWhisperUser(ctx context.Context, user *model.WhisperUser) (created bool, err error) {
	return ac.db.CreateSingleDataIfNotExist(ctx, user)
}
