package model

type PermissionGroup struct {
	Provider string   `json:"provider"`
	Models   []string `json:"models"`
}

type PermissionObject struct {
	Model           string  `json:"model"`
	MaxToken        int     `json:"max_token"`
	PromptPrice     float64 `json:"prompt_price"`
	CompletionPrice float64 `json:"completion_price"`
	LastUpdatedAt   int64   `json:"last_updated_at"`
}

type CreateUserRequest struct {
	Email       string            `json:"email" vc:"key:email,required"`
	Role        string            `json:"role"`
	AllowIPs    []string          `json:"allow_ips"`
	Permissions []PermissionGroup `json:"permissions,omitempty"`
}

type CreateUserResponse struct {
	Role   string `json:"role"`
	Email  string `json:"email"`
	ApiKey string `json:"api_key"`
}

type ModifyUserPermissionsRequest struct {
	UserID      int               `json:"user_id"`
	Permissions []PermissionGroup `json:"permissions"`
}

type ModifyUserPermissionsResponse struct {
	Added   []PermissionGroup `json:"added,omitempty"`
	Removed []PermissionGroup `json:"removed,omitempty"`
}

type GetUserResponse struct {
	UserID      int                `json:"user_id"`
	Email       string             `json:"email"`
	Role        string             `json:"role"`
	ApiKey      string             `json:"api_key"`
	Balance     float64            `json:"balance"`
	Permissions []PermissionObject `json:"permissions"`
	AllowIPs    []string           `json:"allow_ips,omitempty"`
}

type ListUserResponse struct {
	Offset int               `json:"offset"`
	Users  []GetUserResponse `json:"users"`
}

type BaseResponse struct {
	ErrorCode    int    `json:"error_code"`
	ErrorMessage string `json:"error_message,omitempty"`
	RequestID    string `json:"request_id,omitempty"`
	Data         any    `json:"data,omitempty"`
}

type ModifyUserBalanceRequest struct {
	UserID int     `json:"user_id"`
	Amount float64 `json:"amount"`
}

type ModifyUserBalanceResponse struct {
	UserID  int     `json:"user_id"`
	Balance float64 `json:"balance"`
}

type GetUserBalanceResponse struct {
	UserID        int     `json:"user_id"`
	Balance       float64 `json:"balance"`
	LastUpdatedAt int64   `json:"last_updated_at"`
}
