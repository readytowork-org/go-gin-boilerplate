package api_errors

import (
	"net/http"
)

type HttpErrorType int

// these are unsigned integer constants of custom type HttpErrorType
const (
	BadRequest      = HttpErrorType(http.StatusBadRequest)
	StatusOk        = HttpErrorType(http.StatusOK)
	Unauthorized    = HttpErrorType(http.StatusUnauthorized)
	Forbidden       = HttpErrorType(http.StatusForbidden)
	NotFound        = HttpErrorType(http.StatusNotFound)
	Conflict        = HttpErrorType(http.StatusConflict)
	InternalError   = HttpErrorType(http.StatusInternalServerError)
	Unavailable     = HttpErrorType(http.StatusServiceUnavailable)
	TooManyRequests = HttpErrorType(http.StatusTooManyRequests)
)

// ToInt converts HttpErrorType to int
func (h HttpErrorType) ToInt() int {
	return int(h)
}

// ValidationError struct
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ErrorResponse struct
type ErrorResponse struct {
	Message   string        `json:"message"`
	ErrorType HttpErrorType `json:"error_type"`
}
