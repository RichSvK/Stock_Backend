package service

import (
	"backend/internal/entity"
	"backend/internal/model/domainerr"
	"backend/internal/model/query_filter"
	"backend/internal/model/request"
	"backend/internal/model/response"
	"backend/internal/repository"
	"context"
	"fmt"
)

type StockWebService interface {
	GetLinks(ctx context.Context, query *query_filter.GetLinkQuery) (*response.GetLinkResponse, error)
	CreateLink(ctx context.Context, request *request.CreateLinkRequest) (*response.LinkResponse, error)
	UpdateLink(ctx context.Context, request *request.UpdateLinkRequest) (*response.SuccessResponse, error)
	DeleteLink(ctx context.Context, name string) (*response.SuccessResponse, error)
}

type StockWebServiceImpl struct {
	StockWebRepository repository.StockWebRepository
}

func NewStockWebService(repositoryStockWeb repository.StockWebRepository) StockWebService {
	return &StockWebServiceImpl{
		StockWebRepository: repositoryStockWeb,
	}
}

func (service *StockWebServiceImpl) GetLinks(ctx context.Context, query *query_filter.GetLinkQuery) (*response.GetLinkResponse, error) {
	listLink, err := service.StockWebRepository.GetLinks(ctx, query)

	if err != nil {
		return nil, err
	}

	if len(listLink) == 0 {
		return nil, domainerr.ErrLinkNotFound
	}

	response := &response.GetLinkResponse{
		Message: "Link was found",
		Data:    response.MapLinksToResponse(listLink),
	}

	return response, err
}

func (service *StockWebServiceImpl) CreateLink(ctx context.Context, request *request.CreateLinkRequest) (*response.LinkResponse, error) {
	link := &entity.Link{
		URL:         request.URL,
		Name:        request.Name,
		Image:       request.Image,
		Description: request.Description,
		CategoryID:  request.Category,
	}

	err := service.StockWebRepository.CreateLink(ctx, link)
	if err != nil {
		return nil, err
	}

	response := &response.LinkResponse{
		URL:         link.URL,
		Name:        link.Name,
		Image:       link.Image,
		Description: link.Description,
	}

	return response, nil
}

func (service *StockWebServiceImpl) UpdateLink(ctx context.Context, request *request.UpdateLinkRequest) (*response.SuccessResponse, error) {
	updates := make(map[string]any)
	if request.URL != "" {
		updates["url_link"] = request.URL
	}

	if request.Image != "" {
		updates["web_image"] = request.Image
	}

	if request.Description != "" {
		updates["web_description"] = request.Description
	}

	if request.Category != 0 {
		updates["category_id"] = request.Category
	}

	if len(updates) == 0 {
		return nil, domainerr.ErrNoFieldsToUpdate
	}

	err := service.StockWebRepository.UpdateLink(ctx, request.Name, updates)
	if err != nil {
		return nil, err
	}

	response := &response.SuccessResponse{
		Message: fmt.Sprintf("Link %s was updated successfully", request.Name),
	}

	return response, nil
}

func (service *StockWebServiceImpl) DeleteLink(ctx context.Context, name string) (*response.SuccessResponse, error) {
	err := service.StockWebRepository.DeleteLink(ctx, name)
	if err != nil {
		return nil, err
	}

	response := &response.SuccessResponse{
		Message: fmt.Sprintf("Link %s was deleted successfully", name),
	}

	return response, nil
}
