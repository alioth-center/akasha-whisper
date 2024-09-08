package service

import (
	"github.com/alioth-center/akasha-whisper/app/entity"
	"github.com/alioth-center/akasha-whisper/app/global"
	"github.com/alioth-center/akasha-whisper/app/model"
	"github.com/alioth-center/infrastructure/logger"
	"github.com/alioth-center/infrastructure/network/http"
	"github.com/alioth-center/infrastructure/thirdparty/openai"
	"github.com/alioth-center/infrastructure/utils/values"
	"github.com/shopspring/decimal"
	"strings"
)

type ManagementService struct{}

func NewManagementService() *ManagementService { return &ManagementService{} }

func (srv *ManagementService) ListAllClients(ctx http.Context[*entity.ListClientsRequest, *entity.ListClientResponse]) {
	clients, err := global.OpenaiClientDatabaseInstance.ListClients(ctx)
	if err != nil {
		response := http.NewBaseResponse[[]*entity.ClientItem](ctx, []*entity.ClientItem{}, err)
		ctx.SetStatusCode(http.StatusInternalServerError)
		ctx.SetResponse(&response)
		return
	}

	items := make([]*entity.ClientItem, len(clients))
	for i, client := range clients {

		item := &entity.ClientItem{
			ID:       client.ClientID,
			Name:     client.ClientDescription,
			ApiKey:   client.ClientKey,
			Endpoint: client.ClientEndpoint,
			Weight:   client.ClientWeight,
			Balance:  client.ClientBalance,
		}

		if len(client.ClientKey) > 10 {
			// sk-xx****...****xyz
			item.ApiKey = values.BuildStrings(client.ClientKey[:6], strings.Repeat("*", len(client.ClientKey)-10), client.ClientKey[len(client.ClientKey)-4:])
		}

		items[i] = item
	}

	response := http.NewBaseResponse[[]*entity.ClientItem](ctx, items, nil)
	ctx.SetStatusCode(http.StatusOK)
	ctx.SetResponse(&response)
}

func (srv *ManagementService) CreateClient(ctx http.Context[*entity.CreateClientRequest, *entity.CreateClientResponse]) {
	request := ctx.Request()
	client := &model.OpenaiClient{
		Description: request.Name,
		ApiKey:      request.ApiKey,
		Endpoint:    request.Endpoint,
		Weight:      request.Weight,
	}

	// insert into database
	created, createErr := global.OpenaiClientDatabaseInstance.CreateClient(ctx, client)
	if createErr != nil {
		response := http.NewBaseResponse[[]*entity.CreateClientScanModelItem](ctx, []*entity.CreateClientScanModelItem{}, createErr)
		ctx.SetStatusCode(http.StatusInternalServerError)
		ctx.SetResponse(&response)
		return
	}
	if !created {
		response := http.NewBaseResponse[[]*entity.CreateClientScanModelItem](ctx, []*entity.CreateClientScanModelItem{}, http.NewBaseError(http.StatusConflict, "client already exists"))
		ctx.SetStatusCode(http.StatusConflict)
		ctx.SetResponse(&response)
		return
	}

	// initialize client balance
	_, initBalanceErr := global.OpenaiClientBalanceDatabaseInstance.CreateBalanceRecord(ctx, int(client.ID), decimal.Zero, model.OpenaiClientBalanceActionInitial)
	if initBalanceErr != nil {
		response := http.NewBaseResponse[[]*entity.CreateClientScanModelItem](ctx, []*entity.CreateClientScanModelItem{}, initBalanceErr)
		ctx.SetStatusCode(http.StatusInternalServerError)
		ctx.SetResponse(&response)
		return
	}

	// initialize openai client
	openaiClient := openai.NewClient(openai.Config{ApiKey: client.ApiKey, BaseUrl: client.Endpoint}, global.Logger)
	models, listErr := openaiClient.ListModels(ctx, openai.ListModelRequest{})
	if listErr != nil {
		global.Logger.Error(logger.NewFields(ctx).WithMessage("failed to list models when create client").WithData(listErr))
		response := http.NewBaseResponse[[]*entity.CreateClientScanModelItem](ctx, []*entity.CreateClientScanModelItem{}, listErr)
		ctx.SetStatusCode(http.StatusOK)
		ctx.SetResponse(&response)
		return
	}

	items := make([]*entity.CreateClientScanModelItem, len(models.Data))
	for i, m := range models.Data {
		items[i] = &entity.CreateClientScanModelItem{ModelName: m.ID, CreatedAt: m.Created}
	}

	response := http.NewBaseResponse[[]*entity.CreateClientScanModelItem](ctx, items, nil)
	ctx.SetStatusCode(http.StatusOK)
	ctx.SetResponse(&response)
}

func (srv *ManagementService) ListClientModels(ctx http.Context[*entity.ListClientModelRequest, *entity.ListClientModelResponse]) {
	models, queryErr := global.OpenaiModelDatabaseInstance.GetModelsByClientDescription(ctx, ctx.PathParams().GetString("client_name"))
	if queryErr != nil {
		response := http.NewBaseResponse[[]*entity.ModelItem](ctx, []*entity.ModelItem{}, queryErr)
		ctx.SetStatusCode(http.StatusInternalServerError)
		ctx.SetResponse(&response)
		return
	}

	items := make([]*entity.ModelItem, len(models))
	for i, modelItem := range models {
		items[i] = &entity.ModelItem{
			ID:              modelItem.ModelID,
			Name:            modelItem.ModelName,
			MaxTokens:       modelItem.MaxTokens,
			RpmLimit:        modelItem.ModelRpmLimit,
			TpmLimit:        modelItem.ModelTpmLimit,
			PromptPrice:     modelItem.PromptPrice,
			CompletionPrice: modelItem.CompletionPrice,
			LastUpdatedAt:   modelItem.LastUpdatedAt.UnixMilli(),
		}
	}

	response := http.NewBaseResponse[[]*entity.ModelItem](ctx, items, nil)
	ctx.SetStatusCode(http.StatusOK)
	ctx.SetResponse(&response)
}

func (srv *ManagementService) CreateClientModels(ctx http.Context[*entity.CreateClientModelRequest, *entity.CreateClientModelResponse]) {
	request := ctx.Request()

	modelData := make([]*model.OpenaiModel, len(request.Models))
	for i, modelItem := range request.Models {
		modelData[i] = &model.OpenaiModel{
			Model:           modelItem.Name,
			MaxTokens:       modelItem.MaxTokens,
			PromptPrice:     modelItem.PromptPrice,
			CompletionPrice: modelItem.CompletionPrice,
			RpmLimit:        modelItem.RpmLimit,
			TpmLimit:        modelItem.TpmLimit,
		}
	}

	insertErr := global.OpenaiModelDatabaseInstance.CreateOrUpdateModelWithClientDescriptions(ctx, modelData, ctx.PathParams().GetString("client_name"))
	if insertErr != nil {
		response := http.NewBaseResponse[*entity.CreateResponse](ctx, &entity.CreateResponse{Success: false}, insertErr)
		ctx.SetStatusCode(http.StatusInternalServerError)
		ctx.SetResponse(&response)
		return
	}

	response := http.NewBaseResponse[*entity.CreateResponse](ctx, &entity.CreateResponse{Success: true}, nil)
	ctx.SetStatusCode(http.StatusOK)
	ctx.SetResponse(&response)
}

func CheckManagementKey[req any, res any](ctx http.Context[req, *http.BaseResponse[res]]) {
	token := ctx.NormalHeaders().Authorization
	if !CheckManagementKeyAvailable(ctx, token) {
		response := http.NewBaseResponse[res](
			ctx, values.Nil[res](), http.NewBaseError(http.StatusForbidden, "unauthorized"),
		)
		ctx.SetStatusCode(http.StatusForbidden)
		ctx.SetResponse(&response)
		ctx.Abort()
	}
}
