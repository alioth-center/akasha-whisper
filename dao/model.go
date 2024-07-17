package dao

import (
	"context"
	"fmt"
	"github.com/alioth-center/akasha-whisper/model"
	"gorm.io/gorm"
)

func (ac *postgresDatabaseAccessor) GetModelsByClientID(ctx context.Context, clientID int) (result []model.RelatedModelDTO, err error) {
	sql := `select oc.id               as client_id,
       om.id               as model_id,
       om.model            as model_name,
       om.max_tokens       as model_max_tokens,
       om.prompt_price     as model_prompt_price,
       om.completion_price as model_completion_price,
       om.rpm_limit        as model_rpm_limit,
       om.tpm_limit        as model_tpm_limit,
       om.updated_at       as last_updated_at
from openai_clients as oc
         join openai_models as om on om.client_id = oc.id and oc.id = ?`

	err = ac.db.QueryRawWithCtx(ctx, &result, sql, clientID)
	if err != nil {
		return nil, fmt.Errorf("failed to get models by client: %w", err)
	} else if len(result) == 0 {
		return nil, fmt.Errorf("client id %d not exist", clientID)
	}

	return result, nil
}

func (ac *postgresDatabaseAccessor) GetModelsByClientDesc(ctx context.Context, clientName string) (result []model.RelatedModelDTO, err error) {
	sql := `select oc.id               as client_id,
       om.id               as model_id,
       om.model            as model_name,
       om.max_tokens       as model_max_tokens,
       om.prompt_price     as model_prompt_price,
       om.completion_price as model_completion_price,
       om.rpm_limit        as model_rpm_limit,
       om.tpm_limit        as model_tpm_limit,
       om.updated_at       as last_updated_at,
from openai_clients as oc
         join openai_models as om on om.client_id = oc.id and oc.description = ?`

	err = ac.db.QueryRawWithCtx(ctx, &result, sql, clientName)
	if err != nil {
		return nil, fmt.Errorf("failed to get models by client: %w", err)
	} else if len(result) == 0 {
		return nil, fmt.Errorf("client id %d not exist", clientName)
	}

	return result, nil
}

func (ac *postgresDatabaseAccessor) GetModelIDByName(ctx context.Context, clientID int, modelName string) (id int, err error) {
	err = ac.db.QueryGormFunctionWithCtx(ctx, &id, func(db *gorm.DB) *gorm.DB {
		return db.Model(&model.OpenaiModel{}).Select("id").Where("client_id = ? and model = ?", clientID, modelName)
	})
	if err != nil {
		return 0, fmt.Errorf("failed to get model id by name: %w", err)
	}

	return id, nil
}

func (ac *postgresDatabaseAccessor) GetAvailableModelsByApiKey(ctx context.Context, key string) (results []model.RelatedModelDTO, err error) {
	err = ac.db.QueryGormFunctionWithCtx(ctx, &results, func(db *gorm.DB) *gorm.DB {
		return db.Raw(`select om.id               as model_id,
       om.model            as model_name,
       om.max_tokens       as model_max_tokens,
       om.prompt_price     as model_prompt_price,
       om.completion_price as model_completion_price,
       om.rpm_limit        as model_rpm_limit,
       om.tpm_limit        as model_tpm_limit,
       om.updated_at       as last_updated_at
from whisper_users as wu
         join whisper_user_permissions on wu.id = whisper_user_permissions.user_id and wu.api_key = ?
         join openai_models as om on whisper_user_permissions.model_id = om.id`, key)
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get available models by api key: %w", err)
	}

	return results, nil
}

func (ac *postgresDatabaseAccessor) GetAvailableModelsByEmail(ctx context.Context, email string) (results []model.RelatedModelDTO, err error) {
	err = ac.db.QueryGormFunctionWithCtx(ctx, &results, func(db *gorm.DB) *gorm.DB {
		return db.Raw(`select om.id               as model_id,
       om.model            as model_name,
       om.max_tokens       as model_max_tokens,
       om.prompt_price     as model_prompt_price,
       om.completion_price as model_completion_price,
       om.rpm_limit        as model_rpm_limit,
       om.tpm_limit        as model_tpm_limit,
       om.updated_at       as last_updated_at
from whisper_users as wu
         join whisper_user_permissions on wu.id = whisper_user_permissions.user_id and wu.email = ?
         join openai_models as om on whisper_user_permissions.model_id = om.id`, email)
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get available models by user email: %w", err)
	}

	return results, nil
}
