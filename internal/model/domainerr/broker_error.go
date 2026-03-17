package domainerr

import "errors"

var (
	// Broker Errors
	ErrBrokerNotFound    = errors.New("broker not found")
	ErrInvalidBrokerCode = errors.New("invalid broker code")
	ErrBrokerDuplicate = errors.New("duplicate broker data")
)
