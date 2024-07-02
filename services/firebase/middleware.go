package firebase

import (
	"errors"
	"net/http"
	"strings"

	"boilerplate-api/internal/constants"
	"boilerplate-api/internal/json_response"
	"boilerplate-api/internal/types"

	"firebase.google.com/go/auth"
	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
)

type SetClaims func(ctx *gin.Context, claims types.MapString) (int, error)

// AuthMiddleware structure
type AuthMiddleware struct {
	service AuthService
}

// NewFirebaseAuthMiddleware creates new firebase authentication
func NewFirebaseAuthMiddleware(
	service AuthService,
) AuthMiddleware {
	return AuthMiddleware{
		service: service,
	}
}

// HandleAuth Handle handles auth requests
func (f AuthMiddleware) HandleAuth(setClaims ...SetClaims) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := f.getTokenFromHeader(c)

		if err != nil {
			c.JSON(http.StatusUnauthorized, json_response.ErrorMsg{Error: err.Error()})
			c.Abort()
			return
		}

		sentry.ConfigureScope(func(scope *sentry.Scope) {
			scope.SetUser(sentry.User{ID: token.UID})
		})

		for _, setClaim := range setClaims {
			httpCode, err := setClaim(c, token.Claims)
			if err != nil {
				c.JSON(httpCode, json_response.ErrorMsg{Error: err.Error()})
				c.Abort()
				return
			}
		}

		c.Set(constants.UID, token.UID)
		c.Next()
	}
}

// HandleUserAuth Handle handles auth requests
func (f AuthMiddleware) HandleUserAuth() gin.HandlerFunc {
	return f.HandleAuth(func(c *gin.Context, claims types.MapString) (int, error) {
		role := claims[constants.Roles.Key].(constants.Role)
		if role != constants.Roles.User {
			return http.StatusUnauthorized, errors.New("unauthorized request")
		}
		c.Set(constants.Roles.Key, role.ToString())

		userIdKey := constants.Claims.UserId.ToString()
		c.Set(userIdKey, int64(claims[userIdKey].(float64)))

		return -1, nil
	})
}

// HandleAdminAuth handles middleware for roles
func (f AuthMiddleware) HandleAdminAuth(allowedRoles ...constants.Role) gin.HandlerFunc {
	return f.HandleAuth(func(c *gin.Context, claims types.MapString) (int, error) {
		role := claims[constants.Roles.Key].(constants.Role)
		if len(allowedRoles) > 0 {
			if !f.checkRoles(role, allowedRoles) {
				return http.StatusUnauthorized, errors.New("unauthorized request")
			}
		}
		c.Set(constants.Roles.Key, role.ToString())

		adminIdKey := constants.Claims.AdminId.ToString()
		c.Set(adminIdKey, int64(claims[adminIdKey].(float64)))

		return -1, nil
	})
}

// getTokenFromHeader gets token from header
func (f AuthMiddleware) getTokenFromHeader(c *gin.Context) (*auth.Token, error) {
	header := c.GetHeader(constants.Headers.Authorization.ToString())
	idToken := strings.TrimSpace(strings.Replace(header, "Bearer", "", 1))

	token, err := f.service.VerifyToken(idToken)
	if err != nil {
		return nil, err
	}

	return token, nil
}

// checkRoles check if role is allowed
func (f AuthMiddleware) checkRoles(role constants.Role, allowedRoles []constants.Role) bool {
	for _, allowed := range allowedRoles {
		if role == allowed {
			return true
		}
	}
	return false
}
