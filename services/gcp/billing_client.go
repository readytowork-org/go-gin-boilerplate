package gcp

import (
	"context"
	"go.uber.org/zap"
	"google.golang.org/api/cloudbilling/v1"
	"google.golang.org/api/option"
)

type BillingClientConfig struct {
	logger       *zap.SugaredLogger
	clientOption *option.ClientOption
}

type BillingClient struct {
	*cloudbilling.APIService
}

// NewGCPBillingClient creates a new gcp billing api client
func NewGCPBillingClient(config BillingClientConfig) BillingClient {
	billingClient, err := cloudbilling.NewService(context.Background(), *config.clientOption)
	if err != nil {
		config.logger.Panic("Failed to create cloud billing api client: %v \n", err)
	}

	return BillingClient{
		billingClient,
	}
}
