package utility

import (
	"go.uber.org/fx"
)

var Module = fx.Module(
	"utility",
	fx.Options(
		fx.Provide(NewService),
		fx.Provide(NewController),
		fx.Invoke(SetupRoutes),
	),
)
