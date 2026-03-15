package domainerr

import (
	"errors"
	"strings"

	"github.com/go-sql-driver/mysql"
)

func IsDuplicateError(err error) bool {
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		return mysqlErr.Number == 1062
	}
	return strings.Contains(err.Error(), "duplicate key")
}