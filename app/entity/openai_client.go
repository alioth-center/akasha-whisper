package entity

import (
	"github.com/alioth-center/infrastructure/network/http"
	"github.com/shopspring/decimal"
)

type ListClientsRequest = http.NoBody

type ListClientResponse = http.BaseResponse[[]*ClientItem]

type ClientItem struct {
	ID       int             `json:"id"`
	Name     string          `json:"name"`
	ApiKey   string          `json:"api_key"`
	Endpoint string          `json:"endpoint"`
	Weight   int             `json:"weight"`
	Balance  decimal.Decimal `json:"balance"`
}

type CreateClientRequest struct {
	Name     string `json:"name"`
	ApiKey   string `json:"api_key"`
	Endpoint string `json:"endpoint"`
	Weight   int    `json:"weight"`
}

type CreateClientResponse = http.BaseResponse[[]*CreateClientModelItem]

type CreateClientModelItem struct {
	ModelName string `json:"model_name"`
	CreatedAt int64  `json:"created_at"`
}
