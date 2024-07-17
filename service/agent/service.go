package agent

import (
	"github.com/alioth-center/akasha-whisper/dao"
	"github.com/alioth-center/infrastructure/logger"
	"github.com/alioth-center/infrastructure/network/http"
	"github.com/alioth-center/infrastructure/thirdparty/openai"
	"github.com/alioth-center/infrastructure/utils/concurrency"
)

type AkashaAgentSrv interface {
	CompleteChat(ctx http.Context[*openai.CompleteChatRequestBody, *openai.CompleteChatResponseBody])
	ListModel(ctx http.Context[*openai.ListModelRequest, *openai.ListModelResponseBody])
}

func NewAkashaAgentSrv(db dao.DatabaseAccessor, log logger.Logger) AkashaAgentSrv {
	return &akashaAgentSrvImpl{db: db, log: log, clients: concurrency.NewMap[int, openai.Client]()}
}

type akashaAgentSrvImpl struct {
	clients concurrency.Map[int, openai.Client]
	db      dao.DatabaseAccessor
	log     logger.Logger
}
