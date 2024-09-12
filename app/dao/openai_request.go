package dao

import (
	"context"

	"github.com/alioth-center/akasha-whisper/app/model"
	"github.com/alioth-center/infrastructure/database"
)

type OpenaiRequestDatabaseAccessor struct {
	db database.DatabaseV2
}

func NewOpenaiRequestDatabaseAccessor(db database.DatabaseV2) *OpenaiRequestDatabaseAccessor {
	return &OpenaiRequestDatabaseAccessor{db: db}
}

func (ac *OpenaiRequestDatabaseAccessor) CreateOpenaiRequestRecord(ctx context.Context, request *model.OpenaiRequest) (err error) {
	_, err = ac.db.CreateSingleDataIfNotExist(ctx, request)
	return err
}
