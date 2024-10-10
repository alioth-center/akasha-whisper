package service

import (
	"context"
	"strings"

	"github.com/alioth-center/akasha-whisper/app/global"
	"github.com/alioth-center/akasha-whisper/app/model/dto"
	"github.com/alioth-center/infrastructure/logger"
	"github.com/alioth-center/infrastructure/thirdparty/openai"
	"github.com/alioth-center/infrastructure/utils/network"
	"github.com/alioth-center/infrastructure/utils/values"
	"github.com/pandodao/tokenizer-go"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

func CheckApiKeyAvailable(ctx context.Context, key string) (exist bool, allowIPs string, err error) {
	token := strings.TrimPrefix(key, "Bearer ")

	// check from bloom filter first
	if token == "" || !global.BearerTokenBloomFilterInstance.CheckKey(token) {
		return false, "", nil
	}

	// check from database
	return global.WhisperUserDatabaseInstance.CheckWhisperUserApiKey(ctx, token)
}

func CheckAllowIP(_ context.Context, ip string, allowIPs []string) bool {
	if len(allowIPs) == 0 || (len(allowIPs) == 1 && allowIPs[0] == "") {
		return true
	}

	return len(values.FilterArray(allowIPs, func(cidr string) bool {
		if network.IsValidIP(cidr) {
			return ip == cidr
		}

		return network.IPInCIDR(ip, cidr)
	})) > 0
}

func GetAvailableClient(ctx context.Context, key string, modelName string, promptToken int64, endpoint string) (client openai.Client, metadata *dto.AvailableClientDTO, err error) {
	token := strings.TrimPrefix(key, "Bearer ")

	clients, queryErr := global.OpenaiClientDatabaseInstance.GetAvailableClients(ctx, modelName, token, endpoint)
	if queryErr != nil {
		global.Logger.Info(logger.NewFields(ctx).WithMessage("query available clients failed").WithData(queryErr))
		return nil, nil, queryErr
	}

	// filter clients, only return clients that have enough balance
	values.FilterArray(clients, func(client *dto.AvailableClientDTO) bool {
		promptPrice := client.ModelPromptPrice.Mul(decimal.NewFromInt(promptToken)).Div(decimal.NewFromInt(global.Config.App.PriceTokenUnit))
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
			global.Logger.Error(logger.NewFields(ctx).WithMessage("query client secret failed").WithData(map[string]any{"metadata": effectiveClient, "error": querySecretErr}))
			return nil, nil, querySecretErr
		}

		openaiClientConfig := openai.Config{
			ApiKey:  secret.ClientKey,
			BaseUrl: secret.ClientEndpoint,
		}
		openaiClient = openai.NewClient(openaiClientConfig, global.Logger)
		global.OpenaiClientCacheInstance.Set(effectiveClient.ClientID, openaiClient)
		global.OpenaiClientSecretsCacheInstance.Set(effectiveClient.ClientID, &openaiClientConfig)
		global.Logger.Info(logger.NewFields(ctx).WithMessage("openai client initialized").WithData(map[string]any{"metadata": effectiveClient, "client": openaiClient}))
	}

	return openaiClient, effectiveClient, nil
}

func CalculatePromptToken(inputs ...string) (promptToken int64) {
	for _, input := range inputs {
		promptToken += int64(tokenizer.MustCalToken(input))
	}

	return promptToken
}

func CheckManagementKeyAvailable(_ context.Context, key string) bool {
	token := strings.TrimPrefix(key, "Bearer ")

	if token == "" || global.Config.App.ManagementToken != token {
		// management_token not match or empty
		return false
	}

	return true
}

var ErrorNoAvailableClient = errors.New("no available client")
