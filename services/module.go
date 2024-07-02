package services

import (
	"boilerplate-api/services/aws"
	"boilerplate-api/services/firebase"
	"boilerplate-api/services/gcp"
	"go.uber.org/fx"
)

var Module = fx.Options(
	firebase.Module,
	aws.Module,
	gcp.Module,
	fx.Provide(NewStripeService),
	fx.Provide(NewGmailService),
	fx.Provide(NewTwilioService),
)
