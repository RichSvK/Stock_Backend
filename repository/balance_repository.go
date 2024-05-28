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
	err := db.WithContext(ctx).Model(&entity.Stock{}).Create(&stock).Error
	return err
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
