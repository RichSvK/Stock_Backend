package controller

import (
	"backend/service"
	"context"
	"time"

	"github.com/gin-gonic/gin"
)

type StockWebController interface {
	GetLinks(c *gin.Context)
}

type StockWebControllerImpl struct {
	StockWebService service.StockWebService
}

func NewStockWebController(stockWebService service.StockWebService) StockWebController {
	return &StockWebControllerImpl{
		StockWebService: stockWebService,
	}
}

func (controller *StockWebControllerImpl) GetLinks(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()
	status, output := controller.StockWebService.GetLinks(ctx, c.Query("category"))
	c.JSON(status, output)
}
