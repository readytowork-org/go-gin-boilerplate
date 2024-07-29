package gcp

import (
	"cloud.google.com/go/billing/budgets/apiv1"
	"context"
	"google.golang.org/api/option"
)

type budgetClientLogger interface {
	Panic(args ...interface{})
}

type BudgetClientConfig struct {
	logger       budgetClientLogger
	clientOption *option.ClientOption
}
type BudgetClient struct {
	*budgets.BudgetClient
}

func NewGCPBudgetClient(clientConfig BudgetClientConfig) BudgetClient {
	budgetClient, err := budgets.NewBudgetClient(context.Background(), *clientConfig.clientOption)

	if err != nil {
		clientConfig.logger.Panic("Failed to create cloud budget api client: %v \n", err)
	}
	return BudgetClient{
		budgetClient,
	}
}
