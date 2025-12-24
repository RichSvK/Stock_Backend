package service

import (
	domain_error "backend/model/error"
	"backend/model/web/response"
	"backend/repository"
	"context"
)

type StockWebService interface {
	GetLinks(ctx context.Context, categoryID string) (*response.GetLinkResponse, error)
}

type StockWebServiceImpl struct {
	StockWebRepository repository.StockWebRepository
}

func NewStockWebService(repositoryStockWeb repository.StockWebRepository) StockWebService {
	return &StockWebServiceImpl{
		StockWebRepository: repositoryStockWeb,
	}
}

func (service *StockWebServiceImpl) GetLinks(ctx context.Context, categoryID string) (*response.GetLinkResponse, error) {
	listLink, err := service.StockWebRepository.GetLinks(ctx, categoryID)

	if err != nil {
		return nil, err
	}

	if len(listLink) == 0 {
		return nil, domain_error.ErrLinkNotFound
	}

	response := &response.GetLinkResponse{
		Message: "Link was found",
		Data:    response.MapLinksToResponse(listLink),
	}

	return response, err
}
