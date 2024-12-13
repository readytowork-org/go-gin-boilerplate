package services

import (
	"boilerplate-api/lib/config"
	"boilerplate-api/services/aws"
	"go.uber.org/fx"
)

var Module = fx.Options(
	aws.Module,
	// StripeService provider
	fx.Provide(
		func(
			env config.Env,
			logger config.Logger,
		) StripeService {
			return NewStripeService(
				StripeConfig{
					stripeSecretKey: env.StripeSecretKey,
					stripeProductID: env.StripeProductID,
					logger:          logger.SugaredLogger,
				},
			)
		},
	),
	// GmailService provider
	fx.Provide(
		func(
			env config.Env,
			logger config.Logger,
		) GmailService {
			return NewGmailService(
				GmailConfig{
					clientID:     env.MailClientID,
					clientSecret: env.MailClientSecret,
					accessToken:  env.MailAccesstoken,
					refreshToken: env.MailRefreshToken,
					hostURL:      env.HOST,
					logger:       logger.SugaredLogger,
				},
			)
		},
	),
	// TwilioService provider
	fx.Provide(
		func(
			env config.Env,
			logger config.Logger,
		) TwilioService {
			return NewTwilioService(
				TwilioService{
					baseURL:   env.TwilioBaseURL,
					smsFrom:   env.TwilioSMSFrom,
					sID:       env.TwilioSID,
					authToken: env.TwilioAuthToken,
					logger:    logger.SugaredLogger,
				},
			)
		},
	),
)
