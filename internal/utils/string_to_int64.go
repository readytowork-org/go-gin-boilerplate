package utils

import (
	"boilerplate-api/internal/api_errors"
	"strconv"
)

func StringToInt64(stringData string) (int64, *api_errors.ErrorResponse) {
	var intID int64

	if stringData == "" {
		return intID, &api_errors.ErrorResponse{
			ErrorType: api_errors.BadRequest,
			Message:   "Failed to convert ID into int64",
		}
	}

	intID, err := strconv.ParseInt(stringData, 10, 64)
	if err != nil {
		return intID, &api_errors.ErrorResponse{
			ErrorType: api_errors.BadRequest,
			Message:   "Invalid ID",
		}
	}

	return intID, nil
}
