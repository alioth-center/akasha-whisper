package common

import (
	"github.com/alioth-center/akasha-whisper/model"
	"github.com/alioth-center/infrastructure/network/http"
	"github.com/alioth-center/infrastructure/trace"
)

func ErrorUserNotFound[req any](ctx http.Context[req, *model.BaseResponse]) {
	ctx.SetStatusCode(http.StatusNotFound)
	ctx.SetResponse(&model.BaseResponse{
		ErrorCode:    http.ErrorCodeResourceNotFound,
		ErrorMessage: "user not found",
		RequestID:    trace.GetTid(ctx),
	})
}

func ErrorUserAlreadyExists[req any](ctx http.Context[req, *model.BaseResponse]) {
	ctx.SetStatusCode(http.StatusConflict)
	ctx.SetResponse(&model.BaseResponse{
		ErrorCode:    model.ErrorCodeResourceConflict,
		ErrorMessage: "user already exists",
		RequestID:    trace.GetTid(ctx),
	})
}

func ErrorProviderNotFound[req any](ctx http.Context[req, *model.BaseResponse]) {
	ctx.SetStatusCode(http.StatusNotFound)
	ctx.SetResponse(&model.BaseResponse{
		ErrorCode:    http.ErrorCodeResourceNotFound,
		ErrorMessage: "provider not found",
		RequestID:    trace.GetTid(ctx),
	})
}

func ErrorModelNotFound[req any](ctx http.Context[req, *model.BaseResponse]) {
	ctx.SetStatusCode(http.StatusNotFound)
	ctx.SetResponse(&model.BaseResponse{
		ErrorCode:    http.ErrorCodeResourceNotFound,
		ErrorMessage: "model not found",
		RequestID:    trace.GetTid(ctx),
	})
}
