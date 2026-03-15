package repository

import (
	"backend/internal/entity"
	"backend/internal/model/query_filter"
	"context"

	"gorm.io/gorm"
)

type BrokerRepository interface {
	GetBrokers(ctx context.Context, query query_filter.BrokerQuery) ([]entity.Broker, error)
}

type BrokerRepositoryImpl struct {
	DB *gorm.DB
}

func NewBrokerRepository(db *gorm.DB) BrokerRepository {
	return &BrokerRepositoryImpl{
		DB: db,
	}
}

func (repository *BrokerRepositoryImpl) GetBrokers(ctx context.Context, query query_filter.BrokerQuery) ([]entity.Broker, error) {
	var listBroker []entity.Broker = nil
	db := repository.DB.WithContext(ctx)

	if query.Code != "" {
		db = db.Where("broker_code = ?", query.Code)
	}

	err := db.Find(&listBroker).Error
	return listBroker, err
}
