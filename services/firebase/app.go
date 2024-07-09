package firebase

import (
	"context"
	"firebase.google.com/go"
	"go.uber.org/zap"
	"google.golang.org/api/option"
)

// AppConfig structure
type AppConfig struct {
	logger *zap.SugaredLogger
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
