package aws

import (
	"boilerplate-api/lib/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	"go.uber.org/fx"
)

// Module aws module
var Module = fx.Module(
	"aws", fx.Options(
		fx.Provide(
			func(
				logger config.Logger,
				env config.Env,
			) aws.Config {
				return NewAWSConfig(
					AuthConfig{
						logger:    logger.SugaredLogger,
						accessKey: env.AwsAccessKey,
						secretKey: env.AwsSecretKey,
					},
				)
			},
		),
		fx.Provide(
			func(
				logger config.Logger,
				env config.Env,
			) S3BucketService {
				return NewS3BucketService(
					S3BucketConfig{
						s3bucket: env.AwsS3Bucket,
						s3Region: env.AwsS3Region,
					},
				)
			},
		),
	),
)
