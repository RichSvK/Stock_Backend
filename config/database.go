package config

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitializeDBConnection() *gorm.DB {
	var err error

	dbUser := os.Getenv("DB_USER")
	dbPass := url.QueryEscape(os.Getenv("DB_PASSWORD"))
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	if dbUser == "" || dbPass == "" || dbHost == "" || dbPort == "" || dbName == "" {
		log.Fatal("Missing required database environment variables")
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=True",
		dbUser,
		dbPass,
		dbHost,
		dbPort,
		dbName,
	)

	PoolDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		// Logger:                 logger.Default.LogMode(logger.Info), // Use for debug
		Logger:                 logger.Default.LogMode(logger.Silent),
		PrepareStmt:            true,
		SkipDefaultTransaction: true,
	})

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	SqlDB, err := PoolDB.DB()
	if err != nil {
		log.Fatalf("Failed to get database instance: %v", err)
	}

	if err := SqlDB.Ping(); err != nil {
		log.Fatal("Failed to connect")
		return nil
	}

	SqlDB.SetMaxIdleConns(3)
	SqlDB.SetMaxOpenConns(10)
	SqlDB.SetConnMaxIdleTime(5 * time.Minute)
	SqlDB.SetConnMaxLifetime(20 * time.Minute)
	log.Println("Database connection established successfully")
	return PoolDB
}
