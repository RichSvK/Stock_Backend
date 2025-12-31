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

	// Get all IPO data
	router.GET("/ipo", ipoController.GetAllIpo)

	// Get IPO by dynamic conditions
	router.POST("/ipo/condition", ipoController.GetIpoByCondition)
}
