package middlewares

import (
	"net/http"

	"boilerplate-api/internal/auth"
	"boilerplate-api/internal/config"
	"boilerplate-api/internal/constants"
	"boilerplate-api/internal/json_response"

	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
)

type JWTAuthMiddleWare struct {
	jwtService auth.JWTAuthService
	logger     config.Logger
	env        config.Env
}

func NewJWTAuthMiddleWare(
	jwtService auth.JWTAuthService,
	logger config.Logger,
	env config.Env,
) JWTAuthMiddleWare {
	return JWTAuthMiddleWare{
		jwtService: jwtService,
		logger:     logger,
		env:        env,
	}
}

// Handle user with jwt using this middleware
func (m JWTAuthMiddleWare) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the token from the request header
		header := c.GetHeader(constants.Headers.Authorization.ToString())

		tokenString, err := m.jwtService.GetTokenFromHeader(header)
		if err != nil {
			m.logger.Error("Error getting token from header: ", err.Message)
			c.JSON(
				http.StatusUnauthorized, json_response.Error[string]{
					Error:   err.Message,
					Message: "Error getting token from header",
				},
			)
			c.Abort()
			return
		}

		// Parsing and Verifying token
		parsedToken, parseErr := m.jwtService.ParseAndVerifyToken(tokenString, m.env.JwtAccessSecret)
		if parseErr != nil {
			m.logger.Error("Error parsing token: ", parseErr.Message)
			c.JSON(
				http.StatusUnauthorized, json_response.Error[string]{
					Error:   parseErr.Message,
					Message: "Failed to parse and verify token",
				},
			)
			c.Abort()
			return
		}
		// Retrieve claims
		claims, claimsError := m.jwtService.RetrieveClaims(parsedToken)
		if claimsError != nil {
			m.logger.Error("Error retrieving claims: ", claimsError.Message)
			c.JSON(
				http.StatusUnauthorized, json_response.Error[string]{
					Error:   claimsError.Message,
					Message: "Failed to retrieve claims from token",
				},
			)
			c.Abort()
			return
		}
		// ser user to the scope
		sentry.ConfigureScope(
			func(scope *sentry.Scope) {
				scope.SetUser(sentry.User{ID: claims.ID})
			},
		)
		// Can set anything in the request context and passes the request to the next handler.
		c.Set(constants.UserID, claims.ID)
		c.Next()
	}
}
