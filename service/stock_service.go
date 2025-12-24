package service

import (
	domain_error "backend/model/error"
	"backend/model/web/response"
	"backend/repository"
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
		return nil, domain_error.ErrInternalServer
	}

	if len(listStock) == 0 {
		return nil, domain_error.ErrStockNotFound
	}

	response := &response.SearchStockResponse{
		Message: "Stock found",
		Data:    listStock,
	}

	return response, err
}
