package controller

import (
	"backend/internal/helper"
	"backend/internal/model/domainerr"
	"backend/internal/model/request"
	"backend/internal/model/response"
	"backend/internal/service"
	"context"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type BalanceController interface {
	Upload(c *gin.Context)
	ExportBalanceController(c *gin.Context)
	GetBalanceChart(c *gin.Context)
	GetScriptlessChange(c *gin.Context)
	GetBalanceChange(c *gin.Context)
}

type BalanceControllerImpl struct {
	BalanceService service.BalanceService
	Validator      *validator.Validate
}

func NewBalanceController(balanceService service.BalanceService, validate *validator.Validate) BalanceController {
	return &BalanceControllerImpl{
		BalanceService: balanceService,
		Validator:      validate,
	}
}

func (controller *BalanceControllerImpl) Upload(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 10<<20) // 10 MB

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(MapBalanceErrorToCode(domainerr.ErrUploadFailed), response.FailedResponse{
			Message: domainerr.ErrUploadFailed.Error(),
		})
		return
	}

	filePath := filepath.Join("resource", file.Filename)
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(MapBalanceErrorToCode(domainerr.ErrSaveFile), response.FailedResponse{
			Message: domainerr.ErrSaveFile.Error(),
		})
		return
	}

	defer func() {
		if err := helper.RemoveFile(filePath); err != nil {
			log.Printf("Warning: failed to remove temporary file %s: %v", filePath, err)
		}
	}()

	result, err := controller.BalanceService.Create(ctx, file.Filename)
	if err != nil {
		c.JSON(MapBalanceErrorToCode(err), response.FailedResponse{
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, result)
}

func (controller *BalanceControllerImpl) ExportBalanceController(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	stockCode := c.Param("code")
	if err := controller.BalanceService.ExportCode(ctx, stockCode); err != nil {
		c.JSON(MapBalanceErrorToCode(err), response.FailedResponse{
			Message: err.Error(),
		})
		return
	}

	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", "attachment; filename=\""+stockCode+".csv\"")
	filePath := "./resource/" + stockCode + ".csv"
	c.File(filePath)

	if err := helper.RemoveFile(filePath); err != nil {
		log.Printf("Warning: failed to remove temporary file %s: %v", filePath, err)
	}
}

func (controller *BalanceControllerImpl) GetBalanceChart(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	result, err := controller.BalanceService.GetBalanceData(ctx, c.Param("code"))
	if err != nil {
		c.JSON(MapBalanceErrorToCode(err), response.FailedResponse{
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (controller *BalanceControllerImpl) GetScriptlessChange(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	startTime := c.Query("start_time")
	endTime := c.Query("end_time")

	if startTime == "" || endTime == "" {
		c.JSON(http.StatusBadRequest, response.FailedResponse{
			Message: "Missing required parameters: start_time and end_time are required",
		})
		return
	}

	start, err := time.Parse("2006-01-02", startTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.FailedResponse{
			Message: fmt.Sprintf("Invalid start_time format '%s': expected YYYY-MM-DD", startTime),
		})
		return
	}

	end, err := time.Parse("2006-01-02", endTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.FailedResponse{
			Message: fmt.Sprintf("Invalid end_time format '%s': expected YYYY-MM-DD", startTime),
		})
		return
	}

	result, err := controller.BalanceService.GetScriptlessChange(ctx, start, end)
	if err != nil {
		c.JSON(MapBalanceErrorToCode(err), response.FailedResponse{
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (controller *BalanceControllerImpl) GetBalanceChange(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	var query request.GetBalanceChangeQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, response.FailedResponse{
			Message: "Invalid bad request",
		})
		return
	}

	if err := controller.Validator.Struct(query); err != nil {
		c.JSON(MapBalanceErrorToCode(domainerr.ErrPaginationInvalid), response.FailedResponse{
			Message: helper.ValidationError(err),
		})
		return
	}

	result, err := controller.BalanceService.GetBalanceChangeData(ctx, query.Type, query.Change, query.Page)
	if err != nil {
		c.JSON(MapBalanceErrorToCode(err), response.FailedResponse{
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

// Map Balance Error to HTTP Status Code
func MapBalanceErrorToCode(err error) int {
	switch err {
	case domainerr.ErrBalanceNotFound:
		return http.StatusNotFound
	case domainerr.ErrPaginationInvalid, domainerr.ErrInvalidShareholderType:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
