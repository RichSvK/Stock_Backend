package route

import (
	"backend/controller"
	"backend/repository"
	"backend/service"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

func StockWebRoute(router *gin.Engine, db *gorm.DB, validate *validator.Validate) {
	stockWebRepository := repository.NewStockWebRepository(db)
	stockWebService := service.NewStockWebService(stockWebRepository)
	stockWebController := controller.NewStockWebController(stockWebService, validate)

	router.GET("/links", stockWebController.GetLinks)
}
