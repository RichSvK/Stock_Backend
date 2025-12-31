package controller

import (
	"backend/helper"
	domain_error "backend/model/error"
	"backend/model/web/output"
	"backend/model/web/request"
	"backend/service"
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
		c.JSON(MapBalanceErrorToCode(domain_error.ErrUploadFailed), output.FailResponse{
			Message: domain_error.ErrUploadFailed.Error(),
		})
		return
	}

	filePath := filepath.Join("Resource", file.Filename)
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(MapBalanceErrorToCode(domain_error.ErrSaveFile), output.FailResponse{
			Message: domain_error.ErrSaveFile.Error(),
		})
		return
	}

	defer func() {
		if err := helper.RemoveFile(filePath); err != nil {
			log.Printf("Warning: failed to remove temporary file %s: %v", filePath, err)
		}
	}()

	response, err := controller.BalanceService.Create(ctx, file.Filename)
	if err != nil {
		c.JSON(MapBalanceErrorToCode(err), output.FailResponse{
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, response)
}

func (controller *BalanceControllerImpl) ExportBalanceController(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	stockCode := c.Param("code")
	if err := controller.BalanceService.ExportCode(ctx, stockCode); err != nil {
		c.JSON(MapBalanceErrorToCode(err), output.FailResponse{
			Message: err.Error(),
		})
		return
	}

	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", "attachment; filename=\""+stockCode+".csv\"")
	filePath := "./Resource/" + stockCode + ".csv"
	c.File(filePath)

	if err := helper.RemoveFile(filePath); err != nil {
		log.Printf("Warning: failed to remove temporary file %s: %v", filePath, err)
	}
}

func (controller *BalanceControllerImpl) GetBalanceChart(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	response, err := controller.BalanceService.GetBalanceData(ctx, c.Param("code"))
	if err != nil {
		c.JSON(MapBalanceErrorToCode(err), output.FailResponse{
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (controller *BalanceControllerImpl) GetScriptlessChange(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	startTime := c.Query("start_time")
	endTime := c.Query("end_time")

	if startTime == "" || endTime == "" {
		c.JSON(http.StatusBadRequest, output.FailResponse{
			Message: "Missing required parameters: start_time and end_time are required",
		})
		return
	}

	start, err := time.Parse("2006-01-02", startTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, output.FailResponse{
			Message: fmt.Sprintf("Invalid start_time format '%s': expected YYYY-MM-DD", startTime),
		})
		return
	}

	end, err := time.Parse("2006-01-02", endTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, output.FailResponse{
			Message: fmt.Sprintf("Invalid end_time format '%s': expected YYYY-MM-DD", startTime),
		})
		return
	}

	response, err := controller.BalanceService.GetScriptlessChange(ctx, start, end)
	if err != nil {
		c.JSON(MapBalanceErrorToCode(err), output.FailResponse{
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (controller *BalanceControllerImpl) GetBalanceChange(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	var query request.GetBalanceChangeQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, output.FailResponse{
			Message: "Invalid bad request",
		})
		return
	}

	if err := controller.Validator.Struct(query); err != nil {
		c.JSON(MapBalanceErrorToCode(domain_error.ErrPaginationInvalid), output.FailResponse{
			Message: helper.ValidationError(err),
		})
		return
	}

	response, err := controller.BalanceService.GetBalanceChangeData(ctx, query.Type, query.Change, query.Page)
	if err != nil {
		c.JSON(MapBalanceErrorToCode(err), output.FailResponse{
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// Map Balance Error to HTTP Status Code
func MapBalanceErrorToCode(err error) int {
	switch err {
	case domain_error.ErrBalanceNotFound:
		return http.StatusNotFound
	case domain_error.ErrPaginationInvalid, domain_error.ErrInvalidShareholderType:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
