package repository

import (
	"backend/internal/entity"
	"backend/internal/model/domainerr"
	"backend/internal/model/query_filter"
	"backend/internal/model/request"
	"context"
	"fmt"
	"strconv"

	"gorm.io/gorm"
)

type IpoRepository interface {
	GetIpo(ctx context.Context, query query_filter.GetIpoQuery) ([]entity.Ipo, error)
	FindByCondition(ctx context.Context, request []request.Filter) ([]entity.Ipo, error)
}

type IpoRepositoryImpl struct {
	DB *gorm.DB
}

func NewIpoRepository(db *gorm.DB) IpoRepository {
	return &IpoRepositoryImpl{
		DB: db,
	}
}

func (repository *IpoRepositoryImpl) GetIpo(ctx context.Context, query query_filter.GetIpoQuery) ([]entity.Ipo, error) {
	var listStock []entity.Ipo
	var err error
	column := "id.stock_code AS stock_code, price, ipo_shares, listed_shares, equity, warrant, nominal, mcb, is_affiliated, is_acceleration, is_new, lock_up, subscribed_stock, GROUP_CONCAT(uw_code) AS all_underwriter, GROUP_CONCAT(uw_shares) AS all_shares, (price * ipo_shares) AS amount"

	db := repository.DB.WithContext(ctx)
	if query.Code != "" {
		db = db.Where("id.stock_code = ?", query.Code)
	}

	err = db.Table("ipo_detail id").
		Select(column).
		Joins("JOIN stock_ipo s ON id.stock_code = s.stock_code").
		Group("id.stock_code").
		Scan(&listStock).
		Error

	if err != nil {
		return nil, domainerr.ErrInternalServer
	}

	return listStock, nil
}

func (repository *IpoRepositoryImpl) FindByCondition(ctx context.Context, request []request.Filter) ([]entity.Ipo, error) {
	var listStock []entity.Ipo
	db_query := repository.DB.WithContext(ctx)

	for _, filter := range request {
		condition := fmt.Sprintf("%s %s ?", filter.FilterName, filter.Symbol)
		if filter.FilterType == "number" {
			value, _ := strconv.ParseUint(filter.FilterValue, 10, 64)
			db_query = db_query.Where(condition, value)
		} else {
			db_query = db_query.Where(condition, filter.FilterValue)
		}
	}

	query := "id.stock_code AS stock_code, price, ipo_shares, listed_shares, equity, warrant, nominal, mcb, is_affiliated, is_acceleration, is_new, lock_up, subscribed_stock, GROUP_CONCAT(uw_code) AS all_underwriter, GROUP_CONCAT(uw_shares) AS all_shares, (price * ipo_shares) AS amount"
	err := db_query.
		Table("ipo_detail id").
		Select(query).
		Joins("JOIN stock_ipo s ON id.stock_code = s.stock_code").
		Group("id.stock_code").
		Scan(&listStock).
		Error

	if err != nil {
		return nil, domainerr.ErrInternalServer
	}

	return listStock, nil
}