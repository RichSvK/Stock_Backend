package repository

import (
	"backend/internal/entity"
	"context"

	"gorm.io/gorm"
)

type StockWebRepository interface {
	GetLinks(ctx context.Context, categoryID string) ([]entity.Link, error)
}

type StockWebRepositoryImpl struct {
	DB *gorm.DB
}

func NewStockWebRepository(db *gorm.DB) StockWebRepository {
	return &StockWebRepositoryImpl{
		DB: db,
	}
}

func (repository *StockWebRepositoryImpl) GetLinks(ctx context.Context, categoryID string) ([]entity.Link, error) {
	db := repository.DB.WithContext(ctx)

	// Jika categoryID diberikan, tambahkan kondisi WHERE ke query.
	if categoryID != "" {
		db = db.Where("category_id = ?", categoryID)
	}

	var listLink []entity.Link

	err := db.
		Model(&entity.Link{}).
		Select("url_link, web_name, web_image, web_description").
		Find(&listLink).
		Error

	return listLink, err
}
