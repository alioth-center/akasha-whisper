package manager

import (
	"github.com/alioth-center/akasha-whisper/global"
	"github.com/alioth-center/akasha-whisper/model"
	"github.com/alioth-center/akasha-whisper/service/common"
	"github.com/alioth-center/infrastructure/network/http"
	"github.com/alioth-center/infrastructure/trace"
	"strings"
)

func (srv *akashaManagerSrvImpl) CreateUser(ctx http.Context[*model.CreateUserRequest, *model.BaseResponse]) {
	token := common.GetApiKey(ctx.HeaderParams())

	// check api token permissions, only admin token can create system role user
	var (
		userEntity  model.WhisperUserDTO
		permissions []model.UserPermissionItem

		apiKey = srv.generateUserApiKey()
	)
	for _, permission := range ctx.Request().Permissions {
		permissions = append(permissions, model.UserPermissionItem{Desc: permission.Provider, Models: permission.Models})
	}

	if token != global.Config.AdminToken {
		// not admin token, create user with user role, only system role user can create user role user
		userItem, checkErr := srv.db.GetUserByApiKey(ctx, token)
		if checkErr != nil {
			ctx.SetStatusCode(http.StatusUnauthorized)
			ctx.SetResponse(&model.BaseResponse{
				ErrorCode:    http.ErrorCodeResourceNotFound,
				ErrorMessage: checkErr.Error(),
				RequestID:    trace.GetTid(ctx),
			})
			return
		} else if userItem.Role != "system" {
			ctx.SetStatusCode(http.StatusUnauthorized)
			ctx.SetResponse(&model.BaseResponse{
				ErrorCode:    http.ErrorCodeResourceNotFound,
				ErrorMessage: "permission denied",
				RequestID:    trace.GetTid(ctx),
			})
			return
		}

		userEntity = model.WhisperUserDTO{
			Email:    ctx.Request().Email,
			AllowIPs: ctx.Request().AllowIPs,
			Balance:  global.Config.DefaultBalance,
			ApiKey:   apiKey,
			Role:     "user",
		}
	} else {
		// admin token, can create system role user
		userEntity = model.WhisperUserDTO{
			Email:    ctx.Request().Email,
			AllowIPs: ctx.Request().AllowIPs,
			Role:     ctx.Request().Role,
			Balance:  global.Config.DefaultBalance,
			ApiKey:   apiKey,
		}
	}

	// create user
	createErr := srv.db.AddWhisperUser(ctx, userEntity, permissions...)
	if createErr != nil {
		ctx.SetStatusCode(http.StatusInternalServerError)
		ctx.SetResponse(&model.BaseResponse{
			ErrorCode:    http.ErrorCodeInternalErrorOccurred,
			ErrorMessage: createErr.Error(),
			RequestID:    trace.GetTid(ctx),
		})
		return
	}

	ctx.SetStatusCode(http.StatusOK)
	response := &model.CreateUserResponse{
		Role:   userEntity.Role,
		Email:  userEntity.Email,
		ApiKey: userEntity.ApiKey,
	}
	ctx.SetResponse(&model.BaseResponse{Data: response})
}

func (srv *akashaManagerSrvImpl) generateUserApiKey() string {
	return "aw_" + strings.ReplaceAll(trace.GetTid(trace.NewContext()), "-", "")
}
