package global

import (
	"context"
	"path/filepath"
	"time"

	"github.com/alioth-center/akasha-whisper/app/dao"
	"github.com/alioth-center/akasha-whisper/app/model"
	"github.com/alioth-center/infrastructure/config"
	acdb "github.com/alioth-center/infrastructure/database"
	"github.com/alioth-center/infrastructure/database/mysql"
	"github.com/alioth-center/infrastructure/database/postgres"
	"github.com/alioth-center/infrastructure/database/sqlite"
	"github.com/alioth-center/infrastructure/logger"
	"github.com/alioth-center/infrastructure/thirdparty/openai"
	"github.com/alioth-center/infrastructure/trace"
	"github.com/alioth-center/infrastructure/utils/concurrency"
	"github.com/alioth-center/infrastructure/utils/values"
	"github.com/bits-and-blooms/bloom/v3"
)

var syncModels = []any{
	&model.OpenaiClient{}, &model.OpenaiClientBalance{}, &model.OpenaiModel{}, &model.OpenaiRequest{},
	&model.WhisperUser{}, &model.WhisperUserBalance{}, &model.WhisperUserPermission{},
}

func init() {
	// initialize background context
	ctx := trace.NewContext()

	// read config first
	initializeConfig()

	// initialize logger
	initializeLogger()

	// initialize databases
	initializeDatabase()

	// initialize cache
	initializeCache()

	// initialize bloom filter
	initializeBloomFilter(ctx)
}

func initializeConfig() {
	readErr := config.LoadConfig(&Config, "./config/config.yaml")
	if readErr != nil {
		panic(readErr)
	}
}

func initializeLogger() {
	if Config.Logger.LogFileSuffix == "" {
		Config.Logger.LogFileSuffix = "_akasha_whisper_logs_stdout.jsonl"
	}
	switch {
	case !Config.Logger.LogToFile:
		Logger = logger.NewCustomLoggerWithOpts(
			logger.WithLevelOpts(logger.Level(Config.Logger.LogLevel)),
		)
	case !Config.Logger.LogSplit:
		Logger = logger.NewCustomLoggerWithOpts(
			logger.WithLevelOpts(logger.Level(Config.Logger.LogLevel)),
			logger.WithFileWriterOpts(filepath.Join(Config.Logger.LogDirectory, Config.Logger.LogFileSuffix)),
		)
	default:
		Logger = logger.NewCustomLoggerWithOpts(
			logger.WithLevelOpts(logger.Level(Config.Logger.LogLevel)),
			logger.WithCustomWriterOpts(
				logger.NewTimeBasedRotationFileWriter(Config.Logger.LogDirectory, func(time time.Time) (filename string) {
					return values.BuildStrings(time.Format("2006-01-02"), Config.Logger.LogFileSuffix)
				}),
			),
		)
	}
}

func initializeDatabase() {
	var database acdb.DatabaseV2
	switch Config.Database.Driver {
	case postgres.DriverName:
		pgCfg := postgres.Config{
			Host:      Config.Database.Host,
			Port:      Config.Database.Port,
			Username:  Config.Database.Username,
			Password:  Config.Database.Password,
			Database:  Config.Database.Database,
			EnableSSL: Config.Database.SSL,
		}
		pgDB, initErr := postgres.NewWithLogger(pgCfg, Logger, syncModels...)
		if initErr != nil {
			panic(initErr)
		}

		database = pgDB
	case mysql.DriverName:
		mysqlCfg := mysql.Config{
			Server:   Config.Database.Host,
			Port:     Config.Database.Port,
			Username: Config.Database.Username,
			Password: Config.Database.Password,
			Database: Config.Database.Database,
		}
		mysqlDB, initErr := mysql.NewWithLogger(mysqlCfg, Logger, syncModels...)
		if initErr != nil {
			panic(initErr)
		}

		database = mysqlDB
	case sqlite.DriverName:
		sqliteCfg := sqlite.Config{
			Database: "./data/akasha_whisper.db",
		}
		sqliteDB, initErr := sqlite.NewWithLogger(sqliteCfg, Logger, syncModels...)
		if initErr != nil {
			panic(initErr)
		}

		database = sqliteDB
	default:
		panic("unsupported database driver")
	}

	DatabaseInstance = database
	OpenaiClientDatabaseInstance = dao.NewOpenaiClientDatabaseAccessor(database)
	OpenaiClientBalanceDatabaseInstance = dao.NewOpenaiClientBalanceDatabaseAccessor(database)
	OpenaiModelDatabaseInstance = dao.NewOpenaiModelDatabaseAccessor(database)
	OpenaiRequestDatabaseInstance = dao.NewOpenaiRequestDatabaseAccessor(database)
	WhisperUserDatabaseInstance = dao.NewWhisperUserDatabaseAccessor(database)
	WhisperUserBalanceDatabaseInstance = dao.NewWhisperUserBalanceDatabaseAccessor(database)
	WhisperUserPermissionDatabaseInstance = dao.NewWhisperUserPermissionDatabaseAccessor(database)

	dao.LoadRawSqlList(Config.Database.Driver)
}

func initializeCache() {
	OpenaiClientCacheInstance = concurrency.NewMap[int, openai.Client]()
}

func initializeBloomFilter(ctx context.Context) {
	if !Config.BloomFilter.Enable {
		BearerTokenBloomFilterInstance = dao.NewBearerTokenBloomFilter(nil)
		return
	}

	filter := bloom.NewWithEstimates(uint(Config.BloomFilter.FilterSize), Config.BloomFilter.FalseRate)
	BearerTokenBloomFilterInstance = dao.NewBearerTokenBloomFilter(filter)

	// load all tokens from database
	Logger.Info(logger.NewFields(ctx).WithMessage("bloom filter enable, loading whisper user api keys"))
	tokens, loadErr := WhisperUserDatabaseInstance.ListWhisperUserApiKeys(ctx)
	if loadErr != nil {
		panic(loadErr)
	}
	Logger.Info(logger.NewFields(ctx).WithMessage("whisper user api keys loaded").WithData(map[string]any{"key_count": len(tokens)}))

	// add all tokens to bloom filter
	BearerTokenBloomFilterInstance.AddKeys(tokens...)
}
