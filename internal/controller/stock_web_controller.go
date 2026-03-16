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
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type StockWebController interface {
	GetLinks(c *gin.Context)
	CreateLink(c *gin.Context)
	UpdateLink(c *gin.Context)
	DeleteLink(c *gin.Context)
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

	var query query_filter.GetLinkQuery
	_ = c.ShouldBindQuery(&query)
	if err := controller.Validator.Struct(query); err != nil {
		c.JSON(http.StatusBadRequest, response.FailedResponse{
			Message: helper.ValidationError(err),
		})
		return
	}

	result, err := controller.StockWebService.GetLinks(ctx, &query)
	if err != nil {
		c.JSON(MapWebListErrorToCode(err), response.FailedResponse{
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (controller *StockWebControllerImpl) CreateLink(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	var request request.CreateLinkRequest
	_ = c.ShouldBindJSON(&request)
	if err := controller.Validator.Struct(request); err != nil {
		c.JSON(http.StatusBadRequest, response.FailedResponse{
			Message: helper.ValidationError(err),
		})
		return
	}

	result, err := controller.StockWebService.CreateLink(ctx, &request)
	if err != nil {
		c.JSON(MapWebListErrorToCode(err), response.FailedResponse{
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (controller *StockWebControllerImpl) UpdateLink(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	var request request.UpdateLinkRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, response.FailedResponse{
			Message: "invalid request body",
		})
		return
	}

	if err := controller.Validator.Struct(request); err != nil {
		c.JSON(http.StatusBadRequest, response.FailedResponse{
			Message: helper.ValidationError(err),
		})
		return
	}

	result, err := controller.StockWebService.UpdateLink(ctx, &request)
	if err != nil {
		c.JSON(MapWebListErrorToCode(err), response.FailedResponse{
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (controller *StockWebControllerImpl) DeleteLink(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	name := c.Param("name")
	result, err := controller.StockWebService.DeleteLink(ctx, name)
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
	case domainerr.ErrNoFieldsToUpdate:
		return http.StatusBadRequest
	case domainerr.ErrLinkNotFound:
		return http.StatusNotFound
	case domainerr.ErrLinkAlreadyExists:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
