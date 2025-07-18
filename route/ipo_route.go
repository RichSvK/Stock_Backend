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

	router.GET("/ipo", ipoController.GetAllIpo)
	router.GET("/ipo/value", ipoController.GetIpoByValue)
	router.GET("/ipo/underwriter/:underwriter", ipoController.GetIpoByUnderwriter)
	router.POST("/ipo/condition", ipoController.GetIpoByCondition)
}
