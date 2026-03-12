package route

import (
	"backend/internal/controller"
	"backend/internal/repository"
	"backend/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

func StockWebRoute(router *gin.Engine, db *gorm.DB, validate *validator.Validate) {
	stockWebRepository := repository.NewStockWebRepository(db)
	stockWebService := service.NewStockWebService(stockWebRepository)
	stockWebController := controller.NewStockWebController(stockWebService, validate)

	router.GET("/api/v1/links", stockWebController.GetLinks)
}
