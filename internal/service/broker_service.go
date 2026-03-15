package service

import (
	"backend/internal/model/domainerr"
	"backend/internal/model/query_filter"
	"backend/internal/model/response"
	"backend/internal/repository"
	"context"
)

type BrokerService interface {
	GetBrokers(ctx context.Context, query query_filter.BrokerQuery) (*response.GetBrokerResponse, error)
}

type BrokerServiceImpl struct {
	BrokerRepository repository.BrokerRepository
}

func NewBrokerService(repositoryBroker repository.BrokerRepository) BrokerService {
	return &BrokerServiceImpl{
		BrokerRepository: repositoryBroker,
	}
}

func (service *BrokerServiceImpl) GetBrokers(ctx context.Context, query query_filter.BrokerQuery) (*response.GetBrokerResponse, error) {
	listBroker, err := service.BrokerRepository.GetBrokers(ctx, query)
	if err != nil {
		return nil, domainerr.ErrInternalServer
	}

	if len(listBroker) == 0 {
		return nil, domainerr.ErrBrokerNotFound
	}

	response := &response.GetBrokerResponse{
		Message: "Brokers data found",
		Data:    response.ToBrokerResponses(listBroker),
	}

	return response, err
}
