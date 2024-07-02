package aws

import (
	"boilerplate-api/internal/config"
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
)

// NewAWSConfig creates new config instance from default aws profile in ~/.aws/credentials file
func NewAWSConfig(logger config.Logger, env config.Env) aws.Config {
	cfg, err := awsConfig.LoadDefaultConfig(
		context.TODO(),
		awsConfig.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				env.AwsAccessKey,
				env.AwsSecretKey,
				""),
		),
	)
	if err != nil {
		logger.Panic("Unable to load aws configuration from .env file")
	}
	logger.Info("✅ AWS config created.")
	return cfg
}
