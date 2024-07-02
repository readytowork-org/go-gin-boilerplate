package api

import (
	"boilerplate-api/api/admin"
	"boilerplate-api/api/docs"
	"boilerplate-api/api/user"
	"go.uber.org/fx"
)

var Module = fx.Module("api",
	fx.Options(
		docs.Module,
		admin.Module,
		user.Module,
	),
)
