package dao

import (
	"context"

	"github.com/alioth-center/akasha-whisper/app/model"
	"github.com/alioth-center/akasha-whisper/app/model/dto"
	"github.com/alioth-center/infrastructure/database"
	"github.com/pkg/errors"
)

type OpenaiClientDatabaseAccessor struct {
	db database.DatabaseV2
}

func NewOpenaiClientDatabaseAccessor(db database.DatabaseV2) *OpenaiClientDatabaseAccessor {
	return &OpenaiClientDatabaseAccessor{db: db}
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
	sql := rawSqlList[RawsqlOpenaiClientGetAvailableClients]
	template := &dto.GetAvailableClientCND{ModelName: modelName, UserApiKey: token}
	if queryErr := ac.db.ExecuteRawSqlTemplateQuery(ctx, &clients, sql, template); queryErr != nil {
		return nil, errors.Wrap(queryErr, "get available clients failed")
	}

	return clients, nil
}

func (ac *OpenaiClientDatabaseAccessor) GetClientSecret(ctx context.Context, clientID int) (result *dto.ClientSecretDTO, err error) {
	result = new(dto.ClientSecretDTO)
	sql := rawSqlList[RawsqlOpenaiClientGetClientSecrets]
	template := &dto.GetClientSecretCND{ClientIDs: []int{clientID}}
	if queryErr := ac.db.ExecuteRawSqlTemplateQuery(ctx, &result, sql, template); queryErr != nil {
		return nil, errors.Wrap(queryErr, "get client secret failed")
	}

	return result, nil
}

func (ac *OpenaiClientDatabaseAccessor) GetClientSecrets(ctx context.Context, clientIDs ...int) (result []*dto.ClientSecretDTO, err error) {
	if len(clientIDs) == 0 {
		return []*dto.ClientSecretDTO{}, nil
	}

	result = make([]*dto.ClientSecretDTO, 0, len(clientIDs))
	sql := rawSqlList[RawsqlOpenaiClientGetClientSecrets]
	template := &dto.GetClientSecretCND{ClientIDs: clientIDs}
	if queryErr := ac.db.ExecuteRawSqlTemplateQuery(ctx, &result, sql, template); queryErr != nil {
		return nil, errors.Wrap(queryErr, "get client secrets failed")
	}

	return result, nil
}

func (ac *OpenaiClientDatabaseAccessor) ListClients(ctx context.Context) (clients []*dto.ListClientDTO, err error) {
	result := make([]*dto.ListClientDTO, 0)
	sql := rawSqlList[RawsqlOpenaiClientListClients]
	if queryErr := ac.db.GetGormCore(ctx).Raw(sql).Scan(&result).Error; queryErr != nil {
		return nil, errors.Wrap(queryErr, "list clients failed")
	}

	return result, nil
}

func (ac *OpenaiClientDatabaseAccessor) CreateClient(ctx context.Context, client *model.OpenaiClient) (created bool, err error) {
	return ac.db.CreateSingleDataIfNotExist(ctx, client)
}
