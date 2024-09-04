package dao

import (
	"embed"
	"fmt"
	"github.com/alioth-center/akasha-whisper/global"
	"github.com/alioth-center/infrastructure/logger"
	"github.com/alioth-center/infrastructure/trace"
	"github.com/alioth-center/infrastructure/utils/values"
)

var (
	//go:embed rawsql/*.sql
	rawSqlEmbedding embed.FS

	rawSqlList = map[RawsqlKey]string{}
)

type RawsqlKey string

const (
	RawsqlOpenaiClientGetAvailableClients RawsqlKey = "openai_client.get_available_clients.sql"
	RawsqlOpenaiClientGetClientSecrets    RawsqlKey = "openai_client.get_client_secrets.sql"
)

var (
	rawSqlNames = []RawsqlKey{
		RawsqlOpenaiClientGetAvailableClients,
		RawsqlOpenaiClientGetClientSecrets,
	}
)

func init() {
	// load rawsql list
	for _, name := range rawSqlNames {
		content, readErr := rawSqlEmbedding.ReadFile(values.BuildStrings("rawsql/", string(name)))
		if readErr != nil {
			panic(readErr)
		}

		rawSqlList[name] = string(content)
	}
}

func QueryTest() {
	// init db
	clientAccessor := NewOpenaiClientDatabaseAccessor(global.DatabaseV2)
	ctx := trace.NewContext()

	// get available clients
	clients, err := clientAccessor.GetAvailableClients(ctx, "gpt-4o", "114514")
	if err != nil {
		fmt.Println(err)
	}
	for _, client := range clients {
		logger.Info(logger.NewFields(ctx).WithData(client))
	}

	// get client secret
	secrets, err := clientAccessor.GetClientSecret(ctx, 1, 2, 3)
	if err != nil {
		fmt.Println(err)
	}
	for _, secret := range secrets {
		logger.Info(logger.NewFields(ctx).WithData(secret))
	}

	//// create client
	//_, createErr := clientAccessor.CreateClient(ctx, &model.OpenaiClient{
	//	Description: "chatanywhere",
	//	ApiKey:      "sk-zJhQ8H2MklG9r2r03gZIJkBf6bWOZTy9BstMIiGnanOD4m9c",
	//	Endpoint:    "https://api.chatanywhere.org/v1/",
	//	Weight:      1,
	//})
	//if createErr != nil {
	//	fmt.Println(createErr)
	//}
	//
	//// create models
	modelAccessor := NewOpenaiModelDatabaseAccessor(global.DatabaseV2)
	models, err := modelAccessor.GetAvailableModelsByApiKey(ctx, "114514")
	if err != nil {
		fmt.Println(err)
	}

	for _, model := range models {
		logger.Info(logger.NewFields(ctx).WithData(model))
	}
	//createErr = modelAccessor.CreateOrUpdateModelWithClientDescriptions(ctx, &model.OpenaiModel{
	//	Model:           "gpt-4o-mini",
	//	MaxTokens:       128000,
	//	PromptPrice:     0.003555000,
	//	CompletionPrice: 0.010665000,
	//}, "openai_origin", "chatanywhere")
	//if createErr != nil {
	//	fmt.Println(createErr)
	//}

	// add balance
	//clientBalanceAccessor := NewOpenaiClientBalanceDatabaseAccessor(global.DatabaseV2)
	//after, err := clientBalanceAccessor.CreateBalanceRecord(ctx, 1, decimal.NewFromFloat(114), model.OpenaiClientBalanceActionGift)
	//if err != nil {
	//	fmt.Println(err)
	//}

	//logger.Info(logger.NewFields(ctx).WithData(after))
}
