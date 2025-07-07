package repository

import (
	"backend/config"
	"backend/model/entity"
	"context"
	"time"
)

type BalanceRepository interface {
	Create(ctx context.Context, stock []entity.Stock) error
	GetBalanceStock(ctx context.Context, code string) ([]entity.Stock, error)
	GetScriptlessChange(ctx context.Context, startDate time.Time, endDate time.Time) ([]entity.Stock, error)
}

type BalanceRepositoryImpl struct{}

func NewBalanceRepository() BalanceRepository {
	return &BalanceRepositoryImpl{}
}

func (repository *BalanceRepositoryImpl) Create(ctx context.Context, stock []entity.Stock) error {
	db := config.GetDatabaseInstance()
	tx := db.WithContext(ctx).Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Create(&stock).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func (repository *BalanceRepositoryImpl) GetBalanceStock(ctx context.Context, code string) ([]entity.Stock, error) {
	db := config.GetDatabaseInstance()

	var listStock []entity.Stock = nil
	err := db.WithContext(ctx).
		Model(&entity.Stock{}).
		Select("stock.*").
		Where("code = ?", code).
		Order("Date DESC").
		Limit(6).
		Scan(&listStock).
		Error

	return listStock, err
}

func (repository *BalanceRepositoryImpl) GetScriptlessChange(ctx context.Context, startDate time.Time, endDate time.Time) ([]entity.Stock, error) {
	db := config.GetDatabaseInstance()

	var listStock []entity.Stock
	startMonth := int(startDate.Month())
	startYear := startDate.Year()
	endMonth := int(endDate.Month())
	endYear := endDate.Year()

	err := db.WithContext(ctx).
		Model(&entity.Stock{}).
		Select("stock.*").
		Where("(MONTH(date) = ? AND YEAR(date) = ?) OR (MONTH(date) = ? AND YEAR(date) = ?)", startMonth, startYear, endMonth, endYear).
		Order("code ASC").
		Order("Date ASC").
		Scan(&listStock).
		Error

	return listStock, err
}
