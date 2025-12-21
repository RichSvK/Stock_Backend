package route

import (
	"backend/controller"
	"backend/repository"
	"backend/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func StockRoute(router *gin.Engine, db *gorm.DB) {
	stockRepository := repository.NewStockRepository(db)
	stockService := service.NewStockService(stockRepository)
	stockController := controller.NewStockController(stockService)
	router.GET("/stock", stockController.SearchStock)
}
