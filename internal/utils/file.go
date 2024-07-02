package utils

import (
	"errors"
	"path/filepath"
)

func GetFileName(filename string) (string, error) {
	if filename == "" {
		return "", errors.New("filename cannot be empty")
	}

	fileExtension := filepath.Ext(filename)
	return GenerateRandomFileName() + fileExtension, nil
}
