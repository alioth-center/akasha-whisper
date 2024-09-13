package dao

import (
	"context"
	"time"

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
		} else if errors.Is(queryErr, gorm.ErrRecordNotFound) && action != model.WhisperUserBalanceActionInitial {
			return errors.New("user balance not initialized")
		}

		if action == model.WhisperUserBalanceActionInitial {
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

func (ac *WhisperUserBalanceDatabaseAccessor) ListBalanceRecords(ctx context.Context, userID int, start, end time.Time, page int, offset int) (records []*model.WhisperUserBalance, err error) {
	list := make([]*model.WhisperUserBalance, 0, page)

	if queryErr := ac.db.GetGormCore(ctx).
		Model(&model.WhisperUserBalance{}).
		Where(model.WhisperUserBalanceCols.UserID, userID).
		Where(clause.Gte{Column: model.WhisperUserBalanceCols.CreatedAt, Value: start}).
		Where(clause.Lte{Column: model.WhisperUserBalanceCols.CreatedAt, Value: end}).
		Order(clause.OrderByColumn{Column: clause.Column{Name: model.WhisperUserBalanceCols.CreatedAt}, Desc: true}).
		Offset(offset * page).
		Limit(page).
		Find(&list).
		Error; queryErr != nil {
		return nil, queryErr
	}

	return list, nil
}

func (ac *WhisperUserBalanceDatabaseAccessor) BatchCreateBalanceRecord(ctx context.Context, userID []int, changeAmount decimal.Decimal, action model.EnumWhisperUserBalanceAction, reason string) error {
	return ac.db.GetGormCore(ctx).Transaction(func(tx *gorm.DB) error {
		users := int64(0)
		if queryErr := tx.WithContext(ctx).
			Model(&model.WhisperUser{}).
			Where(model.WhisperUserCols.ID, userID).
			Select(model.WhisperUserCols.ID).
			Count(&users).
			Error; queryErr != nil {
			return queryErr
		}
		if users != int64(len(userID)) {
			return errors.New("data not consistent")
		}

		records := make([]model.WhisperUserBalance, 0, len(userID))
		for _, id := range userID {
			records = append(records, model.WhisperUserBalance{
				UserID:              int64(id),
				BalanceChangeAmount: changeAmount,
				Action:              action,
				Reason:              reason,
			})
		}

		if createErr := tx.WithContext(ctx).
			Model(&model.WhisperUserBalance{}).
			CreateInBatches(records, 100).Error; createErr != nil {
			return createErr
		}

		return nil
	})
}
