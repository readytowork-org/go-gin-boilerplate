package gcp

import "go.uber.org/fx"

// Module aws module
var Module = fx.Module("gcp", fx.Options(
	fx.Provide(NewGCPBillingClient),
	//fx.Provide(NewGCPBucketClient),
	fx.Provide(NewGCPBudgetClient),
	//fx.Provide(NewStorageBucketService),
	fx.Provide(NewGCPBillingService),
))
