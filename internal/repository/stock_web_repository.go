package repository

import (
	"backend/internal/entity"
	"backend/internal/model/domainerr"
	"backend/internal/model/query_filter"
	"context"
	"log"

	"gorm.io/gorm"
)

type StockWebRepository interface {
	GetLinks(ctx context.Context, query *query_filter.GetLinkQuery) ([]entity.Link, error)
	CreateLink(ctx context.Context, link *entity.Link) error
	UpdateLink(ctx context.Context, name string, updates map[string]any) error
	DeleteLink(ctx context.Context, name string) error
}

type StockWebRepositoryImpl struct {
	DB *gorm.DB
}

func NewStockWebRepository(db *gorm.DB) StockWebRepository {
	return &StockWebRepositoryImpl{
		DB: db,
	}
}

func (repository *StockWebRepositoryImpl) GetLinks(ctx context.Context, query *query_filter.GetLinkQuery) ([]entity.Link, error) {
	db := repository.DB.WithContext(ctx)

	// Add filter if provided
	if query.Name != "" {
		db = db.Where("web_name LIKE ?", query.Name+"%")
	}

	if query.CategoryID != "" {
		db = db.Where("category_id = ?", query.CategoryID)
	}

	var listLink []entity.Link

	err := db.Select("url_link, web_name, web_image, web_description").
		Find(&listLink).
		Error

	if err != nil {
		return nil, domainerr.ErrInternalServer
	}

	return listLink, nil
}

func (repository *StockWebRepositoryImpl) CreateLink(ctx context.Context, link *entity.Link) error {
	db := repository.DB.WithContext(ctx)

	err := db.Create(link).Error
	if err != nil {
		if domainerr.IsDuplicateError(err) {
			return domainerr.ErrLinkAlreadyExists
		}
		return err
	}

	return nil
}

func (repository *StockWebRepositoryImpl) UpdateLink(ctx context.Context, name string, updates map[string]any) error {
	db := repository.DB.WithContext(ctx)

	result := db.Model(&entity.Link{}).Where("web_name = ?", name).
		Updates(updates)

	if result.Error != nil {
		log.Println(result.Error)
		return domainerr.ErrInternalServer
	}

	if result.RowsAffected == 0 {
		return domainerr.ErrLinkNotFound
	}

	return nil
}

func (repository *StockWebRepositoryImpl) DeleteLink(ctx context.Context, name string) error {
	db := repository.DB.WithContext(ctx)

	result := db.Where("web_name = ?", name).Delete(&entity.Link{})
	if result.Error != nil {
		return domainerr.ErrInternalServer
	}

	if result.RowsAffected == 0 {
		return domainerr.ErrLinkNotFound
	}

	return nil
}
