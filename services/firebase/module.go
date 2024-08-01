package firebase

import (
	"boilerplate-api/internal/config"
	"go.uber.org/fx"
	"google.golang.org/api/option"
)

// Module firebase module
var Module = fx.Module("firebase", fx.Options(
	fx.Provide(func(
		logger config.Logger) AppService {
		opt := option.WithCredentialsFile("./path/to/your/serviceAccountKey.json")
		appConfig := AppConfig{
			logger: logger.SugaredLogger,
			opt:    &opt,
		}
		return NewFirebaseApp(appConfig)
	}),
	fx.Provide(func(
		logger config.Logger,
		appService AppService) AuthService {
		return NewFirebaseAuthService(
			AuthConfig{
				logger: logger.SugaredLogger,
				app:    appService.App,
			})
	}),

	fx.Provide(func(
		logger config.Logger,
		appService AppService) StoreClientService {
		return NewFireStoreClient(
			StoreClientConfig{
				logger: logger.SugaredLogger,
				app:    appService.App,
			})
	}),
	fx.Provide(func(
		logger config.Logger,
		appService AppService) CMClientService {
		return NewFirebaseCMClient(
			CMClientConfig{
				logger: logger.SugaredLogger,
				app:    appService.App,
			})
	}),
	fx.Provide(NewFirebaseAuthMiddleware),
))
