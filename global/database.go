package global

import (
	"github.com/alioth-center/akasha-whisper/app/model"
	"github.com/alioth-center/infrastructure/database"
	"github.com/alioth-center/infrastructure/database/postgres"
)

var (
	Database   database.Database
	DatabaseV2 database.DatabaseV2

	syncModels = []any{
		&model.OpenaiClient{},
		&model.OpenaiModel{},
		&model.OpenaiRequest{},
		&model.WhisperUser{},
		&model.OpenaiClientBalance{},
		&model.WhisperUserPermission{},
		&model.WhisperUserBalance{},
	}
)

func initDatabase() {
	db, initErr := postgres.NewPostgresDb(Config.Database, syncModels...)
	if initErr != nil {
		panic(initErr)
	}

	v2, initV2Err := postgres.NewPostgresSQLv2(Config.Database, syncModels...)
	if initV2Err != nil {
		panic(initV2Err)
	}

	DatabaseV2 = v2
	Database = db
}
