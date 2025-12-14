package route

import (
	"backend/controller"
	"backend/repository"
	"backend/service"

	"github.com/gin-gonic/gin"
)

func StockRoute(router *gin.Engine) {
	stockRepository := repository.NewStockRepository()
	stockService := service.NewStockService(stockRepository)
	stockController := controller.NewStockController(stockService)
	router.GET("/stock", stockController.SearchStock)
}
