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

	balanceGroup := router.Group("/api/v1/balances")
	balanceGroup.GET("/:code/export", balanceController.ExportBalanceController)
	balanceGroup.GET("/:code", balanceController.GetBalanceChart)
	balanceGroup.POST("/upload", balanceController.Upload)

	protectedGroup := router.Group("/api/v1/auth/balances")
	protectedGroup.GET("/scriptless", balanceController.GetScriptlessChange)
	protectedGroup.GET("/change", balanceController.GetBalanceChange)
}
