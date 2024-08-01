package seeds

import (
	"context"
	"time"

	"boilerplate-api/internal/config"
	"boilerplate-api/services/gcp"
)

// ProjectBudgetSeed  Budget setup seed
type ProjectBudgetSeed struct {
	Seed
	logger        config.Logger
	budgetService gcp.BillingService
	env           config.Env
}

// NewProjectBudgetSeed creates budget if set on environment variable
func NewProjectBudgetSeed(
	logger config.Logger,
	budgetService gcp.BillingService,
	env config.Env,
) ProjectBudgetSeed {
	return ProjectBudgetSeed{
		logger:        logger,
		budgetService: budgetService,
		env:           env,
	}
}

func (c ProjectBudgetSeed) getContext() context.Context {

	// Create a context.Background()
	ctx := context.Background()

	// Create a context.WithCancel() to create a cancellable context
	defer context.WithCancel(ctx)

	// Create a context.WithTimeout() to create a context with a timeout
	timeout := 5 * time.Second
	defer context.WithTimeout(ctx, timeout)

	// Create a context.WithDeadline() to create a context with a deadline
	deadline := time.Now().Add(10 * time.Second)
	defer context.WithDeadline(ctx, deadline)

	return ctx
}

// Run the seed data
func (c ProjectBudgetSeed) Run() {
	c.logger.Info("🌱 seeding  budget alert related setup...")

	if c.env.SetBudget == 1 {
		ctx := c.getContext()
		_, err := c.budgetService.CreateOrUpdateBudget(ctx)

		if err != nil {
			c.logger.Info("There is an error setting up budget alert ")
		} else {
			c.logger.Info("budget alert setup successfully")

		}

	}

}
