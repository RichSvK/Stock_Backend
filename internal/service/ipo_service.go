package service

import (
	"backend/internal/model/domainerr"
	"backend/internal/model/request"
	"backend/internal/model/response"
	"backend/internal/repository"
	"context"
)

type IpoService interface {
	GetIpoAll(ctx context.Context) (*response.GetIpoResponse, error)
	GetIpoByCondition(ctx context.Context, request []request.Filter) (*response.GetIpoResponse, error)
}

type IpoServiceImpl struct {
	IpoRepository repository.IpoRepository
}

func NewIpoService(repositoryIPO repository.IpoRepository) IpoService {
	return &IpoServiceImpl{
		IpoRepository: repositoryIPO,
	}
}

func (service *IpoServiceImpl) GetIpoAll(ctx context.Context) (*response.GetIpoResponse, error) {
	listIPO, err := service.IpoRepository.GetAllIpo(ctx)
	if err != nil {
		return nil, err
	}

	if len(listIPO) == 0 {
		return nil, domainerr.ErrIpoDataNotFound
	}

	response := &response.GetIpoResponse{
		Message: "IPO data found",
		Data:    response.ToIpoResponses(listIPO),
	}
	return response, err
}

func (service *IpoServiceImpl) GetIpoByCondition(ctx context.Context, request []request.Filter) (*response.GetIpoResponse, error) {
	listIPO, err := service.IpoRepository.FindByCondition(ctx, request)

	if err != nil {
		return nil, err
	}

	if len(listIPO) == 0 {
		return nil, domainerr.ErrIpoDataNotFound
	}

	response := &response.GetIpoResponse{
		Message: "IPO data found",
		Data:    response.ToIpoResponses(listIPO),
	}

	return response, err
}
