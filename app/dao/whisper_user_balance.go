package dao

import (
	"context"
	"github.com/alioth-center/akasha-whisper/app/model"
	"github.com/alioth-center/infrastructure/database"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type WhisperUserBalanceDatabaseAccessor struct {
	db database.DatabaseV2
}

func NewWhisperUserBalanceDatabaseAccessor(db database.DatabaseV2) *WhisperUserBalanceDatabaseAccessor {
	return &WhisperUserBalanceDatabaseAccessor{db: db}
}

func (ac *WhisperUserBalanceDatabaseAccessor) CreateBalanceRecord(ctx context.Context, userID int, changeAmount decimal.Decimal, action model.EnumWhisperUserBalanceAction, reason ...string) (after decimal.Decimal, err error) {
	recordReason := ""
	if len(reason) > 0 {
		recordReason = reason[0]
	}

	execErr := ac.db.GetGormCore(ctx).Transaction(func(tx *gorm.DB) error {
		receiver := &model.WhisperUserBalance{}
		if queryErr := tx.WithContext(ctx).
			Model(&model.WhisperUserBalance{}).
			Where(model.WhisperUserBalanceCols.UserID, userID).
			Select(model.WhisperUserBalanceCols.BalanceRemaining).
			Order(clause.OrderByColumn{Column: clause.Column{Name: model.WhisperUserBalanceCols.CreatedAt}, Desc: true}).
			First(receiver).
			Error; queryErr != nil && !errors.Is(queryErr, gorm.ErrRecordNotFound) {
			return queryErr
		} else if errors.Is(queryErr, gorm.ErrRecordNotFound) {
			receiver.BalanceRemaining = decimal.Zero
		}

		after = receiver.BalanceRemaining.Add(changeAmount)
		record := &model.WhisperUserBalance{
			UserID:              int64(userID),
			BalanceChangeAmount: changeAmount,
			BalanceRemaining:    after,
			Action:              action,
			Reason:              recordReason,
		}

		if createErr := tx.WithContext(ctx).Create(record).Error; createErr != nil {
			return createErr
		}

		return nil
	})
	if execErr != nil {
		return decimal.Zero, execErr
	}

	return after, nil
}
