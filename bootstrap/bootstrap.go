package bootstrap

import (
	"context"

	"boilerplate-api/api"
	"boilerplate-api/cli"
	"boilerplate-api/database/seeds"
	"boilerplate-api/docs"
	"boilerplate-api/internal"
	"boilerplate-api/internal/config"
	"boilerplate-api/internal/router"
	"boilerplate-api/internal/utils"
	"boilerplate-api/services"
	"go.uber.org/fx"
)

// Module exported for initializing application
var Module = fx.Options(
	internal.Module,
	seeds.Module,
	cli.Module,
	services.Module,
	api.Module,
	fx.Supply(config.EnvPath(".env")),
	fx.Invoke(bootstrap),
)

func bootstrap(
	lifecycle fx.Lifecycle,
	handler router.Router,
	env config.Env,
	logger config.Logger,
	database config.Database,
	cliApp cli.Application,
	migrations config.Migrations,
) {

	appStop := func(context.Context) error {
		logger.Info("Stopping Application")
		conn, _ := database.DB.DB()
		_ = conn.Close()
		return nil
	}

	if utils.IsCli() {
		lifecycle.Append(fx.Hook{
			OnStart: func(context.Context) error {
				logger.Info("Starting cli Application")
				logger.Info("------- (CLI) ------")
				go cliApp.Start()
				return nil
			},
			OnStop: appStop,
		})

		return
	}

	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			logger.Info("Starting Application")
			logger.Info("------------------------")
			logger.Info("------ Gin Skeleton 📺 ------")
			logger.Info("------------------------")

			go func() {
				if env.Environment != "production" && env.HOST != "" {
					logger.Info("Setting Swagger Host...")
					docs.SwaggerInfo.Host = env.HOST
				}

				if env.Environment == "development" || env.Environment == "production" {
					logger.Info("Migrating DB schema...")
					migrations.MigrateUp()
				}
				if env.ServerPort == "" {
					_ = handler.Run()
				} else {
					_ = handler.Run(":" + env.ServerPort)
				}
			}()
			return nil
		},
		OnStop: appStop,
	})
}
