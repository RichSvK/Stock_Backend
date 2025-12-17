package service

import (
	"backend/helper"
	"backend/model/web/request"
	"backend/repository"
	"context"
)

type IpoService interface {
	GetIpoAll(ctx context.Context) (int, any)
	GetIpoByCondition(ctx context.Context, request []request.Filter) (int, any)
}

type IpoServiceImpl struct {
	IpoRepository repository.IpoRepository
}

func NewIpoService(repositoryIPO repository.IpoRepository) IpoService {
	return &IpoServiceImpl{
		IpoRepository: repositoryIPO,
	}
}

func (service *IpoServiceImpl) GetIpoAll(ctx context.Context) (int, any) {
	listIPO, err := service.IpoRepository.GetAllIpo(ctx)
	if err != nil {
		return 500, helper.ToFailedResponse(500, "Failed to get IPO data")
	}

	if len(listIPO) == 0 {
		return 404, helper.ToFailedResponse(404, "IPO data not found")
	}

	return 200, helper.ToWebResponse(200, "IPO data found", helper.ToIpoResponses(listIPO))
}

func (service *IpoServiceImpl) GetIpoByCondition(ctx context.Context, request []request.Filter) (int, any) {
	listIPO, err := service.IpoRepository.FindByCondition(ctx, request)

	if err != nil {
		return 500, helper.ToFailedResponse(500, "Failed to get IPO data")
	}

	if len(listIPO) == 0 {
		return 404, helper.ToFailedResponse(404, "IPO data not found")
	}

	return 200, helper.ToWebResponse(200, "IPO data found", helper.ToIpoResponses(listIPO))
}
