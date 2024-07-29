package firebase

import (
	"context"
	"firebase.google.com/go"
	"google.golang.org/api/option"
)

type appConfigLogger interface {
	Info(args ...interface{})
	Fatalf(template string, args ...interface{})
}

// AppConfig structure
type AppConfig struct {
	logger appConfigLogger
	opt    *option.ClientOption
}

// AppService structure
type AppService struct {
	*firebase.App
}

// NewFirebaseApp creates new firebase app instance
func NewFirebaseApp(config AppConfig) AppService {
	app, err := firebase.NewApp(context.Background(), nil, *config.opt)
	if err != nil {
		config.logger.Fatalf("Firebase NewApp: %v", err)
	}

	config.logger.Info("✅ Firebase app initialized.")
	return AppService{
		App: app,
	}
}
