package route

import (
	"backend/controller"
	"backend/repository"
	"backend/service"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

func BalanceRoute(router *gin.Engine, db *gorm.DB, validate *validator.Validate) {
	balanceRepository := repository.NewBalanceRepository(db)
	balanceService := service.NewBalanceService(balanceRepository)
	balanceController := controller.NewBalanceController(balanceService, validate)

	router.GET("/balance/export/:code", balanceController.ExportBalanceController)
	router.GET("/balance/:code", balanceController.GetBalanceChart)
	router.POST("/balance/upload", balanceController.Upload)

	router.GET("/api/auth/balance/scriptless", balanceController.GetScriptlessChange)
	router.GET("/api/auth/balance/change", balanceController.GetBalanceChange)
}
