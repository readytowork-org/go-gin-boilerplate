package utils

import (
	"boilerplate-api/internal/api_errors"
	"strconv"
)

func StringToInt64(stringData string) (int64, error) {
	var intID int64
	var err error
	if stringData != "" {
		intID, err = strconv.ParseInt(stringData, 10, 64)
		if err != nil {
			err = api_errors.BadRequest.Wrap(err, "Invalid ID")
		}
	}

	if stringData == "" {
		err = api_errors.BadRequest.Wrap(err, "Failed to convert ID into int64")
	}
	return intID, err
}
