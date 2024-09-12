package entity

import (
	"github.com/alioth-center/infrastructure/network/http"
	"github.com/shopspring/decimal"
)

type ListWhisperUserBalanceLogsRequest = http.NoBody

type ListWhisperUserBalanceLogsResponse = http.BaseResponse[[]*WhisperUserBalanceLog]

type WhisperUserBalanceLog struct {
	ID           int             `json:"id,omitempty"`
	ChangeAmount decimal.Decimal `json:"change_amount"`
	Remaining    decimal.Decimal `json:"remaining"`
	Action       string          `json:"action"`
	Reason       string          `json:"reason"`
	CreatedAt    string          `json:"created_at"`
}

type ModifyWhisperUserBalanceRequest struct {
	ChangeAmount decimal.Decimal `json:"change_amount" vc:"key:change_amount,required"`
	Action       string          `json:"action" vc:"key:action,required"`
	Reason       string          `json:"reason" vc:"key:reason,required"`
}

type ModifyWhisperUserBalanceResponse = http.BaseResponse[*WhisperUserBalanceLog]

type BatchModifyWhisperUserBalanceRequest struct {
	Users        []int           `json:"users" vc:"key:users,required"`
	ChangeAmount decimal.Decimal `json:"change_amount" vc:"key:change_amount,required"`
	Action       string          `json:"action" vc:"key:action,required"`
	Reason       string          `json:"reason" vc:"key:reason,required"`
}

type BatchModifyWhisperUserBalanceResult struct {
	Success bool `json:"success"`
}

type BatchModifyWhisperUserBalanceResponse = http.BaseResponse[*BatchModifyWhisperUserBalanceResult]
