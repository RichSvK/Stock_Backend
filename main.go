package main

import (
	"backend/config"
	"backend/model/entity"
	"backend/route"
	"log"
	"os"
)

func main() {
	// Initialize local environment variables
	config.InitEnvironment()

	db := config.InitializeDBConnection()
	defer func() {
		sqlDB, _ := db.DB()
		_ = sqlDB.Close()
	}()

	// Create Resource folder if not exists
	config.MakeFolder("Resource")

	// AutoMigrate database tables (commented out for production)
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

	router := route.SetupRouter(db)

	// Use PORT env or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
