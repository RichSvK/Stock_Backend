package controller

import (
	"backend/internal/helper"
	"backend/internal/model/domainerr"
	"backend/internal/model/query_filter"
	"backend/internal/model/response"
	"backend/internal/service"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type BrokerController interface {
	GetBrokers(c *gin.Context)
}

type BrokerControllerImpl struct {
	BrokerService service.BrokerService
	Validator     *validator.Validate
}

func NewBrokerController(brokerService service.BrokerService, validate *validator.Validate) BrokerController {
	return &BrokerControllerImpl{
		BrokerService: brokerService,
		Validator:     validate,
	}
}

func (controller *BrokerControllerImpl) GetBrokers(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()

	var query query_filter.BrokerQuery
	_ = c.ShouldBindQuery(&query)

	// Broker code is optional. If provided then validate it
	if query.Code != "" {
		if err := controller.Validator.Struct(&query); err != nil {
			c.JSON(http.StatusBadRequest, response.FailedResponse{
				Message: helper.ValidationError(err),
			})
			return
		}
	}

	result, err := controller.BrokerService.GetBrokers(ctx, query)
	if err != nil {
		c.JSON(MapBrokerErrorToCode(err), response.FailedResponse{
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

// Map Broker Error to HTTP Status Code
func MapBrokerErrorToCode(err error) int {
	switch err {
	case domainerr.ErrBrokerNotFound:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}
