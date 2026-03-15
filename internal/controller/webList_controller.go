package controller

import (
	"backend/internal/model/domainerr"
	"backend/internal/model/response"
	"backend/internal/service"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type StockWebController interface {
	GetLinks(c *gin.Context)
}

type StockWebControllerImpl struct {
	StockWebService service.StockWebService
	Validator       *validator.Validate
}

func NewStockWebController(stockWebService service.StockWebService, validate *validator.Validate) StockWebController {
	return &StockWebControllerImpl{
		StockWebService: stockWebService,
		Validator:       validate,
	}
}

func (controller *StockWebControllerImpl) GetLinks(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	category := c.Query("category")
	result, err := controller.StockWebService.GetLinks(ctx, category)
	if err != nil {
		c.JSON(MapWebListErrorToCode(err), response.FailedResponse{
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

// Map Web List Error to HTTP Status Code
func MapWebListErrorToCode(err error) int {
	switch err {
	case domainerr.ErrLinkNotFound:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}
