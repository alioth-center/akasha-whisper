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

type WhisperUserDatabaseAccessor struct {
	db database.DatabaseV2
}

func NewWhisperUserDatabaseAccessor(db database.DatabaseV2) *WhisperUserDatabaseAccessor {
	return &WhisperUserDatabaseAccessor{db: db}
}

func (ac *WhisperUserDatabaseAccessor) CheckWhisperUserApiKey(ctx context.Context, apiKey string) (exist bool, allowIPs string, err error) {
	user := new(model.WhisperUser)
	queryErr := ac.db.GetGormCore(ctx).
		Model(&model.WhisperUser{}).
		Where(model.WhisperUserCols.ApiKey, apiKey).
		Select(model.WhisperUserCols.AllowIps).
		First(user).
		Error
	if queryErr == nil {
		return true, user.AllowIps, nil
	}
	if errors.Is(queryErr, gorm.ErrRecordNotFound) {
		return false, "", nil
	}

	return false, "", queryErr
}

func (ac *WhisperUserDatabaseAccessor) ListWhisperUsers(ctx context.Context, page, limit int) ([]model.WhisperUser, error) {
	users := make([]model.WhisperUser, 0)
	if queryErr := ac.db.GetGormCore(ctx).
		Model(&model.WhisperUser{}).
		Select(model.WhisperUserCols.ID, model.WhisperUserCols.ApiKey, model.WhisperUserCols.Email, model.WhisperUserCols.Language, model.WhisperUserCols.AllowIps).
		Offset(page * limit).
		Limit(limit).
		Scan(&users).
		Error; queryErr != nil {
		return nil, queryErr
	}

	return users, nil
}

func (ac *WhisperUserDatabaseAccessor) ListWhisperUserApiKeys(ctx context.Context) ([]string, error) {
	result := make([]string, 0)
	if queryErr := ac.db.GetGormCore(ctx).
		Model(&model.WhisperUser{}).
		Select(model.WhisperUserCols.ApiKey).
		Scan(&result).
		Error; queryErr != nil {
		return nil, queryErr
	}

	return result, nil
}

func (ac *WhisperUserDatabaseAccessor) CreateWhisperUser(ctx context.Context, user *model.WhisperUser) (created bool, err error) {
	return ac.db.CreateSingleDataIfNotExist(ctx, user)
}

func (ac *WhisperUserDatabaseAccessor) GetWhisperUserInfo(ctx context.Context, userID int) (user *dto.WhisperUserInfo, err error) {
	user = new(dto.WhisperUserInfo)
	return user, ac.db.GetGormCore(ctx).Transaction(func(tx *gorm.DB) error {
		// select om.model as model_name
		// from whisper_users as wu
		//          join whisper_user_permissions as wup on wu.id = wup.user_id
		//          join openai_models as om on wup.model_id = om.id
		// group by om.model
		var models []string
		if queryModelsErr := tx.WithContext(ctx).
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
							Column: database.Column(model.TableNameWhisperUsers, model.WhisperUserCols.ID),
							Value:  userID,
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
			Select(database.Column(model.TableNameOpenaiModels, model.OpenaiModelCols.Model)).
			Group(database.Column(model.TableNameOpenaiModels, model.OpenaiModelCols.Model)).
			Scan(&models).
			Error; queryModelsErr != nil {
			return queryModelsErr
		}

		if queryUserErr := tx.WithContext(ctx).
			Raw(rawSqlList[RawsqlWhisperUserGetUserInfo], userID).
			Scan(&user.UserInfo).
			Error; queryUserErr != nil {
			return queryUserErr
		}

		user.Models = models
		return nil
	})
}

func (ac *WhisperUserDatabaseAccessor) UpdateWhisperUser(ctx context.Context, user *model.WhisperUser) error {
	return ac.db.UpdateDataBySingleCondition(ctx, user, model.WhisperUserCols.ID, user.ID)
}
