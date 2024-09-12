package api

import (
	"github.com/alioth-center/akasha-whisper/app/service"
	"github.com/alioth-center/infrastructure/network/http"
	"github.com/alioth-center/infrastructure/thirdparty/openai"
)

var CompatibleApi compatibleApiImpl

type compatibleApiImpl struct {
	service *service.CompatibleService
}

func (impl compatibleApiImpl) CompleteChat() http.Chain[*openai.CompleteChatRequestBody, *openai.CompleteChatResponseBody] {
	return http.NewChain(impl.service.ChatCompleteAuthorize, impl.service.ChatComplete)
}

func (impl compatibleApiImpl) ListModel() http.Chain[*openai.ListModelRequest, *openai.ListModelResponseBody] {
	return http.NewChain(impl.service.ListModelAuthorize, impl.service.ListModel)
}

func (impl compatibleApiImpl) CreateSpeech() http.Chain[*openai.CreateSpeechRequestBody, *openai.CreateSpeechResponseBody] {
	return http.NewChain(impl.service.CreateSpeechAuthorize, impl.service.CreateSpeech)
}
