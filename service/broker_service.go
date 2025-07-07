package service

import (
	"backend/helper"
	"backend/repository"
	"context"
)

type BrokerService interface {
	GetBrokers(ctx context.Context) (int, any)
}

type BrokerServiceImpl struct {
	BrokerRepository repository.BrokerRepository
}

func NewBrokerService(repositoryBroker repository.BrokerRepository) BrokerService {
	return &BrokerServiceImpl{
		BrokerRepository: repositoryBroker,
	}
}

func (service *BrokerServiceImpl) GetBrokers(ctx context.Context) (int, any) {
	listBroker, err := service.BrokerRepository.GetBrokers(ctx)
	if err != nil {
		return 500, helper.ToFailedResponse(500, "Failed to get broker data")
	}

	if len(listBroker) == 0 {
		return 404, helper.ToFailedResponse(404, "Broker data not found")
	}

	return 200, helper.ToWebResponse(200, "Broker data found", helper.ToBrokerResponses(listBroker))
}
