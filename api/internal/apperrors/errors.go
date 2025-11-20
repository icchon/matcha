package apperrors

import "errors"

var (
	ErrNotFound       = errors.New("resource not found")
	ErrInvalidInput   = errors.New("invalid input provided")
	ErrUnauthorized   = errors.New("unauthorized")
	ErrUnhandled      = errors.New("unhandled error")
	ErrInternalServer = errors.New("internal server error")
	ErrNotImplemented = errors.New("not implemented")
)
