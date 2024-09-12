package entity

import (
	"github.com/alioth-center/infrastructure/network/http"
	"github.com/shopspring/decimal"
)

type ListWhisperUsersRequest = http.NoBody

type ListWhisperUsersResponse = http.BaseResponse[[]*WhisperUserResult]

type GetWhisperUserRequest = http.NoBody

type GetWhisperUserResponse = http.BaseResponse[*WhisperUserInfo]

type CreateWhisperUserRequest struct {
	Email    string   `json:"email" vc:"key:email,required"`
	Language string   `json:"language,omitempty" vc:"key:language"`
	AllowIPs []string `json:"allow_ips,omitempty" vc:"key:allow_ips"`
	Role     string   `json:"role,omitempty" vc:"key:role"`
}

type CreateWhisperUserResponse = http.BaseResponse[*WhisperUserResult]

type UpdateWhisperUserRequest struct {
	Email           string   `json:"email,omitempty" vc:"key:email,required"`
	Language        string   `json:"language,omitempty" vc:"key:language,required"`
	AllowIPs        []string `json:"allow_ips,omitempty" vc:"key:allow_ips"`
	RefreshApiToken bool     `json:"refresh_api_token,omitempty" vc:"key:refresh_api_token"`
}

type UpdateWhisperUserResponse = http.BaseResponse[*WhisperUserResult]

type WhisperUserResult struct {
	ID       int      `json:"id"`
	ApiKey   string   `json:"api_key"`
	Email    string   `json:"email"`
	Language string   `json:"language"`
	AllowIPs []string `json:"allow_ips"`
}

type WhisperUserInfo struct {
	ID              int             `json:"id"`
	Email           string          `json:"email"`
	ApiKey          string          `json:"api_key"`
	Role            string          `json:"role"`
	Language        string          `json:"language"`
	Balance         decimal.Decimal `json:"balance"`
	AvailableModels []string        `json:"available_models"`
	UpdatedAt       string          `json:"updated_at"`
	AllowIPs        []string        `json:"allow_ips,omitempty"`
}
