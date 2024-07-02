package aws

import "go.uber.org/fx"

// Module aws module
var Module = fx.Module("aws", fx.Options(
	fx.Provide(NewAWSConfig),
	fx.Provide(NewS3BucketService),
))
