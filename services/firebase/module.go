package firebase

import "go.uber.org/fx"

// Module firebase module
var Module = fx.Module("firebase", fx.Options(
	fx.Provide(NewFirebaseApp),
	fx.Provide(NewFirestoreClient),
	fx.Provide(NewFirebaseCMClient),
	fx.Provide(NewFirebaseAuthService),
	fx.Provide(NewFirebaseAuthMiddleware),
))
