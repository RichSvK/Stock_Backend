package controller

import (
	"backend/internal/model/domainerr"
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
	GetAllIpo(c *gin.Context)
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

func (controller *IpoControllerImpl) GetAllIpo(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()

	result, err := controller.IpoService.GetIpoAll(ctx)
	if err != nil {
		c.JSON(MapIpoErrorToCode(err), response.FailedResponse{
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (controller *IpoControllerImpl) GetIpoByCondition(c *gin.Context) {
	var request []request.Filter = nil

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

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
