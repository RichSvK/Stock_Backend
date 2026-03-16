package test

import (
	"backend/internal/model/domainerr"
	"backend/internal/model/response"
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUploadBalance(t *testing.T) {
	file, err := os.Open("./resource/ksei_data/2026_02_February.txt")
	require.Nil(t, err)
	defer func() {
		err = file.Close()
		assert.Nil(t, err)
	}()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", "2026_02_February.txt")
	require.Nil(t, err)

	_, err = io.Copy(part, file)
	require.Nil(t, err)

	err = writer.Close()
	require.Nil(t, err)

	headers := map[string]string{
		"Content-Type": writer.FormDataContentType(),
		"Accept":       "application/json",
	}

	path := "/api/v1/balances/import"
	res, status, err := PerformRequest[*response.UploadBalanceResponse](http.MethodPost, path, headers, body)

	require.Nil(t, err)
	assert.Equal(t, http.StatusCreated, status)
	assert.Equal(t, "Data uploaded and inserted successfully", res.Message)
}

func TestUploadBalanceDuplicate(t *testing.T) {
	file, err := os.Open("./resource/ksei_data/2026_02_February.txt")
	require.Nil(t, err)
	defer func() {
		err = file.Close()
		assert.Nil(t, err)
	}()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", "2026_02_February.txt")
	require.Nil(t, err)

	_, err = io.Copy(part, file)
	require.Nil(t, err)

	err = writer.Close()
	require.Nil(t, err)

	headers := map[string]string{
		"Content-Type": writer.FormDataContentType(),
		"Accept":       "application/json",
	}

	path := "/api/v1/balances/import"
	res, status, err := PerformRequest[*response.FailedResponse](http.MethodPost, path, headers, body)

	require.Nil(t, err)
	assert.Equal(t, http.StatusConflict, status)
	assert.Equal(t, "duplicate balance data exists", res.Message)
}

func TestExportBalance(t *testing.T) {
	headers := map[string]string{
		"Accept": "text/csv",
	}

	path := fmt.Sprintf("/api/v1/balances/%s/export", "NOBU")

	_, status, err := PerformRequest[string](http.MethodGet, path, headers, nil)
	require.Nil(t, err)

	assert.Equal(t, http.StatusOK, status)
}

func TestExportBalanceNotFound(t *testing.T) {
	headers := map[string]string{
		"Accept": "text/csv",
	}

	path := fmt.Sprintf("/api/v1/balances/%s/export", "XXXX")

	res, status, err := PerformRequest[*response.FailedResponse](http.MethodGet, path, headers, nil)
	require.Nil(t, err)

	assert.Equal(t, http.StatusNotFound, status)
	assert.Equal(t, "balance not found", res.Message)
}

func TestGetScriptlessChange(t *testing.T) {
	headers := map[string]string{
		"Accept": "application/json",
	}

	path := "/api/v1/auth/balances/scriptless?start_time=2025-09-01&end_time=2026-02-01"
	res, status, err := PerformRequest[*response.GetScriptlessChangeResponse](http.MethodGet, path, headers, nil)
	require.Nil(t, err)

	assert.Equal(t, http.StatusOK, status)
	assert.Equal(t, "Scriptless change data retrieved successfully", res.Message)
	assert.NotEmpty(t, res.Data)
}

func TestGetScriptlessChangeInvalidDate(t *testing.T) {
	headers := map[string]string{
		"Accept": "application/json",
	}

	path := "/api/v1/auth/balances/scriptless?start_time=2026-01-01"
	res, status, err := PerformRequest[*response.FailedResponse](http.MethodGet, path, headers, nil)
	require.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, status)
	assert.Equal(t, "EndTime is required", res.Message)
}

func TestGetScriptlessChangeInvalidStart(t *testing.T) {
	headers := map[string]string{
		"Accept": "application/json",
	}

	path := "/api/v1/auth/balances/scriptless?start_time=01/01/2026&end_time=02/02/2026"
	res, status, err := PerformRequest[*response.FailedResponse](http.MethodGet, path, headers, nil)
	require.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, status)
	assert.Equal(t, "Invalid start_time format expected YYYY-MM-DD", res.Message)
}

func TestGetScriptlessChangeInvalidEnd(t *testing.T) {
	headers := map[string]string{
		"Accept": "application/json",
	}

	path := "/api/v1/auth/balances/scriptless?start_time=2026-01-01&end_time=02/02/2026"
	res, status, err := PerformRequest[*response.FailedResponse](http.MethodGet, path, headers, nil)
	require.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, status)
	assert.Equal(t, "Invalid end_time format expected YYYY-MM-DD", res.Message)
}

func TestGetScriptlessChangeInvalidDateRange(t *testing.T) {
	headers := map[string]string{
		"Accept": "application/json",
	}

	path := "/api/v1/auth/balances/scriptless?start_time=2027-01-01&end_time=2027-02-01"
	res, status, err := PerformRequest[*response.FailedResponse](http.MethodGet, path, headers, nil)
	require.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, status)
	assert.Equal(t, "invalid date range", res.Message)
}

func TestGetScriptlessChangeNotFound(t *testing.T) {
	headers := map[string]string{
		"Accept": "application/json",
	}

	path := "/api/v1/auth/balances/scriptless?start_time=2025-08-01&end_time=2025-10-01"
	res, status, err := PerformRequest[*response.FailedResponse](http.MethodGet, path, headers, nil)
	require.Nil(t, err)

	assert.Equal(t, http.StatusNotFound, status)
	assert.Equal(t, "balance not found", res.Message)
}

func TestGetBalanceChange(t *testing.T) {
	headers := map[string]string{
		"Accept": "application/json",
	}

	path := "/api/v1/auth/balances/changes?type=local_id&change=Decrease&page=1"
	res, status, err := PerformRequest[*response.BalanceChangeResponse](http.MethodGet, path, headers, nil)
	require.Nil(t, err)

	assert.Equal(t, http.StatusOK, status)
	assert.NotEmpty(t, res.HaveNext)
	assert.NotEmpty(t, res.Data)
}

func TestGetBalanceChangeInvalidQuery(t *testing.T) {
	headers := map[string]string{
		"Accept": "application/json",
	}

	path := "/api/v1/auth/balances/changes?type=local_id&change=Dec&page=1"
	res, status, err := PerformRequest[*response.FailedResponse](http.MethodGet, path, headers, nil)
	require.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, status)
	assert.Equal(t, "Change must be one of [Increase Decrease]", res.Message)
}

func TestGetBalanceChangeInvlidType(t *testing.T) {
	headers := map[string]string{
		"Accept": "application/json",
	}

	path := "/api/v1/auth/balances/changes?type=lokal&change=Decrease&page=1"
	res, status, err := PerformRequest[*response.FailedResponse](http.MethodGet, path, headers, nil)
	require.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, status)
	assert.Equal(t, "invalid shareholder type", res.Message)
}

func TestGetBalanceChartData(t *testing.T) {
	headers := map[string]string{
		"Accept": "application/json",
	}

	stock := "NOBU"
	path := "/api/v1/balances/" + stock

	res, status, err := PerformRequest[*response.GetBalanceResponse](http.MethodGet, path, headers, nil)
	require.Nil(t, err)

	assert.Equal(t, http.StatusOK, status)
	assert.Equal(t, fmt.Sprintf("Balance data for stock code %s retrieved successfully", stock), res.Message)
	assert.NotEmpty(t, res.Data)
}

func TestGetBalanceChartDataNotFound(t *testing.T) {
	headers := map[string]string{
		"Accept": "application/json",
	}

	stock := "XXXX"
	path := "/api/v1/balances/" + stock

	res, status, err := PerformRequest[*response.FailedResponse](http.MethodGet, path, headers, nil)
	require.Nil(t, err)

	assert.Equal(t, http.StatusNotFound, status)
	assert.Equal(t, domainerr.ErrBalanceNotFound.Error(), res.Message)
}

func TestSearchStock(t *testing.T) {
	headers := map[string]string{
		"Accept": "application/json",
	}

	path := "/api/v1/stocks?code=BB"
	res, status, err := PerformRequest[*response.SearchStockResponse](http.MethodGet, path, headers, nil)
	require.Nil(t, err)

	assert.Equal(t, http.StatusOK, status)
	assert.Equal(t, "Stocks found", res.Message)
	assert.NotEmpty(t, res.Data)
}

func TestSearchStockInvaliCodeRequired(t *testing.T) {
	headers := map[string]string{
		"Accept": "application/json",
	}

	path := "/api/v1/stocks?code="
	res, status, err := PerformRequest[*response.FailedResponse](http.MethodGet, path, headers, nil)
	require.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, status)
	assert.Equal(t, "Code is required", res.Message)
}

func TestSearchStockInvaliCodeAlpha(t *testing.T) {
	headers := map[string]string{
		"Accept": "application/json",
	}

	path := "/api/v1/stocks?code=123"
	res, status, err := PerformRequest[*response.FailedResponse](http.MethodGet, path, headers, nil)
	require.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, status)
	assert.Equal(t, "Code must contain only alphabetic characters", res.Message)
}

func TestSearchStockInvaliCodeMax(t *testing.T) {
	headers := map[string]string{
		"Accept": "application/json",
	}

	path := "/api/v1/stocks?code=BBCAA"
	res, status, err := PerformRequest[*response.FailedResponse](http.MethodGet, path, headers, nil)
	require.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, status)
	assert.Equal(t, "Code must be less than or equal to 4", res.Message)
}

func TestSearchStockNotFound(t *testing.T) {
	headers := map[string]string{
		"Accept": "application/json",
	}

	path := "/api/v1/stocks?code=XX"
	res, status, err := PerformRequest[*response.FailedResponse](http.MethodGet, path, headers, nil)
	require.Nil(t, err)

	assert.Equal(t, http.StatusNotFound, status)
	assert.Equal(t, "stock not found", res.Message)
}
