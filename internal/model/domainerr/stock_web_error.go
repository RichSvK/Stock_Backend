package domainerr

import "errors"

var (
	// Web List Errors
	ErrLinkNotFound      = errors.New("link not found")
	ErrLinkAlreadyExists = errors.New("link already exists")
	ErrNoFieldsToUpdate  = errors.New("no fields to update")
)
