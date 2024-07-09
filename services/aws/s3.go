package aws

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"go.uber.org/zap"
	"mime/multipart"

	"context"
)

type S3BucketConfig struct {
	logger   *zap.SugaredLogger
	s3bucket string
	s3Region string
	config   aws.Config
}

// S3BucketService handles the file upload functions
type S3BucketService struct {
	client   *s3.Client
	logger   *zap.SugaredLogger
	s3bucket string
	s3Region string
}

// NewS3BucketService initialization for the AWS S3 BucketService struct
func NewS3BucketService(
	bucketConfig S3BucketConfig,
) S3BucketService {
	client := s3.New(s3.Options{Credentials: bucketConfig.config.Credentials, Region: bucketConfig.s3Region})
	bucketConfig.logger.Info("✅  AWS S3 service created")
	return S3BucketService{
		client:   client,
		logger:   bucketConfig.logger,
		s3bucket: bucketConfig.s3bucket,
		s3Region: bucketConfig.s3Region,
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
		Bucket:      aws.String(s.s3bucket),
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
