package config

import (
	"path/filepath"

	"google.golang.org/api/option"
)

func NewGCPClientOption(logger Logger) *option.ClientOption {
	serviceAccountKeyFilePath, err := filepath.Abs("./serviceAccountKey.json")
	if err != nil {
		logger.Panic("Unable to load serviceAccountKey.json file")
	}

	options := option.WithCredentialsFile(serviceAccountKeyFilePath)
	return &options
}
