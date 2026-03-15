package domainerr

import "errors"

var (
	// IPO Errors
	ErrIpoDataNotFound = errors.New("ipo data not found")
	ErrEmptyRequest    = errors.New("request body cannot be empty")
)
