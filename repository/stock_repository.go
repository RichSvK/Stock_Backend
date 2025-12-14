package repository

import (
	"backend/config"
	"context"
)

type StockRepository interface {
	SearchStock(ctx context.Context, stockCode string) ([]string, error)
}

type SearchStockImpl struct{}

func NewStockRepository() StockRepository {
	return &SearchStockImpl{}
}

func (repository *SearchStockImpl) SearchStock(ctx context.Context, stockCode string) ([]string, error) {
	db := config.GetDatabaseInstance()

	var listStock []string = nil
	err := db.WithContext(ctx).
		Table("stock").
		Distinct("code").
		Where("code LIKE ?", stockCode+"%").
		Pluck("code", &listStock).
		Error

	return listStock, err
}
