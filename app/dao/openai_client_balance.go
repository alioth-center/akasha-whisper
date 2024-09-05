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

type OpenaiClientBalanceDatabaseAccessor struct {
	db database.DatabaseV2
}

func NewOpenaiClientBalanceDatabaseAccessor(db database.DatabaseV2) *OpenaiClientBalanceDatabaseAccessor {
	return &OpenaiClientBalanceDatabaseAccessor{db: db}
}

func (ac *OpenaiClientBalanceDatabaseAccessor) CreateBalanceRecord(ctx context.Context, clientID int, changeAmount decimal.Decimal, action model.EnumOpenaiClientBalanceAction, reason ...string) (after decimal.Decimal, err error) {
	recordReason := ""
	if len(reason) > 0 {
		recordReason = reason[0]
	}

	execErr := ac.db.GetGormCore(ctx).Transaction(func(tx *gorm.DB) error {
		receiver := &model.OpenaiClientBalance{}
		if queryErr := tx.WithContext(ctx).
			Model(&model.OpenaiClientBalance{}).
			Where(model.OpenaiClientBalanceCols.ClientID, clientID).
			Select(model.OpenaiClientBalanceCols.BalanceRemaining).
			Order(clause.OrderByColumn{Column: clause.Column{Name: model.OpenaiClientBalanceCols.CreatedAt}, Desc: true}).
			Scan(&receiver).
			Error; queryErr != nil && !errors.Is(queryErr, gorm.ErrRecordNotFound) {
			return queryErr
		} else if receiver.ID == 0 {
			receiver.BalanceRemaining = decimal.Zero
		}

		after = receiver.BalanceRemaining.Add(changeAmount)
		record := &model.OpenaiClientBalance{
			ClientID:            int64(clientID),
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
