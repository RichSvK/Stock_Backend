package controller

import (
	domain_error "backend/model/error"
	"backend/model/web/output"
	"backend/service"
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
	response, err := controller.StockWebService.GetLinks(ctx, category)
	if err != nil {
		c.JSON(MapWebListErrorToCode(err), output.FailResponse{
			Message: err.Error(),
		})
	}

	c.JSON(http.StatusOK, response)
}

// Map Web List Error to HTTP Status Code
func MapWebListErrorToCode(err error) int {
	switch err {
	case domain_error.ErrLinkNotFound:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}
