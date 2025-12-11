package config

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var PoolDB *gorm.DB = nil
var SqlDB *sql.DB = nil

func InitializeDBConnection() {
	var err error

	dbUser := os.Getenv("DB_USER")
	dbPass := url.QueryEscape(os.Getenv("DB_PASSWORD"))
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=True",
		dbUser,
		dbPass,
		dbHost,
		dbPort,
		dbName,
	)

	PoolDB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		// Logger:                 logger.Default.LogMode(logger.Info), // Use for debug
		Logger:                 logger.Default.LogMode(logger.Silent),
		PrepareStmt:            true,
		SkipDefaultTransaction: true,
	})

	if err != nil {
		panic("Error")
	}

	SqlDB, err = PoolDB.DB()
	if err != nil {
		panic("Failed to connect")
	}

	if SqlDB.Ping() != nil {
		log.Fatal("Failed to connect")
		return
	}

	SqlDB.SetMaxIdleConns(3)
	SqlDB.SetMaxOpenConns(10)
	SqlDB.SetConnMaxIdleTime(10 * time.Second)
	SqlDB.SetConnMaxLifetime(20 * time.Second)
}

func GetDatabaseInstance() *gorm.DB {
	return PoolDB
}
