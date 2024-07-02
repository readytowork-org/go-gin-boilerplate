package gcp

import (
	"context"
	"errors"

	"boilerplate-api/internal/config"
	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

type BucketClient struct {
	*storage.Client
}

// NewGCPBucketClient creates a new gcp bucket api client
func NewGCPBucketClient(logger config.Logger, env config.Env, clientOption *option.ClientOption) BucketClient {
	bucketName := env.StorageBucketName
	ctx := context.Background()
	if bucketName == "" {
		logger.Error("Please check your env file for STORAGE_BUCKET_NAME")
	}
	client, err := storage.NewClient(ctx, *clientOption)
	if err != nil {
		logger.Fatal(err.Error())
	}

	bucket := client.Bucket(bucketName)
	_, err = bucket.Attrs(ctx)
	if errors.Is(err, storage.ErrBucketNotExist) {
		logger.Fatalf("Provided bucket %v doesn't exists", bucketName)
	}

	if err != nil {
		logger.Fatalf("Cloud bucket error: %v", err.Error())
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
		logger.Fatalf("Cloud bucket update error: %v", err.Error())
	}
	return BucketClient{
		client,
	}
}
