package route

import (
	"backend/controller"
	"backend/repository"
	"backend/service"

	"github.com/gin-gonic/gin"
)

func BrokerRoute(router *gin.Engine) {
	brokerRepository := repository.NewBrokerRepository()
	brokerService := service.NewBrokerService(brokerRepository)
	brokerController := controller.NewBrokerController(brokerService)
	router.GET("/brokers", brokerController.GetBrokers)
}
