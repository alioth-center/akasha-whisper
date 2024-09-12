package entity

import "github.com/alioth-center/infrastructure/network/http"

type ListWhisperUserPermissionsRequest = http.NoBody

type ListWhisperUserPermissionsResponse = http.BaseResponse[[]*WhisperUserClientPermission]

type WhisperUserClientPermission struct {
	ClientID   int                          `json:"client_id"`
	ClientName string                       `json:"client_name"`
	Models     []WhisperUserModelPermission `json:"models"`
}

type WhisperUserModelPermission struct {
	ModelID   int    `json:"model_id"`
	ModelName string `json:"model_name"`
}

type ModifyWhisperUserPermissionRequest struct {
	Permissions []ModifyWhisperUserClientPermissions `json:"permissions" vc:"key:permissions,required"`
}

type ModifyWhisperUserClientPermissions struct {
	ClientName string   `json:"client_name" vc:"key:client_name,required"`
	Models     []string `json:"models" vc:"key:models,required"`
}

type ModifyWhisperUserPermissionResponse = http.BaseResponse[*ModifyWhisperUserPermissionResult]

type ModifyWhisperUserPermissionResult struct {
	Success bool `json:"success"`
}
