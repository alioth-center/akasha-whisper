package global

import (
	"github.com/alioth-center/akasha-whisper/app/dao"
	"github.com/alioth-center/akasha-whisper/app/model"
	"github.com/alioth-center/infrastructure/config"
	"github.com/alioth-center/infrastructure/database/postgres"
	"github.com/alioth-center/infrastructure/logger"
	"github.com/alioth-center/infrastructure/utils/values"
	"time"
)

var syncModels = []any{
	&model.OpenaiClient{}, &model.OpenaiClientBalance{}, &model.OpenaiModel{}, &model.OpenaiRequest{},
	&model.WhisperUser{}, &model.WhisperUserBalance{}, &model.WhisperUserPermission{},
}

func init() {
	// read config first
	readErr := config.LoadConfig(&Config, "./config/config.yaml")
	if readErr != nil {
		panic(readErr)
	}

	// initialize logger
	Logger = logger.NewCustomLoggerWithOpts(logger.WithCustomWriterOpts(logger.NewTimeBasedRotationFileWriter(Config.LogDir, func(time time.Time) (filename string) {
		return values.BuildStrings(time.Format("2006-01-02"), "_stdout.jsonl")
	})))

	// initialize databases
	database, initErr := postgres.NewPostgresSQLv2(Config.Database, syncModels...)
	if initErr != nil {
		panic(initErr)
	}
	DatabaseInstance = database
	OpenaiClientDatabaseInstance = dao.NewOpenaiClientDatabaseAccessor(database)
	OpenaiClientBalanceDatabaseInstance = dao.NewOpenaiClientBalanceDatabaseAccessor(database)
	OpenaiModelDatabaseInstance = dao.NewOpenaiModelDatabaseAccessor(database)
	OpenaiRequestDatabaseInstance = dao.NewOpenaiRequestDatabaseAccessor(database)
	WhisperUserDatabaseInstance = dao.NewWhisperUserDatabaseAccessor(database)
	WhisperUserBalanceDatabaseInstance = dao.NewWhisperUserBalanceDatabaseAccessor(database)
}
