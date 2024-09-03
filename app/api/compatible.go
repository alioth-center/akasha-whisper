package api

import (
	"github.com/alioth-center/akasha-whisper/service/agent"
	"github.com/alioth-center/infrastructure/network/http"
	"github.com/alioth-center/infrastructure/thirdparty/openai"
)

var CompatibleApi compatibleApiImpl

type compatibleApiImpl struct {
	srv agent.AkashaAgentSrv
}

func (impl compatibleApiImpl) CompleteChat() http.Chain[*openai.CompleteChatRequestBody, *openai.CompleteChatResponseBody] {
	return http.NewChain(impl.srv.CompleteChat)
}

func (impl compatibleApiImpl) ListModel() http.Chain[*openai.ListModelRequest, *openai.ListModelResponseBody] {
	return http.NewChain(impl.srv.ListModel)
}
