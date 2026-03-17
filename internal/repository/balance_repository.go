package repository

import (
	"backend/internal/entity"
	"backend/internal/model/domainerr"
	"backend/internal/model/projection"
	"backend/internal/model/query_filter"
	"context"
	"fmt"
	"log"

	"gorm.io/gorm"
)

type BalanceRepository interface {
	Create(ctx context.Context, stock []entity.Stock) error
	GetBalanceStock(ctx context.Context, code string) ([]entity.Stock, error)
	GetScriptlessChange(ctx context.Context, dateRange query_filter.DateRangeQuery) ([]entity.Stock, error)
	GetBalanceChangeData(ctx context.Context, shareholderType string, change string, page int, dateRange query_filter.DateRangeQuery) ([]projection.BalanceChange, error)
}

type BalanceRepositoryImpl struct {
	DB *gorm.DB
}

func NewBalanceRepository(db *gorm.DB) BalanceRepository {
	return &BalanceRepositoryImpl{
		DB: db,
	}
}

func (repository *BalanceRepositoryImpl) Create(ctx context.Context, stock []entity.Stock) error {
	db := repository.DB.WithContext(ctx)
	tx := db.Begin()

	if tx.Error != nil {
		return domainerr.ErrBalanceCreationFailed
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Create(&stock).Error; err != nil {
		tx.Rollback()

		if domainerr.IsDuplicateError(err) {
			return domainerr.ErrDuplicateBalanceData
		}

		return domainerr.ErrBalanceCreationFailed
	}

	if err := tx.Commit().Error; err != nil {
		return domainerr.ErrBalanceCreationFailed
	}

	return nil
}

func (repository *BalanceRepositoryImpl) GetBalanceStock(ctx context.Context, code string) ([]entity.Stock, error) {
	var listStock []entity.Stock

	err := repository.DB.WithContext(ctx).
		Where("code = ?", code).
		Order("Date DESC").
		Limit(6).
		Find(&listStock).
		Error

	if err != nil {
		return nil, domainerr.ErrInternalServer
	}

	return listStock, nil
}

func (repository *BalanceRepositoryImpl) GetScriptlessChange(ctx context.Context, dateRange query_filter.DateRangeQuery) ([]entity.Stock, error) {
	var listStock []entity.Stock

	err := repository.DB.WithContext(ctx).
		Where("(date >= ? AND date < ?) OR (date >= ? AND date < ?)", dateRange.StartTime, dateRange.StartTimeLast, dateRange.EndTime, dateRange.EndTimeLast).
		Order("code ASC").
		Order("Date ASC").
		Find(&listStock).
		Error

	if err != nil {
		return nil, domainerr.ErrInternalServer
	}

	return listStock, err
}

func quoteIdent(col string) string {
	// Use backticks for MySQL identifiers.
	// Since col is validated from whitelist, this is safe.
	return "`" + col + "`"
}

func (repository *BalanceRepositoryImpl) GetBalanceChangeData(ctx context.Context, shareholderType string, change string, page int, dateRange query_filter.DateRangeQuery) ([]projection.BalanceChange, error) {
	db := repository.DB.WithContext(ctx)

	quotedCol := quoteIdent(shareholderType)

	const pageSize = 11
	offset := (page - 1) * (pageSize - 1)

	var changeExpr string
	var whereCond string
	if change == "Decrease" {
		changeExpr = fmt.Sprintf(`CASE WHEN s2.%[1]s = 0 THEN NULL ELSE ((s2.%[1]s - s.%[1]s) / IFNULL(s2.%[1]s, 1) * 100) END`, quotedCol)
		whereCond = fmt.Sprintf("s.%[1]s < s2.%[1]s", quotedCol)
	} else {
		changeExpr = fmt.Sprintf(`CASE WHEN s2.%[1]s = 0 THEN NULL ELSE ((s.%[1]s - s2.%[1]s) / IFNULL(s2.%[1]s, 1) * 100) END`, quotedCol)
		whereCond = fmt.Sprintf("s.%[1]s > s2.%[1]s", quotedCol)
	}

	query := fmt.Sprintf(`
	SELECT
		s.code AS stock_code,
		s.%[1]s AS current_ownership,
		s2.%[1]s AS previous_ownership,
		%[2]s AS change_percentage
	FROM stock s
	JOIN stock s2 
		ON s.code = s2.code
		AND s2.date = (
			SELECT MAX(date)
			FROM stock
			WHERE date >= ? AND date < ?
		)
		AND s.date = (
			SELECT MAX(date)
			FROM stock
			WHERE date >= ? AND date < ?
		)
	WHERE s2.%[1]s IS NOT NULL AND %[3]s
	ORDER BY change_percentage DESC, s.code
	LIMIT ? OFFSET ?;
`, quotedCol, changeExpr, whereCond)

	var listStock []projection.BalanceChange
	err := db.Raw(query, dateRange.StartTime, dateRange.StartTimeLast, dateRange.EndTime, dateRange.EndTimeLast, pageSize, offset).
		Scan(&listStock).
		Error

	if err != nil {
		log.Println(err.Error())
		return nil, domainerr.ErrInternalServer
	}
	return listStock, nil
}
