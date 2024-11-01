package repository

import (
	"backend/config"
	"backend/model/entity"
	"context"
)

type PostIpoRepository interface {
	GetBalanceFilter(ctx context.Context, condition string) ([]entity.Stock, error)
}

type PostIpoRepositoryImpl struct{}

func (repository *PostIpoRepositoryImpl) GetBalanceFilter(ctx context.Context, condition string) ([]entity.Stock, error) {
	db := config.GetDatabaseInstance()

	var listStock []entity.Stock = nil
	err := db.WithContext(ctx).
		Table("stock_ipo si").
		Select("code, local_is, local_cp, local_pf, local_ib, local_id, local_mf, local_sc, local_fd, local_ot, foreign_is, foreign_cp, foreign_pf, foreign_ib, foreign_id, foreign_mf, foreign_sc, foreign_fd, foreign_ot, GREATEST(local_is, local_cp, local_pf, local_ib, local_id, local_mf, local_sc, local_fd, local_ot, foreign_is, foreign_cp, foreign_pf, foreign_ib, foreign_id, foreign_mf, foreign_sc, foreign_fd, foreign_ot) AS maxhold").
		Joins("JOIN stock s ON si.stock_code = s.code").
		Where("? = maxhold", condition).
		Error

	return listStock, err
}
