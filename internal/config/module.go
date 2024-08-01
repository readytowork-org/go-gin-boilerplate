package config

import (
	"go.uber.org/fx"
	"log"
	"strings"
)

// ENVModule config dependency
var ENVModule = fx.Options(
	fx.Provide(NewEnv),
	fx.Provide(NewDSNConfig),
)

// BaseModule base config
var BaseModule = fx.Options(
	fx.Provide(NewDatabase),
	fx.Provide(NewMigrations),
	fx.Provide(GetLogger),
)

// Module main config
var Module = fx.Options(
	ENVModule,
	BaseModule,
	fx.Provide(NewGCPClientOption),
)

// TestENVModule required for test env
var TestENVModule = fx.Module("ENV", fx.Options(
	ENVModule,
	fx.Invoke(validateTestEnv),
))

func validateTestEnv(envPath EnvPath, env Env) {
	if !strings.Contains(envPath.ToString(), "test") {
		log.Fatalf("Error :: Trying to test with incorrect .env. Make sure to use .test.env")
		return
	}

	if !(strings.Contains(env.Environment, "test")) {
		log.Fatalf("Error :: Testing in incorrect environment. Make sure env is set to test")
		return
	}

	if !strings.Contains(env.DBName, "test") {
		log.Fatalf("Error :: Testing with incorrect db. Make sure to use test database")
		return
	}
}
