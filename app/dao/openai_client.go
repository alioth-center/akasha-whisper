package dao

import (
	"context"
	"github.com/alioth-center/akasha-whisper/app/model"
	"github.com/alioth-center/akasha-whisper/app/model/dto"
	"github.com/alioth-center/infrastructure/database"
	"github.com/pkg/errors"
	"gorm.io/gorm/clause"
)

type OpenaiClientDatabaseAccessor struct {
	db database.DatabaseV2
}

func (ac *OpenaiClientDatabaseAccessor) CheckAllClientsExist(ctx context.Context, clients []string) (notExistClientDescriptions []string, err error) {
	if len(clients) == 0 {
		return []string{}, nil
	}

	clientDTOs := make([]dto.ClientCheckDTO, 0, len(clients))
	needFields := []string{model.OpenaiClientCols.ID, model.OpenaiClientCols.Description}

	// select id, description from openai_client where description in ${clients}
	if queryErr := ac.db.GetDataBySingleCondition(ctx, &clientDTOs, model.OpenaiClientCols.Description, clients, needFields...); queryErr != nil {
		return nil, errors.Wrap(queryErr, "check all clients exist failed")
	}

	// mapping dto(s) to a map
	mapping := map[string]struct{}{}
	for _, client := range clientDTOs {
		mapping[client.Name] = struct{}{}
	}

	// check if all clients exist
	for _, client := range clients {
		if _, exist := mapping[client]; !exist {
			notExistClientDescriptions = append(notExistClientDescriptions, client)
		}
	}

	// return not exist clients
	return notExistClientDescriptions, nil
}

func (ac *OpenaiClientDatabaseAccessor) GetAvailableClients(ctx context.Context, modelName, token string) (clients []*dto.AvailableClientDTO, err error) {
	clients = make([]*dto.AvailableClientDTO, 0)
	needFields := []string{
		database.ColumnAlias(model.TableNameOpenaiClients, model.OpenaiClientCols.ID, "client_id"),
		database.ColumnAlias(model.TableNameOpenaiClients, model.OpenaiClientCols.Balance, "client_balance"),
		database.ColumnAlias(model.TableNameOpenaiClients, model.OpenaiClientCols.Weight, "client_weight"),
		database.ColumnAlias(model.TableNameOpenaiClients, model.OpenaiClientCols.Description, "client_description"),
		database.ColumnAlias(model.TableNameWhisperUsers, model.WhisperUserCols.ID, "user_id"),
		database.ColumnAlias(model.TableNameWhisperUsers, model.WhisperUserCols.Balance, "user_balance"),
		database.ColumnAlias(model.TableNameWhisperUsers, model.WhisperUserCols.Role, "user_role"),
		database.ColumnAlias(model.TableNameOpenaiModels, model.OpenaiModelCols.Model, "model_name"),
		database.ColumnAlias(model.TableNameOpenaiModels, model.OpenaiModelCols.ID, "model_id"),
		database.ColumnAlias(model.TableNameOpenaiModels, model.OpenaiModelCols.MaxTokens, "model_max_tokens"),
		database.ColumnAlias(model.TableNameOpenaiModels, model.OpenaiModelCols.PromptPrice, "model_prompt_price"),
		database.ColumnAlias(model.TableNameOpenaiModels, model.OpenaiModelCols.CompletionPrice, "model_completion_price"),
	}

	// select oc.id               as client_id,
	//        oc.balance          as client_balance,
	//        oc.weight           as client_weight,
	//        oc.description      as client_description,
	//        wu.id               as user_id,
	//        wu.balance          as user_balance,
	//        wu.role             as user_role,
	//        om.model            as model_name,
	//        om.id               as model_id,
	//        om.max_tokens       as model_max_tokens,
	//        om.prompt_price     as model_prompt_price,
	//        om.completion_price as model_completion_price
	// from whisper_users as wu
	//          join whisper_user_permissions as wup on wu.id = wup.user_id and wu.api_key = ${api_key}
	//          join openai_models as om on wup.model_id = om.id and om.model = ${model}
	//          join openai_clients as oc on om.client_id = oc.id and oc.balance > 0
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
						Column: database.Column(model.TableNameWhisperUsers, model.WhisperUserCols.ApiKey),
						Value:  token,
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
					clause.Eq{
						Column: database.Column(model.TableNameOpenaiModels, model.OpenaiModelCols.Model),
						Value:  modelName,
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
						Column: database.Column(model.TableNameOpenaiModels, model.OpenaiModelCols.ClientID),
						Value:  clause.Column{Table: model.TableNameOpenaiClients, Name: model.OpenaiClientCols.ID},
					},
					clause.Gt{
						Column: database.Column(model.TableNameOpenaiClients, model.OpenaiClientCols.Balance),
						Value:  0,
					},
				},
			},
		}).
		Select(needFields).
		Scan(&clients).
		Error; queryErr != nil {
		return nil, errors.Wrap(queryErr, "get available clients failed")
	}

	return clients, nil
}

func (ac *OpenaiClientDatabaseAccessor) GetClientSecret(ctx context.Context, clientID int) (result *dto.ClientSecretDTO, err error) {
	result = new(dto.ClientSecretDTO)
	needFields := []string{
		model.OpenaiClientCols.ID, model.OpenaiClientCols.ApiKey, model.OpenaiClientCols.Endpoint,
		model.OpenaiClientCols.Weight, model.OpenaiClientCols.Balance,
	}

	// select id, api_key, endpoint, weight, balance from openai_client where id = ${client_id}
	if queryErr := ac.db.GetDataBySingleCondition(ctx, result, model.OpenaiClientCols.ID, clientID, needFields...); queryErr != nil {
		return nil, errors.Wrap(queryErr, "get client secret failed")
	}

	return result, nil
}
