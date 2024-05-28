package repository

import (
	"backend/config"
	"backend/helper"
	"backend/model/entity"
	"context"
)

type IpoRepository interface {
	GetAllIpo(ctx context.Context) ([]entity.Ipo, error)
	FindByValue(ctx context.Context, value int, underwriter string) ([]entity.Ipo, error)
	FindByUnderwriter(ctx context.Context, underwriter string) ([]entity.Ipo, error)
}

type IpoRepositoryImpl struct{}

func NewIpoRepository() IpoRepository {
	return &IpoRepositoryImpl{}
}

func (repository *IpoRepositoryImpl) GetAllIpo(ctx context.Context) ([]entity.Ipo, error) {
	var listStock []entity.Ipo = nil
	var err error = nil
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

func (repository *IpoRepositoryImpl) FindByValue(ctx context.Context, value int, underwriter string) ([]entity.Ipo, error) {
	var listStock []entity.Ipo = nil
	db := config.GetDatabaseInstance()
	if underwriter != "ALL" {
		db = db.Joins("JOIN (SELECT s.stock_code AS stock_code, price, ipo_shares, listed_shares, equity, warrant, nominal, mcb, is_affiliated, is_acceleration, is_new, lock_up, subscribed_stock FROM ipo_detail ids JOIN stock_ipo s ON s.stock_code = ids.stock_code WHERE uw_code = ?) ts ON ts.stock_code = id.stock_code", underwriter)
	} else {
		db = db.Joins("JOIN stock_ipo s ON id.stock_code = s.stock_code")
	}

	query := "id.stock_code AS stock_code, price, ipo_shares, listed_shares, equity, warrant, nominal, mcb, is_affiliated, is_acceleration, is_new, lock_up, subscribed_stock, GROUP_CONCAT(uw_code) AS all_underwriter, GROUP_CONCAT(uw_shares) AS all_shares, (price * ipo_shares) AS amount"
	err := db.Table("ipo_detail id").
		WithContext(ctx).
		Select(query).
		Where(helper.GetAmountCondition(value)).
		Group("id.stock_code").
		Scan(&listStock).
		Error
	return listStock, err
}

func (repository *IpoRepositoryImpl) FindByUnderwriter(ctx context.Context, underwriter string) ([]entity.Ipo, error) {
	var listStock []entity.Ipo = nil
	db := config.GetDatabaseInstance()

	query := "id.stock_code AS stock_code, price, ipo_shares, listed_shares, equity, warrant, nominal, mcb, is_affiliated, is_acceleration, is_new, lock_up, subscribed_stock, GROUP_CONCAT(uw_code) AS all_underwriter, GROUP_CONCAT(uw_shares) AS all_shares, (price * ipo_shares) AS amount"
	err := db.Table("ipo_detail id").
		WithContext(ctx).
		Select(query).
		Joins("JOIN (SELECT s.stock_code AS stock_code, price, ipo_shares, listed_shares, equity, warrant, nominal, mcb, is_affiliated, is_acceleration, is_new, lock_up, subscribed_stock FROM ipo_detail ids JOIN stock_ipo s ON s.stock_code = ids.stock_code WHERE uw_code = ?) ts ON ts.stock_code = id.stock_code", underwriter).
		Group("id.stock_code").
		Scan(&listStock).
		Error

	return listStock, err
}
