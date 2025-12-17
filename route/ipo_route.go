package route

import (
	"backend/controller"
	"backend/repository"
	"backend/service"

	"github.com/gin-gonic/gin"
)

func IpoRoute(router *gin.Engine) {
	ipoRepository := repository.NewIpoRepository()
	ipoService := service.NewIpoService(ipoRepository)
	ipoController := controller.NewIpoController(ipoService)

	// Get all IPO data
	router.GET("/ipo", ipoController.GetAllIpo)

	// Get IPO by dynamic conditions
	router.POST("/ipo/condition", ipoController.GetIpoByCondition)
}
