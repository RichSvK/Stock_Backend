package test

import (
	"backend/internal/model/response"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
)

func ClearTable(tableName string) error {
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	query := fmt.Sprintf(`DELETE FROM %s`, tableName)
	if err := tx.Exec(query).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete from table: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to clear database: %w", err)
	}

	return nil
}

func ExecuteSQLFile(filePath string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read SQL file: %w", err)
	}

	statements := strings.Split(string(content), ";")

	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}

		if err := db.Exec(stmt).Error; err != nil {
			return fmt.Errorf("failed to execute SQL file: %w", err)
		}
	}

	return nil
}

func PerformRequest[T any](method string, path string, httpHeader map[string]string, body any) (T, int, error) {
	var response T

	var bodyReader io.Reader
	switch v := body.(type) {
	case io.Reader:
		bodyReader = v
	case nil:
		bodyReader = nil
	default:
		jsonRequest, err := json.Marshal(v)
		if err != nil {
			return response, 0, err
		}
		bodyReader = strings.NewReader(string(jsonRequest))
	}

	req := httptest.NewRequest(method, path, bodyReader)
	for key, value := range httpHeader {
		req.Header.Set(key, value)
	}

	recorder := httptest.NewRecorder()
	app.ServeHTTP(recorder, req)

	if recorder.Body.Len() > 0 {
		contentType := recorder.Header().Get("Content-Type")
		if strings.Contains(contentType, "application/json") {
			err := json.Unmarshal(recorder.Body.Bytes(), &response)
			if err != nil {
				return response, 0, err
			}
		}
	}

	return response, recorder.Code, nil
}

func InsertTestStockData(filepath string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("failed to open test data file: %w", err)
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", "2026_02_February.txt")
	if err != nil {
		return fmt.Errorf("failed to create form file: %w", err)
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return fmt.Errorf("failed to copy file content: %w", err)
	}

	err = file.Close()
	if err != nil {
		return fmt.Errorf("failed to close: %w", err)
	}

	err = writer.Close()
	if err != nil {
		return fmt.Errorf("failed to close writer: %w", err)
	}

	headers := map[string]string{
		"Content-Type": writer.FormDataContentType(),
		"Accept":       "application/json",
	}

	path := "/api/v1/balances/import"
	_, _, err = PerformRequest[*response.UploadBalanceResponse](http.MethodPost, path, headers, body)
	if err != nil {
		return err
	}

	return nil
}
