package dao

import (
	"context"

	"github.com/alioth-center/akasha-whisper/app/model"
	"github.com/alioth-center/akasha-whisper/app/model/dto"
	"github.com/alioth-center/infrastructure/database"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type WhisperUserPermissionDatabaseAccessor struct {
	db database.DatabaseV2
}

func NewWhisperUserPermissionDatabaseAccessor(db database.DatabaseV2) *WhisperUserPermissionDatabaseAccessor {
	return &WhisperUserPermissionDatabaseAccessor{db: db}
}

func (ac *WhisperUserPermissionDatabaseAccessor) CreatePermissionRecord(ctx context.Context, userID, clientID int, models ...string) (created []string, err error) {
	if len(models) == 0 {
		return []string{}, nil
	}

	if executeErr := ac.db.GetGormCore(ctx).Transaction(func(tx *gorm.DB) error {
		// query if the model exists
		modelDTOs := make([]model.OpenaiModel, 0, len(models))
		if queryErr := tx.WithContext(ctx).
			Model(&model.OpenaiModel{}).
			Where(model.OpenaiModelCols.ClientID, clientID).
			Where(model.OpenaiModelCols.Model, models).
			Select(model.OpenaiModelCols.Model, model.OpenaiModelCols.ID).
			Scan(&modelDTOs).
			Error; queryErr != nil {
			return queryErr
		}

		// insert into database
		permissions := make([]model.WhisperUserPermission, len(models))
		for i := range models {
			permissions[i] = model.WhisperUserPermission{
				UserID:  int64(userID),
				ModelID: modelDTOs[i].ID,
			}
			created = append(created, modelDTOs[i].Model)
		}

		return tx.CreateInBatches(permissions, 100).Error
	}); executeErr != nil {
		return nil, executeErr
	}

	return created, nil
}

func (ac *WhisperUserPermissionDatabaseAccessor) SyncPermissions(ctx context.Context, userID int, permissions map[string][]string) error {
	return ac.db.GetGormCore(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. query clients exists, if not, return error
		clientNames := make([]string, 0, len(permissions))
		clientNamesCond := make([]any, 0, len(permissions)) // clause.IN.Values use []any
		for clientName := range permissions {
			clientNames = append(clientNames, clientName)
			clientNamesCond = append(clientNamesCond, clientName)
		}
		var queryClients []model.OpenaiClient
		if queryClientsErr := tx.Model(&model.OpenaiClient{}).
			Where(model.OpenaiClientCols.Description, clientNames).
			Select(model.OpenaiClientCols.ID, model.OpenaiClientCols.Description).
			Scan(&queryClients).
			Error; queryClientsErr != nil {
			return queryClientsErr
		}
		if len(queryClients) != len(clientNames) {
			return gorm.ErrRecordNotFound
		}

		// additional: mapping client name to client id
		clientsMapping := map[int64]string{}
		for _, client := range queryClients {
			clientsMapping[client.ID] = client.Description
		}

		// 2. appends `WHERE client.client_name = ${client}` and `WHERE model.model_name in (${models})` to query conditions

		// select
		//		openai_clients.id as client_id,
		//		openai_clients.description as client_name,
		//		openai_models.id as model_id,
		//		openai_models.model as model_name
		// from openai_clients
		// join openai_models on openai_clients.id = openai_models.client_id and openai_clients.description in (${client_names})
		// where openai_clients = ${client} and openai_models.model in (${models})
		// 	or openai_clients.id = ${client_id} and openai_models.model in (${models}) ...
		modifyQuery := tx.Model(&model.OpenaiClient{}).Joins("?",
			&clause.Join{
				Type:  clause.InnerJoin,
				Table: clause.Table{Name: model.TableNameOpenaiModels},
				ON: clause.Where{
					Exprs: []clause.Expression{
						clause.Eq{
							Column: clause.Column{Table: model.TableNameOpenaiClients, Name: model.OpenaiClientCols.ID},
							Value:  clause.Column{Table: model.TableNameOpenaiModels, Name: model.OpenaiModelCols.ClientID},
						},
						clause.IN{
							Column: clause.Column{Table: model.TableNameOpenaiClients, Name: model.OpenaiClientCols.Description},
							Values: clientNamesCond,
						},
					},
				},
			},
		).Select(
			database.ColumnAlias(model.TableNameOpenaiClients, model.OpenaiClientCols.ID, "client_id"),
			database.ColumnAlias(model.TableNameOpenaiClients, model.OpenaiClientCols.Description, "client_name"),
			database.ColumnAlias(model.TableNameOpenaiModels, model.OpenaiModelCols.ID, "model_id"),
			database.ColumnAlias(model.TableNameOpenaiModels, model.OpenaiModelCols.Model, "model_name"),
		)
		for clientName, models := range permissions {
			var modelsCond []any
			for _, modelName := range models {
				modelsCond = append(modelsCond, modelName)
			}
			modifyQuery = modifyQuery.Where(clause.Eq{
				Column: clause.Column{Table: model.TableNameOpenaiClients, Name: model.OpenaiClientCols.Description},
				Value:  clientName,
			}).Where(clause.IN{
				Column: clause.Column{Table: model.TableNameOpenaiModels, Name: model.OpenaiModelCols.Model},
				Values: modelsCond,
			})
		}

		// 3. execute modify query and original permissions query
		var originalPermissions, modifyPermissions []dto.ClientModelDTO
		if queryOriginalPermissionsErr := tx.Model(&model.WhisperUserPermission{}).
			Joins("?", &clause.Join{
				Type:  clause.InnerJoin,
				Table: clause.Table{Name: model.TableNameOpenaiModels},
				ON: clause.Where{
					Exprs: []clause.Expression{
						clause.Eq{
							Column: clause.Column{Table: model.TableNameWhisperUserPermissions, Name: model.WhisperUserPermissionCols.ModelID},
							Value:  clause.Column{Table: model.TableNameOpenaiModels, Name: model.OpenaiModelCols.ID},
						},
					},
				},
			}).
			Joins("?", &clause.Join{
				Type:  clause.InnerJoin,
				Table: clause.Table{Name: model.TableNameOpenaiClients},
				ON: clause.Where{
					Exprs: []clause.Expression{
						clause.Eq{
							Column: clause.Column{Table: model.TableNameOpenaiModels, Name: model.OpenaiModelCols.ClientID},
							Value:  clause.Column{Table: model.TableNameOpenaiClients, Name: model.OpenaiClientCols.ID},
						},
					},
				},
			}).
			Select(
				database.ColumnAlias(model.TableNameOpenaiClients, model.OpenaiClientCols.ID, "client_id"),
				database.ColumnAlias(model.TableNameOpenaiClients, model.OpenaiClientCols.Description, "client_name"),
				database.ColumnAlias(model.TableNameOpenaiModels, model.OpenaiModelCols.ID, "model_id"),
				database.ColumnAlias(model.TableNameOpenaiModels, model.OpenaiModelCols.Model, "model_name"),
			).
			Where(model.WhisperUserPermissionCols.UserID, userID).
			Scan(&originalPermissions).
			Error; queryOriginalPermissionsErr != nil {
			return queryOriginalPermissionsErr
		}

		if queryModifyPermissionsErr := modifyQuery.Scan(&modifyPermissions).Error; queryModifyPermissionsErr != nil {
			return queryModifyPermissionsErr
		}

		// 4. mapping query result to map[client_name]map[model_name]*dto.ClientModelDTO{}
		currentPermissions := map[string]map[string]*dto.ClientModelDTO{}
		for _, permission := range originalPermissions {
			if currentPermissions[permission.ClientName] == nil {
				currentPermissions[permission.ClientName] = map[string]*dto.ClientModelDTO{}
			}

			currentPermissions[permission.ClientName][permission.ModelName] = &permission
		}

		// 5. statistics the difference between permissions and query result
		insertPermissions := map[string]map[string]*dto.ClientModelDTO{}
		for _, permission := range modifyPermissions {
			// currentPermissions[clientName] not exist, means need to insert model
			if currentPermissions[permission.ClientName] == nil {
				if insertPermissions[permission.ClientName] == nil {
					insertPermissions[permission.ClientName] = map[string]*dto.ClientModelDTO{}
				}

				insertPermissions[permission.ClientName][permission.ModelName] = &permission
				continue
			}

			// currentPermissions[clientName][modelName] not exist, means need to insert model
			if _, ok := currentPermissions[permission.ClientName][permission.ModelName]; !ok {
				if insertPermissions[permission.ClientName] == nil {
					insertPermissions[permission.ClientName] = map[string]*dto.ClientModelDTO{}
				}

				insertPermissions[permission.ClientName][permission.ModelName] = &permission
				continue
			}

			// currentPermissions[clientName][modelName] exist, remove it from currentPermissions
			// after scanning, the nodes left in currentPermissions are the nodes that need to delete
			delete(currentPermissions[permission.ClientName], permission.ModelName)
		}

		// 6. insert and delete the difference
		deleteModels := make([]int, 0, len(currentPermissions))
		for _, models := range insertPermissions {
			for _, modelDTO := range models {
				deleteModels = append(deleteModels, modelDTO.ModelID)
			}
		}
		if deleteErr := tx.Model(&model.WhisperUserPermission{}).
			Where(model.WhisperUserPermissionCols.UserID, userID).
			Where(model.WhisperUserPermissionCols.ModelID, deleteModels).
			Delete(&model.WhisperUserPermission{}).
			Error; deleteErr != nil {
			return deleteErr
		}

		insertPermissionsDTOs := make([]model.WhisperUserPermission, 0, len(insertPermissions))
		for _, models := range insertPermissions {
			for _, modelDTO := range models {
				insertPermissionsDTOs = append(insertPermissionsDTOs, model.WhisperUserPermission{
					UserID:  int64(userID),
					ModelID: int64(modelDTO.ModelID),
				})
			}
		}
		if insertErr := tx.Model(&model.WhisperUserPermission{}).
			CreateInBatches(insertPermissionsDTOs, 100).Error; insertErr != nil {
			return insertErr
		}

		return nil
	})
}
