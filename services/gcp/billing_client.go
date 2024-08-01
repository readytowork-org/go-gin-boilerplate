package gcp

import (
	"context"
	"google.golang.org/api/cloudbilling/v1"
	"google.golang.org/api/option"
)

type billingClientLogger interface {
	Panic(args ...interface{})
}

type BillingClientConfig struct {
	logger       billingClientLogger
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
