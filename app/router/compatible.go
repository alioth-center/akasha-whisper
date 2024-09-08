package router

import (
	"github.com/alioth-center/akasha-whisper/app/api"
	"github.com/alioth-center/infrastructure/network/http"
	"github.com/alioth-center/infrastructure/thirdparty/openai"
)

var OpenAiCompatibleRouterGroup = []http.EndPointInterface{
	http.NewEndPointBuilder[*openai.CompleteChatRequestBody, *openai.CompleteChatResponseBody]().
		SetNecessaryHeaders("Authorization").
		SetHandlerChain(api.CompatibleApi.CompleteChat()).
		SetAllowMethods(http.POST).
		SetRouter(http.NewRouter("/chat/completions")).
		Build(),
	http.NewEndPointBuilder[*openai.ListModelRequest, *openai.ListModelResponseBody]().
		SetNecessaryHeaders("Authorization").
		SetHandlerChain(api.CompatibleApi.ListModel()).
		SetAllowMethods(http.GET).
		SetRouter(http.NewRouter("/models")).
		Build(),
}
