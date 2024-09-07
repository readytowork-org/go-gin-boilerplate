package admin

import (
	"boilerplate-api/api/admin/user"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"admin",
	fx.Options(
		user.Module,
		//gcp_billing.Module,
		//utility.Module,
	),
)
