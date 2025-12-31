package repository

import (
	"backend/model/entity"
	domain_error "backend/model/error"
	query_model "backend/model/query"
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type BalanceRepository interface {
	Create(ctx context.Context, stock []entity.Stock) error
	GetBalanceStock(ctx context.Context, code string) ([]entity.Stock, error)
	GetScriptlessChange(ctx context.Context, startDate time.Time, endDate time.Time) ([]entity.Stock, error)
	GetBalanceChangeData(ctx context.Context, shareholderType string, change string, page int, startDate string, endDate string) ([]query_model.BalanceChange, error)
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

func (repository *BalanceRepositoryImpl) GetScriptlessChange(ctx context.Context, startDate time.Time, endDate time.Time) ([]entity.Stock, error) {
	db := repository.DB

	var listStock []entity.Stock
	startMonth := int(startDate.Month())
	startYear := startDate.Year()
	endMonth := int(endDate.Month())
	endYear := endDate.Year()

	err := db.WithContext(ctx).
		Model(&entity.Stock{}).
		Select("stock.*").
		Where("(MONTH(date) = ? AND YEAR(date) = ?) OR (MONTH(date) = ? AND YEAR(date) = ?)", startMonth, startYear, endMonth, endYear).
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

func (repository *BalanceRepositoryImpl) GetBalanceChangeData(ctx context.Context, shareholderType string, change string, page int, startDate string, endDate string) ([]query_model.BalanceChange, error) {
	db := repository.DB

	quotedCol := quoteIdent(shareholderType)
	prevCol := "prev_" + shareholderType

	const pageSize = 11
	offset := (page - 1) * (pageSize - 1)

	var changeExpr string
	var whereCond string
	if change == "Decrease" {
		changeExpr = fmt.Sprintf(`CASE WHEN t.%[2]s = 0 THEN NULL ELSE ((t.%[2]s - t.%[1]s) / t.%[2]s * 100) END`, quotedCol, prevCol)
		whereCond = fmt.Sprintf("t.%s < t.%s", quotedCol, prevCol)
	} else {
		changeExpr = fmt.Sprintf(`CASE WHEN t.%[2]s = 0 THEN NULL ELSE ((t.%[1]s - t.%[2]s) / t.%[2]s * 100) END`, quotedCol, prevCol)
		whereCond = fmt.Sprintf("t.%s > t.%s", quotedCol, prevCol)
	}

	query := fmt.Sprintf(`
		SELECT
			t.code AS stock_code,
			t.%[1]s AS current_ownership,
			t.%[2]s AS previous_ownership,
			%[3]s AS change_percentage
		FROM (
			SELECT
				s.*,
				LAG(s.%[1]s) OVER (PARTITION BY s.code ORDER BY s.date) AS %[2]s
			FROM stock s
			WHERE DATE_FORMAT(s.date, '%%Y-%%m') IN (?, ?)
		) t
		WHERE t.%[2]s IS NOT NULL AND %[4]s
		ORDER BY change_percentage DESC, t.code
		LIMIT ? OFFSET ?;`, quotedCol, prevCol, changeExpr, whereCond)

	var listStock []query_model.BalanceChange
	err := db.WithContext(ctx).
		Raw(query, startDate, endDate, pageSize, offset).
		Scan(&listStock).
		Error

	if err != nil {
		return nil, domain_error.ErrInternalServer
	}
	return listStock, nil
}
