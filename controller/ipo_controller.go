package controller

import (
	"backend/service"
	"context"
	"time"

	"github.com/gin-gonic/gin"
)

type IpoController interface {
	GetAllIpo(c *gin.Context)
	GetIpoByUnderwriter(c *gin.Context)
	GetIpoByValue(c *gin.Context)
}

type IpoControllerImpl struct {
	IpoService service.IpoService
}

func NewIpoController(ipoService service.IpoService) IpoController {
	return &IpoControllerImpl{
		IpoService: ipoService,
	}
}

func (controller *IpoControllerImpl) GetAllIpo(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()
	status, response := controller.IpoService.GetIpoAll(ctx)
	c.JSON(status, response)
}

func (controller *IpoControllerImpl) GetIpoByUnderwriter(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()
	status, response := controller.IpoService.GetIpoByUnderwriter(ctx, c.Param("underwriter"))
	c.JSON(status, response)
}

func (controller *IpoControllerImpl) GetIpoByValue(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()
	status, response := controller.IpoService.GetIpoByValue(ctx, c.Query("value"), c.Query("underwriter"))
	c.JSON(status, response)
}
