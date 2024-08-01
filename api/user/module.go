package user

import (
	"boilerplate-api/api/user/user"
	"go.uber.org/fx"
)

var Module = fx.Module("user",
	fx.Options(
		user.Module,
	),
)
