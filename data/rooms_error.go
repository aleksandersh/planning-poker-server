package data

import "errors"

var (
	ErrUnauthorized    = errors.New("unathorized")
	ErrMissingResource = errors.New("missing resource")
	ErrLimitExceeded   = errors.New("resource limit exceeded")
	ErrIllegalState    = errors.New("resource illegal state")
)
