package api

import (
	"github.com/alioth-center/akasha-whisper/app/entity"
	"github.com/alioth-center/akasha-whisper/app/service"
	"github.com/alioth-center/infrastructure/network/http"
	gin "github.com/gin-gonic/gin"
)

var ManagementApi managementApiImpl

type managementApiImpl struct {
	service *service.ManagementService
}

func (impl managementApiImpl) AuthorizeManagementKey() http.Chain[http.NoBody, http.NoResponse] {
	return http.NewChain(impl.service.AuthorizeManagementKey)
}

func (impl managementApiImpl) Overview() http.Chain[*entity.OverviewRequest, *entity.OverviewResponse] {
	return http.NewChain(impl.service.Overview)
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

func (impl managementApiImpl) ModifyClientBalance() http.Chain[*entity.ModifyOpenaiClientBalanceRequest, *entity.ModifyOpenaiClientBalanceResponse] {
	return http.NewChain(
		service.CheckManagementKey[*entity.ModifyOpenaiClientBalanceRequest, *entity.ModifyOpenaiClientBalanceResult],
		impl.service.ModifyClientBalance,
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

func (impl managementApiImpl) ListWhisperUsers() http.Chain[*entity.ListWhisperUsersRequest, *entity.ListWhisperUsersResponse] {
	return http.NewChain(
		service.CheckManagementKey[*entity.ListWhisperUsersRequest, []*entity.WhisperUserResult],
		impl.service.ListWhisperUsers,
	)
}

func (impl managementApiImpl) CreateWhisperUser() http.Chain[*entity.CreateWhisperUserRequest, *entity.CreateWhisperUserResponse] {
	return http.NewChain(
		service.CheckManagementKey[*entity.CreateWhisperUserRequest, *entity.WhisperUserResult],
		impl.service.CreateWhisperUser,
	)
}

func (impl managementApiImpl) GetWhisperUser() http.Chain[*entity.GetWhisperUserRequest, *entity.GetWhisperUserResponse] {
	return http.NewChain(
		service.CheckManagementKey[*entity.GetWhisperUserRequest, *entity.WhisperUserInfo],
		impl.service.GetWhisperUser,
	)
}

func (impl managementApiImpl) UpdateWhisperUser() http.Chain[*entity.UpdateWhisperUserRequest, *entity.UpdateWhisperUserResponse] {
	return http.NewChain(
		service.CheckManagementKey[*entity.UpdateWhisperUserRequest, *entity.WhisperUserResult],
		impl.service.UpdateWhisperUser,
	)
}

func (impl managementApiImpl) ListWhisperUserBalanceLogs() http.Chain[*entity.ListWhisperUserBalanceLogsRequest, *entity.ListWhisperUserBalanceLogsResponse] {
	return http.NewChain(
		service.CheckManagementKey[*entity.ListWhisperUserBalanceLogsRequest, []*entity.WhisperUserBalanceLog],
		impl.service.ListWhisperUserBalanceLogs,
	)
}

func (impl managementApiImpl) ModifyWhisperUserBalance() http.Chain[*entity.ModifyWhisperUserBalanceRequest, *entity.ModifyWhisperUserBalanceResponse] {
	return http.NewChain(
		service.CheckManagementKey[*entity.ModifyWhisperUserBalanceRequest, *entity.WhisperUserBalanceLog],
		impl.service.ModifyWhisperUserBalance,
	)
}

func (impl managementApiImpl) BatchModifyWhisperUserBalance() http.Chain[*entity.BatchModifyWhisperUserBalanceRequest, *entity.BatchModifyWhisperUserBalanceResponse] {
	return http.NewChain(
		service.CheckManagementKey[*entity.BatchModifyWhisperUserBalanceRequest, *entity.BatchModifyWhisperUserBalanceResult],
		impl.service.BatchModifyWhisperUserBalance,
	)
}

func (impl managementApiImpl) ModifyWhisperUserPermissions() http.Chain[*entity.ModifyWhisperUserPermissionRequest, *entity.ModifyWhisperUserPermissionResponse] {
	return http.NewChain(
		service.CheckManagementKey[*entity.ModifyWhisperUserPermissionRequest, *entity.ModifyWhisperUserPermissionResult],
		impl.service.ModifyWhisperUserPermissions,
	)
}

func (impl managementApiImpl) PreCheckCookie() []gin.HandlerFunc {
	return []gin.HandlerFunc{impl.service.PreCheckCookie}
}
