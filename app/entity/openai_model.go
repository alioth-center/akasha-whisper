package entity

import (
	"github.com/alioth-center/infrastructure/network/http"
	"github.com/shopspring/decimal"
)

type ListClientModelRequest = http.NoBody

type ListClientModelResponse = http.BaseResponse[[]*ModelItem]

type ModelItem struct {
	ID              int             `json:"id"`
	Name            string          `json:"name"`
	MaxTokens       int             `json:"max_tokens"`
	RpmLimit        int             `json:"rpm_limit"`
	TpmLimit        int             `json:"tpm_limit"`
	PromptPrice     decimal.Decimal `json:"prompt_price"`
	CompletionPrice decimal.Decimal `json:"completion_price"`
	LastUpdatedAt   int64           `json:"last_updated_at"`
}

type CreateClientModelRequest struct {
	Models []CreateClientModelItem `json:"models" vc:"key:models,required"`
}

type CreateClientModelItem struct {
	Name            string          `json:"name" vc:"key:name,required"`
	MaxTokens       int             `json:"max_tokens" vc:"key:max_tokens,required"`
	PromptPrice     decimal.Decimal `json:"prompt_price" vc:"key:prompt_price,required"`
	CompletionPrice decimal.Decimal `json:"completion_price" vc:"key:completion_price,required"`
	RpmLimit        int             `json:"rpm_limit,omitempty"`
	TpmLimit        int             `json:"tpm_limit,omitempty"`
}

type CreateResponse struct {
	Success bool `json:"success"`
}

type CreateClientModelResponse = http.BaseResponse[*CreateResponse]
