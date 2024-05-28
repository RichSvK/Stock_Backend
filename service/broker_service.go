package service

import (
	"backend/helper"
	"backend/repository"
	"context"
	"net/http"
)

type BrokerService interface {
	GetBrokers(ctx context.Context) (int, interface{})
}

type BrokerServiceImpl struct {
	BrokerRepository repository.BrokerRepository
}

func NewBrokerService(repositoryBroker repository.BrokerRepository) BrokerService {
	return &BrokerServiceImpl{
		BrokerRepository: repositoryBroker,
	}
}

func (service *BrokerServiceImpl) GetBrokers(ctx context.Context) (int, interface{}) {
	listBroker, err := service.BrokerRepository.GetBrokers(ctx)
	if err != nil {
		return http.StatusInternalServerError, helper.ToFailedResponse(http.StatusInternalServerError, "Failed to get broker data")
	}

	if len(listBroker) == 0 {
		return http.StatusNotFound, helper.ToFailedResponse(http.StatusNotFound, "Broker data not found")
	}

	return http.StatusOK, helper.ToWebResponse(http.StatusOK, "Broker data found", helper.ToBrokerResponses(listBroker))
}
