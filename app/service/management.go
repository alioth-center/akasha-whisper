package service

import (
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/alioth-center/akasha-whisper/app/entity"
	"github.com/alioth-center/akasha-whisper/app/global"
	"github.com/alioth-center/akasha-whisper/app/model"
	"github.com/alioth-center/infrastructure/logger"
	"github.com/alioth-center/infrastructure/network/http"
	"github.com/alioth-center/infrastructure/thirdparty/openai"
	"github.com/alioth-center/infrastructure/utils/generate"
	"github.com/alioth-center/infrastructure/utils/network"
	"github.com/alioth-center/infrastructure/utils/values"
	"github.com/shopspring/decimal"
)

type ManagementService struct{}

func NewManagementService() *ManagementService { return &ManagementService{} }

func (srv *ManagementService) AuthorizeManagementKey(ctx http.Context[http.NoBody, http.NoResponse]) {
	token := ctx.NormalHeaders().Authorization
	if !CheckManagementKeyAvailable(ctx, token) {
		ctx.SetStatusCode(http.StatusForbidden)
		return
	}

	srv.setManagementCookie(ctx)
	ctx.SetStatusCode(http.StatusOK)
}

func (srv *ManagementService) Overview(ctx http.Context[*entity.OverviewRequest, *entity.OverviewResponse]) {
	clients, listClientsErr := global.OpenaiClientDatabaseInstance.ListClients(ctx)
	if listClientsErr != nil {
		response := http.NewBaseResponse(ctx, &entity.OverviewResult{}, listClientsErr)
		ctx.SetStatusCode(http.StatusInternalServerError)
		ctx.SetResponse(&response)
		return
	}

	cutoff := time.Now().AddDate(0, 0, -14)
	clientBalanceLogs, listLogsErr := global.OpenaiClientBalanceDatabaseInstance.StatisticsClientBalance(ctx, cutoff)
	if listLogsErr != nil {
		response := http.NewBaseResponse(ctx, &entity.OverviewResult{}, listLogsErr)
		ctx.SetStatusCode(http.StatusInternalServerError)
		ctx.SetResponse(&response)
		return
	}

	result := &entity.OverviewResult{
		Clients:           make([]entity.ClientItem, len(clients)),
		ClientBalanceLogs: make([]entity.OverviewClientBalanceLog, len(clientBalanceLogs)),
	}
	for i, client := range clients {
		result.Clients[i] = entity.ClientItem{
			ID:       client.ClientID,
			Name:     client.ClientDescription,
			ApiKey:   values.SecretString(client.ClientKey, 6, 4, "*"),
			Endpoint: client.ClientEndpoint,
			Weight:   client.ClientWeight,
			Balance:  client.ClientBalance,
		}
	}
	for i, log := range clientBalanceLogs {
		result.ClientBalanceLogs[i] = entity.OverviewClientBalanceLog{
			ClientName:   log.ClientName,
			TotalRequest: log.RequestCount,
			TotalCost:    log.TotalCost,
			Date:         log.DateDay.Format("20060102"),
		}
	}

	response := http.NewBaseResponse(ctx, result, nil)
	ctx.SetStatusCode(http.StatusOK)
	ctx.SetResponse(&response)
}

func (srv *ManagementService) ListAllClients(ctx http.Context[*entity.ListClientsRequest, *entity.ListClientResponse]) {
	clients, err := global.OpenaiClientDatabaseInstance.ListClients(ctx)
	if err != nil {
		response := http.NewBaseResponse(ctx, []*entity.ClientItem{}, err)
		ctx.SetStatusCode(http.StatusInternalServerError)
		ctx.SetResponse(&response)
		return
	}

	items := make([]*entity.ClientItem, len(clients))
	for i, client := range clients {
		items[i] = &entity.ClientItem{
			ID:       client.ClientID,
			Name:     client.ClientDescription,
			ApiKey:   values.SecretString(client.ClientKey, 6, 4, "*"),
			Endpoint: client.ClientEndpoint,
			Weight:   client.ClientWeight,
			Balance:  client.ClientBalance,
		}
	}

	response := http.NewBaseResponse(ctx, items, nil)
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
		response := http.NewBaseResponse(ctx, []*entity.CreateClientScanModelItem{}, createErr)
		ctx.SetStatusCode(http.StatusInternalServerError)
		ctx.SetResponse(&response)
		return
	}
	if !created {
		response := http.NewBaseResponse(ctx, []*entity.CreateClientScanModelItem{}, http.NewBaseError(http.StatusConflict, "client already exists"))
		ctx.SetStatusCode(http.StatusConflict)
		ctx.SetResponse(&response)
		return
	}
	global.Logger.Info(logger.NewFields(ctx).WithMessage("client created").WithData(client))

	// initialize client balance
	_, initBalanceErr := global.OpenaiClientBalanceDatabaseInstance.CreateBalanceRecord(ctx, int(client.ID), decimal.Zero, model.OpenaiClientBalanceActionInitial)
	if initBalanceErr != nil {
		response := http.NewBaseResponse(ctx, []*entity.CreateClientScanModelItem{}, initBalanceErr)
		ctx.SetStatusCode(http.StatusInternalServerError)
		ctx.SetResponse(&response)
		return
	}
	global.Logger.Info(logger.NewFields(ctx).WithMessage("client balance initialized"))

	// initialize openai client
	openaiClient := openai.NewClient(openai.Config{ApiKey: client.ApiKey, BaseUrl: client.Endpoint}, global.Logger)
	models, listErr := openaiClient.ListModels(ctx, openai.ListModelRequest{})
	if listErr != nil {
		global.Logger.Error(logger.NewFields(ctx).WithMessage("failed to list models when create client").WithData(listErr))
		response := http.NewBaseResponse(ctx, []*entity.CreateClientScanModelItem{}, listErr)
		ctx.SetStatusCode(http.StatusOK)
		ctx.SetResponse(&response)
		return
	}
	global.OpenaiClientCacheInstance.Set(int(client.ID), openaiClient)
	global.Logger.Info(logger.NewFields(ctx).WithMessage("openai client initialized"))

	// list openai supported models
	items := make([]*entity.CreateClientScanModelItem, len(models.Data))
	for i, m := range models.Data {
		items[i] = &entity.CreateClientScanModelItem{ModelName: m.ID, CreatedAt: m.Created}
	}

	response := http.NewBaseResponse(ctx, items, nil)
	ctx.SetStatusCode(http.StatusOK)
	ctx.SetResponse(&response)
}

func (srv *ManagementService) ModifyClientBalance(ctx http.Context[*entity.ModifyOpenaiClientBalanceRequest, *entity.ModifyOpenaiClientBalanceResponse]) {
	client := ctx.PathParams().GetString("client_name")

	request := ctx.Request()
	// check action
	if !srv.checkBalanceChangeAction(request.Action) {
		response := http.NewBaseResponse(ctx, &entity.ModifyOpenaiClientBalanceResult{}, http.NewBaseError(http.StatusBadRequest, "invalid action"))
		ctx.SetStatusCode(http.StatusBadRequest)
		ctx.SetResponse(&response)
		return
	}

	after, modifyErr := global.OpenaiClientBalanceDatabaseInstance.CreateBalanceRecordByName(ctx, client, request.ChangeAmount, request.Action, request.Reason)
	if modifyErr != nil {
		global.Logger.Error(logger.NewFields(ctx).WithMessage("failed to modify client balance").WithData(modifyErr))
		response := http.NewBaseResponse(ctx, &entity.ModifyOpenaiClientBalanceResult{}, modifyErr)
		ctx.SetStatusCode(http.StatusInternalServerError)
		ctx.SetResponse(&response)
		return
	}

	global.Logger.Info(logger.NewFields(ctx).WithMessage("modify client balance").WithData(map[string]any{"client": client, "action": request.Action, "change": request.ChangeAmount, "reason": request.Reason}))
	result := &entity.ModifyOpenaiClientBalanceResult{Remaining: after}
	response := http.NewBaseResponse(ctx, result, nil)
	ctx.SetStatusCode(http.StatusOK)
	ctx.SetResponse(&response)
}

func (srv *ManagementService) ListClientModels(ctx http.Context[*entity.ListClientModelRequest, *entity.ListClientModelResponse]) {
	models, queryErr := global.OpenaiModelDatabaseInstance.GetModelsByClientDescription(ctx, ctx.PathParams().GetString("client_name"))
	if queryErr != nil {
		response := http.NewBaseResponse(ctx, []*entity.ModelItem{}, queryErr)
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

	response := http.NewBaseResponse(ctx, items, nil)
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
		response := http.NewBaseResponse(ctx, &entity.CreateResponse{Success: false}, insertErr)
		ctx.SetStatusCode(http.StatusInternalServerError)
		ctx.SetResponse(&response)
		return
	}

	response := http.NewBaseResponse(ctx, &entity.CreateResponse{Success: true}, nil)
	ctx.SetStatusCode(http.StatusOK)
	ctx.SetResponse(&response)

	ctx.ExtraParams()
}

func (srv *ManagementService) ListWhisperUsers(ctx http.Context[*entity.ListWhisperUsersRequest, *entity.ListWhisperUsersResponse]) {
	page, limit := ctx.QueryParams().GetInt("page"), ctx.QueryParams().GetInt("limit")
	if limit == 0 || limit > 100 {
		limit = 100
	}

	users, err := global.WhisperUserDatabaseInstance.ListWhisperUsers(ctx, page, limit)
	if err != nil {
		response := http.NewBaseResponse(ctx, []*entity.WhisperUserResult{}, err)
		ctx.SetStatusCode(http.StatusInternalServerError)
		ctx.SetResponse(&response)
		return
	}

	items := make([]*entity.WhisperUserResult, len(users))
	for i, user := range users {
		items[i] = &entity.WhisperUserResult{
			ID:       int(user.ID),
			ApiKey:   user.ApiKey,
			Email:    user.Email,
			Language: user.Language,
			AllowIPs: strings.Split(user.AllowIps, ","),
		}
	}

	response := http.NewBaseResponse(ctx, items, nil)
	ctx.SetStatusCode(http.StatusOK)
	ctx.SetResponse(&response)
}

func (srv *ManagementService) CreateWhisperUser(ctx http.Context[*entity.CreateWhisperUserRequest, *entity.CreateWhisperUserResponse]) {
	request := ctx.Request()
	switch request.Role {
	case model.WhisperUserRoleClient, model.WhisperUserRoleUser:
	default:
		response := http.NewBaseResponse(ctx, &entity.WhisperUserResult{}, http.NewBaseError(http.StatusBadRequest, "invalid role"))
		ctx.SetStatusCode(http.StatusBadRequest)
		ctx.SetResponse(&response)
		return
	}

	// create user
	user := &model.WhisperUser{
		Email:    request.Email,
		ApiKey:   generate.RandomBase62WithPrefix("aw-", 64),
		Role:     request.Role,
		Language: request.Language,
		AllowIps: strings.Join(values.FilterArray(request.AllowIPs, func(s string) bool { return network.IsValidIPOrCIDR(s) }), ","),
	}
	created, createErr := global.WhisperUserDatabaseInstance.CreateWhisperUser(ctx, user)
	if createErr != nil {
		response := http.NewBaseResponse(ctx, &entity.WhisperUserResult{}, createErr)
		ctx.SetStatusCode(http.StatusInternalServerError)
		ctx.SetResponse(&response)
		return
	}
	if !created {
		response := http.NewBaseResponse(ctx, &entity.WhisperUserResult{}, http.NewBaseError(http.StatusConflict, "user already exists"))
		ctx.SetStatusCode(http.StatusConflict)
		ctx.SetResponse(&response)
		return
	}

	// initialize user balance
	_, initBalanceErr := global.WhisperUserBalanceDatabaseInstance.CreateBalanceRecord(ctx, int(user.ID), decimal.Zero, model.WhisperUserBalanceActionInitial)
	if initBalanceErr != nil {
		response := http.NewBaseResponse(ctx, &entity.WhisperUserResult{}, initBalanceErr)
		ctx.SetStatusCode(http.StatusInternalServerError)
		ctx.SetResponse(&response)
		return
	}

	// add api-key to bloom filter
	global.BearerTokenBloomFilterInstance.AddKeys(user.ApiKey)

	result := &entity.WhisperUserResult{
		ID:       int(user.ID),
		ApiKey:   user.ApiKey,
		Email:    user.Email,
		Language: user.Language,
		AllowIPs: strings.Split(user.AllowIps, ","),
	}
	response := http.NewBaseResponse(ctx, result, nil)
	ctx.SetStatusCode(http.StatusOK)
	ctx.SetResponse(&response)
}

func (srv *ManagementService) BatchModifyWhisperUserBalance(ctx http.Context[*entity.BatchModifyWhisperUserBalanceRequest, *entity.BatchModifyWhisperUserBalanceResponse]) {
	request := ctx.Request()

	// check action
	if !srv.checkBalanceChangeAction(request.Action) {
		response := http.NewBaseResponse(ctx, &entity.BatchModifyWhisperUserBalanceResult{Success: false}, http.NewBaseError(http.StatusBadRequest, "invalid action"))
		ctx.SetStatusCode(http.StatusBadRequest)
		ctx.SetResponse(&response)
		return
	}

	// check users
	if len(request.Users) == 0 {
		response := http.NewBaseResponse(ctx, &entity.BatchModifyWhisperUserBalanceResult{Success: false}, http.NewBaseError(http.StatusBadRequest, "empty users"))
		ctx.SetStatusCode(http.StatusBadRequest)
		ctx.SetResponse(&response)
		return
	}

	// creat balance records
	if createErr := global.WhisperUserBalanceDatabaseInstance.BatchCreateBalanceRecord(ctx, request.Users, request.ChangeAmount, request.Action, request.Reason); createErr != nil {
		global.Logger.Error(logger.NewFields(ctx).WithMessage("failed to batch modify user balance").WithData(createErr))
		response := http.NewBaseResponse(ctx, &entity.BatchModifyWhisperUserBalanceResult{Success: false}, createErr)
		ctx.SetStatusCode(http.StatusInternalServerError)
		ctx.SetResponse(&response)
		return
	}

	global.Logger.Info(logger.NewFields(ctx).WithMessage("batch modify user balance").WithData(map[string]any{"users": request.Users, "action": request.Action, "change": request.ChangeAmount, "reason": request.Reason}))
	response := http.NewBaseResponse(ctx, &entity.BatchModifyWhisperUserBalanceResult{Success: true}, nil)
	ctx.SetStatusCode(http.StatusOK)
	ctx.SetResponse(&response)
}

func (srv *ManagementService) GetWhisperUser(ctx http.Context[*entity.GetWhisperUserRequest, *entity.GetWhisperUserResponse]) {
	userID := ctx.PathParams().GetInt("user_id")
	user, queryErr := global.WhisperUserDatabaseInstance.GetWhisperUserInfo(ctx, userID)
	if queryErr != nil {
		response := http.NewBaseResponse(ctx, &entity.WhisperUserInfo{}, queryErr)
		ctx.SetStatusCode(http.StatusInternalServerError)
		ctx.SetResponse(&response)
		return
	}

	result := &entity.WhisperUserInfo{
		ID:              user.UserInfo.ID,
		Email:           user.UserInfo.Email,
		ApiKey:          user.UserInfo.ApiKey,
		Role:            user.UserInfo.Role,
		Language:        user.UserInfo.Language,
		Balance:         user.UserInfo.Balance,
		AvailableModels: user.Models,
		UpdatedAt:       user.UserInfo.UpdatedAt.Format(time.RFC3339),
		AllowIPs:        strings.Split(user.UserInfo.AllowIps, ","),
	}

	response := http.NewBaseResponse(ctx, result, nil)
	ctx.SetStatusCode(http.StatusOK)
	ctx.SetResponse(&response)
}

func (srv *ManagementService) UpdateWhisperUser(ctx http.Context[*entity.UpdateWhisperUserRequest, *entity.CreateWhisperUserResponse]) {
	request, userID := ctx.Request(), ctx.PathParams().GetInt("user_id")

	// update user
	user := &model.WhisperUser{
		ID:       int64(userID),
		Email:    request.Email,
		Language: request.Language,
		AllowIps: strings.Join(request.AllowIPs, ","),
	}
	if request.RefreshApiToken {
		user.ApiKey = generate.RandomBase62WithPrefix("aw-", 64)
	}

	updateErr := global.WhisperUserDatabaseInstance.UpdateWhisperUser(ctx, user)
	if updateErr != nil {
		response := http.NewBaseResponse(ctx, &entity.WhisperUserResult{}, updateErr)
		ctx.SetStatusCode(http.StatusInternalServerError)
		ctx.SetResponse(&response)
		return
	}

	// add api-key to bloom filter
	if request.RefreshApiToken {
		global.BearerTokenBloomFilterInstance.AddKeys(user.ApiKey)
	}

	result := &entity.WhisperUserResult{
		ID:       int(user.ID),
		ApiKey:   user.ApiKey,
		Email:    user.Email,
		Language: user.Language,
		AllowIPs: request.AllowIPs,
	}
	response := http.NewBaseResponse(ctx, result, nil)
	ctx.SetStatusCode(http.StatusOK)
	ctx.SetResponse(&response)
}

func (srv *ManagementService) ListWhisperUserBalanceLogs(ctx http.Context[*entity.ListWhisperUserBalanceLogsRequest, *entity.ListWhisperUserBalanceLogsResponse]) {
	user := ctx.PathParams().GetInt("user_id")

	page, offset := ctx.QueryParams().GetInt("page"), ctx.QueryParams().GetInt("offset")
	if page == 0 || page > 100 {
		page = 100
	}

	startStr, endStr := ctx.QueryParams().GetString("start"), ctx.QueryParams().GetString("end")
	start, parseErr := strconv.ParseInt(startStr, 10, 64)
	if parseErr != nil || start == 0 {
		start = time.Now().AddDate(0, 0, -7).UnixMilli()
	}
	end, parseErr := strconv.ParseInt(endStr, 10, 64)
	if parseErr != nil || end == 0 {
		end = time.Now().UnixMilli()
	}

	global.Logger.Info(logger.NewFields(ctx).WithMessage("list user balance params parsed").WithData(map[string]any{"start": time.UnixMilli(start).String(), "end": time.UnixMilli(end).String(), "page": page, "offset": offset}))
	logs, queryErr := global.WhisperUserBalanceDatabaseInstance.ListBalanceRecords(ctx, user, time.UnixMilli(start), time.UnixMilli(end), page, offset)
	if queryErr != nil {
		global.Logger.Error(logger.NewFields(ctx).WithMessage("failed to list user balance logs").WithData(queryErr))
		response := http.NewBaseResponse(ctx, []*entity.WhisperUserBalanceLog{}, queryErr)
		ctx.SetStatusCode(http.StatusInternalServerError)
		ctx.SetResponse(&response)
		return
	}

	items := make([]*entity.WhisperUserBalanceLog, len(logs))
	for i, log := range logs {
		items[i] = &entity.WhisperUserBalanceLog{
			ID:           int(log.ID),
			ChangeAmount: log.BalanceChangeAmount,
			Remaining:    log.BalanceRemaining,
			Action:       log.Action,
			Reason:       log.Reason,
			CreatedAt:    log.CreatedAt.Format(time.RFC3339),
		}
	}

	global.Logger.Info(logger.NewFields(ctx).WithMessage("list user balance logs").WithData(map[string]any{"count": len(items)}))
	response := http.NewBaseResponse(ctx, items, nil)
	ctx.SetStatusCode(http.StatusOK)
	ctx.SetResponse(&response)
}

func (srv *ManagementService) ModifyWhisperUserBalance(ctx http.Context[*entity.ModifyWhisperUserBalanceRequest, *entity.ModifyWhisperUserBalanceResponse]) {
	// check user id
	user := ctx.PathParams().GetInt("user_id")
	if user == 0 {
		response := http.NewBaseResponse(ctx, &entity.WhisperUserBalanceLog{}, http.NewBaseError(http.StatusBadRequest, "invalid user"))
		ctx.SetStatusCode(http.StatusBadRequest)
		ctx.SetResponse(&response)
		return
	}

	request := ctx.Request()
	// check action
	if !srv.checkBalanceChangeAction(request.Action) {
		response := http.NewBaseResponse(ctx, &entity.WhisperUserBalanceLog{}, http.NewBaseError(http.StatusBadRequest, "invalid action"))
		ctx.SetStatusCode(http.StatusBadRequest)
		ctx.SetResponse(&response)
		return
	}

	after, modifyErr := global.WhisperUserBalanceDatabaseInstance.CreateBalanceRecord(ctx, user, request.ChangeAmount, request.Action, request.Reason)
	if modifyErr != nil {
		global.Logger.Error(logger.NewFields(ctx).WithMessage("failed to modify user balance").WithData(modifyErr))
		response := http.NewBaseResponse(ctx, &entity.WhisperUserBalanceLog{}, modifyErr)
		ctx.SetStatusCode(http.StatusInternalServerError)
		ctx.SetResponse(&response)
		return
	}

	log := &entity.WhisperUserBalanceLog{
		ChangeAmount: request.ChangeAmount,
		Remaining:    after,
		Action:       request.Action,
		Reason:       request.Reason,
		CreatedAt:    time.Now().Format(time.RFC3339),
	}

	global.Logger.Info(logger.NewFields(ctx).WithMessage("modify user balance").WithData(log))
	response := http.NewBaseResponse(ctx, log, nil)
	ctx.SetStatusCode(http.StatusOK)
	ctx.SetResponse(&response)
}

func (srv *ManagementService) ModifyWhisperUserPermissions(ctx http.Context[*entity.ModifyWhisperUserPermissionRequest, *entity.ModifyWhisperUserPermissionResponse]) {
	user := ctx.PathParams().GetInt("user_id")
	request := ctx.Request()

	modifyMap := map[string][]string{}
	for _, permission := range request.Permissions {
		modifyMap[permission.ClientName] = permission.Models
	}

	global.Logger.Info(logger.NewFields(ctx).WithMessage("sync user permissions").WithData(map[string]any{"user": user, "permissions": modifyMap}))
	syncErr := global.WhisperUserPermissionDatabaseInstance.SyncPermissions(ctx, user, modifyMap)
	if syncErr != nil {
		global.Logger.Error(logger.NewFields(ctx).WithMessage("failed to sync user permissions").WithData(syncErr))
		response := http.NewBaseResponse(ctx, &entity.ModifyWhisperUserPermissionResult{}, syncErr)
		ctx.SetStatusCode(http.StatusInternalServerError)
		ctx.SetResponse(&response)
		return
	}

	global.Logger.Info(logger.NewFields(ctx).WithMessage("sync user permissions success").WithData(map[string]any{"user": user}))
	response := http.NewBaseResponse(ctx, &entity.ModifyWhisperUserPermissionResult{Success: true}, nil)
	ctx.SetStatusCode(http.StatusOK)
	ctx.SetResponse(&response)
}

func (srv *ManagementService) checkBalanceChangeAction(action string) bool {
	switch action {
	case model.WhisperUserBalanceActionConsumption, model.WhisperUserBalanceActionRecharge, model.WhisperUserBalanceActionGift,
		model.WhisperUserBalanceActionSpecial, model.WhisperUserBalanceActionInitial:
		return true
	default:
		return false
	}
}

func CheckManagementKey[req any, res any](ctx http.Context[req, *http.BaseResponse[res]]) {
	token := ctx.NormalHeaders().Authorization
	if !CheckManagementKeyAvailable(ctx, token) {
		response := http.NewBaseResponse(
			ctx, values.Nil[res](), http.NewBaseError(http.StatusForbidden, "unauthorized"),
		)
		ctx.SetStatusCode(http.StatusForbidden)
		ctx.SetResponse(&response)
		ctx.Abort()
	}
}

func (srv *ManagementService) setManagementCookie(ctx http.Context[http.NoBody, http.NoResponse]) {
	key := generate.RandomBase62(64)
	for exist, _ := global.LoginCookieCacheInstance.ExistKey(ctx, values.BuildStrings("login_token:", key)); exist; key = generate.RandomBase62(64) {
	}

	_ = global.LoginCookieCacheInstance.StoreJsonEX(ctx, values.BuildStrings("login_token:", key), &entity.LoginToken{IP: ctx.ClientIP(), CreatedAt: time.Now()}, time.Hour)
	ctx.SetResponseSetCookie(http.NewBasicCookie(global.Config.App.LoginTokenKey, key))
}

func (srv *ManagementService) PreCheckCookie(ctx *gin.Context) {
	if cookie, _ := ctx.Cookie(global.Config.App.LoginTokenKey); cookie != "" {
		loginToken := &entity.LoginToken{}
		if exist, err := global.LoginCookieCacheInstance.LoadJson(ctx, values.BuildStrings("login_token:", cookie), loginToken); err == nil && exist {
			if loginToken.IP == ctx.ClientIP() {
				_ = global.LoginCookieCacheInstance.StoreJsonEX(ctx, values.BuildStrings("login_token:", cookie), loginToken, time.Hour)

				// write back to request header to pass header params check
				ctx.Request.Header.Set(http.HeaderAuthorization, "Bearer "+global.Config.App.ManagementToken)
				return
			}
		}

		// unauthorized, login token invalid, clear cookie
		cookie := http.NewBasicCookie(global.Config.App.LoginTokenKey, "")
		ctx.SetCookie(cookie.Name, cookie.Value, cookie.MaxAge, cookie.Path, cookie.Domain, cookie.Secure, cookie.HttpOnly)
	}
}
