package docs

import (
	"go.uber.org/fx"
)

var Module = fx.Module("docs",
	fx.Options(
		fx.Invoke(SetupRoutes),
	),
)
