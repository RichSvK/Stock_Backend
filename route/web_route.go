package route

import (
	"backend/controller"
	"backend/repository"
	"backend/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func StockWebRoute(router *gin.Engine, db *gorm.DB) {
	stockWebRepository := repository.NewStockWebRepository(db)
	stockWebService := service.NewStockWebService(stockWebRepository)
	stockWebController := controller.NewStockWebController(stockWebService)

	router.GET("/links", stockWebController.GetLinks)
}
