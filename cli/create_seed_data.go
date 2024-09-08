package cli

import (
	"boilerplate-api/database/seeds/faker"
	"boilerplate-api/internal/config"
)

// CreateSeedData command
type CreateSeedData struct {
	logger   config.Logger
	database *config.Database
	//seeds  seeds.Seeds
}

// NewCreateSeedData creates instance of admin user
func NewCreateSeedData(
	logger config.Logger,
	database *config.Database,
//seeds seeds.Seeds,
) CreateSeedData {
	return CreateSeedData{
		logger:   logger,
		database: database,
		//seeds:  seeds,
	}
}

// Run runs command
func (c CreateSeedData) Run() {
	c.logger.Info("🌱 Creating seed data...")
	//c.seeds.Run()
	_ = faker.NewFaker(c.database.DB, c.logger, faker.Config{})
}

// Name return name of command
func (c CreateSeedData) Name() string {
	return "CREATE_SEED_DATA"
}
