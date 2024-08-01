package main

import (
	"boilerplate-api/bootstrap"
	"github.com/joho/godotenv"
	"go.uber.org/fx"
)

//	@title						Demo API
//	@version					1.0
//	@description				An API in Go using Gin framework
//	@host						localhost:8000
//	@securityDefinitions.apikey	Bearer
//	@in							header
//	@name						Authorization
//	@description				Description for what is this security definition being used
func main() {
	_ = godotenv.Load()
	fx.New(bootstrap.Module).Run()
}
