package middlewares

import (
	"go.uber.org/fx"
)

// Module Middleware exported
var Module = fx.Options(
	fx.Provide(NewDBTransactionMiddleware),
	fx.Provide(NewRateLimitMiddleware),
	fx.Provide(NewJWTAuthMiddleWare),
)
