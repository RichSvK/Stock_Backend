package repository

import (
	"backend/config"
	"backend/model/entity"
	"backend/model/web/request"
	"context"
	"fmt"
	"strconv"
)

type IpoRepository interface {
	GetAllIpo(ctx context.Context) ([]entity.Ipo, error)
	FindByCondition(ctx context.Context, request []request.Filter) ([]entity.Ipo, error)
}

type IpoRepositoryImpl struct{}

func NewIpoRepository() IpoRepository {
	return &IpoRepositoryImpl{}
}

func (repository *IpoRepositoryImpl) GetAllIpo(ctx context.Context) ([]entity.Ipo, error) {
	var listStock []entity.Ipo = nil
	var err error
	db := config.GetDatabaseInstance()
	query := "id.stock_code AS stock_code, price, ipo_shares, listed_shares, equity, warrant, nominal, mcb, is_affiliated, is_acceleration, is_new, lock_up, subscribed_stock, GROUP_CONCAT(uw_code) AS all_underwriter, GROUP_CONCAT(uw_shares) AS all_shares, (price * ipo_shares) AS amount"
	err = db.Table("ipo_detail id").
		WithContext(ctx).
		Select(query).
		Joins("JOIN stock_ipo s ON id.stock_code = s.stock_code").
		Group("id.stock_code").
		Scan(&listStock).
		Error
	return listStock, err
}

func (repository *IpoRepositoryImpl) FindByCondition(ctx context.Context, request []request.Filter) ([]entity.Ipo, error) {
	var listStock []entity.Ipo = nil
	db := config.GetDatabaseInstance()

	db_query := db.Table("ipo_detail id")
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
		WithContext(ctx).
		Select(query).
		Joins("JOIN stock_ipo s ON id.stock_code = s.stock_code").
		Group("id.stock_code").
		Scan(&listStock).
		Error

	return listStock, err
}
