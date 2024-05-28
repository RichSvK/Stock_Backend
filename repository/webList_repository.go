package repository

import (
	"backend/config"
	"backend/model/entity"
	"context"
)

type StockWebRepository interface {
	GetLinks(ctx context.Context, categoryID string) ([]entity.Link, error)
}

type StockWebRepositoryImpl struct{}

func NewStockWebRepository() StockWebRepository {
	return &StockWebRepositoryImpl{}
}

func (repository *StockWebRepositoryImpl) GetLinks(ctx context.Context, categoryID string) ([]entity.Link, error) {
	db := config.GetDatabaseInstance()

	// Jika categoryID diberikan, tambahkan kondisi WHERE ke query.
	if categoryID != "" {
		db = db.Where("category_id = ?", categoryID)
	}

	var listLink []entity.Link = nil

	err := db.WithContext(ctx).
		Model(&entity.Link{}).
		Select("url_link, web_name, web_image, web_description").
		Scan(&listLink).
		Error

	return listLink, err
}
