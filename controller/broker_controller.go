package controller

import (
	"backend/service"
	"context"
	"time"

	"github.com/gin-gonic/gin"
)

type BrokerController interface {
	GetBrokers(c *gin.Context)
}

type BrokerControllerImpl struct {
	BrokerService service.BrokerService
}

func NewBrokerController(brokerService service.BrokerService) BrokerController {
	return &BrokerControllerImpl{
		BrokerService: brokerService,
	}
}

func (controller *BrokerControllerImpl) GetBrokers(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()
	status, response := controller.BrokerService.GetBrokers(ctx)
	c.JSON(status, response)
}
