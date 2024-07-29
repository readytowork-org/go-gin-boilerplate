package aws

import (
	"context"
	"go.uber.org/zap"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
)

type AuthConfig struct {
	logger    *zap.SugaredLogger
	accessKey string
	secretKey string
}

// NewAWSConfig creates new config instance from default aws profile in ~/.aws/credentials file
func NewAWSConfig(config AuthConfig) aws.Config {
	cfg, err := awsConfig.LoadDefaultConfig(
		context.TODO(),
		awsConfig.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				config.accessKey,
				config.secretKey,
				""),
		),
	)
	if err != nil {
		config.logger.Panic("Unable to load aws configuration from provided resources")
	}
	config.logger.Info("✅ AWS config created.")
	return cfg
}
