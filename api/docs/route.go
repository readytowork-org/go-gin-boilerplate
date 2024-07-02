package docs

import (
	"boilerplate-api/internal/config"
	"boilerplate-api/internal/router"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRoutes(
	logger config.Logger,
	router router.Router,
	env config.Env,
) {
	if env.Environment != "production" {
		logger.Info(" Setting up Docs routes")
		swagger := router.Group("/docs")
		{
			swagger.GET("/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		}
	}
}
