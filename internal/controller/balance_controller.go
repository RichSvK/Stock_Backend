package controller

import (
	"backend/internal/helper"
	"backend/internal/model/domainerr"
	"backend/internal/model/query_filter"
	"backend/internal/model/request"
	"backend/internal/model/response"
	"backend/internal/service"
	"context"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type BalanceController interface {
	Upload(c *gin.Context)
	ExportBalanceData(c *gin.Context)
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

func (controller *BalanceControllerImpl) ExportBalanceData(c *gin.Context) {
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

	var query query_filter.ScriptlessChangeQuery
	_ = c.ShouldBindQuery(&query)
	if err := controller.Validator.Struct(query); err != nil {
		c.JSON(http.StatusBadRequest, response.FailedResponse{
			Message: helper.ValidationError(err),
		})
		return
	}

	start, err := time.Parse("2006-01-02", query.StartTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.FailedResponse{
			Message: "Invalid start_time format expected YYYY-MM-DD",
		})
		return
	}

	end, err := time.Parse("2006-01-02", query.EndTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.FailedResponse{
			Message: "Invalid end_time format expected YYYY-MM-DD",
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
	_ = c.ShouldBindQuery(&query)
	if err := controller.Validator.Struct(query); err != nil {
		c.JSON(http.StatusBadRequest, response.FailedResponse{
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
	case domainerr.ErrInvalidShareholderType, domainerr.ErrInvalidDateRange:
		return http.StatusBadRequest
	case domainerr.ErrDuplicateBalanceData:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
