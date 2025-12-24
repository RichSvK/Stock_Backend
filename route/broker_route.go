package route

import (
	"backend/controller"
	"backend/repository"
	"backend/service"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

func BrokerRoute(router *gin.Engine, db *gorm.DB, validate *validator.Validate) {
	brokerRepository := repository.NewBrokerRepository(db)
	brokerService := service.NewBrokerService(brokerRepository)
	brokerController := controller.NewBrokerController(brokerService, validate)

	router.GET("/brokers", brokerController.GetBrokers)
}
