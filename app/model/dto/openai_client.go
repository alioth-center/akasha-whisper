package dto

import (
	"strconv"
	"strings"

	"github.com/alioth-center/infrastructure/utils/values"
	"github.com/shopspring/decimal"
)

type ClientCheckDTO struct {
	ID   int    `gorm:"column:id"`
	Name string `gorm:"column:description"`
}

type GetAvailableClientCND struct {
	UserApiKey string
	ModelName  string
}

func (cnd *GetAvailableClientCND) ParseTemplate(tmpl string) string {
	condition := map[string]string{
		"user_api_key": cnd.UserApiKey,
		"model_name":   cnd.ModelName,
	}

	return values.NewRawSqlTemplateWithMap(tmpl, condition).Parse()
}

type GetClientSecretCND struct {
	ClientIDs []int
}

func (cnd *GetClientSecretCND) ParseTemplate(tmpl string) string {
	conditions := make([]string, 0, len(cnd.ClientIDs))
	for _, clientID := range cnd.ClientIDs {
		conditions = append(conditions, strconv.Itoa(clientID))
	}

	condition := map[string]string{
		"openai_client_id": strings.Join(conditions, ","),
	}

	return values.NewRawSqlTemplateWithMap(tmpl, condition).Parse()
}

type AvailableClientDTO struct {
	ClientID             int             `gorm:"column:client_id"`
	ClientWeight         int64           `gorm:"column:client_weight"`
	ClientBalance        decimal.Decimal `gorm:"column:client_balance"`
	UserID               int             `gorm:"column:user_id"`
	UserBalance          decimal.Decimal `gorm:"column:user_balance"`
	UserRole             string          `gorm:"column:user_role"`
	ModelID              int             `gorm:"column:model_id"`
	ModelName            string          `gorm:"column:model_name"`
	ModelMaxToken        int             `gorm:"column:model_max_token"`
	ModelPromptPrice     decimal.Decimal `gorm:"column:model_prompt_price"`
	ModelCompletionPrice decimal.Decimal `gorm:"column:model_completion_price"`
}

type ClientSecretDTO struct {
	ClientID       int             `gorm:"column:id"`
	ClientKey      string          `gorm:"column:api_key"`
	ClientEndpoint string          `gorm:"column:endpoint"`
	ClientWeight   int             `gorm:"column:weight"`
	ClientBalance  decimal.Decimal `gorm:"column:balance"`
}

type ListClientDTO struct {
	ClientID          int             `gorm:"column:id"`
	ClientDescription string          `gorm:"column:description"`
	ClientKey         string          `gorm:"column:api_key"`
	ClientEndpoint    string          `gorm:"column:endpoint"`
	ClientWeight      int             `gorm:"column:weight"`
	ClientBalance     decimal.Decimal `gorm:"column:balance"`
}
