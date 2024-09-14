package service

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	nh "net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/sse"
	"github.com/gin-gonic/gin"

	"github.com/alioth-center/akasha-whisper/app/global"
	"github.com/alioth-center/akasha-whisper/app/model"
	"github.com/alioth-center/infrastructure/logger"
	"github.com/alioth-center/infrastructure/network/http"
	"github.com/alioth-center/infrastructure/thirdparty/openai"
	"github.com/alioth-center/infrastructure/trace"
	"github.com/alioth-center/infrastructure/utils/values"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type CompatibleService struct{}

func NewCompatibleService() *CompatibleService { return &CompatibleService{} }

func (srv *CompatibleService) ChatComplete(ctx http.Context[*openai.CompleteChatRequestBody, *openai.CompleteChatResponseBody]) {
	apiKey, request := ctx.NormalHeaders().Authorization, ctx.Request()

	// calculate prompt token
	inputMessages := make([]string, len(request.Messages))
	for i, message := range request.Messages {
		inputMessages[i] = message.Content
	}
	promptToken := CalculatePromptToken(inputMessages...)

	// get available openai client
	client, metadata, getErr := GetAvailableClient(ctx, apiKey, request.Model, promptToken)
	if getErr != nil && errors.Is(getErr, ErrorNoAvailableClient) {
		ctx.SetStatusCode(http.StatusForbidden)
		ctx.SetResponse(srv.buildErrorChatCompleteResponse(ctx, "no available client"))
		ctx.Abort()
		return
	} else if getErr != nil {
		ctx.SetStatusCode(http.StatusInternalServerError)
		ctx.SetResponse(srv.buildErrorChatCompleteResponse(ctx, values.BuildStrings("internal server error: ", getErr.Error())))
		ctx.Abort()
		return
	}

	// complete chat
	response, executeErr := client.CompleteChat(ctx, openai.CompleteChatRequest{
		Body: openai.CompleteChatRequestBody{
			Model:            request.Model,
			Messages:         request.Messages,
			Temperature:      request.Temperature,
			TopP:             request.TopP,
			N:                request.N,
			Stream:           false,
			MaxTokens:        min(request.MaxTokens, global.Config.App.MaxToken),
			PresencePenalty:  request.PresencePenalty,
			FrequencyPenalty: request.FrequencyPenalty,
			User:             request.User,
			LogitBias:        request.LogitBias,
			LogProbs:         request.LogProbs,
			TopLogProbs:      request.TopLogProbs,
			ResponseFormat:   request.ResponseFormat,
			Seed:             request.Seed,
			ServiceTier:      request.ServiceTier,
			Tools:            request.Tools,
			ToolChoice:       request.ToolChoice,
		},
	})
	if executeErr != nil {
		ctx.SetStatusCode(http.StatusInternalServerError)
		ctx.SetResponse(srv.buildErrorChatCompleteResponse(ctx, values.BuildStrings("internal server error: ", executeErr.Error())))
		ctx.Abort()
		return
	}

	// consume success, update balances
	promptCostAmount := metadata.ModelPromptPrice.Mul(decimal.NewFromInt(int64(response.Usage.PromptTokens))).Div(decimal.NewFromInt(global.Config.App.PriceTokenUnit))
	completionCostAmount := metadata.ModelCompletionPrice.Mul(decimal.NewFromInt(int64(response.Usage.CompletionTokens))).Div(decimal.NewFromInt(global.Config.App.PriceTokenUnit))
	balanceCost := promptCostAmount.Add(completionCostAmount).Mul(decimal.NewFromInt(-1))
	_, updateClientBalanceErr := global.OpenaiClientBalanceDatabaseInstance.CreateBalanceRecord(ctx, metadata.ClientID, balanceCost, model.OpenaiClientBalanceActionConsumption)
	_, updateUserBalanceErr := global.WhisperUserBalanceDatabaseInstance.CreateBalanceRecord(ctx, metadata.UserID, balanceCost, model.WhisperUserBalanceActionConsumption)
	updateRequestErr := global.OpenaiRequestDatabaseInstance.CreateOpenaiRequestRecord(ctx, &model.OpenaiRequest{
		ClientID:             int64(metadata.ClientID),
		ModelID:              int64(metadata.ModelID),
		UserID:               int64(metadata.UserID),
		RequestIP:            ctx.ExtraParams().GetString(http.RemoteIPKey),
		RequestID:            response.ID,
		TraceID:              trace.GetTid(ctx),
		PromptTokenUsage:     response.Usage.PromptTokens,
		CompletionTokenUsage: response.Usage.CompletionTokens,
		BalanceCost:          balanceCost.Abs(),
	})
	for _, err := range []error{updateClientBalanceErr, updateUserBalanceErr, updateRequestErr} {
		if err != nil {
			global.Logger.Error(logger.NewFields(ctx).WithMessage("update response result failed").WithData(err))
		}
	}

	// return openai response
	ctx.SetStatusCode(http.StatusOK)
	ctx.SetResponse(&response)
}

func (srv *CompatibleService) StreamingChatCompletion(ctx *gin.Context) {
	// check api key available
	apiKey := ctx.GetHeader(http.HeaderAuthorization)
	exist, allowIPs, err := CheckApiKeyAvailable(ctx, apiKey)
	if err != nil {
		global.Logger.Error(logger.NewFields(ctx).WithMessage("check api key available failed").WithData(err))
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, srv.buildErrorChatCompleteResponse(ctx, values.BuildStrings("internal server error: ", err.Error())))
		return
	}
	if !exist {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, srv.buildErrorChatCompleteResponse(ctx, "unauthorized"))
		return
	}

	// check allow ip
	if !CheckAllowIP(ctx, ctx.ClientIP(), strings.Split(allowIPs, ",")) {
		ctx.AbortWithStatusJSON(http.StatusForbidden, srv.buildErrorChatCompleteResponse(ctx, "ip forbidden"))
		return
	}

	request := &openai.CompleteChatRequestBody{}
	bindErr := ctx.ShouldBindJSON(request)
	if bindErr != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": bindErr.Error()})
	}
	writeBack, _ := json.Marshal(request)
	ctx.Request.Body = io.NopCloser(bytes.NewBuffer(writeBack))

	if request.Stream {
		defer ctx.Abort()
		inputMessages := make([]string, len(request.Messages))
		for i, message := range request.Messages {
			inputMessages[i] = message.Content
		}
		promptToken := CalculatePromptToken(inputMessages...)

		// get available openai client
		_, metadata, getErr := GetAvailableClient(ctx, apiKey, request.Model, promptToken)
		if getErr != nil && errors.Is(getErr, ErrorNoAvailableClient) {
			ctx.AbortWithStatusJSON(http.StatusForbidden, srv.buildErrorChatCompleteResponse(ctx, "no available client"))
			return
		} else if getErr != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, srv.buildErrorChatCompleteResponse(ctx, values.BuildStrings("internal server error: ", getErr.Error())))
			return
		}

		// get client secrets
		secrets, exist := global.OpenaiClientSecretsCacheInstance.Get(metadata.ClientID)
		if !exist {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, srv.buildErrorChatCompleteResponse(ctx, values.BuildStrings("internal server error: ", "client secrets not found")))
			return
		}

		// complete chat
		requestURL, parseErr := url.JoinPath(secrets.BaseUrl, "/chat/completions")
		if parseErr != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, srv.buildErrorChatCompleteResponse(ctx, values.BuildStrings("internal server error: ", parseErr.Error())))
			return
		}
		global.Logger.Info(logger.NewFields(ctx).WithMessage("request openai chat completion").WithData(requestURL))
		body, _ := json.Marshal(request)
		openaiRequest, buildErr := nh.NewRequest(http.POST, requestURL, bytes.NewBuffer(body))
		if buildErr != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, srv.buildErrorChatCompleteResponse(ctx, values.BuildStrings("internal server error: ", buildErr.Error())))
			return
		}

		openaiRequest.Header.Set(http.HeaderAuthorization, values.BuildStrings("Bearer ", secrets.ApiKey))
		openaiRequest.Header.Set(http.HeaderContentType, http.ContentTypeJson)

		response, executeErr := httpClient.Do(openaiRequest)
		if executeErr != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, srv.buildErrorChatCompleteResponse(ctx, values.BuildStrings("internal server error: ", executeErr.Error())))
			return
		}

		if response.StatusCode != http.StatusOK {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, srv.buildErrorChatCompleteResponse(ctx, values.BuildStrings("internal server error: ", "response status code not 200, ", strconv.Itoa(response.StatusCode))))
			return
		}

		ctx.Header(http.HeaderContentType, "text/event-stream")
		ctx.Header("Cache-Control", "no-cache")
		ctx.Header("Transfer-Encoding", "chunked")
		ctx.Header("Connection", "keep-alive")

		// parse streaming response
		reader := bufio.NewReader(response.Body)
		for {
			line, readErr := reader.ReadString('\n')
			line = strings.TrimSuffix(line, "\n")
			line = strings.TrimPrefix(line, "data: ")
			if readErr != nil && errors.Is(readErr, io.EOF) {
				break
			} else if readErr != nil {
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, srv.buildErrorChatCompleteResponse(ctx, values.BuildStrings("internal server error: ", readErr.Error())))
			}

			if len(strings.TrimSpace(line)) > 0 {
				fmt.Println(line)
				e := sse.Encode(ctx.Writer, sse.Event{Data: line})
				if e != nil {
					fmt.Println("error", e)
				}
				ctx.Writer.Flush()
				fmt.Println("flush")
			}
		}
	}
}

func (srv *CompatibleService) ListModel(ctx http.Context[*openai.ListModelRequest, *openai.ListModelResponseBody]) {
	apiKey := strings.TrimPrefix(ctx.NormalHeaders().Authorization, "Bearer ")

	models, queryErr := global.OpenaiModelDatabaseInstance.GetAvailableModelsByApiKey(ctx, apiKey)
	if queryErr != nil {
		global.Logger.Error(logger.NewFields(ctx).WithMessage("query available models failed").WithData(queryErr))
		ctx.SetStatusCode(http.StatusInternalServerError)
		ctx.SetResponse(&openai.ListModelResponseBody{})
		ctx.Abort()
		return
	}

	set := map[string]openai.ModelObject{}
	for _, dto := range models {
		set[dto.ModelName] = openai.ModelObject{ID: dto.ModelName, Created: dto.LastUpdatedAt.Unix(), Object: "model", OwnedBy: "openai"}
	}

	response := &openai.ListModelResponseBody{Object: "list", Data: make([]openai.ModelObject, 0, len(set))}
	for _, value := range set {
		response.Data = append(response.Data, value)
	}

	ctx.SetStatusCode(http.StatusOK)
	ctx.SetResponse(response)
}

func (srv *CompatibleService) CreateSpeechAuthorize(ctx http.Context[*openai.CreateSpeechRequestBody, *openai.CreateSpeechResponseBody]) {
	// check api key available
	exist, allowIPs, err := CheckApiKeyAvailable(ctx, ctx.NormalHeaders().Authorization)
	if err != nil {
		global.Logger.Error(logger.NewFields(ctx).WithMessage("check api key available failed").WithData(err))
		ctx.SetStatusCode(http.StatusInternalServerError)
		ctx.SetResponse(&openai.CreateSpeechResponseBody{})
		ctx.Abort()
		return
	}

	if !exist {
		ctx.SetStatusCode(http.StatusUnauthorized)
		ctx.SetResponse(&openai.CreateSpeechResponseBody{})
		ctx.Abort()
		return
	}

	// check allow ip
	if !CheckAllowIP(ctx, ctx.ClientIP(), strings.Split(allowIPs, ",")) {
		ctx.SetStatusCode(http.StatusForbidden)
		ctx.SetResponse(&openai.CreateSpeechResponseBody{})
		ctx.Abort()
	}
}

func (srv *CompatibleService) CreateSpeech(ctx http.Context[*openai.CreateSpeechRequestBody, *openai.CreateSpeechResponseBody]) {
	apiKey, request := ctx.NormalHeaders().Authorization, ctx.Request()

	// calculate prompt token
	promptToken := int64(len([]rune(request.Input)))

	// get available openai client
	client, metadata, getErr := GetAvailableClient(ctx, apiKey, request.Model, promptToken)
	if getErr != nil && errors.Is(getErr, ErrorNoAvailableClient) {
		ctx.SetStatusCode(http.StatusForbidden)
		ctx.SetResponse(&openai.CreateSpeechResponseBody{})
		ctx.Abort()
		return
	} else if getErr != nil {
		ctx.SetStatusCode(http.StatusInternalServerError)
		ctx.SetResponse(&openai.CreateSpeechResponseBody{})
		ctx.Abort()
		return
	}

	response, executeErr := client.CreateSpeech(ctx, openai.CreateSpeechRequest{
		Body: openai.CreateSpeechRequestBody{
			Model:          request.Model,
			Input:          request.Input,
			Voice:          request.Voice,
			ResponseFormat: request.ResponseFormat,
			Speed:          request.Speed,
		},
	})
	if executeErr != nil {
		ctx.SetStatusCode(http.StatusInternalServerError)
		ctx.SetResponse(&openai.CreateSpeechResponseBody{})
		ctx.Abort()
		return
	}

	// consume success, update balances
	promptCostAmount := metadata.ModelPromptPrice.Mul(decimal.NewFromInt(promptToken).Div(decimal.NewFromInt(global.Config.App.PriceTokenUnit)))
	balanceCost := promptCostAmount.Mul(decimal.NewFromInt(-1))
	_, updateClientBalanceErr := global.OpenaiClientBalanceDatabaseInstance.CreateBalanceRecord(ctx, metadata.ClientID, balanceCost, model.OpenaiClientBalanceActionConsumption)
	_, updateUserBalanceErr := global.WhisperUserBalanceDatabaseInstance.CreateBalanceRecord(ctx, metadata.UserID, balanceCost, model.WhisperUserBalanceActionConsumption)
	updateRequestErr := global.OpenaiRequestDatabaseInstance.CreateOpenaiRequestRecord(ctx, &model.OpenaiRequest{
		ClientID:             int64(metadata.ClientID),
		ModelID:              int64(metadata.ModelID),
		UserID:               int64(metadata.UserID),
		RequestIP:            ctx.ExtraParams().GetString(http.RemoteIPKey),
		RequestID:            trace.GetTid(ctx),
		TraceID:              trace.GetTid(ctx),
		PromptTokenUsage:     int(promptToken),
		CompletionTokenUsage: 0,
		BalanceCost:          balanceCost.Abs(),
	})
	for _, err := range []error{updateClientBalanceErr, updateUserBalanceErr, updateRequestErr} {
		if err != nil {
			global.Logger.Error(logger.NewFields(ctx).WithMessage("update response result failed").WithData(err))
		}
	}

	// set response file header
	switch request.ResponseFormat {
	case "opus":
		ctx.SetResponseHeader(http.HeaderContentType, "audio/ogg")
	case "aac":
		ctx.SetResponseHeader(http.HeaderContentType, "audio/aac")
	case "flac":
		ctx.SetResponseHeader(http.HeaderContentType, "audio/flac")
	case "pcm":
		ctx.SetResponseHeader(http.HeaderContentType, "audio/L16")
	default:
		ctx.SetResponseHeader(http.HeaderContentType, "audio/mpeg")
	}
	ctx.SetStatusCode(http.StatusOK)
	ctx.SetResponse(&response)
}

func (srv *CompatibleService) ListModelAuthorize(ctx http.Context[*openai.ListModelRequest, *openai.ListModelResponseBody]) {
	// check api key available
	exist, allowIPs, err := CheckApiKeyAvailable(ctx, ctx.NormalHeaders().Authorization)
	if err != nil {
		global.Logger.Error(logger.NewFields(ctx).WithMessage("check api key available failed").WithData(err))
		ctx.SetStatusCode(http.StatusInternalServerError)
		ctx.SetResponse(&openai.ListModelResponseBody{})
		ctx.Abort()
		return
	}
	if !exist {
		ctx.SetStatusCode(http.StatusUnauthorized)
		ctx.SetResponse(&openai.ListModelResponseBody{})
		ctx.Abort()
		return
	}

	// check allow ip
	if !CheckAllowIP(ctx, ctx.ClientIP(), strings.Split(allowIPs, ",")) {
		ctx.SetStatusCode(http.StatusForbidden)
		ctx.SetResponse(&openai.ListModelResponseBody{})
		ctx.Abort()
	}
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
