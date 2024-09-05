package service

import (
	"context"
	"github.com/alioth-center/akasha-whisper/app/global"
	"github.com/alioth-center/akasha-whisper/app/model/dto"
	"github.com/alioth-center/infrastructure/logger"
	"github.com/alioth-center/infrastructure/thirdparty/openai"
	"github.com/alioth-center/infrastructure/utils/values"
	"github.com/pandodao/tokenizer-go"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"strings"
)

func CheckApiKeyAvailable(ctx context.Context, key string) (bool, error) {
	token := strings.TrimPrefix(key, "Bearer ")

	// check from bloom filter first
	if token == "" || !global.BearerTokenBloomFilterInstance.CheckKey(token) {
		return false, nil
	}

	// check from database
	return global.WhisperUserDatabaseInstance.CheckWhisperUserApiKey(ctx, token)
}

func GetAvailableClient(ctx context.Context, key string, modelName string, promptToken int64) (client openai.Client, metadata *dto.AvailableClientDTO, err error) {
	token := strings.TrimPrefix(key, "Bearer ")

	clients, queryErr := global.OpenaiClientDatabaseInstance.GetAvailableClients(ctx, modelName, token)
	if queryErr != nil {
		global.Logger.Info(logger.NewFields(ctx).WithMessage("query available clients failed").WithData(queryErr))
		return nil, nil, queryErr
	}

	// filter clients, only return clients that have enough balance
	values.FilterArray(clients, func(client *dto.AvailableClientDTO) bool {
		promptPrice := client.ModelPromptPrice.Mul(decimal.NewFromInt(promptToken))
		affordable := client.ClientBalance.GreaterThanOrEqual(promptPrice) && client.UserBalance.GreaterThanOrEqual(promptPrice)

		if affordable {
			// update client weight, weight = balance/price * weight
			client.ClientWeight = client.ClientBalance.Div(client.ModelPromptPrice).Mul(decimal.NewFromInt(client.ClientWeight)).IntPart()
		}

		return affordable
	})

	// no available client, return error
	if len(clients) == 0 {
		return nil, nil, ErrorNoAvailableClient
	}

	// sort clients by weight
	clients = values.SortArray(clients, func(a, b *dto.AvailableClientDTO) bool { return a.ClientWeight > b.ClientWeight })

	// get openai client from cache
	effectiveClient := clients[0]
	global.Logger.Info(logger.NewFields(ctx).WithMessage("effective client calculated").WithData(effectiveClient))
	openaiClient, exist := global.OpenaiClientCacheInstance.Get(effectiveClient.ClientID)
	if !exist {
		// lazy initialize openai client
		secret, querySecretErr := global.OpenaiClientDatabaseInstance.GetClientSecret(ctx, effectiveClient.ClientID)
		if querySecretErr != nil {
			global.Logger.Info(logger.NewFields(ctx).WithMessage("query client secret failed").WithData(map[string]any{"metadata": effectiveClient, "error": querySecretErr}))
			return nil, nil, querySecretErr
		}

		openaiClientConfig := openai.Config{
			ApiKey:  secret.ClientKey,
			BaseUrl: secret.ClientEndpoint,
		}
		openaiClient = openai.NewClient(openaiClientConfig, global.Logger)
		global.OpenaiClientCacheInstance.Set(effectiveClient.ClientID, openaiClient)
	}

	return openaiClient, effectiveClient, nil
}

func CalculatePromptToken(inputs ...string) (promptToken int64) {
	for _, input := range inputs {
		promptToken += int64(tokenizer.MustCalToken(input))
	}

	return promptToken
}

func CheckManagementKeyAvailable(ctx context.Context, key string) bool {
	token := strings.TrimPrefix(key, "Bearer ")

	if token == "" || global.Config.App.ManagementToken != token {
		// management_token not match or empty
		return false
	}

	return true
}

var (
	ErrorNoAvailableClient = errors.New("no available client")
)
