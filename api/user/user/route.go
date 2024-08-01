package user

import (
	"boilerplate-api/internal/config"
	"boilerplate-api/internal/middlewares"
	"boilerplate-api/internal/router"
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
