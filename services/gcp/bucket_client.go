package gcp

import (
	"cloud.google.com/go/storage"
	"context"
	"errors"
	"go.uber.org/zap"
	"google.golang.org/api/option"
)

type BucketClientConfig struct {
	logger            *zap.SugaredLogger
	storageBucketName string
	clientOption      *option.ClientOption
}

type BucketClient struct {
	*storage.Client
}

// NewGCPBucketClient creates a new gcp bucket api client
func NewGCPBucketClient(config BucketClientConfig) BucketClient {
	bucketName := config.storageBucketName
	ctx := context.Background()
	if bucketName == "" {
		config.logger.Error("Please check your env file for STORAGE_BUCKET_NAME")
	}
	client, err := storage.NewClient(ctx, *config.clientOption)
	if err != nil {
		config.logger.Fatal(err.Error())
	}

	bucket := client.Bucket(bucketName)
	_, err = bucket.Attrs(ctx)
	if errors.Is(err, storage.ErrBucketNotExist) {
		config.logger.Fatalf("Provided bucket %v doesn't exists", bucketName)
	}

	if err != nil {
		config.logger.Fatalf("Cloud bucket error: %v", err.Error())
	}

	bucketAttrsToUpdate := storage.BucketAttrsToUpdate{
		CORS: []storage.CORS{
			{
				MaxAge:          600,
				Methods:         []string{"PUT", "PATCH", "GET", "POST", "OPTIONS", "DELETE"},
				Origins:         []string{"*"},
				ResponseHeaders: []string{"Content-Type"},
			}},
	}
	if _, err := bucket.Update(ctx, bucketAttrsToUpdate); err != nil {
		config.logger.Fatalf("Cloud bucket update error: %v", err.Error())
	}
	return BucketClient{
		client,
	}
}
