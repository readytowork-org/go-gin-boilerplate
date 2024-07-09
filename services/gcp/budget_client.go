package gcp

import (
	"cloud.google.com/go/billing/budgets/apiv1"
	"context"
	"go.uber.org/zap"
	"google.golang.org/api/option"
)

type BudgetClientConfig struct {
	logger       *zap.SugaredLogger
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
