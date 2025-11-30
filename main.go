package main

import (
	"backend/config"
	"backend/model/entity"
	"backend/route"
	"log"
)

func init() {
	config.InitEnvironment()
	config.InitializeDBConnection()
	config.MakeFolder("Resource")
}

func main() {
	db := config.GetDatabaseInstance()

	// AutoMigrate
	if err := db.AutoMigrate(
		&entity.StockIPO{},
		&entity.Broker{},
		&entity.IPO_Detail{},
		&entity.Stock{},
		&entity.Category{},
		&entity.Link{},
	); err != nil {
		log.Fatalf("migration failed: %v", err)
	}

	router := route.SetupRouter()

	// Run the server
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
