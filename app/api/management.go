package api

import (
	"github.com/alioth-center/akasha-whisper/app/entity"
	"github.com/alioth-center/akasha-whisper/app/service"
	"github.com/alioth-center/infrastructure/network/http"
)

var ManagementApi managementApiImpl

type managementApiImpl struct {
	service *service.ManagementService
}

func (impl managementApiImpl) ListClients() http.Chain[*entity.ListClientsRequest, *entity.ListClientResponse] {
	return http.NewChain(
		service.CheckManagementKey[*entity.ListClientsRequest, []*entity.ClientItem],
		impl.service.ListAllClients,
	)
}

func (impl managementApiImpl) CreateClient() http.Chain[*entity.CreateClientRequest, *entity.CreateClientResponse] {
	return http.NewChain(
		service.CheckManagementKey[*entity.CreateClientRequest, []*entity.CreateClientScanModelItem],
		impl.service.CreateClient,
	)
}

func (impl managementApiImpl) ListClientModels() http.Chain[*entity.ListClientModelRequest, *entity.ListClientModelResponse] {
	return http.NewChain(
		service.CheckManagementKey[*entity.ListClientModelRequest, []*entity.ModelItem],
		impl.service.ListClientModels,
	)
}

func (impl managementApiImpl) CreateClientModels() http.Chain[*entity.CreateClientModelRequest, *entity.CreateClientModelResponse] {
	return http.NewChain(
		service.CheckManagementKey[*entity.CreateClientModelRequest, *entity.CreateResponse],
		impl.service.CreateClientModels,
	)
}
