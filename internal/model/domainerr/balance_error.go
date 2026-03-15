package domainerr

import "errors"

var (
	// Stock Balance Errors
	ErrBalanceNotFound        = errors.New("balance not found")
	ErrBalanceCreationFailed  = errors.New("failed to create stock balance")
	ErrInvalidShareholderType = errors.New("invalid shareholder type")
	ErrDuplicateBalanceData  = errors.New("duplicate balance data exists")
)