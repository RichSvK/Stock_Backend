package route

import (
	"backend/controller"
	"backend/repository"
	"backend/service"

	"github.com/gin-gonic/gin"
)

func BalanceRoute(router *gin.Engine) {
	balanceRepository := repository.NewBalanceRepository()
	balanceService := service.NewBalanceService(balanceRepository)
	balanceController := controller.NewBalanceController(balanceService)
	router.GET("/balance/export/:code", balanceController.ExportBalanceController)
	router.GET("/balance/:code", balanceController.GetBalanceChart)
	router.POST("/balance/upload", balanceController.Upload)

	router.GET("/api/auth/balance/scriptless", balanceController.GetScriptlessChange)
	router.GET("/api/auth/balance/change", balanceController.GetBalanceChange)
}
