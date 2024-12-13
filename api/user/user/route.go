package user

import (
	"boilerplate-api/lib/config"
	"boilerplate-api/lib/middlewares"
	"boilerplate-api/lib/router"
)

// SetupRoutes user routes
func SetupRoutes(
	logger config.Logger,
	router router.Router,
	userController Controller,
	jwtMiddleware middlewares.JWTAuthMiddleWare,
) {
	logger.Info(" Setting up user routes")
	router.V1.GET("/profile", jwtMiddleware.Handle(), userController.GetUserProfile)
}
