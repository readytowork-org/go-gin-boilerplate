package aws

import (
	"boilerplate-api/internal/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"mime/multipart"

	"context"
)

// S3BucketService handles the file upload functions
type S3BucketService struct {
	client *s3.Client
	logger config.Logger
	env    config.Env
}

// NewS3BucketService initialization for the AWS S3 BucketService struct
func NewS3BucketService(
	logger config.Logger,
	config aws.Config,
	env config.Env,
) S3BucketService {
	client := s3.New(s3.Options{Credentials: config.Credentials, Region: env.AwsS3Region})
	logger.Info("✅  AWS S3 service created")
	return S3BucketService{
		client: client,
		logger: logger,
		env:    env,
	}
}

// UploadToS3 uploads the file to the aws s3 bucket
func (s S3BucketService) UploadToS3(
	file multipart.File,
	fileHeader *multipart.FileHeader,
	fileName string,
) (string, error) {
	uploader := manager.NewUploader(s.client)
	result, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(s.env.AwsS3Bucket),
		Key:         aws.String(fileName),
		Body:        file,
		ContentType: aws.String(fileHeader.Header.Get("content-type")),
		ACL:         types.ObjectCannedACLPublicRead,
	})
	if err != nil {
		s.logger.Fatalf("aws s3 cloud bucket upload error: %v", err.Error())
		return "", err
	}
	return result.Location, nil
}
