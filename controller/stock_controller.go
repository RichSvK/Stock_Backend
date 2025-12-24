package controller

import (
	"backend/helper"
	domain_error "backend/model/error"
	"backend/model/web/output"
	"backend/model/web/request"
	"backend/service"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type StockController interface {
	SearchStock(c *gin.Context)
}

type StockControllerImpl struct {
	StockService service.StockService
	Validator    *validator.Validate
}

func NewStockController(stockService service.StockService, validate *validator.Validate) StockController {
	return &StockControllerImpl{
		StockService: stockService,
		Validator:    validate,
	}
}

func (controller *StockControllerImpl) SearchStock(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	var query request.SearchStockQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, output.FailResponse{
			Message: err.Error(),
		})
		return
	}

	if err := controller.Validator.Struct(query); err != nil {
		c.JSON(http.StatusBadRequest, output.FailResponse{
			Message: helper.ValidationError(err),
		})
		return
	}

	response, err := controller.StockService.SearchStock(ctx, query.Code)
	if err != nil {
		c.JSON(MapStockErrorToCode(err), output.FailResponse{
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// Map Stock Error to HTTP Status Code
func MapStockErrorToCode(err error) int {
	switch err {
	case domain_error.ErrStockNotFound:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}
