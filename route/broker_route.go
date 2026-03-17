package route

import (
	"backend/internal/controller"
	"backend/internal/repository"
	"backend/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

func BrokerRoute(router *gin.Engine, db *gorm.DB, validate *validator.Validate) {
	brokerRepository := repository.NewBrokerRepository(db)
	brokerService := service.NewBrokerService(brokerRepository)
	brokerController := controller.NewBrokerController(brokerService, validate)

	brokerRoute := router.Group("/api/v1/brokers")
	brokerRoute.GET("", brokerController.GetBrokers)
	brokerRoute.POST("", brokerController.CreateBroker)
	brokerRoute.PUT("", brokerController.UpdateBroker)
	brokerRoute.DELETE("/:code", brokerController.DeleteBroker)
}
