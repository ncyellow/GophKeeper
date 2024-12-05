package api

import "errors"

// List of errors that the client can generate when interacting with the server
// While some errors in the http client come as StatusCode with no error message.
// Then we generate an error, and if a grpc error occurs, we propagate it using the error format
var (
	FmtErrInternalServer    = "server unavailable, please try later: %w"
	FmtErrServerTimout      = "server unavailable, please try later: %w"
	FmtErrDeserialization   = "deserialization error: %w"
	FmtErrRequestPrepare    = "failed to prepare http request: %w"
	FmtErrUserAlreadyExists = "a user with this login is already registered: %w"
	FmtErrUserNotFound      = "a user with this login was not found: %w"
	FmtErrAlreadyExists     = "ID with this identifier is already registered: %w"
	FmtErrNotFound          = "record with this identifier was not found: %w"
	FmtErrSerialization     = "serialization error: %w"

	ErrSerialization     = errors.New("serialization error")
	ErrAuthRequire       = errors.New("authorization required")
	ErrUserAlreadyExists = errors.New("a user with this login is already registered")
	ErrInternalServer    = errors.New("server unavailable, please try later")
	ErrUserNotFound      = errors.New("a user with this login was not found")
	ErrAlreadyExists     = errors.New("ID with this identifier is already registered")
	ErrNotFound          = errors.New("record with this identifier was not found")
)
