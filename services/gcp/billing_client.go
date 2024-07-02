package gcp

import (
	"context"

	"boilerplate-api/internal/config"
	"google.golang.org/api/option"

	"google.golang.org/api/cloudbilling/v1"
)

type BillingClient struct {
	*cloudbilling.APIService
}

// NewGCPBillingClient creates a new gcp billing api client
func NewGCPBillingClient(logger config.Logger, clientOption *option.ClientOption) BillingClient {
	billingClient, err := cloudbilling.NewService(context.Background(), *clientOption)
	if err != nil {
		logger.Panic("Failed to create cloud billing api client: %v \n", err)
	}

	return BillingClient{
		billingClient,
	}
}
