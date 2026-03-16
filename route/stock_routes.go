package route

import (
	"backend/internal/controller"
	"backend/internal/repository"
	"backend/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

func StockRoute(router *gin.Engine, db *gorm.DB, validate *validator.Validate) {
	stockRepository := repository.NewStockRepository(db)
	stockService := service.NewStockService(stockRepository)
	stockController := controller.NewStockController(stockService, validate)

	router.GET("/api/v1/stocks", stockController.SearchStock)
}