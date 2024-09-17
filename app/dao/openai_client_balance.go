package dao

import (
	"context"
	"time"

	"github.com/alioth-center/akasha-whisper/app/model/dto"

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
			Select(model.OpenaiClientBalanceCols.BalanceRemaining, model.OpenaiClientBalanceCols.ID).
			Order(clause.OrderByColumn{Column: clause.Column{Name: model.OpenaiClientBalanceCols.CreatedAt}, Desc: true}).
			Scan(&receiver).
			Error; queryErr != nil {
			return queryErr
		} else if receiver.ID == 0 && action != model.OpenaiClientBalanceActionInitial {
			return errors.New("client balance not initialized")
		}

		if action == model.OpenaiClientBalanceActionInitial {
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

func (ac *OpenaiClientBalanceDatabaseAccessor) CreateBalanceRecordByName(ctx context.Context, clientName string, changeAmount decimal.Decimal, action model.EnumOpenaiClientBalanceAction, reason string) (after decimal.Decimal, err error) {
	execErr := ac.db.GetGormCore(ctx).Transaction(func(tx *gorm.DB) error {
		var clientID int64
		if queryErr := tx.WithContext(ctx).
			Model(&model.OpenaiClient{}).
			Where(model.OpenaiClientCols.Description, clientName).
			Select(model.OpenaiClientCols.ID).
			Scan(&clientID).
			Error; queryErr != nil {
			return queryErr
		}
		if clientID == 0 {
			return errors.New("client not found")
		}

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
			Reason:              reason,
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

func (ac *OpenaiClientBalanceDatabaseAccessor) StatisticsClientBalance(ctx context.Context, startDate time.Time) (result []*dto.OpenaiClientBalanceStatisticsDTO, err error) {
	result = make([]*dto.OpenaiClientBalanceStatisticsDTO, 0)
	sql := rawSqlList[RawsqlOpenaiClientBalanceStatistics]
	if queryErr := ac.db.GetGormCore(ctx).Raw(sql, startDate).Scan(&result).Error; queryErr != nil {
		return nil, errors.Wrap(queryErr, "statistics client balance failed")
	}

	return result, nil
}
