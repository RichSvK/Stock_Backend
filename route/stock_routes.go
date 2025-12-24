package route

import (
	"backend/controller"
	"backend/repository"
	"backend/service"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

func StockRoute(router *gin.Engine, db *gorm.DB, validate *validator.Validate) {
	stockRepository := repository.NewStockRepository(db)
	stockService := service.NewStockService(stockRepository)
	stockController := controller.NewStockController(stockService, validate)

	router.GET("/stock", stockController.SearchStock)
}
