package router

import (
	"github.com/alioth-center/akasha-whisper/model"
	"github.com/alioth-center/akasha-whisper/service"
	"github.com/alioth-center/infrastructure/network/http"
)

var managementRouter = http.NewRouter("management")

var ManagementRouterGroup = []http.EndPointInterface{
	http.NewEndPointBuilder[*model.CreateUserRequest, *model.BaseResponse]().
		SetNecessaryHeaders("Authorization").
		SetHandlerChain(http.NewChain(service.ManagerService.CreateUser)).
		SetAllowMethods(http.POST).
		SetRouter(managementRouter.Group("/users")).
		Build(),
}
