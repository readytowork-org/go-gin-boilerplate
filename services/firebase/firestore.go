package firebase

import (
	"boilerplate-api/internal/config"
	"cloud.google.com/go/firestore"
	"context"
	"firebase.google.com/go"
)

// NewFirestoreClient creates new firestore client
func NewFirestoreClient(logger config.Logger, app *firebase.App) *firestore.Client {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	firestoreClient, err := app.Firestore(ctx)
	if err != nil {
		logger.Fatalf("Firestore client: %v", err)
	}

	return firestoreClient
}
