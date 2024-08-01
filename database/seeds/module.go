package seeds

import (
	"boilerplate-api/internal/config"
	"go.uber.org/fx"
)

// Seed db seed
type Seed interface {
	Run()
}

// Module exports seed module
var Module = fx.Module("seeds",
	fx.Options(
		fx.Provide(
			fx.Annotate(
				NewAdminSeed,
				fx.As(new(Seed)),
				fx.ResultTags(`group:"seeds"`),
			),
		),
		fx.Provide(
			fx.Annotate(
				NewProjectBudgetSeed,
				fx.As(new(Seed)),
				fx.ResultTags(`group:"seeds"`),
			),
		),
		fx.Invoke(
			fx.Annotate(
				SetupSeeds,
				fx.ParamTags(`group:"seeds"`),
			),
		),
	),
)

// SetupSeeds creates new seeds
func SetupSeeds(seeds []Seed, logger config.Logger) {
	logger.Info("🌱 seeding data...")
	for _, seed := range seeds {
		seed.Run()
	}
}
