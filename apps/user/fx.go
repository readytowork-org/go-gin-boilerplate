package user

import (
	"go.uber.org/fx"
)

// Module exported for initializing application
var Module = fx.Options(
	fx.Provide(ControllerConstuctor),
	fx.Provide(ServiceConstuctor),
	fx.Provide(RepositoryConstuctor),
	fx.Provide(NewUserValidator),
)
