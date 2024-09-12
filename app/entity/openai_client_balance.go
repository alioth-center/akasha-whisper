package entity

import (
	"github.com/alioth-center/akasha-whisper/app/model"
	"github.com/alioth-center/infrastructure/network/http"
	"github.com/shopspring/decimal"
)

type ModifyOpenaiClientBalanceRequest struct {
	ChangeAmount decimal.Decimal                     `json:"change_amount" vc:"key:change_amount,required"`
	Action       model.EnumOpenaiClientBalanceAction `json:"action" vc:"key:action,required"`
	Reason       string                              `json:"reason" vc:"key:reason,required"`
}

type ModifyOpenaiClientBalanceResult struct {
	Remaining decimal.Decimal `json:"remaining"`
}

type ModifyOpenaiClientBalanceResponse = http.BaseResponse[*ModifyOpenaiClientBalanceResult]
