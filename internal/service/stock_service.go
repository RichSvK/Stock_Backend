package service

import (
	"backend/internal/model/domainerr"
	"backend/internal/model/response"
	"backend/internal/repository"
	"context"
)

type StockService interface {
	SearchStock(ctx context.Context, code string) (*response.SearchStockResponse, error)
}

type StockServiceImpl struct {
	StockRepository repository.StockRepository
}

func NewStockService(stockRepository repository.StockRepository) StockService {
	return &StockServiceImpl{
		StockRepository: stockRepository,
	}
}

func (service *StockServiceImpl) SearchStock(ctx context.Context, code string) (*response.SearchStockResponse, error) {
	listStock, err := service.StockRepository.SearchStock(ctx, code)
	if err != nil {
		return nil, domainerr.ErrInternalServer
	}

	if len(listStock) == 0 {
		return nil, domainerr.ErrStockNotFound
	}

	response := &response.SearchStockResponse{
		Message: "Stocks found",
		Data:    listStock,
	}

	return response, err
}
