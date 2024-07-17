package dao

import (
	"context"
	"fmt"
	"github.com/alioth-center/akasha-whisper/model"
	"github.com/alioth-center/infrastructure/database/extension/orm"
	"gorm.io/gorm"
	"strings"
	"time"
)

type postgresDatabaseAccessor struct {
	db orm.Extended
}

func (ac *postgresDatabaseAccessor) GetUserByID(ctx context.Context, id int64) (result model.QueryWhisperUserDTO, err error) {
	err = ac.db.QueryGormFunctionWithCtx(ctx, &result, func(db *gorm.DB) *gorm.DB {
		return db.Model(&model.WhisperUser{}).Where("id =?", id).Limit(1).Scan(&result)
	})
	if err != nil {
		return model.QueryWhisperUserDTO{}, fmt.Errorf("failed to get user by id: %w", err)
	}

	return result, nil
}

func (ac *postgresDatabaseAccessor) GetUserByApiKey(ctx context.Context, key string) (result model.QueryWhisperUserDTO, err error) {
	err = ac.db.QueryGormFunctionWithCtx(ctx, &result, func(db *gorm.DB) *gorm.DB {
		return db.Model(&model.WhisperUser{}).Where("api_key = ?", key).Limit(1).Scan(&result)
	})
	if err != nil {
		return model.QueryWhisperUserDTO{}, fmt.Errorf("failed to get user by api key: %w", err)
	}

	return result, nil
}

func (ac *postgresDatabaseAccessor) AddClient(ctx context.Context, client model.ClientDTO, models ...model.ClientModelDTO) (err error) {
	err = ac.db.ExecuteGormTransactionWithCtx(ctx, func(tx *gorm.DB) error {
		// add client
		openaiClient := model.OpenaiClient{
			Description: client.Description,
			ApiKey:      client.ApiKey,
			Endpoint:    client.Endpoint,
			Weight:      client.Weight,
			Balance:     client.Balance,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		if err := tx.Model(&model.OpenaiClient{}).Create(&openaiClient).Error; err != nil {
			return err
		}

		// add client models
		var buffer []model.OpenaiModel
		for _, openaiModel := range models {
			buffer = append(buffer, model.OpenaiModel{
				ClientID:        openaiClient.ID,
				Model:           openaiModel.Name,
				MaxTokens:       openaiModel.MaxTokens,
				PromptPrice:     openaiModel.PromptPrice,
				CompletionPrice: openaiModel.CompletionPrice,
				RpmLimit:        openaiModel.RpmLimit,
				TpmLimit:        openaiModel.TpmLimit,
				CreatedAt:       time.Now(),
				UpdatedAt:       time.Now(),
			})
		}
		if len(buffer) == 0 {
			return nil
		}

		return tx.Model(&model.OpenaiModel{}).CreateInBatches(buffer, 100).Error
	})

	if err != nil {
		return fmt.Errorf("failed to add client: %w", err)
	}

	return nil
}

func (ac *postgresDatabaseAccessor) AddClientModels(ctx context.Context, clientID int, models []model.ClientModelDTO) (err error) {
	err = ac.db.ExecuteGormTransactionWithCtx(ctx, func(tx *gorm.DB) error {
		// check the given client_id exist
		var count int64
		if err := tx.Model(&model.OpenaiClient{}).Where("id = ?", clientID).Count(&count).Error; err != nil {
			return err
		} else if count == 0 {
			return fmt.Errorf("client id %d not exist", clientID)
		}

		// add client models
		var buffer []model.OpenaiModel
		for _, openaiModel := range models {
			buffer = append(buffer, model.OpenaiModel{
				ClientID:        int64(clientID),
				Model:           openaiModel.Name,
				MaxTokens:       openaiModel.MaxTokens,
				PromptPrice:     openaiModel.PromptPrice,
				CompletionPrice: openaiModel.CompletionPrice,
				RpmLimit:        openaiModel.RpmLimit,
				TpmLimit:        openaiModel.TpmLimit,
				CreatedAt:       time.Now(),
				UpdatedAt:       time.Now(),
			})
		}
		return tx.Model(&model.OpenaiModel{}).CreateInBatches(buffer, 100).Error
	})

	if err != nil {
		return fmt.Errorf("failed to add client models: %w", err)
	}

	return nil
}

func (ac *postgresDatabaseAccessor) AddRequestRecord(ctx context.Context, request model.RequestRecordDTO) (err error) {
	err = ac.db.ExecuteGormTransactionWithCtx(ctx, func(tx *gorm.DB) error {
		// check the given client_id exist
		var count int64
		if err := tx.Model(&model.OpenaiClient{}).Where("id = ?", request.ClientID).Count(&count).Error; err != nil {
			return err
		} else if count == 0 {
			return fmt.Errorf("client id %d does not exist", request.ClientID)
		}

		// record the request
		recordErr := tx.Model(&model.OpenaiRequest{}).Create(&model.OpenaiRequest{
			ClientID:             int64(request.ClientID),
			ModelID:              int64(request.ModelID),
			UserID:               int64(request.UserID),
			RequestIP:            request.RequestIP,
			PromptTokenUsage:     request.PromptTokenUsage,
			CompletionTokenUsage: request.CompletionTokenUsage,
			BalanceCost:          request.BalanceCost,
			CreatedAt:            time.Now(),
		}).Error
		if recordErr != nil {
			return recordErr
		}

		// update the client balance
		updateClientBalanceResult := tx.Exec("UPDATE openai_clients SET balance = balance - ? WHERE id = ?", request.BalanceCost, request.ClientID)
		if updateClientBalanceResult.Error != nil {
			return updateClientBalanceResult.Error
		} else if updateClientBalanceResult.RowsAffected == 0 {
			return fmt.Errorf("client balance update failed, no rows affected for client id %d", request.ClientID)
		}

		// update the user balance
		updateUserBalanceResult := tx.Exec("UPDATE whisper_users SET balance = balance - ? WHERE id = ?", request.BalanceCost, request.UserID)
		if updateUserBalanceResult.Error != nil {
			return updateUserBalanceResult.Error
		} else if updateUserBalanceResult.RowsAffected == 0 {
			return fmt.Errorf("user balance update failed, no rows affected for user id %d", request.UserID)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to add request record: %w", err)
	}

	return nil
}

func (ac *postgresDatabaseAccessor) AddWhisperUser(ctx context.Context, user model.WhisperUserDTO, permissions ...model.UserPermissionItem) (err error) {
	err = ac.db.ExecuteGormTransactionWithCtx(ctx, func(tx *gorm.DB) error {
		// add whisper user
		whisperUser := model.WhisperUser{
			Email:     user.Email,
			ApiKey:    user.ApiKey,
			Balance:   user.Balance,
			Role:      user.Role,
			AllowIps:  strings.Join(user.AllowIPs, ","),
			CreatedAt: time.Now(),
		}
		if err := tx.Model(&model.WhisperUser{}).Create(&whisperUser).Error; err != nil {
			return err
		}

		// no given permissions, return
		if len(permissions) == 0 {
			return nil
		}

		// check the given permissions exist
		type tempBuffer struct {
			ID          int    `gorm:"column:id"`
			Description string `gorm:"column:description"`
		}
		var (
			descriptions []string
			clientBuffer []tempBuffer
			modelIDs     = map[string][]int64{}
			idsMap       = map[string]int{}
		)
		for _, permission := range permissions {
			descriptions = append(descriptions, permission.Desc)
		}

		if queryClientErr := tx.Model(&model.OpenaiClient{}).
			Where("description in ?", descriptions).Scan(&clientBuffer).Error; queryClientErr != nil {
			return queryClientErr
		} else if len(clientBuffer) != len(permissions) {
			return fmt.Errorf("not all permissions exist")
		}
		for _, client := range clientBuffer {
			idsMap[client.Description] = client.ID
		}

		fmt.Printf("idsMap: %v\n", idsMap)

		for _, permission := range permissions {
			// check the given models exist
			modelIDs[permission.Desc] = []int64{}
			var resultBuffer []int64
			if queryModelErr := tx.Model(&model.OpenaiModel{}).Select("id").
				Where("model in ? and client_id = ?", permission.Models, idsMap[permission.Desc]).Scan(&resultBuffer).Error; queryModelErr != nil {
				return queryModelErr
			} else if len(resultBuffer) != len(permission.Models) {
				return fmt.Errorf("not all models exist")
			}

			modelIDs[permission.Desc] = resultBuffer
		}

		// add whisper user permissions
		var buffer []model.WhisperUserPermission
		for _, permission := range permissions {
			for j := range permission.Models {
				buffer = append(buffer, model.WhisperUserPermission{
					UserID:  whisperUser.ID,
					ModelID: modelIDs[permission.Desc][j],
				})
			}
		}

		return tx.Model(&model.WhisperUserPermission{}).CreateInBatches(buffer, 100).Error
	})
	if err != nil {
		return fmt.Errorf("failed to add whisper user: %w", err)
	}

	return nil
}
