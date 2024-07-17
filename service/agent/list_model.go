package agent

import (
	"github.com/alioth-center/infrastructure/network/http"
	"github.com/alioth-center/infrastructure/thirdparty/openai"
	"github.com/alioth-center/infrastructure/utils/values"
	"strings"
)

func (srv *akashaAgentSrvImpl) ListModel(ctx http.Context[*openai.ListModelRequest, *openai.ListModelResponseBody]) {
	token := strings.TrimPrefix(ctx.HeaderParams().GetString("Authorization"), "Bearer ")
	if token == "" {
		// empty token, return empty message
		ctx.SetStatusCode(401)
		ctx.SetResponse(&openai.ListModelResponseBody{Object: "list", Data: []openai.ModelObject{}})
		return
	}

	userItem, checkErr := srv.db.GetUserByApiKey(ctx, token)
	if checkErr != nil {
		// check error, return empty message
		ctx.SetStatusCode(401)
		ctx.SetResponse(&openai.ListModelResponseBody{Object: "list", Data: []openai.ModelObject{}})
		return
	}

	// return the models
	var (
		modelMapping = map[string]struct{}{}
		modelObjects []openai.ModelObject
	)

	// query user/system owned models
	if userItem.Role != "user" {
		// system role, return client owned models
		models, queryErr := srv.db.GetAvailableModelsByApiKey(ctx, token)
		if queryErr != nil {
			// query error, return empty message
			ctx.SetStatusCode(401)
			ctx.SetResponse(&openai.ListModelResponseBody{Object: "list", Data: []openai.ModelObject{}})
			return
		}

		// unique the model with model name
		for _, modelItem := range models {
			if _, exist := modelMapping[modelItem.ModelName]; !exist {
				modelMapping[modelItem.ModelName] = struct{}{}
				modelObjects = append(modelObjects, openai.ModelObject{
					ID:      modelItem.ModelName,
					Created: modelItem.LastUpdatedAt.Unix(),
					Object:  "model",
					OwnedBy: "openai",
				})
			}
		}
	} else {
		// user role, return user permitted models
		models, queryErr := srv.db.GetAvailableModelsByEmail(ctx, userItem.Email)
		if queryErr != nil {
			// query error, return empty message
			ctx.SetStatusCode(401)
			ctx.SetResponse(&openai.ListModelResponseBody{Object: "list", Data: []openai.ModelObject{}})
			return
		}

		// unique the model with model name
		for _, modelItem := range models {
			if _, exist := modelMapping[modelItem.ModelName]; !exist {
				modelMapping[modelItem.ModelName] = struct{}{}
				modelObjects = append(modelObjects, openai.ModelObject{
					ID:      modelItem.ModelName,
					Created: modelItem.LastUpdatedAt.Unix(),
					Object:  "model",
					OwnedBy: "openai",
				})
			}
		}
	}

	// sort the model by created time and return
	ctx.SetStatusCode(200)
	ctx.SetResponse(&openai.ListModelResponseBody{
		Object: "list",
		Data:   values.SortArray(modelObjects, func(a, b openai.ModelObject) bool { return a.Created > b.Created }),
	})
}
