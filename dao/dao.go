package dao

import (
	"context"
	"github.com/alioth-center/akasha-whisper/model"
	"github.com/alioth-center/infrastructure/database"
	"github.com/alioth-center/infrastructure/database/extension/orm"
)

var accessor DatabaseAccessor

type DatabaseAccessor interface {
	GetUserByID(ctx context.Context, id int64) (result model.QueryWhisperUserDTO, err error)
	GetUserByApiKey(ctx context.Context, key string) (result model.QueryWhisperUserDTO, err error)
	GetAvailableClient(ctx context.Context, model, token string) (results []model.AvailableClientDTO, err error)
	GetAvailableUserClient(ctx context.Context, model, email string) (results []model.AvailableClientDTO, err error)
	GetClientSecret(ctx context.Context, clientID int) (result model.ClientSecretDTO, err error)
	GetModelsByClientID(ctx context.Context, clientID int) (result []model.RelatedModelDTO, err error)
	GetModelsByClientDesc(ctx context.Context, clientName string) (result []model.RelatedModelDTO, err error)
	GetModelIDByName(ctx context.Context, clientID int, modelName string) (id int, err error)
	GetAvailableModelsByApiKey(ctx context.Context, key string) (results []model.RelatedModelDTO, err error)
	GetAvailableModelsByEmail(ctx context.Context, email string) (results []model.RelatedModelDTO, err error)
	AddClient(ctx context.Context, client model.ClientDTO, models ...model.ClientModelDTO) (err error)
	AddClientModels(ctx context.Context, clientID int, models []model.ClientModelDTO) (err error)
	AddRequestRecord(ctx context.Context, request model.RequestRecordDTO) (err error)
	AddWhisperUser(ctx context.Context, user model.WhisperUserDTO, permissions ...model.UserPermissionItem) (err error)
}

// NewDatabaseAccessor initializes and returns an instance of DatabaseAccessor.
// This function ensures a singleton pattern is applied to the DatabaseAccessor instance.
// If an instance already exists, it returns the existing instance instead of creating a new one.
//
// Parameters:
//   - db: A database.Database object required for database operations.
//
// Returns:
//   - An instance of DatabaseAccessor that provides various database access methods.
func NewDatabaseAccessor(db database.Database) DatabaseAccessor {
	if accessor == nil {
		accessor = &postgresDatabaseAccessor{db: orm.NewExtension().InitializeExtension(db)}
	}

	return accessor
}
