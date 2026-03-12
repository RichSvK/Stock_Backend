package route

import (
	"backend/controller"
	"backend/repository"
	"backend/service"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

func IpoRoute(router *gin.Engine, db *gorm.DB, validate *validator.Validate) {
	ipoRepository := repository.NewIpoRepository(db)
	ipoService := service.NewIpoService(ipoRepository)
	ipoController := controller.NewIpoController(ipoService, validate)

	ipoGroup := router.Group("/api/v1/ipo")
	// Get all IPO data
	ipoGroup.GET("", ipoController.GetAllIpo)

	// Get IPO by dynamic conditions
	ipoGroup.POST("/condition", ipoController.GetIpoByCondition)
}
