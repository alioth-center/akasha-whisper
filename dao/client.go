package dao

import (
	"context"
	"fmt"
	"github.com/alioth-center/akasha-whisper/model"
	"github.com/alioth-center/infrastructure/utils/values"
	"gorm.io/gorm"
)

func (ac *postgresDatabaseAccessor) CheckAllClientsExist(ctx context.Context, clients []string) (notExist []string, err error) {
	var dto []model.CheckDTO
	checkErr := ac.db.QueryGormFunctionWithCtx(ctx, &dto, func(db *gorm.DB) *gorm.DB {
		return db.Model(&model.OpenaiClient{}).Select("id, description as name").Where("description in ?", clients).Scan(&dto)
	})
	if checkErr != nil {
		return nil, fmt.Errorf("failed to check all clients exist: %w", checkErr)
	}

	// mapping the dto to a map
	mapping := map[string]struct{}{}
	for _, checkDTO := range dto {
		mapping[checkDTO.Name] = struct{}{}
	}

	// check the clients not exist
	for _, client := range clients {
		if _, ok := mapping[client]; !ok {
			notExist = append(notExist, client)
		}
	}

	return notExist, nil
}

func (ac *postgresDatabaseAccessor) GetAvailableClient(ctx context.Context, mdl, tk string) (results []model.AvailableClientDTO, err error) {
	sql := values.NewRawSqlTemplateWithMap(`select oc.id               as client_id,
       oc.balance          as client_balance,
       oc.weight           as client_weight,
       wu.id               as user_id,
       wu.balance          as user_balance,
       wu.allow_ips        as user_allow_ips,
       wu.role             as user_role,
       om.model            as model_name,
       om.id               as model_id,
       om.max_tokens       as model_max_tokens,
       om.prompt_price     as model_prompt_price,
       om.completion_price as model_completion_price
from whisper_users as wu
         join whisper_user_permissions as wup on wu.id = wup.user_id
         join openai_models as om on wup.model_id = om.id and om.model = '${model}'
         join openai_clients as oc on om.client_id = oc.id and oc.balance > 0
where wu.api_key = '${api_key}'`, map[string]string{"model": mdl, "api_key": tk}).Parse()

	err = ac.db.QueryRawWithCtx(ctx, &results, sql)
	if err != nil {
		return nil, fmt.Errorf("failed to get available client: %w", err)
	} else if len(results) == 0 {
		return nil, fmt.Errorf("no available client")
	}

	return results, nil
}

func (ac *postgresDatabaseAccessor) GetAvailableUserClient(ctx context.Context, mdl, email string) (results []model.AvailableClientDTO, err error) {
	sql := values.NewRawSqlTemplateWithMap(`select oc.id               as client_id,
       oc.balance          as client_balance,
       oc.weight           as client_weight,
       wu.id               as user_id,
       wu.balance          as user_balance,
       wu.allow_ips        as user_allow_ips,
       wu.role             as user_role,
       om.model            as model_name,
       om.id               as model_id,
       om.max_tokens       as model_max_tokens,
       om.prompt_price     as model_prompt_price,
       om.completion_price as model_completion_price
from whisper_users as wu
         join whisper_user_permissions as wup on wu.id = wup.user_id
         join openai_models as om on wup.model_id = om.id and om.model = '${model}'
         join openai_clients as oc on om.client_id = oc.id and oc.balance > 0
where wu.email = '${email}'`, map[string]string{"model": mdl, "email": email}).Parse()

	err = ac.db.QueryRawWithCtx(ctx, &results, sql)
	if err != nil {
		return nil, fmt.Errorf("failed to get available client: %w", err)
	} else if len(results) == 0 {
		return nil, fmt.Errorf("no available client")
	}

	return results, nil
}

func (ac *postgresDatabaseAccessor) GetClientSecret(ctx context.Context, clientID int) (result model.ClientSecretDTO, err error) {
	sql := `select id       as client_id,
       api_key  as client_key,
       endpoint as client_endpoint,
       weight   as client_weight,
       balance  as client_balance
from openai_clients
where id = ?`

	err = ac.db.QueryRawWithCtx(ctx, &result, sql, clientID)
	if err != nil {
		return model.ClientSecretDTO{}, fmt.Errorf("failed to get client secret: %w", err)
	} else if result.ClientID == 0 {
		return model.ClientSecretDTO{}, fmt.Errorf("client id %d not exist", clientID)
	}

	return result, nil
}
