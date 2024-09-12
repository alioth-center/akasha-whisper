package dto

import (
	"time"

	"github.com/shopspring/decimal"
)

type WhisperUserInfo struct {
	UserInfo WhisperUserInfoDTO
	Models   []string
}

type WhisperUserInfoDTO struct {
	ID        int             `gorm:"column:id"`
	Email     string          `gorm:"column:email"`
	ApiKey    string          `gorm:"column:api_key"`
	Role      string          `gorm:"column:role"`
	Language  string          `gorm:"column:language"`
	AllowIps  string          `gorm:"column:allow_ips"`
	UpdatedAt time.Time       `gorm:"column:updated_at"`
	Balance   decimal.Decimal `gorm:"column:balance"`
}
