package repository

import (
	"backend/config"
	"backend/model/entity"
	"context"
)

type BalanceRepository interface {
	Create(ctx context.Context, stock []entity.Stock) error
	GetBalanceStock(ctx context.Context, code string) ([]entity.Stock, error)
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
