package service

import (
	"github.com/alioth-center/akasha-whisper/dao"
	"github.com/alioth-center/akasha-whisper/global"
	"github.com/alioth-center/akasha-whisper/model"
	"github.com/alioth-center/infrastructure/network/http"
	"github.com/alioth-center/infrastructure/trace"
	"strings"
)

// CheckRequestHasAdminPermission creates a middleware handler to check if the request has admin-level permissions.
// This function leverages type parameters to work with any request type, returning a middleware handler
// that intercepts HTTP requests to perform authorization checks. It validates the API key provided in the
// Authorization header against a predefined admin token stored in the global configuration.
//
// The middleware performs the following checks:
//  1. Extracts the API key from the Authorization header, removing the "Bearer " prefix. If the API key is missing,
//     it sets the response status code to http.StatusUnauthorized and provides an error message indicating an
//     empty API key. The request processing is then aborted.
//  2. Compares the extracted API key with the global admin token. If they do not match, it sets the response status
//     code to http.StatusUnauthorized and provides an error message indicating permission denied. The request
//     processing is then aborted.
//
// Parameters:
//   - req: The request type parameter, allowing this function to be used with any request type.
//
// Returns:
//   - A middleware handler function that takes an http.Context[req, *model.BaseResponse] and performs the
//     authorization checks described above.
func CheckRequestHasAdminPermission[req any]() http.Handler[req, *model.BaseResponse] {
	return func(ctx http.Context[req, *model.BaseResponse]) {
		token := strings.TrimPrefix(ctx.HeaderParams().GetString("Authorization"), "Bearer ")
		switch {
		case token == "":
			ctx.SetStatusCode(http.StatusUnauthorized)
			ctx.SetResponse(&model.BaseResponse{
				ErrorCode:    model.ErrorCodeUnauthorized,
				ErrorMessage: "empty api key",
				RequestID:    trace.GetTid(ctx),
			})
			ctx.Abort()
			return
		case token != global.Config.AdminToken:
			ctx.SetStatusCode(http.StatusUnauthorized)
			ctx.SetResponse(&model.BaseResponse{
				ErrorCode:    model.ErrorCodeUnauthorized,
				ErrorMessage: "permission denied",
				RequestID:    trace.GetTid(ctx),
			})
			ctx.Abort()
			return
		}
	}
}

// CheckRequestHasSystemPermission creates a middleware handler to check if the request has system-level permissions.
// This function is a generic function that can work with any request type. It returns a middleware handler
// that can be used to intercept HTTP requests and perform authorization checks based on the API key provided
// in the Authorization header. The middleware ensures that the API key belongs to a user with a "system" role.
// If the database connection is not established, it responds with an internal server error. If the API key is
// missing or invalid, or if the user does not have the "system" role, it responds with an unauthorized error.
//
// The middleware uses the following steps:
//  1. Checks if the global database connection is nil. If so, it sets the response status code to
//     http.StatusInternalServerError and provides an error message indicating an invalid database connection.
//     The request processing is then aborted.
//  2. Extracts the API key from the Authorization header. If the API key is missing, it sets the response status
//     code to http.StatusUnauthorized and provides an error message indicating an empty API key. The request
//     processing is then aborted.
//  3. Attempts to retrieve the user associated with the API key from the database. If the API key is invalid
//     or the user cannot be retrieved, it sets the response status code to http.StatusUnauthorized and provides
//     an error message indicating an invalid API key. The request processing is then aborted.
//  4. Checks if the retrieved user has a role of "system". If not, it sets the response status code to
//     http.StatusUnauthorized and provides an error message indicating permission denied. The request processing
//     is then aborted.
//
// Parameters:
//   - req: The request type parameter, allowing this function to be used with any request type.
//
// Returns:
//   - A middleware handler function that takes an http.Context[req, *model.BaseResponse] and performs the
//     authorization checks described above.
func CheckRequestHasSystemPermission[req any]() http.Handler[req, *model.BaseResponse] {
	return func(ctx http.Context[req, *model.BaseResponse]) {
		if global.Database == nil {
			ctx.SetStatusCode(http.StatusInternalServerError)
			ctx.SetResponse(&model.BaseResponse{
				ErrorCode:    http.ErrorCodeInternalErrorOccurred,
				ErrorMessage: "database connection invalid",
				RequestID:    trace.GetTid(ctx),
			})
			ctx.Abort()
			return
		}

		token := strings.TrimPrefix(ctx.HeaderParams().GetString("Authorization"), "Bearer ")
		if token == "" {
			ctx.SetStatusCode(http.StatusUnauthorized)
			ctx.SetResponse(&model.BaseResponse{
				ErrorCode:    http.ErrorCodeMissingRequiredHeader,
				ErrorMessage: "empty api key",
				RequestID:    trace.GetTid(ctx),
			})
			ctx.Abort()
			return
		}

		user, queryErr := dao.NewDatabaseAccessor(global.Database).GetUserByApiKey(ctx, token)
		if queryErr != nil {
			ctx.SetStatusCode(http.StatusUnauthorized)
			ctx.SetResponse(&model.BaseResponse{
				ErrorCode:    model.ErrorCodeUnauthorized,
				ErrorMessage: "invalid api key",
				RequestID:    trace.GetTid(ctx),
			})
			ctx.Abort()
			return
		}

		if user.Role != "system" {
			ctx.SetStatusCode(http.StatusUnauthorized)
			ctx.SetResponse(&model.BaseResponse{
				ErrorCode:    model.ErrorCodeUnauthorized,
				ErrorMessage: "permission denied",
				RequestID:    trace.GetTid(ctx),
			})
			ctx.Abort()
			return
		}
	}
}

// CheckRequestHasUserPermission creates a middleware handler to check if the request has user-level permissions.
// This function is designed to be generic, allowing it to work with any request type. It returns a middleware
// handler that intercepts HTTP requests to perform authorization checks based on the API key provided in the
// Authorization header. The middleware ensures that the API key belongs to a user with a "user" role. If the
// database connection is not established, it responds with an internal server error. If the API key is missing
// or invalid, or if the user does not have the "user" role, it responds with an unauthorized error.
//
// The middleware performs the following checks:
//  1. Verifies the global database connection is not nil. If it is, it sets the response status code to
//     http.StatusInternalServerError and provides an error message indicating an invalid database connection.
//     The request processing is then aborted.
//  2. Extracts the API key from the Authorization header. If the API key is missing, it sets the response status
//     code to http.StatusUnauthorized and provides an error message indicating an empty API key. The request
//     processing is then aborted.
//  3. Attempts to retrieve the user associated with the API key from the database. If the API key is invalid
//     or the user cannot be retrieved, it sets the response status code to http.StatusUnauthorized and provides
//     an error message indicating an invalid API key. The request processing is then aborted.
//  4. Checks if the retrieved user has a role of "user". If not, it sets the response status code to
//     http.StatusUnauthorized and provides an error message indicating permission denied. The request processing
//     is then aborted.
//
// Parameters:
//   - req: The request type parameter, allowing this function to be used with any request type.
//
// Returns:
//   - A middleware handler function that takes an http.Context[req, *model.BaseResponse] and performs the
//     authorization checks described above.
func CheckRequestHasUserPermission[req any]() http.Handler[req, *model.BaseResponse] {
	return func(ctx http.Context[req, *model.BaseResponse]) {
		if global.Database == nil {
			ctx.SetStatusCode(http.StatusInternalServerError)
			ctx.SetResponse(&model.BaseResponse{
				ErrorCode:    http.ErrorCodeInternalErrorOccurred,
				ErrorMessage: "database connection invalid",
				RequestID:    trace.GetTid(ctx),
			})
			ctx.Abort()
			return
		}

		token := strings.TrimPrefix(ctx.HeaderParams().GetString("Authorization"), "Bearer ")
		if token == "" {
			ctx.SetStatusCode(http.StatusUnauthorized)
			ctx.SetResponse(&model.BaseResponse{
				ErrorCode:    http.ErrorCodeMissingRequiredHeader,
				ErrorMessage: "empty api key",
				RequestID:    trace.GetTid(ctx),
			})
			ctx.Abort()
			return
		}

		user, queryErr := dao.NewDatabaseAccessor(global.Database).GetUserByApiKey(ctx, token)
		if queryErr != nil {
			ctx.SetStatusCode(http.StatusUnauthorized)
			ctx.SetResponse(&model.BaseResponse{
				ErrorCode:    model.ErrorCodeUnauthorized,
				ErrorMessage: "invalid api key",
				RequestID:    trace.GetTid(ctx),
			})
			ctx.Abort()
			return
		}

		if user.Role != "user" {
			ctx.SetStatusCode(http.StatusUnauthorized)
			ctx.SetResponse(&model.BaseResponse{
				ErrorCode:    model.ErrorCodeUnauthorized,
				ErrorMessage: "permission denied",
				RequestID:    trace.GetTid(ctx),
			})
			ctx.Abort()
			return
		}
	}
}
