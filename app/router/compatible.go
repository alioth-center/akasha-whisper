package router

import (
	"github.com/alioth-center/akasha-whisper/app/api"
	"github.com/alioth-center/infrastructure/network/http"
	"github.com/alioth-center/infrastructure/thirdparty/openai"
)

var OpenAiCompatibleRouterGroup = []http.EndPointInterface{
	http.NewEndPointBuilder[*openai.CompleteChatRequestBody, *openai.CompleteChatResponseBody]().
		SetNecessaryHeaders("Authorization").
		SetGinMiddlewares(api.CompatibleApi.StreamingCompleteChat()...).
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
	// yet have some problem which cannot return audio file correctly
	// http.NewEndPointBuilder[*openai.CreateSpeechRequestBody, *openai.CreateSpeechResponseBody]().
	// 	SetNecessaryHeaders("Authorization").
	//	SetHandlerChain(api.CompatibleApi.CreateSpeech()).
	//	SetAllowMethods(http.POST).
	//	SetRouter(http.NewRouter("/audio/speech")).
	//	Build(),
}
