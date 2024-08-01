package cli

import (
	"boilerplate-api/internal/config"
)

// CreateSeedData command
type CreateSeedData struct {
	logger config.Logger
	//seeds  seeds.Seeds
}

// NewCreateSeedData creates instance of admin user
func NewCreateSeedData(
	logger config.Logger,
//seeds seeds.Seeds,
) CreateSeedData {
	return CreateSeedData{
		logger: logger,
		//seeds:  seeds,
	}
}

// Run runs command
func (c CreateSeedData) Run() {
	c.logger.Info("🌱 Creating seed data...")
	//c.seeds.Run()
}

// Name return name of command
func (c CreateSeedData) Name() string {
	return "CREATE_SEED_DATA"
}
