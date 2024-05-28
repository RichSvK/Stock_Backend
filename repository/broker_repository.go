package repository

import (
	"backend/config"
	"backend/model/entity"
	"context"
)

type BrokerRepository interface {
	GetBrokers(ctx context.Context) ([]entity.Broker, error)
}

type BrokerRepositoryImpl struct{}

func NewBrokerRepository() BrokerRepository {
	return &BrokerRepositoryImpl{}
}

func (repository *BrokerRepositoryImpl) GetBrokers(ctx context.Context) ([]entity.Broker, error) {
	var listBroker []entity.Broker = nil
	db := config.GetDatabaseInstance()

	err := db.Model(&entity.Broker{}).
		WithContext(ctx).
		Scan(&listBroker).
		Error
	return listBroker, err
}
