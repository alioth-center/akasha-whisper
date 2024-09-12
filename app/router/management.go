package router

import (
	"github.com/alioth-center/akasha-whisper/app/api"
	"github.com/alioth-center/akasha-whisper/app/entity"
	"github.com/alioth-center/infrastructure/network/http"
)

var managementRouter = http.NewRouter("management")

var ManagementRouterGroup = []http.EndPointInterface{
	http.NewEndPointBuilder[*entity.ListClientsRequest, *entity.ListClientResponse]().
		SetNecessaryHeaders(http.HeaderAuthorization).
		SetHandlerChain(api.ManagementApi.ListClients()).
		SetAllowMethods(http.GET).
		SetRouter(managementRouter.Group("clients")).
		Build(),
	http.NewEndPointBuilder[*entity.CreateClientRequest, *entity.CreateClientResponse]().
		SetNecessaryHeaders(http.HeaderAuthorization).
		SetHandlerChain(api.ManagementApi.CreateClient()).
		SetAllowMethods(http.POST).
		SetRouter(managementRouter.Group("clients")).
		Build(),
	http.NewEndPointBuilder[*entity.ModifyOpenaiClientBalanceRequest, *entity.ModifyOpenaiClientBalanceResponse]().
		SetNecessaryHeaders(http.HeaderAuthorization).
		SetNecessaryParams("client_name").
		SetHandlerChain(api.ManagementApi.ModifyClientBalance()).
		SetAllowMethods(http.POST).
		SetRouter(managementRouter.Group("client/:client_name/balance")).
		Build(),
	http.NewEndPointBuilder[*entity.ListClientModelRequest, *entity.ListClientModelResponse]().
		SetNecessaryHeaders(http.HeaderAuthorization).
		SetNecessaryParams("client_name").
		SetHandlerChain(api.ManagementApi.ListClientModels()).
		SetAllowMethods(http.GET).
		SetRouter(managementRouter.Group("client/:client_name/models")).
		Build(),
	http.NewEndPointBuilder[*entity.CreateClientModelRequest, *entity.CreateClientModelResponse]().
		SetNecessaryHeaders(http.HeaderAuthorization).
		SetNecessaryParams("client_name").
		SetHandlerChain(api.ManagementApi.CreateClientModels()).
		SetAllowMethods(http.POST).
		SetRouter(managementRouter.Group("client/:client_name/models")).
		Build(),
	http.NewEndPointBuilder[*entity.ListWhisperUsersRequest, *entity.ListWhisperUsersResponse]().
		SetNecessaryHeaders(http.HeaderAuthorization).
		SetHandlerChain(api.ManagementApi.ListWhisperUsers()).
		SetAllowMethods(http.GET).
		SetRouter(managementRouter.Group("users")).
		Build(),
	http.NewEndPointBuilder[*entity.CreateWhisperUserRequest, *entity.CreateWhisperUserResponse]().
		SetNecessaryHeaders(http.HeaderAuthorization).
		SetHandlerChain(api.ManagementApi.CreateWhisperUser()).
		SetAllowMethods(http.POST).
		SetRouter(managementRouter.Group("users")).
		Build(),
	http.NewEndPointBuilder[*entity.BatchModifyWhisperUserBalanceRequest, *entity.BatchModifyWhisperUserBalanceResponse]().
		SetNecessaryHeaders(http.HeaderAuthorization).
		SetHandlerChain(api.ManagementApi.BatchModifyWhisperUserBalance()).
		SetAllowMethods(http.POST).
		SetRouter(managementRouter.Group("users/balance")).
		Build(),
	http.NewEndPointBuilder[*entity.GetWhisperUserRequest, *entity.GetWhisperUserResponse]().
		SetNecessaryHeaders(http.HeaderAuthorization).
		SetNecessaryParams("user_id").
		SetHandlerChain(api.ManagementApi.GetWhisperUser()).
		SetAllowMethods(http.GET).
		SetRouter(managementRouter.Group("user/:user_id")).
		Build(),
	http.NewEndPointBuilder[*entity.UpdateWhisperUserRequest, *entity.UpdateWhisperUserResponse]().
		SetNecessaryHeaders(http.HeaderAuthorization).
		SetNecessaryParams("user_id").
		SetHandlerChain(api.ManagementApi.UpdateWhisperUser()).
		SetAllowMethods(http.PUT).
		SetRouter(managementRouter.Group("user/:user_id")).
		Build(),
	http.NewEndPointBuilder[*entity.ListWhisperUserBalanceLogsRequest, *entity.ListWhisperUserBalanceLogsResponse]().
		SetNecessaryHeaders(http.HeaderAuthorization).
		SetNecessaryParams("user_id").
		SetAdditionalQueries("page", "offset", "start", "end").
		SetHandlerChain(api.ManagementApi.ListWhisperUserBalanceLogs()).
		SetAllowMethods(http.GET).
		SetRouter(managementRouter.Group("user/:user_id/balance_logs")).
		Build(),
	http.NewEndPointBuilder[*entity.ModifyWhisperUserBalanceRequest, *entity.ModifyWhisperUserBalanceResponse]().
		SetNecessaryHeaders(http.HeaderAuthorization).
		SetNecessaryParams("user_id").
		SetHandlerChain(api.ManagementApi.ModifyWhisperUserBalance()).
		SetAllowMethods(http.POST).
		SetRouter(managementRouter.Group("user/:user_id/balance")).
		Build(),
	http.NewEndPointBuilder[*entity.ModifyWhisperUserPermissionRequest, *entity.ModifyWhisperUserPermissionResponse]().
		SetNecessaryHeaders(http.HeaderAuthorization).
		SetHandlerChain(api.ManagementApi.ModifyWhisperUserPermissions()).
		SetNecessaryParams("user_id").
		SetAllowMethods(http.PUT).
		SetRouter(managementRouter.Group("user/:user_id/permissions")).
		Build(),
}
