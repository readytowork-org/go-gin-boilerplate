package lib

import (
	"boilerplate-api/lib/auth"
	"boilerplate-api/lib/config"
	"boilerplate-api/lib/middlewares"
	"boilerplate-api/lib/request_validator"
	"boilerplate-api/lib/router"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"internal",
	config.Module,
	middlewares.Module,
	fx.Options(
		fx.Provide(
			router.NewRouter,
			request_validator.NewValidator,
			auth.NewJWTAuthService,
		),
	),
)
