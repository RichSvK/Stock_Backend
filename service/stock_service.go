package service

import (
	"backend/helper"
	"backend/repository"
	"context"
)

type StockService interface {
	SearchStock(ctx context.Context, code string) (int, any)
}

type StockServiceImpl struct {
	StockRepository repository.StockRepository
}

func NewStockService(stockRepository repository.StockRepository) StockService {
	return &StockServiceImpl{
		StockRepository: stockRepository,
	}
}

func (service *StockServiceImpl) SearchStock(ctx context.Context, code string) (int, any) {
	listStock, err := service.StockRepository.SearchStock(ctx, code)
	if err != nil {
		return 500, "Failed to search stock"
	}

	if len(listStock) == 0 {
		return 404, "Stock not found"
	}

	return 200, helper.ToSearchStockResponses(listStock)
}
