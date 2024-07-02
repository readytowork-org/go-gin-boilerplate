package firebase

import (
	"boilerplate-api/internal/config"
	"context"
	"firebase.google.com/go"
	"firebase.google.com/go/messaging"
)

// NewFirebaseCMClient creates new firebase cloud messaging client
func NewFirebaseCMClient(logger config.Logger, app *firebase.App) *messaging.Client {
	ctx := context.Background()
	messagingClient, err := app.Messaging(ctx)
	if err != nil {
		logger.Fatalf("Firebase messaing: %v", err)
	}
	return messagingClient
}
