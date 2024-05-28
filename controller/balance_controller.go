package controller

import (
	"backend/helper"
	"backend/model/web/output"
	"backend/service"
	"context"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

type BalanceController interface {
	Upload(c *gin.Context)
	ExportBalanceController(c *gin.Context)
	GetBalanceChart(c *gin.Context)
}

type BalanceControllerImpl struct {
	BalanceService service.BalanceService
}

func NewBalanceController(balanceService service.BalanceService) BalanceController {
	return &BalanceControllerImpl{
		BalanceService: balanceService,
	}
}

func (controller *BalanceControllerImpl) Upload(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Second)
	defer cancel()
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 10<<20) // 10 MB

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, output.FailedResponse{})
		return
	}

	filePath := filepath.Join("Resource", file.Filename)
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving the file"})
		return
	}

	code, output := controller.BalanceService.Create(ctx, file.Filename)
	c.JSON(code, output)
	err = helper.RemoveFile(filePath)
	helper.PanicIfError(err)
}

func (controller *BalanceControllerImpl) ExportBalanceController(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Second)
	defer cancel()

	stockCode := c.Query("code")
	status, output := controller.BalanceService.ExportCode(ctx, stockCode)
	if status != http.StatusOK {
		c.Header("Content-Type", "application/json")
		c.JSON(status, output)
		return
	}

	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", "attachment; filename=\""+stockCode+".csv\"")
	filePath := "./Resource/" + stockCode + ".csv"
	c.File(filePath)
	err := helper.RemoveFile(filePath)
	helper.PanicIfError(err)
}

func (controller *BalanceControllerImpl) GetBalanceChart(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Second)
	defer cancel()
	status, output := controller.BalanceService.GetBalanceData(ctx, c.Param("code"))
	c.JSON(status, output)
}
