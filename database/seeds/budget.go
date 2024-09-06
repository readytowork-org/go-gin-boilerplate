package seeds

import (
	"context"

	"boilerplate-api/internal/config"
	"boilerplate-api/internal/utils"
)

type IGcpBillingService interface {
	CreateOrUpdateBudget(ctx context.Context) (interface{}, error)
}

// ProjectBudgetSeed  Budget setup seed
type ProjectBudgetSeed struct {
	Seed
	logger        config.Logger
	budgetService IGcpBillingService
	env           config.Env
}

// NewProjectBudgetSeed creates budget if set on environment variable
func NewProjectBudgetSeed(
	logger config.Logger,
	budgetService IGcpBillingService,
	env config.Env,
) ProjectBudgetSeed {
	return ProjectBudgetSeed{
		logger:        logger,
		budgetService: budgetService,
		env:           env,
	}
}

// Run the seed data
func (c ProjectBudgetSeed) Run() {
	c.logger.Info("🌱 seeding  budget alert related setup...")

	if c.env.SetBudget == 1 {
		ctx := utils.GetContext()
		_, err := c.budgetService.CreateOrUpdateBudget(ctx)

		if err != nil {
			c.logger.Info("There is an error setting up budget alert ")
		} else {
			c.logger.Info("budget alert setup successfully")
		}
	}
}
