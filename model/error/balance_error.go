package domain_error

import "errors"

var (
	// Stock Balance Errors
	ErrBalanceNotFound        = errors.New("balance not found")
	ErrBalanceCreationFailed  = errors.New("failed to create stock balance")
	ErrInvalidShareholderType = errors.New("invalid shareholder type")
)