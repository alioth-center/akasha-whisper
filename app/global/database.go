package global

import (
	"github.com/alioth-center/akasha-whisper/app/dao"
	"github.com/alioth-center/infrastructure/database"
)

var (
	DatabaseInstance database.DatabaseV2

	OpenaiClientDatabaseInstance          *dao.OpenaiClientDatabaseAccessor
	OpenaiClientBalanceDatabaseInstance   *dao.OpenaiClientBalanceDatabaseAccessor
	OpenaiModelDatabaseInstance           *dao.OpenaiModelDatabaseAccessor
	OpenaiRequestDatabaseInstance         *dao.OpenaiRequestDatabaseAccessor
	WhisperUserDatabaseInstance           *dao.WhisperUserDatabaseAccessor
	WhisperUserBalanceDatabaseInstance    *dao.WhisperUserBalanceDatabaseAccessor
	WhisperUserPermissionDatabaseInstance *dao.WhisperUserPermissionDatabaseAccessor
)
