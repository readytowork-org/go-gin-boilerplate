package firebase

import (
	"context"

	"boilerplate-api/internal/config"
	"firebase.google.com/go"
	"google.golang.org/api/option"
)

// NewFirebaseApp creates new firebase app instance
func NewFirebaseApp(logger config.Logger, opt *option.ClientOption) *firebase.App {
	app, err := firebase.NewApp(context.Background(), nil, *opt)
	if err != nil {
		logger.Fatalf("Firebase NewApp: %v", err)
	}

	logger.Info("✅ Firebase app initialized.")
	return app
}
