package api

import (
	"github.com/alioth-center/akasha-whisper/global"
	"github.com/alioth-center/akasha-whisper/model"
	"github.com/alioth-center/akasha-whisper/service"
	"github.com/alioth-center/infrastructure/network/http"
	"github.com/alioth-center/infrastructure/thirdparty/openai"
)

// BindChatCompletion binds the endpoint for chat completion
//
// Router: POST /chat/completions
func BindChatCompletion() {
	global.Engine.AddEndPoints(http.NewEndPointBuilder[*openai.CompleteChatRequestBody, *openai.CompleteChatResponseBody]().
		SetNecessaryHeaders("Authorization").
		SetHandlerChain(http.NewChain(service.AgentService.CompleteChat)).
		SetAllowMethods(http.POST).
		SetRouter(http.NewRouter("/chat/completions")).
		Build())
}

// BindListModel binds the endpoint for listing models
//
// Router: GET /models
func BindListModel() {
	global.Engine.AddEndPoints(http.NewEndPointBuilder[*openai.ListModelRequest, *openai.ListModelResponseBody]().
		SetNecessaryHeaders("Authorization").
		SetHandlerChain(http.NewChain(service.AgentService.ListModel)).
		SetAllowMethods(http.GET).
		SetRouter(http.NewRouter("/models")).
		Build())
}

// BindCreateUser binds the endpoint for creating user
//
// Router: POST /manage/users
func BindCreateUser() {
	global.Engine.AddEndPoints(http.NewEndPointBuilder[*model.CreateUserRequest, *model.BaseResponse]().
		SetNecessaryHeaders("Authorization").
		SetHandlerChain(http.NewChain(service.ManagerService.CreateUser)).
		SetAllowMethods(http.POST).
		SetRouter(http.NewRouter("/manage/users")).
		Build())
}
