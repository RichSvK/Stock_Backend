package test

import (
	"backend/internal/model/response"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetBroker(t *testing.T) {
	err := ExecuteSQLFile("./resource/sql_data/insert_broker.sql")
	require.Nil(t, err)

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