package config

import (
	"database/sql"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var PoolDB *gorm.DB = nil
var SqlDB *sql.DB = nil

func InitializeDBConnection() {
	dsn := os.Getenv("DB_URL")
	var err error

	PoolDB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger:                 logger.Default.LogMode(logger.Info),
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
