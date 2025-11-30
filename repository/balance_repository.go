package repository

import (
	"backend/config"
	"backend/model/entity"
	"backend/model/web/response"
	"context"
	"fmt"
	"strconv"
	"time"
)

type BalanceRepository interface {
	Create(ctx context.Context, stock []entity.Stock) error
	GetBalanceStock(ctx context.Context, code string) ([]entity.Stock, error)
	GetScriptlessChange(ctx context.Context, startDate time.Time, endDate time.Time) ([]entity.Stock, error)
	GetBalanceChangeData(ctx context.Context, shareholderType string, change string, page string, startDate string, endDate string) ([]response.BalanceChangeData, error)
}

type BalanceRepositoryImpl struct{}

func NewBalanceRepository() BalanceRepository {
	return &BalanceRepositoryImpl{}
}

func (repository *BalanceRepositoryImpl) Create(ctx context.Context, stock []entity.Stock) error {
	db := config.GetDatabaseInstance()
	tx := db.WithContext(ctx).Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Create(&stock).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func (repository *BalanceRepositoryImpl) GetBalanceStock(ctx context.Context, code string) ([]entity.Stock, error) {
	db := config.GetDatabaseInstance()

	var listStock []entity.Stock = nil
	err := db.WithContext(ctx).
		Model(&entity.Stock{}).
		Select("stock.*").
		Where("code = ?", code).
		Order("Date DESC").
		Limit(6).
		Scan(&listStock).
		Error

	return listStock, err
}

func (repository *BalanceRepositoryImpl) GetScriptlessChange(ctx context.Context, startDate time.Time, endDate time.Time) ([]entity.Stock, error) {
	db := config.GetDatabaseInstance()

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

	return listStock, err
}

func quoteIdent(col string) string {
	// Use backticks for MySQL identifiers.
	// Since col is validated from whitelist, this is safe.
	return "`" + col + "`"
}

func (repository *BalanceRepositoryImpl) GetBalanceChangeData(ctx context.Context, shareholderType string, change string, page string, startDate string, endDate string) ([]response.BalanceChangeData, error) {
	db := config.GetDatabaseInstance()
	var listStock []response.BalanceChangeData

	var AllowedColumns = map[string]bool{
		"local_is": true,
		"local_cp": true,
		"local_pf": true,
		"local_ib": true,
		"local_id": true,
		"local_mf": true,
		"local_sc": true,
		"local_fd": true,
		"local_ot": true,

		"foreign_is": true,
		"foreign_cp": true,
		"foreign_pf": true,
		"foreign_ib": true,
		"foreign_id": true,
		"foreign_mf": true,
		"foreign_sc": true,
		"foreign_fd": true,
		"foreign_ot": true,
	}

	if !AllowedColumns[shareholderType] {
		return nil, fmt.Errorf("invalid column name: %s", shareholderType)
	}

	quotedCol := quoteIdent(shareholderType)
	prevCol := "prev_" + shareholderType

	const pageSize = 11
	var offset int
	if page != "" {
		if pNum, err := strconv.Atoi(page); err == nil && pNum > 0 {
			offset = (pNum - 1) * (pageSize - 1)
		}
	}

	var changeExpr string
	var whereCond string
	if change == "Decrease" {
		changeExpr = fmt.Sprintf(
			`CASE WHEN t.%[2]s = 0 THEN NULL ELSE ((t.%[2]s - t.%[1]s) / t.%[2]s * 100) END`,
			quotedCol, prevCol)
		whereCond = fmt.Sprintf("t.%s < t.%s", quotedCol, prevCol)
	} else {
		changeExpr = fmt.Sprintf(
			`CASE WHEN t.%[2]s = 0 THEN NULL ELSE ((t.%[1]s - t.%[2]s) / t.%[2]s * 100) END`,
			quotedCol, prevCol)
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

	if err := db.WithContext(ctx).
		Raw(query, startDate, endDate, pageSize, offset).
		Scan(&listStock).Error; err != nil {
		return nil, err
	}
	return listStock, nil
}
