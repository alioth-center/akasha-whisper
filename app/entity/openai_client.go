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
	Name     string `json:"name" vc:"key:name,required"`
	ApiKey   string `json:"api_key" vc:"key:api_key,required"`
	Endpoint string `json:"endpoint" vc:"key:endpoint,required"`
	Weight   int    `json:"weight" vc:"key:weight,required"`
}

type CreateClientResponse = http.BaseResponse[[]*CreateClientScanModelItem]

type CreateClientScanModelItem struct {
	ModelName string `json:"model_name"`
	CreatedAt int64  `json:"created_at"`
}
