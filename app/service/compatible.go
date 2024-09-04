package service

import (
	"context"
	"github.com/alioth-center/infrastructure/network/http"
	"github.com/alioth-center/infrastructure/thirdparty/openai"
	"github.com/alioth-center/infrastructure/trace"
	"strings"
	"time"
)

type CompatibleService struct{}

func NewCompatibleService() *CompatibleService {
	return &CompatibleService{}
}

func (srv *CompatibleService) ChatCompleteAuthorize(ctx http.Context[*openai.CompleteChatRequestBody, *openai.CompleteChatResponseBody]) {
	if strings.TrimPrefix(ctx.NormalHeaders().Authorization, "Bearer ") == "" {
		ctx.SetStatusCode(http.StatusUnauthorized)
		ctx.SetResponse(srv.buildErrorChatCompleteResponse(ctx, "unauthorized"))
		ctx.Abort()
	}
}

func (srv *CompatibleService) ChatComplete(ctx http.Context[*openai.CompleteChatResponseBody, *openai.CompleteChatResponseBody]) {

}

func (srv *CompatibleService) buildErrorChatCompleteResponse(ctx context.Context, content string) *openai.CompleteChatResponseBody {
	return &openai.CompleteChatResponseBody{
		ID:      trace.GetTid(ctx),
		Object:  "chat.completion",
		Created: time.Now().UnixMilli(),
		Choices: []openai.ReplyChoiceObject{{Index: 0, Message: openai.ChatMessageObject{Role: openai.ChatRoleEnumAssistant, Content: content}, FinishReason: "error"}},
		Usage:   openai.UsageObject{},
		Model:   "akasha-whisper",
	}
}
