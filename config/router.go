package config

import (
	"go.uber.org/fx"
)

// / Module exports dependency to container
var RouterModule = fx.Options(
	fx.Provide(RoutersConstructor),
)

// Routes contains multiple routes
type Routes []Route

// Route interface
type Route interface {
	Setup()
}

// NewRoutes sets up routes
func RoutersConstructor() Routes {
	return Routes{}
}

// Setup all the route
func (r Routes) Setup() {
	for _, route := range r {
		route.Setup()
	}
}
