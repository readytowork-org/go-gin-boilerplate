package middlewares

import (
	"boilerplate-api/internal/api_errors"
	"boilerplate-api/internal/constants"
	"boilerplate-api/internal/json_response"
	"boilerplate-api/services/firebase"
	"firebase.google.com/go/auth"
	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type SetClaims func(ctx *gin.Context, claims map[string]interface{}) *api_errors.ErrorResponse

// FirebaseAuthMiddleware structure
type FirebaseAuthMiddleware struct {
	service firebase.AuthService
}

// NewFirebaseAuthMiddleware creates new firebase authentication
func NewFirebaseAuthMiddleware(
	service firebase.AuthService,
) FirebaseAuthMiddleware {
	return FirebaseAuthMiddleware{
		service: service,
	}
}

// HandleAuth Handle handles auth requests
func (f FirebaseAuthMiddleware) HandleAuth(setClaims ...SetClaims) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := f.getTokenFromHeader(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, json_response.Error[string]{Error: err.Message})
			c.Abort()
			return
		}

		sentry.ConfigureScope(func(scope *sentry.Scope) {
			scope.SetUser(sentry.User{ID: token.UID})
		})

		for _, setClaim := range setClaims {
			err = setClaim(c, token.Claims)
			if err != nil {
				c.JSON(err.ErrorType.ToInt(), json_response.Error[string]{Error: err.Message})
				c.Abort()
				return
			}
		}

		c.Set(constants.UID, token.UID)
		c.Next()
	}
}

// HandleUserAuth Handle handles auth requests
func (f FirebaseAuthMiddleware) HandleUserAuth() gin.HandlerFunc {
	return f.HandleAuth(func(c *gin.Context, claims map[string]interface{}) *api_errors.ErrorResponse {
		role := claims[constants.Roles.Key].(constants.Role)
		if role != constants.Roles.User {
			return &api_errors.ErrorResponse{
				ErrorType: http.StatusUnauthorized,
				Message:   "unauthorized request",
			}
		}
		c.Set(constants.Roles.Key, role.ToString())

		userIdKey := constants.Claims.UserId.ToString()
		c.Set(userIdKey, int64(claims[userIdKey].(float64)))

		return nil
	})
}

// HandleAdminAuth handles middleware for roles
func (f FirebaseAuthMiddleware) HandleAdminAuth(allowedRoles ...constants.Role) gin.HandlerFunc {
	return f.HandleAuth(func(c *gin.Context, claims map[string]interface{}) *api_errors.ErrorResponse {
		role := claims[constants.Roles.Key].(constants.Role)
		if len(allowedRoles) > 0 {
			if !f.checkRoles(role, allowedRoles) {
				return &api_errors.ErrorResponse{
					ErrorType: http.StatusUnauthorized,
					Message:   "unauthorized request",
				}
			}
		}
		c.Set(constants.Roles.Key, role.ToString())

		adminIdKey := constants.Claims.AdminId.ToString()
		c.Set(adminIdKey, int64(claims[adminIdKey].(float64)))

		return nil
	})
}

// getTokenFromHeader gets token from header
func (f FirebaseAuthMiddleware) getTokenFromHeader(c *gin.Context) (*auth.Token, *api_errors.ErrorResponse) {
	header := c.GetHeader(constants.Headers.Authorization.ToString())
	idToken := strings.TrimSpace(strings.Replace(header, "Bearer", "", 1))

	token, err := f.service.VerifyToken(idToken)
	if err != nil {
		return nil, &api_errors.ErrorResponse{
			Message: err.Message,
		}
	}

	return token, nil
}

// checkRoles check if role is allowed
func (f FirebaseAuthMiddleware) checkRoles(role constants.Role, allowedRoles []constants.Role) bool {
	for _, allowed := range allowedRoles {
		if role == allowed {
			return true
		}
	}
	return false
}
