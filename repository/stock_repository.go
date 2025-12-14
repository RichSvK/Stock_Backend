package repository

import (
	"backend/config"
	"context"
	"fmt"
	"time"
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

	start := time.Now()

	var listStock []string = nil
	err := db.WithContext(ctx).
		Table("stock").
		Distinct("code").
		Where("code LIKE ?", stockCode+"%").
		Pluck("code", &listStock).
		Error

	elapsed := time.Since(start)

	fmt.Println("Time for query: " + elapsed.String())
	return listStock, err
}
