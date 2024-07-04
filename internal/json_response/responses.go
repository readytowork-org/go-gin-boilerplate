package json_response

type Message struct {
	Msg string `json:"message" validate:"required"`
} // @name Message

type Error[T any] struct {
	Message string `json:"message" validate:"required"`
	Error   T      `json:"error"`
} // @name ApiError

type Data[T any] struct {
	Data T `json:"data" validate:"required"`
} // @name Data

type DataCount[T any] struct {
	Data  []T   `json:"data" validate:"required"`
	Count int64 `json:"count" validate:"required"`
} // @name DataCount
