package main

import (
	"backend/config"
	"backend/route"
	"log"
	"os"
	// _ "net/http/pprof"
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
	// if err := db.AutoMigrate(
	// 	&entity.StockIPO{},
	// 	&entity.Broker{},
	// 	&entity.IPO_Detail{},
	// 	&entity.Stock{},
	// 	&entity.Category{},
	// 	&entity.Link{},
	// ); err != nil {
	// 	log.Fatalf("migration failed: %v", err)
	// }

	// Commented out pprof for production
	// go func() {
	// 	log.Println("pprof listening on :6060")
	// 	log.Println(http.ListenAndServe("localhost:6060", nil))
	// }()

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
