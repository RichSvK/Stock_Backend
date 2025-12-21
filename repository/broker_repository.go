package repository

import (
	"backend/model/entity"
	"context"

	"gorm.io/gorm"
)

type BrokerRepository interface {
	GetBrokers(ctx context.Context) ([]entity.Broker, error)
}

type BrokerRepositoryImpl struct {
	DB *gorm.DB
}

func NewBrokerRepository(db *gorm.DB) BrokerRepository {
	return &BrokerRepositoryImpl{
		DB: db,
	}
}

func (repository *BrokerRepositoryImpl) GetBrokers(ctx context.Context) ([]entity.Broker, error) {
	var listBroker []entity.Broker = nil
	db := repository.DB

	err := db.Model(&entity.Broker{}).
		WithContext(ctx).
		Scan(&listBroker).
		Error
	return listBroker, err
}
