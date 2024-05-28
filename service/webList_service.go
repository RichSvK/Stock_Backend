package service

import (
	"backend/helper"
	"backend/repository"
	"context"
	"net/http"
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
		return http.StatusInternalServerError, helper.ToFailedResponse(http.StatusInternalServerError, "Failed to get link data")
	}

	if len(listLink) == 0 {
		return http.StatusNotFound, helper.ToFailedResponse(http.StatusNotFound, "Link not found")
	}

	return http.StatusOK, helper.ToWebResponse(http.StatusOK, "Link was found", listLink)
}

func (service *StockWebServiceImpl) GetLinks(ctx context.Context) (int, interface{}) {
	listLink, err := service.StockWebRepository.GetLinks(ctx, "")

	if err != nil {
		return http.StatusInternalServerError, helper.ToFailedResponse(http.StatusInternalServerError, "Failed to get link data")
	}

	if len(listLink) == 0 {
		return http.StatusNotFound, helper.ToFailedResponse(http.StatusNotFound, "Link not found")
	}

	return http.StatusOK, helper.ToWebResponse(http.StatusOK, "Link was found", listLink)
}
