package entity

import "github.com/alioth-center/infrastructure/network/http"

type ListOpenaiModelRequest = http.NoBody

type ListOpenaiModelResponse = http.BaseResponse[[]*ModelItem]

type ModelItem struct{}
