package firebase

import (
	"context"
	"firebase.google.com/go"
	"firebase.google.com/go/messaging"
)

type loggerCMClient interface {
	Fatalf(template string, args ...interface{})
}

type CMClientConfig struct {
	logger loggerCMClient
	app    *firebase.App
}

type CMClientService struct {
	*messaging.Client
}

// NewFirebaseCMClient creates new firebase cloud messaging client
func NewFirebaseCMClient(config CMClientConfig) CMClientService {
	ctx := context.Background()
	messagingClient, err := config.app.Messaging(ctx)
	if err != nil {
		config.logger.Fatalf("Firebase messaing: %v", err)
	}
	return CMClientService{
		Client: messagingClient,
	}
}
