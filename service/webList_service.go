package service

import (
	"backend/helper"
	"backend/repository"
	"context"
)

type StockWebService interface {
	GetLinkReference(ctx context.Context, categoryID string) (int, interface{})
	GetLinks(ctx context.Context) (int, interface{})
}

type StockWebServiceImpl struct {
	StockWebRepository repository.StockWebRepository
}

func NewStockWebService(repositoryStockWeb repository.StockWebRepository) StockWebService {
	return &StockWebServiceImpl{
		StockWebRepository: repositoryStockWeb,
	}
}

func (service *StockWebServiceImpl) GetLinkReference(ctx context.Context, categoryID string) (int, interface{}) {
	listLink, err := service.StockWebRepository.GetLinks(ctx, categoryID)

	if err != nil {
		return 500, helper.ToFailedResponse(500, "Failed to get link data")
	}

	if len(listLink) == 0 {
		return 404, helper.ToFailedResponse(404, "Link not found")
	}

	return 200, helper.ToWebResponse(200, "Link was found", listLink)
}

func (service *StockWebServiceImpl) GetLinks(ctx context.Context) (int, interface{}) {
	listLink, err := service.StockWebRepository.GetLinks(ctx, "")

	if err != nil {
		return 500, helper.ToFailedResponse(500, "Failed to get link data")
	}

	if len(listLink) == 0 {
		return 404, helper.ToFailedResponse(404, "Link not found")
	}

	return 200, helper.ToWebResponse(200, "Link was found", listLink)
}
