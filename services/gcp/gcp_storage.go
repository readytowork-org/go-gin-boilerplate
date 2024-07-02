package gcp

import (
	"boilerplate-api/internal/config"
	"boilerplate-api/internal/utils"
	"bytes"
	"context"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"mime/multipart"
	"net/url"
	"strings"
	"time"

	"cloud.google.com/go/storage"
)

// StorageBucketService the file upload/download functions
type StorageBucketService struct {
	client BucketClient
	logger config.Logger
	env    config.Env
}

// NewStorageBucketService for the StorageBucketService struct
func NewStorageBucketService(
	logger config.Logger,
	env config.Env,
	client BucketClient,
) StorageBucketService {
	return StorageBucketService{
		client: client,
		logger: logger,
		env:    env,
	}
}

func (service StorageBucketService) GetImageUrl(ctx context.Context, image multipart.File, imageFileHeader *multipart.FileHeader) (uploadedUrl string, err error) {
	if imageFileHeader != nil && image != nil {
		fileName, _ := utils.GetFileName(imageFileHeader.Filename)
		originalFileName := "images/" + fileName
		uploadedUrl, err = service.UploadFile(ctx, image, originalFileName)
		if err != nil {
			return uploadedUrl, err
		}
	}
	return uploadedUrl, nil
}

// UploadFile uploads the file to the cloud storage
func (service StorageBucketService) UploadFile(
	ctx context.Context,
	file multipart.File,
	fileName string,
) (string, error) {
	bucketName := service.env.StorageBucketName

	if bucketName == "" {
		service.logger.Fatal("Please check your env file for StorageBucketName")
	}

	_, err := service.client.Bucket(bucketName).Attrs(ctx)

	if errors.Is(err, storage.ErrBucketNotExist) {
		service.logger.Fatalf("provided bucket %v doesn't exists", bucketName)
	}

	if err != nil {
		service.logger.Fatalf("cloud bucket error: %v", err.Error())
	}

	wc := service.client.Bucket(bucketName).Object(fileName).NewWriter(ctx)
	wc.ContentType = "application/octet-stream"

	if _, err := io.Copy(wc, file); err != nil {
		return "", err
	}

	if err := wc.Close(); err != nil {
		return "", err
	}

	return fileName, nil
}

// UploadBinary the binary to the cloud storage
func (service StorageBucketService) UploadBinary(
	ctx context.Context,
	file []byte,
	fileName string,
) (string, error) {

	var bucketName = service.env.StorageBucketName

	if bucketName == "" {
		service.logger.Fatal("Please check your env file for StorageBucketName")
	}

	_, err := service.client.Bucket(bucketName).Attrs(ctx)

	if err == storage.ErrBucketNotExist {
		service.logger.Fatalf("provided bucket %v doesn't exists", bucketName)
	}

	if err != nil {
		service.logger.Fatalf("cloud bucket error: %v", err.Error())
	}

	wc := service.client.Bucket(bucketName).Object(fileName).NewWriter(ctx)
	wc.ContentType = "application/octet-stream"

	if _, err := io.Copy(wc, bytes.NewReader(file)); err != nil {
		return "", err
	}

	if err := wc.Close(); err != nil {
		return "", err
	}

	u, err := url.ParseRequestURI("/" + bucketName + "/" + wc.Attrs().Name)

	if err != nil {
		return "", err
	}

	path := u.EscapedPath()
	path = strings.Replace(path, "/"+bucketName, "", 1)
	path = strings.Replace(path, "/", "", 1)

	return path, nil

}

// RemoveObject removes the file from the storage bucket
func (service StorageBucketService) RemoveObject(objectName string) error {

	bucketName := service.env.StorageBucketName
	if bucketName == "" {
		service.logger.Fatal("Please check your env file for StorageBucketName")
	}
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	objectToDelete := service.client.Bucket(bucketName).Object(objectName)
	attrs, err := objectToDelete.Attrs(ctx)
	if err != nil {
		return fmt.Errorf("Object(%v).Attrs: %v", objectToDelete, err)
	}
	if err != nil {
		return fmt.Errorf("object.Attrs: %v", err)
	}
	objectToDelete = objectToDelete.If(storage.Conditions{GenerationMatch: attrs.Generation})

	err = objectToDelete.Delete(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (service StorageBucketService) UploadThumbnailFile(ctx context.Context,
	file image.Image,
	fileName string, extension string) (string, error) {

	var bucketName = service.env.StorageBucketName
	if bucketName == "" {
		service.logger.Fatal("Please check your env file for StorageBucketName")
	}

	_, err := service.client.Bucket(bucketName).Attrs(ctx)
	if errors.Is(err, storage.ErrBucketNotExist) {
		service.logger.Fatalf("provided bucket %v doesn't exists", bucketName)
	}
	if err != nil {
		service.logger.Fatalf("cloud bucket error: %v", err.Error())
	}

	wc := service.client.Bucket(bucketName).Object(fileName).NewWriter(ctx)
	wc.ContentType = "application/octet-stream"

	if extension == "jpg" || extension == "jpeg" {
		err = jpeg.Encode(wc, file, nil)
	} else {
		err = png.Encode(wc, file)
	}

	if err != nil {
		return "", err
	}

	if err := wc.Close(); err != nil {
		return "", err
	}

	u, err := url.ParseRequestURI("/" + bucketName + "/" + wc.Attrs().Name)
	if err != nil {
		return "", err
	}

	path := u.EscapedPath()
	path = strings.Replace(path, "/"+bucketName, "", 1)
	path = strings.Replace(path, "/", "", 1)

	return path, nil

}
