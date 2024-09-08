package router

import (
	"github.com/alioth-center/akasha-whisper/app/api"
	"github.com/alioth-center/akasha-whisper/app/entity"
	"github.com/alioth-center/infrastructure/network/http"
)

var managementRouter = http.NewRouter("management")

var ManagementRouterGroup = []http.EndPointInterface{
	//http.NewEndPointBuilder[*model.CreateUserRequest, *model.BaseResponse]().
	//	SetNecessaryHeaders("Authorization").
	//	SetHandlerChain(http.NewChain(service.ManagerService.CreateUser)).
	//	SetAllowMethods(http.POST).
	//	SetRouter(managementRouter.Group("/users")).
	//	Build(),
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
}
