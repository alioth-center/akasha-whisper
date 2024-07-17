package global

import (
	"github.com/alioth-center/akasha-whisper/model"
	"github.com/alioth-center/infrastructure/database"
	"github.com/alioth-center/infrastructure/database/postgres"
)

var (
	Database database.Database

	syncModels = []any{
		&model.OpenaiClient{},
		&model.OpenaiModel{},
		&model.OpenaiRequest{},
		&model.WhisperUser{},
		&model.WhisperUserPermission{},
	}
)

func initDatabase() {
	db, initErr := postgres.NewPostgresDb(Config.Database, syncModels...)
	if initErr != nil {
		panic(initErr)
	}

	Database = db
}
