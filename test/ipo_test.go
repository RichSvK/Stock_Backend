package test

import (
	"backend/internal/model/domainerr"
	"backend/internal/model/request"
	"backend/internal/model/response"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetIpo(t *testing.T) {
	header := map[string]string{
		"Accept": "application/json",
	}

	code := "STRK"
	path := "/api/v1/ipo?code=" + code
	res, status, err := PerformRequest[*response.GetIpoResponse](http.MethodGet, path, header, nil)
	require.Nil(t, err)

	assert.Equal(t, http.StatusOK, status)
	assert.Equal(t, "IPO data found", res.Message)
	assert.NotEmpty(t, res.Data)
}

func TestGetIpoInvalidCode(t *testing.T) {
	header := map[string]string{
		"Accept": "application/json",
	}
	
	code := "ADMR1"
	path := "/api/v1/ipo?code=" + code
	res, status, err := PerformRequest[*response.FailedResponse](http.MethodGet, path, header, nil)
	require.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, status)
	assert.Equal(t, "Code must contain only alphabetic characters", res.Message)
}

func TestGetIpoByCondition(t *testing.T) {
	header := map[string]string{
		"Content-Type": "application/json",
		"Accept":       "application/json",
	}

	requestBody := []request.Filter{
		{
			FilterName:  "uw_code",
			FilterValue: "FZ",
			Symbol:      "=",
			FilterType:  "string",
		},
	}

	path := "/api/v1/ipo/condition"
	res, status, err := PerformRequest[*response.GetIpoResponse](http.MethodPost, path, header, requestBody)
	require.Nil(t, err)

	assert.Equal(t, http.StatusOK, status)
	assert.Equal(t, "IPO data found", res.Message)
	assert.NotEmpty(t, res.Data)
}

func TestGetIpoEmptyFilter(t *testing.T) {
	header := map[string]string{
		"Accept":       "application/json",
	}

	path := "/api/v1/ipo/condition"
	res, status, err := PerformRequest[*response.FailedResponse](http.MethodPost, path, header, nil)
	require.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, status)
	assert.Equal(t, domainerr.ErrInvalidRequestBody.Error(), res.Message)
}

func TestGetIpoInvalidCondition(t *testing.T) {
	header := map[string]string{
		"Content-Type": "application/json",
		"Accept":       "application/json",
	}

	requestBody := []request.Filter{
		{
			FilterName:  "uw_code",
			FilterValue: "KI",
			Symbol:      "",
			FilterType:  "string",
		},
	}

	path := "/api/v1/ipo/condition"
	res, status, err := PerformRequest[*response.FailedResponse](http.MethodPost, path, header, requestBody)
	require.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, status)
	assert.Equal(t, "Symbol is required", res.Message)
}

func TestGetIpoNotFoundn(t *testing.T) {
	header := map[string]string{
		"Content-Type": "application/json",
		"Accept":       "application/json",
	}

	requestBody := []request.Filter{
		{
			FilterName:  "uw_code",
			FilterValue: "XX",
			Symbol:      "=",
			FilterType:  "string",
		},
	}

	path := "/api/v1/ipo/condition"
	res, status, err := PerformRequest[*response.FailedResponse](http.MethodPost, path, header, requestBody)
	require.Nil(t, err)

	assert.Equal(t, http.StatusNotFound, status)
	assert.Equal(t, "ipo data not found", res.Message)
}

func TestGetIpoNotFound(t *testing.T) {
	header := map[string]string{
		"Accept": "application/json",
	}

	path := "/api/v1/ipo?code=XXXX"
	res, status, err := PerformRequest[*response.FailedResponse](http.MethodGet, path, header, nil)
	require.Nil(t, err)

	assert.Equal(t, http.StatusNotFound, status)
	assert.Equal(t, "ipo data not found", res.Message)
}
