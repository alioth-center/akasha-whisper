package agent

import (
	"context"
	"fmt"
	"github.com/alioth-center/akasha-whisper/global"
	"github.com/alioth-center/akasha-whisper/model"
	"github.com/alioth-center/infrastructure/logger"
	"github.com/alioth-center/infrastructure/network/http"
	"github.com/alioth-center/infrastructure/thirdparty/openai"
	"github.com/alioth-center/infrastructure/trace"
	"github.com/alioth-center/infrastructure/utils/values"
	"github.com/pandodao/tokenizer-go"
	"strings"
	"time"
)

func (srv *akashaAgentSrvImpl) CompleteChat(ctx http.Context[*openai.CompleteChatRequestBody, *openai.CompleteChatResponseBody]) {
	token, request := strings.TrimPrefix(ctx.HeaderParams().GetString("Authorization"), "Bearer "), ctx.Request()

	// calculate the token of the input strings, abort if token exceed the limit
	var inputs []string
	for _, message := range request.Messages {
		inputs = append(inputs, message.Content)
	}
	tokens := srv.calculateToken(inputs...)
	if global.Config.MaxToken > 0 && tokens > global.Config.MaxToken {
		// token limit exceed, return error message
		ctx.SetStatusCode(400)
		ctx.SetResponse(srv.processCompletionMessage(ctx, "token limit exceed"))
		return
	}

	// get available client with its client info
	client, clientInfo := srv.getClient(ctx, token, request.User, request.Model, tokens)
	if client == nil {
		// no available client, return error message
		ctx.SetStatusCode(400)
		ctx.SetResponse(srv.processCompletionMessage(ctx, "no available client"))
		return
	}

	// execute the request
	response, execErr := client.CompleteChat(openai.CompleteChatRequest{
		Body: openai.CompleteChatRequestBody{
			Model:            request.Model,
			Messages:         request.Messages,
			Temperature:      request.Temperature,
			TopP:             request.TopP,
			N:                request.N,
			PresencePenalty:  request.PresencePenalty,
			FrequencyPenalty: request.FrequencyPenalty,
			LogitBias:        request.LogitBias,
			LogProbs:         request.LogProbs,
			TopLogProbs:      request.TopLogProbs,
			ResponseFormat:   request.ResponseFormat,
			Seed:             request.Seed,
			ServiceTier:      request.ServiceTier,
			Tools:            request.Tools,
			ToolChoice:       request.ToolChoice,
			Stream:           false,
			MaxTokens:        global.Config.MaxToken,
		},
	})
	if execErr != nil {
		// execute error, return error message
		ctx.SetStatusCode(500)
		ctx.SetResponse(srv.processCompletionMessage(ctx, execErr.Error()))
		return
	}

	// request success, update the client balance
	updateErr := srv.updateBalance(ctx, clientInfo, response, ctx.ExtraParams().GetString(http.RemoteIPKey))
	if updateErr != nil {
		srv.log.Error(logger.NewFields(ctx).WithMessage("update balance error").WithData(updateErr))
	}

	// return the response
	ctx.SetStatusCode(200)
	ctx.SetResponse(&response)
}

// processCompletionMessage is a helper function to generate an error completion message
func (srv *akashaAgentSrvImpl) processCompletionMessage(ctx context.Context, content string) *openai.CompleteChatResponseBody {
	return &openai.CompleteChatResponseBody{
		ID:      trace.GetTid(ctx),
		Object:  "error",
		Created: time.Now().UnixMilli(),
		Choices: []openai.ReplyChoiceObject{
			{Index: 0, Message: openai.ChatMessageObject{Role: "system", Content: content}, FinishReason: "error"},
		},
		Usage: openai.UsageObject{},
		Model: "akasha-whisper",
	}
}

// getClient is a helper function to get the client with the api key and user email
func (srv *akashaAgentSrvImpl) getClient(ctx context.Context, apiKey, userEmail, modelName string, tokens int) (openai.Client, model.AvailableClientDTO) {
	userDTO, checkErr := srv.db.GetUserByApiKey(ctx, apiKey)
	if checkErr != nil {
		return nil, model.AvailableClientDTO{}
	}

	var (
		isUser       = userDTO.Role == "user"
		clientRecord []model.AvailableClientDTO
		queryErr     error
	)
	if !isUser {
		// not user api key, use user's email to query the real client
		clientRecord, queryErr = srv.db.GetAvailableUserClient(ctx, modelName, userEmail)
		if queryErr != nil {
			return nil, model.AvailableClientDTO{}
		}
	} else {
		// user api key, get the client directly
		clientRecord, queryErr = srv.db.GetAvailableClient(ctx, modelName, apiKey)
		if queryErr != nil {
			return nil, model.AvailableClientDTO{}
		}
	}

	// calculate client/user balance affordability
	clientRecord = values.FilterArray(clientRecord, func(dto model.AvailableClientDTO) bool {
		clientAffordable := dto.ClientBalance >= float64(tokens)*dto.ModelPromptPrice
		userAffordable := dto.UserBalance >= float64(tokens)*dto.ModelPromptPrice

		return clientAffordable && userAffordable
	})

	// no available client, return nil
	if len(clientRecord) == 0 {
		return nil, model.AvailableClientDTO{}
	}

	// calculate the weight of the client and sort the client by weight
	for _, dto := range clientRecord {
		// W = balance/price * weight
		dto.ClientWeight = int(dto.ClientBalance / dto.ModelPromptPrice * float64(dto.ClientWeight))
	}
	fixed := values.SortArray(clientRecord, func(a, b model.AvailableClientDTO) bool { return a.ClientWeight > b.ClientWeight })

	// client is nil, lazy init the client
	client, exist := srv.clients.Get(fixed[0].ClientID)
	if client == nil || !exist {
		secret, getSecretErr := srv.db.GetClientSecret(ctx, fixed[0].ClientID)
		if getSecretErr != nil || secret.ClientEndpoint == "" || secret.ClientKey == "" {
			// get secret error, return nil
			return nil, model.AvailableClientDTO{}
		}

		client = openai.NewClient(openai.Config{ApiKey: secret.ClientKey, BaseUrl: secret.ClientEndpoint}, srv.log)
		srv.clients.Set(fixed[0].ClientID, client)
	}

	return client, fixed[0]
}

// updateBalance is a helper function to update the balance of the client and user
func (srv *akashaAgentSrvImpl) updateBalance(ctx context.Context, info model.AvailableClientDTO, response openai.CompleteChatResponseBody, ip string) error {
	// calculate the token of the response
	cost := float64(response.Usage.PromptTokens)*info.ModelPromptPrice +
		float64(response.Usage.CompletionTokens)*info.ModelCompletionPrice

	// update balance
	recordErr := srv.db.AddRequestRecord(ctx, model.RequestRecordDTO{
		ClientID:             info.ClientID,
		ModelID:              info.ModelID,
		UserID:               info.UserID,
		PromptTokenUsage:     response.Usage.PromptTokens,
		CompletionTokenUsage: response.Usage.CompletionTokens,
		RequestIP:            ip,
		BalanceCost:          cost,
	})
	if recordErr != nil {
		return fmt.Errorf("failed to record request: %w", recordErr)
	}

	return nil
}

// calculateToken is a helper function to calculate the token of the input strings
func (srv *akashaAgentSrvImpl) calculateToken(inputs ...string) (tokens int) {
	for _, input := range inputs {
		tokens += tokenizer.MustCalToken(input)
	}

	return tokens
}
