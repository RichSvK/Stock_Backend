package domainerr

import "errors"

var (
	// General Errors
	ErrInternalServer     = errors.New("internal server error")
	ErrInvalidRequestBody = errors.New("invalid request body")

	// File err
	ErrSaveFile       = errors.New("failed to save file")
	ErrUploadFailed   = errors.New("file upload failed")
	ErrFailedWriteCSV = errors.New("failed to write CSV file")

	// Date Errors
	ErrInvalidDateRange = errors.New("invalid date range")
)
