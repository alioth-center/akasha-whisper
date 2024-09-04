package dao

import (
	"context"
	"github.com/alioth-center/akasha-whisper/app/model"
	"github.com/alioth-center/akasha-whisper/app/model/dto"
	"github.com/alioth-center/infrastructure/database"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type OpenaiModelDatabaseAccessor struct {
	db database.DatabaseV2
}

func NewOpenaiModelDatabaseAccessor(db database.DatabaseV2) *OpenaiModelDatabaseAccessor {
	return &OpenaiModelDatabaseAccessor{db: db}
}

func (ac *OpenaiModelDatabaseAccessor) GetModelsByClientID(ctx context.Context, clientID int) (result []*dto.RelatedModelDTO, err error) {
	result = make([]*dto.RelatedModelDTO, 0)
	needFields := []string{
		database.ColumnAlias(model.TableNameOpenaiClients, model.OpenaiClientCols.ID, "client_id"),
		database.ColumnAlias(model.TableNameOpenaiModels, model.OpenaiModelCols.ID, "model_id"),
		database.ColumnAlias(model.TableNameOpenaiModels, model.OpenaiModelCols.Model, "model_name"),
		database.ColumnAlias(model.TableNameOpenaiModels, model.OpenaiModelCols.MaxTokens, "model_max_tokens"),
		database.ColumnAlias(model.TableNameOpenaiModels, model.OpenaiModelCols.PromptPrice, "model_prompt_price"),
		database.ColumnAlias(model.TableNameOpenaiModels, model.OpenaiModelCols.CompletionPrice, "model_completion_price"),
		database.ColumnAlias(model.TableNameOpenaiModels, model.OpenaiModelCols.RpmLimit, "model_rpm_limit"),
		database.ColumnAlias(model.TableNameOpenaiModels, model.OpenaiModelCols.TpmLimit, "model_tpm_limit"),
		database.ColumnAlias(model.TableNameOpenaiModels, model.OpenaiModelCols.UpdatedAt, "last_updated_at"),
	}

	// select oc.id               as client_id,
	//        om.id               as model_id,
	//        om.model            as model_name,
	//        om.max_tokens       as model_max_tokens,
	//        om.prompt_price     as model_prompt_price,
	//        om.completion_price as model_completion_price,
	//        om.rpm_limit        as model_rpm_limit,
	//        om.tpm_limit        as model_tpm_limit,
	//        om.updated_at       as last_updated_at
	// from openai_clients as oc
	//          join openai_models as om on om.client_id = oc.id and oc.id = ${client_id}
	if queryErr := ac.db.GetGormCore(ctx).
		Model(&model.OpenaiClient{}).
		Joins("?", &clause.Join{
			Type:  clause.InnerJoin,
			Table: clause.Table{Name: model.TableNameOpenaiModels},
			ON: clause.Where{
				Exprs: []clause.Expression{
					clause.Eq{
						Column: database.Column(model.TableNameOpenaiClients, model.OpenaiClientCols.ID),
						Value:  clause.Column{Table: model.TableNameOpenaiModels, Name: model.OpenaiModelCols.ClientID},
					},
					clause.Eq{
						Column: clause.Column{Table: model.TableNameOpenaiClients, Name: model.OpenaiClientCols.ID},
						Value:  clientID,
					},
				},
			},
		}).
		Select(needFields).
		Scan(&result).
		Error; queryErr != nil {
		return nil, errors.Wrap(queryErr, "failed to get models by client id")
	}

	return result, nil
}

func (ac *OpenaiModelDatabaseAccessor) GetModelsByClientDescription(ctx context.Context, description string) (result []*dto.RelatedModelDTO, err error) {
	result = make([]*dto.RelatedModelDTO, 0)
	needFields := []string{
		database.ColumnAlias(model.TableNameOpenaiClients, model.OpenaiClientCols.ID, "client_id"),
		database.ColumnAlias(model.TableNameOpenaiModels, model.OpenaiModelCols.ID, "model_id"),
		database.ColumnAlias(model.TableNameOpenaiModels, model.OpenaiModelCols.Model, "model_name"),
		database.ColumnAlias(model.TableNameOpenaiModels, model.OpenaiModelCols.MaxTokens, "model_max_tokens"),
		database.ColumnAlias(model.TableNameOpenaiModels, model.OpenaiModelCols.PromptPrice, "model_prompt_price"),
		database.ColumnAlias(model.TableNameOpenaiModels, model.OpenaiModelCols.CompletionPrice, "model_completion_price"),
		database.ColumnAlias(model.TableNameOpenaiModels, model.OpenaiModelCols.RpmLimit, "model_rpm_limit"),
		database.ColumnAlias(model.TableNameOpenaiModels, model.OpenaiModelCols.TpmLimit, "model_tpm_limit"),
		database.ColumnAlias(model.TableNameOpenaiModels, model.OpenaiModelCols.UpdatedAt, "last_updated_at"),
	}

	// select oc.id               as client_id,
	//        om.id               as model_id,
	//        om.model            as model_name,
	//        om.max_tokens       as model_max_tokens,
	//        om.prompt_price     as model_prompt_price,
	//        om.completion_price as model_completion_price,
	//        om.rpm_limit        as model_rpm_limit,
	//        om.tpm_limit        as model_tpm_limit,
	//        om.updated_at       as last_updated_at
	// from openai_clients as oc
	//          join openai_models as om on om.client_id = oc.id and oc.description = ${description}
	if queryErr := ac.db.GetGormCore(ctx).
		Model(&model.OpenaiClient{}).
		Joins("?", &clause.Join{
			Type:  clause.InnerJoin,
			Table: clause.Table{Name: model.TableNameOpenaiModels},
			ON: clause.Where{
				Exprs: []clause.Expression{
					clause.Eq{
						Column: database.Column(model.TableNameOpenaiClients, model.OpenaiClientCols.ID),
						Value:  clause.Column{Table: model.TableNameOpenaiModels, Name: model.OpenaiModelCols.ClientID},
					},
					clause.Eq{
						Column: clause.Column{Table: model.TableNameOpenaiClients, Name: model.OpenaiClientCols.Description},
						Value:  description,
					},
				},
			},
		}).
		Select(needFields).
		Scan(&result).
		Error; queryErr != nil {
		return nil, errors.Wrap(queryErr, "failed to get models by client id")
	}

	return result, nil
}

func (ac *OpenaiModelDatabaseAccessor) GetModelIDByName(ctx context.Context, modelName string) (modelID int, err error) {
	result := new(model.OpenaiModel)

	// select id from openai_models where model = ${modelName}
	if queryErr := ac.db.GetDataBySingleCondition(ctx, result, model.OpenaiModelCols.Model, modelName, model.OpenaiClientCols.ID); queryErr != nil {
		return 0, errors.Wrap(queryErr, "failed to get model id by model name")
	}

	return int(result.ID), nil
}

func (ac *OpenaiModelDatabaseAccessor) GetAvailableModelsByApiKey(ctx context.Context, key string) (result []*dto.RelatedModelDTO, err error) {
	result = make([]*dto.RelatedModelDTO, 0)
	needFields := []string{
		database.ColumnAlias(model.TableNameOpenaiModels, model.OpenaiModelCols.ClientID, "client_id"),
		database.ColumnAlias(model.TableNameOpenaiModels, model.OpenaiModelCols.ID, "model_id"),
		database.ColumnAlias(model.TableNameOpenaiModels, model.OpenaiModelCols.Model, "model_name"),
		database.ColumnAlias(model.TableNameOpenaiModels, model.OpenaiModelCols.MaxTokens, "model_max_tokens"),
		database.ColumnAlias(model.TableNameOpenaiModels, model.OpenaiModelCols.PromptPrice, "model_prompt_price"),
		database.ColumnAlias(model.TableNameOpenaiModels, model.OpenaiModelCols.CompletionPrice, "model_completion_price"),
		database.ColumnAlias(model.TableNameOpenaiModels, model.OpenaiModelCols.RpmLimit, "model_rpm_limit"),
		database.ColumnAlias(model.TableNameOpenaiModels, model.OpenaiModelCols.TpmLimit, "model_tpm_limit"),
		database.ColumnAlias(model.TableNameOpenaiModels, model.OpenaiModelCols.UpdatedAt, "last_updated_at"),
	}

	// select om.client_id		  as client_id,
	//        om.id               as model_id,
	//        om.model            as model_name,
	//        om.max_tokens       as model_max_tokens,
	//        om.prompt_price     as model_prompt_price,
	//        om.completion_price as model_completion_price,
	//        om.rpm_limit        as model_rpm_limit,
	//        om.tpm_limit        as model_tpm_limit,
	//        om.updated_at       as last_updated_at
	// from whisper_users as wu
	//          join whisper_user_permissions on wu.id = whisper_user_permissions.user_id and wu.api_key = ${key}
	//          join openai_models as om on whisper_user_permissions.model_id = om.id`
	if queryErr := ac.db.GetGormCore(ctx).
		Model(&model.WhisperUser{}).
		Joins("?", &clause.Join{
			Type:  clause.InnerJoin,
			Table: clause.Table{Name: model.TableNameWhisperUserPermissions},
			ON: clause.Where{
				Exprs: []clause.Expression{
					clause.Eq{
						Column: database.Column(model.TableNameWhisperUsers, model.WhisperUserCols.ID),
						Value:  clause.Column{Table: model.TableNameWhisperUserPermissions, Name: model.WhisperUserPermissionCols.UserID},
					},
					clause.Eq{
						Column: clause.Column{Table: model.TableNameWhisperUsers, Name: model.WhisperUserCols.ApiKey},
						Value:  key,
					},
				},
			},
		}).
		Joins("?", &clause.Join{
			Type:  clause.InnerJoin,
			Table: clause.Table{Name: model.TableNameOpenaiModels},
			ON: clause.Where{
				Exprs: []clause.Expression{
					clause.Eq{
						Column: database.Column(model.TableNameWhisperUserPermissions, model.WhisperUserPermissionCols.ModelID),
						Value:  clause.Column{Table: model.TableNameOpenaiModels, Name: model.OpenaiModelCols.ID},
					},
				},
			},
		}).
		Select(needFields).
		Scan(&result).
		Error; queryErr != nil {
		return nil, errors.Wrap(queryErr, "failed to get available models by api key")
	}

	return result, nil
}

func (ac *OpenaiModelDatabaseAccessor) CreateOrUpdateModel(ctx context.Context, modelData *model.OpenaiModel, clientIDs ...int) (err error) {
	updates := make([]*model.OpenaiModel, 0, len(clientIDs))
	for _, client := range clientIDs {
		updates = append(updates, &model.OpenaiModel{
			ClientID:        int64(client),
			Model:           modelData.Model,
			MaxTokens:       modelData.MaxTokens,
			PromptPrice:     modelData.PromptPrice,
			CompletionPrice: modelData.CompletionPrice,
			RpmLimit:        modelData.RpmLimit,
			TpmLimit:        modelData.TpmLimit,
		})
	}

	indexKeys := []string{model.OpenaiModelCols.ClientID, model.OpenaiModelCols.Model}
	updateKeys := []string{model.OpenaiModelCols.MaxTokens, model.OpenaiModelCols.PromptPrice, model.OpenaiModelCols.CompletionPrice, model.OpenaiModelCols.RpmLimit, model.OpenaiModelCols.TpmLimit}

	return ac.db.CreateDataOnDuplicateKeyUpdate(ctx, updates, indexKeys, updateKeys)
}

func (ac *OpenaiModelDatabaseAccessor) CreateOrUpdateModelWithClientDescriptions(ctx context.Context, modelData *model.OpenaiModel, descriptions ...string) (err error) {
	return ac.db.GetGormCore(ctx).Transaction(func(tx *gorm.DB) error {
		clientIDs := make([]int, 0, len(descriptions))
		queryErr := tx.WithContext(ctx).
			Model(&model.OpenaiClient{}).
			Select(model.OpenaiClientCols.ID).
			Where(model.OpenaiClientCols.Description, descriptions).
			Scan(&clientIDs).Error
		if queryErr != nil {
			return queryErr
		}

		updates := make([]*model.OpenaiModel, 0, len(clientIDs))
		for _, client := range clientIDs {
			updates = append(updates, &model.OpenaiModel{
				ClientID:        int64(client),
				Model:           modelData.Model,
				MaxTokens:       modelData.MaxTokens,
				PromptPrice:     modelData.PromptPrice,
				CompletionPrice: modelData.CompletionPrice,
				RpmLimit:        modelData.RpmLimit,
				TpmLimit:        modelData.TpmLimit,
			})
		}

		indexKeys := []string{model.OpenaiModelCols.ClientID, model.OpenaiModelCols.Model}
		updateKeys := []string{model.OpenaiModelCols.MaxTokens, model.OpenaiModelCols.PromptPrice, model.OpenaiModelCols.CompletionPrice, model.OpenaiModelCols.RpmLimit, model.OpenaiModelCols.TpmLimit}

		duplicatedColumns := make([]clause.Column, len(indexKeys))
		for i, key := range indexKeys {
			duplicatedColumns[i] = clause.Column{Name: key}
		}

		return tx.WithContext(ctx).Model(&model.OpenaiModel{}).Clauses(clause.OnConflict{
			Columns:   duplicatedColumns,
			DoUpdates: clause.AssignmentColumns(updateKeys),
		}).Create(updates).Error
	})
}
