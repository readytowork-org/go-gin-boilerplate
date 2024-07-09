package firebase

import (
	"cloud.google.com/go/firestore"
	"context"
	"firebase.google.com/go"
	"go.uber.org/zap"
)

type StoreClientConfig struct {
	logger *zap.SugaredLogger
	app    *firebase.App
}

type StoreClientService struct {
	*firestore.Client
}

// NewFireStoreClient creates new firestore client
func NewFireStoreClient(config StoreClientConfig) StoreClientService {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	firestoreClient, err := config.app.Firestore(ctx)
	if err != nil {
		config.logger.Fatalf("Firestore client: %v", err)
	}

	return StoreClientService{Client: firestoreClient}
}
