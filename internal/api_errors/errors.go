package api_errors

import (
	"fmt"
	"net/http"
	"os"

	"github.com/pkg/errors"
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

// ErrorContext struct
type ErrorContext struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ErrResponse[T any] struct {
	Message string `json:"message" validate:"required"`
	Error   string `json:"error" validate:"required"`
	Errors  T      `json:"errors"`
} // @name ErrResponse

type responseError struct {
	httpErrorType HttpErrorType
	originalError error
	customMessage CustomMessage
	context       []ErrorContext
}

// CustomMessage struct
type CustomMessage struct {
	Message string `json:"message"`
}

// Error() returns the origin message
func (error responseError) Error() string {
	return error.originalError.Error()
}

// ----------------- HTTPErrorType ----------------------

// New func :
func (errorType HttpErrorType) New(msg string) error {
	return responseError{
		httpErrorType: errorType,
		originalError: errors.New(msg),
	}
}

// Newf func :
func (errorType HttpErrorType) Newf(msg string, args ...interface{}) error {
	return responseError{
		httpErrorType: errorType,
		originalError: fmt.Errorf(msg, args...),
	}
}

// Wrap func:
func (errorType HttpErrorType) Wrap(err error, msg string) error {
	return errorType.Wrapf(err, msg)
}

// Wrapf func:
func (errorType HttpErrorType) Wrapf(err error, msg string, args ...interface{}) error {
	return responseError{
		httpErrorType: errorType,
		originalError: errors.Wrapf(err, msg, args...),
		customMessage: CustomMessage{
			Message: msg,
		},
	}
}

// GetErrorType func
func GetErrorType(err error) HttpErrorType {
	if responseErr, ok := err.(responseError); ok {
		return responseErr.httpErrorType
	}
	return InternalError
}

// ---------- Ends  HttpErrorType ------------

// AddErrorContext func:
func AddErrorContext(err error, field, message string) error {
	context := ErrorContext{Field: field, Message: message}
	if responseErr, ok := err.(responseError); ok {
		return responseError{
			httpErrorType: responseErr.httpErrorType,
			originalError: responseErr.originalError,
			customMessage: responseErr.customMessage,
			context:       append(responseErr.context, context),
		}
	}

	re := &responseError{}
	re.httpErrorType = InternalError
	re.context = append(re.context, context)
	re.originalError = err
	return re
}

// AddErrorContextBlock func:
func AddErrorContextBlock(err error, errorCtx []ErrorContext) error {
	if responseErr, ok := err.(responseError); ok {
		return responseError{
			httpErrorType: responseErr.httpErrorType,
			originalError: responseErr.originalError,
			customMessage: responseErr.customMessage,
			context:       append(responseErr.context, errorCtx...),
		}
	}

	re := &responseError{}
	re.httpErrorType = InternalError
	re.context = append(re.context, errorCtx...)
	re.originalError = err
	return re
}

// GetErrorContext func
func GetErrorContext(err error) []ErrorContext {
	if responseErr, ok := err.(responseError); ok && responseErr.context != nil {
		return responseErr.context
	}
	return nil
}

func SetCustomMessage(err error, msg string) error {
	customMessage := CustomMessage{Message: msg}
	if responseErr, ok := err.(responseError); ok {
		return responseError{
			httpErrorType: responseErr.httpErrorType,
			originalError: responseErr.originalError,
			customMessage: customMessage,
			context:       responseErr.context,
		}
	}
	return responseError{
		httpErrorType: InternalError,
		originalError: err,
		customMessage: customMessage,
	}
}

func GetCustomMessage(err error) string {
	emptyStruct := CustomMessage{}
	if responseErr, ok := err.(responseError); ok && responseErr.customMessage != emptyStruct {
		return responseErr.customMessage.Message
	}
	return ""
}

// HandleError func
func HandleError(err error) (status int, response ErrResponse[[]ErrorContext]) {
	_status := GetErrorType(err)

	errorContext := GetErrorContext(err)
	customMessage := GetCustomMessage(err)
	if os.Getenv("ENV") != "production" {
		response.Error = err.Error()
	}
	response.Message = customMessage

	if errorContext != nil {
		response.Errors = errorContext
	}
	return int(_status), response
}
