package route

import (
	"backend/controller"
	"backend/repository"
	"backend/service"

	"github.com/gin-gonic/gin"
)

func StockWebRoute(router *gin.Engine) {
	stockWebRepository := repository.NewStockWebRepository()
	stockWebService := service.NewStockWebService(stockWebRepository)
	stockWebController := controller.NewStockWebController(stockWebService)
	router.GET("/links", stockWebController.GetLinks)
}
