package admin

import (
	"boilerplate-api/api/admin/auth"
	"boilerplate-api/api/admin/user"
	"go.uber.org/fx"
)

var Module = fx.Module("admin",
	fx.Options(
		auth.Module,
		user.Module,
		//gcp_billing.Module,
		//utility.Module,
	),
)
