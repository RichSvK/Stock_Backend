package repository

import (
	"backend/model/entity"
	domain_error "backend/model/error"
	query_model "backend/model/query"
	"context"
	"fmt"

	"gorm.io/gorm"
)

type BalanceRepository interface {
	Create(ctx context.Context, stock []entity.Stock) error
	GetBalanceStock(ctx context.Context, code string) ([]entity.Stock, error)
	GetScriptlessChange(ctx context.Context, dateRange query_model.DateRangeQuery) ([]entity.Stock, error)
	GetBalanceChangeData(ctx context.Context, shareholderType string, change string, page int, dateRange query_model.DateRangeQuery) ([]query_model.BalanceChange, error)
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
	db := repository.DB
	tx := db.WithContext(ctx).Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Create(&stock).Error; err != nil {
		tx.Rollback()
		return domain_error.ErrBalanceCreationFailed
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return domain_error.ErrBalanceCreationFailed
	}

	return nil
}

func (repository *BalanceRepositoryImpl) GetBalanceStock(ctx context.Context, code string) ([]entity.Stock, error) {
	db := repository.DB

	var listStock []entity.Stock = nil
	err := db.WithContext(ctx).
		Model(&entity.Stock{}).
		Select("stock.*").
		Where("code = ?", code).
		Order("Date DESC").
		Limit(6).
		Scan(&listStock).
		Error

	if err != nil {
		return nil, domain_error.ErrInternalServer
	}

	if len(listStock) == 0 {
		return nil, domain_error.ErrBalanceNotFound
	}

	return listStock, err
}

func (repository *BalanceRepositoryImpl) GetScriptlessChange(ctx context.Context, dateRange query_model.DateRangeQuery) ([]entity.Stock, error) {
	db := repository.DB

	var listStock []entity.Stock

	err := db.WithContext(ctx).
		Model(&entity.Stock{}).
		Select("stock.*").
		Where("(date >= ? AND date < ?) OR (date >= ? AND date < ?)", dateRange.StartTime, dateRange.StartTimeLast, dateRange.EndTime, dateRange.EndTimeLast).
		Order("code ASC").
		Order("Date ASC").
		Scan(&listStock).
		Error

	if err != nil {
		return nil, domain_error.ErrInternalServer
	}

	if len(listStock) == 0 {
		return nil, domain_error.ErrBalanceNotFound
	}

	return listStock, err
}

func quoteIdent(col string) string {
	// Use backticks for MySQL identifiers.
	// Since col is validated from whitelist, this is safe.
	return "`" + col + "`"
}

func (repository *BalanceRepositoryImpl) GetBalanceChangeData(ctx context.Context, shareholderType string, change string, page int, dateRange query_model.DateRangeQuery) ([]query_model.BalanceChange, error) {
	db := repository.DB

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
		AND s2.date >= ? AND s2.date < ? 
		AND s.date >= ? AND s.date < ?
	WHERE s2.%[1]s IS NOT NULL AND %[3]s
	ORDER BY change_percentage DESC, s.code
	LIMIT ? OFFSET ?;
`, quotedCol, changeExpr, whereCond)

	var listStock []query_model.BalanceChange
	err := db.WithContext(ctx).
		Raw(query, dateRange.StartTime, dateRange.StartTimeLast, dateRange.EndTime, dateRange.EndTimeLast, pageSize, offset).
		Scan(&listStock).
		Error

	if err != nil {
		fmt.Println(err.Error())
		return nil, domain_error.ErrInternalServer
	}
	return listStock, nil
}
