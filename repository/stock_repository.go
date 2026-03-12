package repository

import (
	"context"

	"gorm.io/gorm"
)

type StockRepository interface {
	SearchStock(ctx context.Context, stockCode string) ([]string, error)
}

type StockRepositoryImpl struct {
	DB *gorm.DB
}

func NewStockRepository(db *gorm.DB) StockRepository {
	return &StockRepositoryImpl{
		DB: db,
	}
}

func (repository *StockRepositoryImpl) SearchStock(ctx context.Context, stockCode string) ([]string, error) {
	db := repository.DB

	var listStock []string = nil
	err := db.WithContext(ctx).
		Table("stock").
		Distinct("code").
		Where("code LIKE ?", stockCode+"%").
		Pluck("code", &listStock).
		Error

	return listStock, err
}
