package swagger

import (
	"go.uber.org/fx"
)

var Module = fx.Module(
	"swagger",
	fx.Options(
		fx.Invoke(SetupRoutes),
	),
)
