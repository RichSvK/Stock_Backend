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

func TestCreateBroker(t *testing.T) {
	request := request.CreateBrokerRequest{
		Code: "RI",
		Name: "PT Testing Sekuritas",
	}

	headers := map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/json",
	}

	path := "/api/v1/brokers"
	res, status, err := PerformRequest[*response.BrokerResponse](http.MethodPost, path, headers, request)
	require.Nil(t, err)

	assert.Equal(t, http.StatusCreated, status)
	assert.Equal(t, request.Code, res.Code)
	assert.Equal(t, request.Name, res.Name)
}

func TestCreateBrokerInvalidBody(t *testing.T) {
	headers := map[string]string{
		"Accept": "application/json",
	}

	path := "/api/v1/brokers"
	res, status, err := PerformRequest[*response.FailedResponse](http.MethodPost, path, headers, nil)
	require.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, status)
	assert.Equal(t, domainerr.ErrInvalidRequestBody.Error(), res.Message)
}

func TestCreateBrokerInvalidBroker(t *testing.T) {
	request := request.CreateBrokerRequest{
		Code: "",
		Name: "",
	}

	headers := map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/json",
	}

	path := "/api/v1/brokers"
	res, status, err := PerformRequest[*response.FailedResponse](http.MethodPost, path, headers, request)
	require.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, status)
	assert.Equal(t, "Code is required", res.Message)
}

func TestCreateBrokerDuplicate(t *testing.T) {
	request := request.CreateBrokerRequest{
		Code: "FZ",
		Name: "PT Testing Sekuritas",
	}

	headers := map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/json",
	}

	path := "/api/v1/brokers"
	res, status, err := PerformRequest[*response.FailedResponse](http.MethodPost, path, headers, request)
	require.Nil(t, err)

	assert.Equal(t, http.StatusConflict, status)
	assert.Equal(t, "duplicate broker data", res.Message)
}

func TestUpdateBroker(t *testing.T) {
	request := request.UpdateBrokerRequest{
		Code: "RI",
		Name: "PT RI Sekuritas",
	}

	headers := map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/json",
	}

	path := "/api/v1/brokers"
	res, status, err := PerformRequest[*response.SuccessResponse](http.MethodPut, path, headers, request)
	require.Nil(t, err)

	assert.Equal(t, http.StatusOK, status)
	assert.Equal(t, "Success update RI broker", res.Message)
}

func TestUpdateBrokerInvalid(t *testing.T) {
	headers := map[string]string{
		"Accept": "application/json",
	}

	path := "/api/v1/brokers"
	res, status, err := PerformRequest[*response.FailedResponse](http.MethodPut, path, headers, nil)
	require.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, status)
	assert.Equal(t, domainerr.ErrInvalidRequestBody.Error(), res.Message)
}

func TestUpdateBrokerNotFound(t *testing.T) {
	request := request.UpdateBrokerRequest{
		Code: "XX",
		Name: "PT XX Sekuritas",
	}

	headers := map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/json",
	}

	path := "/api/v1/brokers"
	res, status, err := PerformRequest[*response.FailedResponse](http.MethodPut, path, headers, request)
	require.Nil(t, err)

	assert.Equal(t, http.StatusNotFound, status)
	assert.Equal(t, domainerr.ErrBrokerNotFound.Error(), res.Message)
}

func TestDeleteBroker(t *testing.T) {
	headers := map[string]string{
		"Accept": "application/json",
	}

	code := "RI"
	path := "/api/v1/brokers/" + code
	res, status, err := PerformRequest[*response.SuccessResponse](http.MethodDelete, path, headers, nil)
	require.Nil(t, err)

	assert.Equal(t, http.StatusOK, status)
	assert.Equal(t, "Success delete RI broker", res.Message)
}

func TestDeleteBrokerInvalid(t *testing.T) {
	headers := map[string]string{
		"Accept": "application/json",
	}

	code := "R"
	path := "/api/v1/brokers/" + code
	res, status, err := PerformRequest[*response.FailedResponse](http.MethodDelete, path, headers, nil)
	require.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, status)
	assert.Equal(t, "code must be exactly 2 alphabetic characters", res.Message)
}

func TestDeleteBrokerNotFound(t *testing.T) {
	headers := map[string]string{
		"Accept": "application/json",
	}

	code := "XX"
	path := "/api/v1/brokers/" + code
	res, status, err := PerformRequest[*response.FailedResponse](http.MethodDelete, path, headers, nil)
	require.Nil(t, err)

	assert.Equal(t, http.StatusNotFound, status)
	assert.Equal(t, domainerr.ErrBrokerNotFound.Error(), res.Message)
}

func TestGetBroker(t *testing.T) {
	headers := map[string]string{
		"Accept": "application/json",
	}

	path := "/api/v1/brokers"
	res, status, err := PerformRequest[*response.GetBrokerResponse](http.MethodGet, path, headers, nil)
	require.Nil(t, err)

	assert.Equal(t, http.StatusOK, status)
	assert.Equal(t, "Brokers data found", res.Message)
	assert.NotEmpty(t, res.Data)
}

func TestGetBrokerNotFound(t *testing.T) {
	headers := map[string]string{
		"Accept": "application/json",
	}

	broker := "XX"
	path := "/api/v1/brokers?code=" + broker
	res, status, err := PerformRequest[*response.FailedResponse](http.MethodGet, path, headers, nil)
	require.Nil(t, err)

	assert.Equal(t, http.StatusNotFound, status)
	assert.Equal(t, "broker not found", res.Message)
}

func TestGetBrokerInvalidBroker(t *testing.T) {
	headers := map[string]string{
		"Accept": "application/json",
	}

	broker := "X"
	path := "/api/v1/brokers?code=" + broker
	res, status, err := PerformRequest[*response.FailedResponse](http.MethodGet, path, headers, nil)
	require.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, status)
	assert.Equal(t, "Code must be exactly 2 characters long", res.Message)
}
