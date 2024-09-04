package model

import (
	"github.com/shopspring/decimal"
	"time"
)

type EnumWhisperUserBalanceAction = string

const (
	WhisperUserBalanceActionConsumption EnumWhisperUserBalanceAction = "consumption" // 1. 消费：Consumption - Natural consumption in the system
	WhisperUserBalanceActionRecharge    EnumWhisperUserBalanceAction = "recharge"    // 2. 充值：Recharge - Actively recharging in the system
	WhisperUserBalanceActionGift        EnumWhisperUserBalanceAction = "gift"        // 3. 赠送：Gift - Given according to rules
	WhisperUserBalanceActionSpecial     EnumWhisperUserBalanceAction = "special"     // 4. 特殊：Special - Manually changed by an administrator
	WhisperUserBalanceActionInitial     EnumWhisperUserBalanceAction = "initial"     // 5. 初始：Initial - Initial balance when the user is created
)

// WhisperUserBalance whisper user balance record
//
// Reference: https://docs.alioth.center/akasha-whisper-database.html#whisper-user-balance
type WhisperUserBalance struct {
	ID                  int64           `gorm:"column:id;type:integer;autoIncrement:true;primaryKey:true"`
	UserID              int64           `gorm:"column:user_id;type:integer;not null;index:idx_client_id;index:idx_scan"`
	BalanceChangeAmount decimal.Decimal `gorm:"column:balance_change_amount;type:decimal(16,8);not null;default:0"`
	BalanceRemaining    decimal.Decimal `gorm:"column:balance_remaining;type:decimal(16,8);not null;default:0;index:idx_balance_remaining"`
	Action              string          `gorm:"column:action;type:varchar(32);not null;index:idx_action"`
	Reason              string          `gorm:"column:reason;type:varchar(255);not null;default:''"`
	CreatedAt           time.Time       `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP;index:idx_created_at;index:idx_scan"`
}

func (m WhisperUserBalance) TableName() string {
	return TableNameWhisperUserBalance
}
