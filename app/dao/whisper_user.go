package dao

import "github.com/alioth-center/infrastructure/database"

type WhisperUserDatabaseAccessor struct {
	db database.DatabaseV2
}

func NewWhisperUserDatabaseAccessor(db database.DatabaseV2) *WhisperUserDatabaseAccessor {
	return &WhisperUserDatabaseAccessor{db: db}
}
