package repository

import (
	"backend/internal/entity"
	"backend/internal/model/domainerr"
	"backend/internal/model/query_filter"
	"backend/internal/model/request"
	"context"
	"log"

	"gorm.io/gorm"
)

type BrokerRepository interface {
	GetBrokers(ctx context.Context, query *query_filter.BrokerQuery) ([]entity.Broker, error)
	CreateBroker(ctx context.Context, broker *entity.Broker) error
	UpdateBroker(ctx context.Context, request *request.UpdateBrokerRequest) error
	DeleteBroker(ctx context.Context, code string) error
}

type BrokerRepositoryImpl struct {
	DB *gorm.DB
}

func NewBrokerRepository(db *gorm.DB) BrokerRepository {
	return &BrokerRepositoryImpl{
		DB: db,
	}
}

func (repository *BrokerRepositoryImpl) GetBrokers(ctx context.Context, query *query_filter.BrokerQuery) ([]entity.Broker, error) {
	var listBroker []entity.Broker = nil
	db := repository.DB.WithContext(ctx)

	if query.Code != "" {
		db = db.Where("broker_code = ?", query.Code)
	}

	err := db.Find(&listBroker).Error
	return listBroker, err
}

func (repository *BrokerRepositoryImpl) CreateBroker(ctx context.Context, broker *entity.Broker) error {
	db := repository.DB.WithContext(ctx)

	if err := db.Create(broker).Error; err != nil {
		log.Println(err.Error())
		if domainerr.IsDuplicateError(err) {
			return domainerr.ErrBrokerDuplicate
		}
		return domainerr.ErrInternalServer
	}

	return nil
}

func (repository *BrokerRepositoryImpl) UpdateBroker(ctx context.Context, request *request.UpdateBrokerRequest) error {
	result := repository.DB.WithContext(ctx).
		Model(&entity.Broker{}).
		Where("broker_code = ?", request.Code).
		Update("broker_name", request.Name)

	if result.Error != nil {
		return domainerr.ErrInternalServer
	}

	if result.RowsAffected == 0 {
		return domainerr.ErrBrokerNotFound
	}

	return nil
}

func (repository *BrokerRepositoryImpl) DeleteBroker(ctx context.Context, code string) error {
	result := repository.DB.WithContext(ctx).
		Where("broker_code = ?", code).
		Delete(&entity.Broker{})

	if result.Error != nil {
		log.Println(result.Error)
		return domainerr.ErrInternalServer
	}

	if result.RowsAffected == 0 {
		return domainerr.ErrBrokerNotFound
	}

	return nil
}
