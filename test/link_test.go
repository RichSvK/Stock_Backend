package test

import (
	"backend/internal/model/response"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLinkNotFound(t *testing.T) {
	headers := map[string]string{
		"Accept": "application/json",
	}

	path := "/api/v1/links"
	res, status, err := PerformRequest[*response.FailedResponse](http.MethodGet, path, headers, nil)
	require.Nil(t, err)

	assert.Equal(t, http.StatusNotFound, status)
	assert.Equal(t, "link not found", res.Message)
}

func TestGetLinks(t *testing.T) {
	err := ExecuteSQLFile("./resource/sql_data/insert_category.sql")
	require.Nil(t, err)

	err = ExecuteSQLFile("./resource/sql_data/insert_link.sql")
	require.Nil(t, err)

	headers := map[string]string{
		"Accept": "application/json",
	}

	path := "/api/v1/links?category=1"
	res, status, err := PerformRequest[*response.GetLinkResponse](http.MethodGet, path, headers, nil)
	require.Nil(t, err)

	assert.Equal(t, http.StatusOK, status)
	assert.Equal(t, "Link was found", res.Message)
	assert.NotEmpty(t, res.Data)
}