package gcp

import (
	"context"

	"boilerplate-api/internal/config"
	"cloud.google.com/go/billing/budgets/apiv1"
	"google.golang.org/api/option"
)

type BudgetClient struct {
	*budgets.BudgetClient
}

func NewGCPBudgetClient(logger config.Logger, clientOption *option.ClientOption) BudgetClient {
	budgetClient, err := budgets.NewBudgetClient(context.Background(), *clientOption)

	if err != nil {
		logger.Panic("Failed to create cloud budget api client: %v \n", err)
	}
	return BudgetClient{
		budgetClient,
	}
}
