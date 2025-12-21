package repository

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type StockRepository interface {
	SearchStock(ctx context.Context, stockCode string) ([]string, error)
}

type SearchStockImpl struct {
	DB *gorm.DB
}

func NewStockRepository(db *gorm.DB) StockRepository {
	return &SearchStockImpl{
		DB: db,
	}
}

func (repository *SearchStockImpl) SearchStock(ctx context.Context, stockCode string) ([]string, error) {
	db := repository.DB

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
