package test

import (
	"backend/internal/model/request"
	"backend/internal/model/response"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateLink(t *testing.T) {
	request := request.CreateLinkRequest{
		URL:         "https://www.example.com",
		Name:        "Test Link",
		Image:       "https://www.example.com/image.jpg",
		Description: "This is an example description",
		Category:    1,
	}

	headers := map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/json",
	}

	path := "/api/v1/links"
	res, status, err := PerformRequest[*response.LinkResponse](http.MethodPost, path, headers, request)
	require.Nil(t, err)

	assert.Equal(t, http.StatusOK, status)
	assert.Equal(t, request.URL, res.URL)
	assert.Equal(t, request.Name, res.Name)
	assert.Equal(t, request.Image, res.Image)
	assert.Equal(t, request.Description, res.Description)
}

func TestCreateLinkBadRequest(t *testing.T) {
	request := request.CreateLinkRequest{}

	headers := map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/json",
	}

	path := "/api/v1/links"
	res, status, err := PerformRequest[*response.FailedResponse](http.MethodPost, path, headers, request)
	require.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, status)
	assert.Equal(t, "URL is required", res.Message)
}

func TestCreateLinkDuplicate(t *testing.T) {
	request := request.CreateLinkRequest{
		URL:         "https://www.example.com",
		Name:        "Test Link",
		Image:       "https://www.example.com/image.jpg",
		Description: "This is an example link",
		Category:    1,
	}

	headers := map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/json",
	}

	path := "/api/v1/links"
	res, status, err := PerformRequest[*response.FailedResponse](http.MethodPost, path, headers, request)
	require.Nil(t, err)

	assert.Equal(t, http.StatusConflict, status)
	assert.Equal(t, "link already exists", res.Message)
}

func TestLinkNotFound(t *testing.T) {
	headers := map[string]string{
		"Accept": "application/json",
	}

	path := "/api/v1/links?name=Invalid%20Link"
	res, status, err := PerformRequest[*response.FailedResponse](http.MethodGet, path, headers, nil)
	require.Nil(t, err)

	assert.Equal(t, http.StatusNotFound, status)
	assert.Equal(t, "link not found", res.Message)
}

func TestGetLinks(t *testing.T) {
	headers := map[string]string{
		"Accept": "application/json",
	}

	categoryId := 1
	name := "Test Link"
	path := fmt.Sprintf("/api/v1/links?category_id=%d&name=%s", categoryId, strings.ReplaceAll(name, " ", "%20"))
	res, status, err := PerformRequest[*response.GetLinkResponse](http.MethodGet, path, headers, nil)
	require.Nil(t, err)

	assert.Equal(t, http.StatusOK, status)
	assert.Equal(t, "Link was found", res.Message)
	assert.NotEmpty(t, res.Data)
}

func TestUpdateLink(t *testing.T) {
	request := request.UpdateLinkRequest{
		Name:        "Test Link",
		URL:         "https://www.test.com",
		Description: "This is an updated test link",
		Image: "update.png",
		Category: 2,
	}

	headers := map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/json",
	}

	path := "/api/v1/links"
	res, status, err := PerformRequest[*response.SuccessResponse](http.MethodPatch, path, headers, request)
	require.Nil(t, err)

	assert.Equal(t, http.StatusOK, status)
	assert.Equal(t, "Link Test Link was updated successfully", res.Message)
}

func TestUpdateLinkBadRequest(t *testing.T) {
	headers := map[string]string{
		"Accept": "application/json",
	}

	path := "/api/v1/links"
	res, status, err := PerformRequest[*response.FailedResponse](http.MethodPatch, path, headers, nil)
	require.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, status)
	assert.Equal(t, "invalid request body", res.Message)
}

func TestUpdateLinkFieldNotExist(t *testing.T) {
	request := request.UpdateLinkRequest{
		Name: "Test Link",
	}

	headers := map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/json",
	}

	path := "/api/v1/links"
	res, status, err := PerformRequest[*response.FailedResponse](http.MethodPatch, path, headers, request)
	require.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, status)
	assert.Equal(t, "no fields to update", res.Message)
}

func TestUpdateLinkFieldNotFound(t *testing.T) {
	request := request.UpdateLinkRequest{
		Name: "Invalid Link",
		URL:  "https://www.test.com",
	}

	headers := map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/json",
	}

	path := "/api/v1/links"
	res, status, err := PerformRequest[*response.FailedResponse](http.MethodPatch, path, headers, request)
	require.Nil(t, err)

	assert.Equal(t, http.StatusNotFound, status)
	assert.Equal(t, "link not found", res.Message)
}

func TestDelete(t *testing.T) {
	name := "Test Link"
	headers := map[string]string{
		"Accept": "application/json",
	}

	path := "/api/v1/links/" + strings.ReplaceAll(name, " ", "%20")
	res, status, err := PerformRequest[*response.SuccessResponse](http.MethodDelete, path, headers, nil)
	require.Nil(t, err)

	assert.Equal(t, http.StatusOK, status)
	assert.Equal(t, "Link Test Link was deleted successfully", res.Message)
}

func TestDeleteNotFound(t *testing.T) {
	name := "InvalidLink"
	headers := map[string]string{
		"Accept": "application/json",
	}

	path := "/api/v1/links/" + name
	res, status, err := PerformRequest[*response.FailedResponse](http.MethodDelete, path, headers, nil)
	require.Nil(t, err)

	assert.Equal(t, http.StatusNotFound, status)
	assert.Equal(t, "link not found", res.Message)
}