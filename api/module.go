package api

import (
	"boilerplate-api/api/admin"
	"boilerplate-api/api/auth"
	"boilerplate-api/api/swagger"
	"boilerplate-api/api/user"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"api",
	fx.Options(
		swagger.Module,
		admin.Module,
		user.Module,
		auth.Module,
	),
)
