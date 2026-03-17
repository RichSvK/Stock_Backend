package controller

import (
	"backend/internal/helper"
	"backend/internal/model/domainerr"
	"backend/internal/model/query_filter"
	"backend/internal/model/request"
	"backend/internal/model/response"
	"backend/internal/service"
	"context"
	"net/http"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type BrokerController interface {
	GetBrokers(c *gin.Context)
	CreateBroker(c *gin.Context)
	UpdateBroker(c *gin.Context)
	DeleteBroker(c *gin.Context)
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
		if err := controller.Validator.Struct(query); err != nil {
			c.JSON(http.StatusBadRequest, response.FailedResponse{
				Message: helper.ValidationError(err),
			})
			return
		}
	}

	result, err := controller.BrokerService.GetBrokers(ctx, &query)
	if err != nil {
		c.JSON(MapBrokerErrorToCode(err), response.FailedResponse{
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (controller *BrokerControllerImpl) CreateBroker(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()

	var request request.CreateBrokerRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, response.FailedResponse{
			Message: domainerr.ErrInvalidRequestBody.Error(),
		})
		return
	}

	if err := controller.Validator.Struct(request); err != nil {
		c.JSON(http.StatusBadRequest, response.FailedResponse{
			Message: helper.ValidationError(err),
		})
		return
	}

	result, err := controller.BrokerService.CreateBroker(ctx, &request)
	if err != nil {
		c.JSON(MapBrokerErrorToCode(err), response.FailedResponse{
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, result)
}

func (controller *BrokerControllerImpl) UpdateBroker(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	var request request.UpdateBrokerRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, response.FailedResponse{
			Message: domainerr.ErrInvalidRequestBody.Error(),
		})
		return
	}

	if err := controller.Validator.Struct(request); err != nil {
		c.JSON(http.StatusBadRequest, response.FailedResponse{
			Message: helper.ValidationError(err),
		})
		return
	}

	result, err := controller.BrokerService.UpdateBroker(ctx, &request)
	if err != nil {
		c.JSON(MapBrokerErrorToCode(err), response.FailedResponse{
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (controller *BrokerControllerImpl) DeleteBroker(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	code := c.Param("code")
	codeRegex := regexp.MustCompile("^[A-Z]{2}$")
	if !codeRegex.MatchString(code) {
		c.JSON(http.StatusBadRequest, response.FailedResponse{
			Message: "code must be exactly 2 alphabetic characters",
		})
		return
	}

	result, err := controller.BrokerService.DeleteBroker(ctx, code)
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
	case domainerr.ErrBrokerDuplicate:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
