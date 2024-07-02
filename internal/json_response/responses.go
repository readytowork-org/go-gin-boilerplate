package json_response

import (
	"boilerplate-api/internal/api_errors"
)

type Message struct {
	Msg string `json:"message" validate:"required"`
} // @name Message

type Error struct {
	Error api_errors.ErrResponse[[]api_errors.ErrorContext] `json:"error" validate:"required"`
} // @name ApiError

type ErrorMsg struct {
	Error string `json:"error" validate:"required"`
} // @name ErrorMessage

type Data[T any] struct {
	Data T `json:"data" validate:"required"`
} // @name Data

type DataCount[T any] struct {
	Data  []T   `json:"data" validate:"required"`
	Count int64 `json:"count" validate:"required"`
} // @name DataCount
