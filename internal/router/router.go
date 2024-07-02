package router

import (
	"fmt"
	"net/http"

	"boilerplate-api/internal/config"

	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Router Gin Router
type Router struct {
	*gin.Engine
	V1 *gin.RouterGroup
}

// NewRouter : all the routes are defined here
func NewRouter(env config.Env, logger config.Logger) Router {
	appEnv := env.Environment

	if appEnv != "local" {
		if err := sentry.Init(sentry.ClientOptions{
			Dsn:              env.SentryDSN,
			Environment:      `Demo-backend-` + env.Environment,
			AttachStacktrace: true,
		}); err != nil {
			fmt.Printf("Sentry initialization failed: %v\n", err)
		}
	}

	// TODO :: after cli config
	// gin.DefaultWriter = logger.GetGinLogger()

	if appEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	httpRouter := gin.Default()

	httpRouter.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "PATCH", "GET", "POST", "OPTIONS", "DELETE"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
	}))

	httpRouter.Use(sentrygin.New(sentrygin.Options{
		Repanic: true,
	}))

	httpRouter.GET("/health-check", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": "Demo  📺 API Up and Running"})
	})

	api := httpRouter.Group("/api")
	v1 := api.Group("/v1")

	return Router{
		Engine: httpRouter,
		V1:     v1,
	}
}
