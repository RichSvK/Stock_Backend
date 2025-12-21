package route

import (
	"backend/controller"
	"backend/repository"
	"backend/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func BrokerRoute(router *gin.Engine, db *gorm.DB) {
	brokerRepository := repository.NewBrokerRepository(db)
	brokerService := service.NewBrokerService(brokerRepository)
	brokerController := controller.NewBrokerController(brokerService)
	router.GET("/brokers", brokerController.GetBrokers)
}
