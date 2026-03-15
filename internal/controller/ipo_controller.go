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

type IpoController interface {
	GetIpo(c *gin.Context)
	GetIpoByCondition(c *gin.Context)
}

type IpoControllerImpl struct {
	IpoService service.IpoService
	Validator  *validator.Validate
}

func NewIpoController(ipoService service.IpoService, validate *validator.Validate) IpoController {
	return &IpoControllerImpl{
		IpoService: ipoService,
		Validator:  validate,
	}
}

func (controller *IpoControllerImpl) GetIpo(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()

	var query query_filter.GetIpoQuery
	_ = c.ShouldBindQuery(&query)
	if query.Code != "" { 
		if err := controller.Validator.Struct(query); err != nil {
			c.JSON(http.StatusBadRequest, response.FailedResponse{
				Message: helper.ValidationError(err),
			})
			return
		}
	}

	result, err := controller.IpoService.GetIpo(ctx, query)
	if err != nil {
		c.JSON(MapIpoErrorToCode(err), response.FailedResponse{
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (controller *IpoControllerImpl) GetIpoByCondition(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	var request []request.Filter = nil

	_ = c.ShouldBindJSON(&request)
	if len(request) == 0 {
		c.JSON(http.StatusBadRequest, response.FailedResponse{
			Message: domainerr.ErrEmptyRequest.Error(),
		})
		return
	}

	for _, filter := range request {
		if err := controller.Validator.Struct(filter); err != nil {
			c.JSON(http.StatusBadRequest, response.FailedResponse{
				Message: helper.ValidationError(err),
			})
			return
		}
	}

	result, err := controller.IpoService.GetIpoByCondition(ctx, request)
	if err != nil {
		c.JSON(MapIpoErrorToCode(err), response.FailedResponse{
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

// Map IPO Errror to HTTP Status Code
func MapIpoErrorToCode(err error) int {
	switch err {
	case domainerr.ErrIpoDataNotFound:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}
