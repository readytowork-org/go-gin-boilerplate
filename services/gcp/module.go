package gcp

import (
	"boilerplate-api/internal/config"

	"go.uber.org/fx"
	"google.golang.org/api/option"
)

// Module aws module
var Module = fx.Module("gcp", fx.Options(
	// BillingClient provider
	fx.Provide(func(logger config.Logger) BillingClient {
		opt := option.WithCredentialsFile("serviceAccountKey.json")
		return NewGCPBillingClient(
			BillingClientConfig{
				clientOption: &opt,
				logger:       logger.SugaredLogger,
			})
	}),
	// BucketClient provider
	fx.Provide(func(
		logger config.Logger) BucketClient {
		opt := option.WithCredentialsFile("serviceAccountKey.json")
		return NewGCPBucketClient(
			BucketClientConfig{
				logger:       logger.SugaredLogger,
				clientOption: &opt,
			})
	}),
	// BudgetClient provider
	fx.Provide(func(logger config.Logger) BudgetClient {
		opt := option.WithCredentialsFile("serviceAccountKey.json")
		return NewGCPBudgetClient(
			BudgetClientConfig{
				clientOption: &opt,
				logger:       logger.SugaredLogger,
			})
	}),
	// StorageBucketService provider
	fx.Provide(func(
		logger config.Logger,
		env config.Env,
		bucket BucketClient) StorageBucketService {
		return NewStorageBucketService(
			StorageBucketService{
				client: bucket.Client,
				logger: logger.SugaredLogger,
			},
		)
	}),
	// BillingService provider
	fx.Provide(func(
		logger config.Logger,
		env config.Env,
		billingClient BillingClient,
		budgetClient BudgetClient) BillingService {
		return NewGCPBillingService(
			BillingService{
				projectName:       env.ProjectName,
				billingAccountID:  env.BillingAccountId,
				budgetDisplayName: env.BudgetDisplayName,
				budgetAmount:      env.BudgetAmount,
				logger:            logger.SugaredLogger,
				gcpBilling:        billingClient,
				budgetClient:      budgetClient,
			},
		)
	}),
))
