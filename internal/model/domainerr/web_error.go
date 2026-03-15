package domainerr

import "errors"

var (
	// General Errors
	ErrInternalServer       = errors.New("internal server error")
	ErrInvalidChangeRequest = errors.New("invalid change request")

	// File err
	ErrSaveFile       = errors.New("failed to save file")
	ErrUploadFailed   = errors.New("file upload failed")
	ErrFailedWriteCSV = errors.New("failed to write CSV file")

	// Date Errors
	ErrInvalidDateRange = errors.New("invalid date range")
)
