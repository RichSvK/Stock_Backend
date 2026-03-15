package test

import (
	"backend/config"
	"backend/internal/entity"
	"backend/route"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var app *gin.Engine
var db *gorm.DB

func init() {
	config.InitEnvironment("test.env")
	db = config.InitializeDBConnection()
	app = route.SetupRouter(db)

	config.MakeFolder("resource")

	if err := db.AutoMigrate(
		&entity.StockIPO{},
		&entity.Broker{},
		&entity.IpoDetail{},
		&entity.Stock{},
		&entity.Category{},
		&entity.Link{},
	); err != nil {
		log.Fatalf("migration failed: %v", err)
	}

	table := []string{"ipo_detail", "stock_ipo", "stock", "broker_underwriter", "link", "category"}
	for _, t := range table {
		if err := ClearTable(t); err != nil {
			log.Fatalf("failed to clear table %s: %v", t, err)
		}
	}

	insertData := []string{
		"./resource/ksei_data/2025_09_September.txt",
		"./resource/ksei_data/2025_11_November.txt",
		"./resource/ksei_data/2025_12_December.txt",
		"./resource/ksei_data/2026_01_January.txt",
	}

	for _, file := range insertData {
		if err := InsertTestStockData(file); err != nil {
			log.Fatalf("failed to insert test stock data from file %s: %v", file, err)
		}
	}

	if	err := ExecuteSQLFile("./resource/start/insert_test.sql"); err != nil {
		log.Fatalf("failed to execute test SQL file: %v", err)
	}
}
