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
	"log"
)

type BrokerService interface {
	GetBrokers(ctx context.Context, query *query_filter.BrokerQuery) (*response.GetBrokerResponse, error)
	CreateBroker(ctx context.Context, request *request.CreateBrokerRequest) (*response.BrokerResponse, error)
	UpdateBroker(ctx context.Context, request *request.UpdateBrokerRequest) (*response.SuccessResponse, error)
	DeleteBroker(ctx context.Context, code string) (*response.SuccessResponse, error)
}

type BrokerServiceImpl struct {
	BrokerRepository repository.BrokerRepository
}

func NewBrokerService(repositoryBroker repository.BrokerRepository) BrokerService {
	return &BrokerServiceImpl{
		BrokerRepository: repositoryBroker,
	}
}

func (service *BrokerServiceImpl) GetBrokers(ctx context.Context, query *query_filter.BrokerQuery) (*response.GetBrokerResponse, error) {
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

func (service *BrokerServiceImpl) CreateBroker(ctx context.Context, request *request.CreateBrokerRequest) (*response.BrokerResponse, error) {
	broker := &entity.Broker{
		BrokerCode: request.Code,
		BrokerName: request.Name,
	}

	err := service.BrokerRepository.CreateBroker(ctx, broker)
	if err != nil {
		return nil, err
	}

	response := &response.BrokerResponse{
		Code: broker.BrokerCode,
		Name: broker.BrokerName,
	}
	return response, nil
}

func (service *BrokerServiceImpl) UpdateBroker(ctx context.Context, request *request.UpdateBrokerRequest) (*response.SuccessResponse, error) {
	err := service.BrokerRepository.UpdateBroker(ctx, request)

	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	response := &response.SuccessResponse{
		Message: fmt.Sprintf("Success update %s broker", request.Code),
	}
	return response, nil
}

func (service *BrokerServiceImpl) DeleteBroker(ctx context.Context, code string) (*response.SuccessResponse, error) {
	err := service.BrokerRepository.DeleteBroker(ctx, code)
	if err != nil {
		return nil, err
	}

	response := &response.SuccessResponse{
		Message: fmt.Sprintf("Success delete %s broker", code),
	}
	return response, nil
}
