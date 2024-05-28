package service

import (
	"backend/helper"
	"backend/repository"
	"context"
	"fmt"
	"net/http"
	"strconv"
)

type IpoService interface {
	GetIpoAll(ctx context.Context) (int, interface{})
	GetIpoByUnderwriter(ctx context.Context, underwriter string) (int, interface{})
	GetIpoByValue(ctx context.Context, value string, underwriter string) (int, interface{})
}

type IpoServiceImpl struct {
	IpoRepository repository.IpoRepository
}

func NewIpoService(repositoryIPO repository.IpoRepository) IpoService {
	return &IpoServiceImpl{
		IpoRepository: repositoryIPO,
	}
}

func (service *IpoServiceImpl) GetIpoAll(ctx context.Context) (int, interface{}) {
	listIPO, err := service.IpoRepository.GetAllIpo(ctx)
	if err != nil {
		return http.StatusInternalServerError, helper.ToFailedResponse(http.StatusInternalServerError, "Failed to get IPO data")
	}

	if len(listIPO) == 0 {
		return http.StatusNotFound, helper.ToFailedResponse(http.StatusNotFound, "IPO data not found")
	}

	return http.StatusOK, helper.ToWebResponse(http.StatusOK, "IPO data found", helper.ToIpoResponses(listIPO))
}

func (service *IpoServiceImpl) GetIpoByUnderwriter(ctx context.Context, underwriter string) (int, interface{}) {
	listIPO, err := service.IpoRepository.FindByUnderwriter(ctx, underwriter)
	if err != nil {
		return http.StatusInternalServerError, helper.ToFailedResponse(http.StatusInternalServerError, "Failed to get IPO data")
	}

	if len(listIPO) == 0 {
		return http.StatusNotFound, helper.ToFailedResponse(http.StatusNotFound, fmt.Sprintf("IPO data with %s underwriter not found", underwriter))
	}

	return http.StatusOK, helper.ToWebResponse(http.StatusOK, "IPO data found", helper.ToIpoResponses(listIPO))
}

func (service *IpoServiceImpl) GetIpoByValue(ctx context.Context, value string, underwriter string) (int, interface{}) {
	values, err := strconv.Atoi(value)
	if err != nil {
		return http.StatusBadRequest, helper.ToFailedResponse(http.StatusBadRequest, "Bad Request")
	}

	listIPO, err := service.IpoRepository.FindByValue(ctx, values, underwriter)
	if err != nil {
		return http.StatusInternalServerError, helper.ToFailedResponse(http.StatusInternalServerError, "Failed to get IPO data")
	}

	if len(listIPO) == 0 {
		return http.StatusNotFound, helper.ToFailedResponse(http.StatusNotFound, fmt.Sprintf("IPO data with %s value not found", value))
	}

	return http.StatusOK, helper.ToWebResponse(http.StatusOK, "IPO data found", helper.ToIpoResponses(listIPO))
}
