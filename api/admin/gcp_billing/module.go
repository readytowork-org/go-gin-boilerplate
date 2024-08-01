package gcp_billing

import (
	"go.uber.org/fx"
)

var Module = fx.Module("gcp_billing",
	fx.Options(
		fx.Provide(NewController),
		fx.Invoke(SetupRoutes),
	),
)
