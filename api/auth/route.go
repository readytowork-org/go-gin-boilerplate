package auth

import (
	"boilerplate-api/lib/config"
	"boilerplate-api/lib/constants"
	"boilerplate-api/lib/middlewares"
	"boilerplate-api/lib/router"
)

// JwtAuthRoutes struct
type JwtAuthRoutes struct {
	logger              config.Logger
	router              router.Router
	jwtController       JwtAuthController
	rateLimitMiddleware middlewares.RateLimitMiddleware
}

// NewJwtAuthRoutes creates new jwt controller
func NewJwtAuthRoutes(
	logger config.Logger,
	router router.Router,
	jwtController JwtAuthController,
	rateLimitMiddleware middlewares.RateLimitMiddleware,
) JwtAuthRoutes {
	return JwtAuthRoutes{
		router:              router,
		logger:              logger,
		jwtController:       jwtController,
		rateLimitMiddleware: rateLimitMiddleware,
	}
}

// SetupRoutes Obtain Jwt Token Routes
func SetupRoutes(
	logger config.Logger,
	router router.Router,
	jwtController JwtAuthController,
	rateLimitMiddleware middlewares.RateLimitMiddleware,
) {
	logger.Info(" Setting up jwt routes")
	jwt := router.V1.Group("/login").Use(
		rateLimitMiddleware.HandleRateLimit(
			constants.LoginRateLimit, constants.LoginPeriod,
		),
	)
	{
		jwt.POST("", jwtController.LoginUserWithJWT)
		jwt.POST("/refresh", jwtController.RefreshJwtToken)
	}
}
