package middlewares

import (
	"boilerplate-api/internal/api_errors"
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
	db         config.Database
}

func NewJWTAuthMiddleWare(
	jwtService auth.JWTAuthService,
	logger config.Logger,
	env config.Env,
	db config.Database,
) JWTAuthMiddleWare {
	return JWTAuthMiddleWare{
		jwtService: jwtService,
		logger:     logger,
		env:        env,
		db:         db,
	}
}

// Handle user with jwt using this middleware
func (m JWTAuthMiddleWare) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the token from the request header
		header := c.GetHeader(constants.Headers.Authorization.ToString())

		tokenString, err := m.jwtService.GetTokenFromHeader(header)
		if err != nil {
			m.logger.Error("Error getting token from header: ", err.Error())
			err = api_errors.Unauthorized.Wrap(err, "Error getting token from header")
			status, errM := api_errors.HandleError(err)
			c.JSON(status, json_response.Error{Error: errM})
			c.Abort()
			return
		}

		// Parsing and Verifying token
		parsedToken, parseErr := m.jwtService.ParseAndVerifyToken(tokenString, m.env.JwtAccessSecret)
		if parseErr != nil {
			m.logger.Error("Error parsing token: ", parseErr.Error())
			err = api_errors.Unauthorized.Wrap(parseErr, "Failed to parse and verify token")
			status, errM := api_errors.HandleError(err)
			c.JSON(status, json_response.Error{Error: errM})
			c.Abort()
			return
		}
		// Retrieve claims
		claims, claimsError := m.jwtService.RetrieveClaims(parsedToken)
		if claimsError != nil {
			m.logger.Error("Error retrieving claims: ", claimsError.Error())
			err = api_errors.Unauthorized.Wrap(claimsError, "Failed to retrieve claims from token")
			status, errM := api_errors.HandleError(err)
			c.JSON(status, json_response.Error{Error: errM})
			c.Abort()
			return
		}
		// ser user to the scope
		sentry.ConfigureScope(func(scope *sentry.Scope) {
			scope.SetUser(sentry.User{ID: claims.ID})
		})
		// Can set anything in the request context and passes the request to the next handler.
		c.Set(constants.UserID, claims.ID)
		c.Next()
	}
}
