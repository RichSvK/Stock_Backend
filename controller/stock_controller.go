package controller

import (
	"backend/service"
	"context"
	"time"

	"github.com/gin-gonic/gin"
)

type StockController interface {
	SearchStock(c *gin.Context)
}

type StockControllerImpl struct {
	StockService service.StockService
}

func NewStockController(stockService service.StockService) StockController {
	return &StockControllerImpl{
		StockService: stockService,
	}
}

func (controller *StockControllerImpl) SearchStock(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	codeParam := c.Query("code")

	status, response := controller.StockService.SearchStock(ctx, codeParam)
	c.JSON(status, response)
}
