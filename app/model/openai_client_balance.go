package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type EnumOpenaiClientBalanceAction = string

const (
	OpenaiClientBalanceActionConsumption EnumOpenaiClientBalanceAction = "consumption" // 1. 消费：Consumption - Natural consumption in the system
	OpenaiClientBalanceActionRecharge    EnumOpenaiClientBalanceAction = "recharge"    // 2. 充值：Recharge - Actively recharging in the system
	OpenaiClientBalanceActionGift        EnumOpenaiClientBalanceAction = "gift"        // 3. 赠送：Gift - Given according to rules
	OpenaiClientBalanceActionSpecial     EnumOpenaiClientBalanceAction = "special"     // 4. 特殊：Special - Manually changed by an administrator
	OpenaiClientBalanceActionInitial     EnumOpenaiClientBalanceAction = "initial"     // 5. 初始：Initial - Initial balance when the account is created
)

// OpenaiClientBalance openai client balance record
//
// Reference: https://docs.alioth.center/akasha-whisper-database.html#openai-client-balance
type OpenaiClientBalance struct {
	ID                  int64           `gorm:"column:id;type:integer;autoIncrement:true;primaryKey:true"`
	ClientID            int64           `gorm:"column:client_id;type:integer;not null;index:idx_client_id;index:idx_scan"`
	BalanceChangeAmount decimal.Decimal `gorm:"column:balance_change_amount;type:decimal(16,8);not null;default:0"`
	BalanceRemaining    decimal.Decimal `gorm:"column:balance_remaining;type:decimal(16,8);not null;default:0;index:idx_balance_remaining"`
	Action              string          `gorm:"column:action;type:varchar(32);not null;index:idx_action"`
	Reason              string          `gorm:"column:reason;type:varchar(255);not null;default:''"`
	CreatedAt           time.Time       `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP;index:idx_created_at;index:idx_scan"`
}

func (m OpenaiClientBalance) TableName() string {
	return TableNameOpenaiClientBalance
}
