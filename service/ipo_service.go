package service

import (
	"backend/helper"
	"backend/model/web/request"
	"backend/repository"
	"context"
	"fmt"
	"strconv"
)

type IpoService interface {
	GetIpoAll(ctx context.Context) (int, any)
	GetIpoByUnderwriter(ctx context.Context, underwriter string) (int, any)
	GetIpoByValue(ctx context.Context, value string, underwriter string) (int, any)
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

func (service *IpoServiceImpl) GetIpoByUnderwriter(ctx context.Context, underwriter string) (int, any) {
	listIPO, err := service.IpoRepository.FindByUnderwriter(ctx, underwriter)
	if err != nil {
		return 500, helper.ToFailedResponse(500, "Failed to get IPO data")
	}

	if len(listIPO) == 0 {
		return 404, helper.ToFailedResponse(404, fmt.Sprintf("IPO data with %s underwriter not found", underwriter))
	}

	return 200, helper.ToWebResponse(200, "IPO data found", helper.ToIpoResponses(listIPO))
}

func (service *IpoServiceImpl) GetIpoByValue(ctx context.Context, value string, underwriter string) (int, any) {
	values, err := strconv.Atoi(value)
	if err != nil {
		return 400, helper.ToFailedResponse(400, "Bad Request")
	}

	listIPO, err := service.IpoRepository.FindByValue(ctx, values, underwriter)
	if err != nil {
		return 500, helper.ToFailedResponse(500, "Failed to get IPO data")
	}

	if len(listIPO) == 0 {
		return 404, helper.ToFailedResponse(404, fmt.Sprintf("IPO data with %s value not found", value))
	}

	return 200, helper.ToWebResponse(200, "IPO data found", helper.ToIpoResponses(listIPO))
}

func (service *IpoServiceImpl) GetIpoByCondition(ctx context.Context, request []request.Filter) (int, any) {
	listIPO, err := service.IpoRepository.FindByCondition(ctx, request)
	fmt.Printf("HELLO Service")
	if err != nil {
		return 500, helper.ToFailedResponse(500, "Failed to get IPO data")
	}

	if len(listIPO) == 0 {
		return 404, helper.ToFailedResponse(404, "IPO data not found")
	}

	return 200, helper.ToWebResponse(200, "IPO data found", helper.ToIpoResponses(listIPO))
}
