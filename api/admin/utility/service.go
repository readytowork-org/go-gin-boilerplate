package utility

import (
	"context"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"

	"boilerplate-api/lib/config"
	"boilerplate-api/lib/constants"
	"boilerplate-api/lib/utils"
)

type GcpStorageBucketService interface {
	UploadFile(ctx context.Context, file io.Reader, fileName string) (string, error)
}

type UploadResponse struct {
	Message    string
	StatusCode int
}

// Response for the util scope
type Response struct {
	Success bool        `json:"success" validate:"required"`
	Message string      `json:"message" validate:"required"`
	Data    string      `json:"data" validate:"required"`
	Path    string      `json:"path" validate:"required"`
	Value   interface{} `json:"attributes" validate:"required"`
}

type Service struct {
	logger config.Logger
	env    config.Env
	bucket GcpStorageBucketService
}

func NewService(
	logger config.Logger,
	env config.Env,
	bucket GcpStorageBucketService,
) Service {
	return Service{
		logger: logger,
		env:    env,
		bucket: bucket,
	}
}

func (s Service) UploadImage(file multipart.File, uploadFile *multipart.FileHeader) (UploadResponse, *Response, error) {
	fileType := uploadFile.Header.Get("Content-Type")
	fileExtension := filepath.Ext(uploadFile.Filename)
	fileName := utils.GenerateRandomFileName() + fileExtension

	originalFileName := "images/original/" + fileName
	switch fileType {
	case "image/png",
		"image/jpg",
		"image/jpeg",
		"image/gif",
		"image/webp":
		{
			uploadedOriginalURL, errs := s.bucket.UploadFile(context.Background(), file, originalFileName)
			if errs != nil {
				s.logger.Error("Error Failed to upload File::", errs.Error())
				return UploadResponse{
					StatusCode: http.StatusBadRequest,
					Message:    "Failed to upload File",
				}, nil, errs
			}

			//upload thumbnail
			thumbnailFileName := "images/thumbnail/" + fileName

			thumbnailImg, err := utils.CreateThumbnail(file, fileType, 200, 0)
			if err != nil {
				s.logger.Error("Error Failed create thumbnail", err.Error())
				return UploadResponse{
					Message:    "Failed to create thumbnail",
					StatusCode: http.StatusBadRequest,
				}, nil, err
			}

			uploadThumbnailUrl, err := s.bucket.UploadFile(context.Background(), thumbnailImg, thumbnailFileName)
			if err != nil {
				s.logger.Error("Error Failed to upload File::", err.Error())
				return UploadResponse{
					Message:    "Failed to upload thumbnail File",
					StatusCode: http.StatusBadRequest,
				}, nil, err
			}

			signedURL, err := utils.GetObjectSignedURL(s.env.StorageBucketName, uploadedOriginalURL)
			if err != nil {
				s.logger.Error("Error Failed to convert signed url:", err.Error())
				return UploadResponse{
					Message:    "Failed to convert signed url",
					StatusCode: http.StatusBadRequest,
				}, nil, err
			}

			return UploadResponse{
					Message:    "Uploaded Successfully",
					StatusCode: http.StatusOK,
				}, &Response{
					Success: true,
					Message: "Uploaded Successfully",
					Data:    signedURL,
					Path:    uploadedOriginalURL,
					Value: map[string]string{
						"original_image_url":   constants.STORAGE_URL + s.env.StorageBucketName + uploadedOriginalURL,
						"original_image_path":  uploadedOriginalURL,
						"thumbnail_image_url":  constants.STORAGE_URL + s.env.StorageBucketName + uploadThumbnailUrl,
						"thumbnail_image_path": uploadThumbnailUrl,
					},
				}, nil
		}
	default:
		originalFileName = "files/" + fileName
		uploadedFileURL, err := s.bucket.UploadFile(context.Background(), file, originalFileName)
		if err != nil {
			s.logger.Error("Error Failed to upload File::", err.Error())
			return UploadResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "Failed to upload File",
			}, nil, err
		}

		response := &Response{
			Success: true,
			Message: "Uploaded Successfully",
			Data:    constants.STORAGE_URL + s.env.StorageBucketName + "/" + uploadedFileURL,
			Path:    uploadedFileURL,
		}

		return UploadResponse{
			StatusCode: http.StatusOK,
			Message:    "Uploaded Successfully",
		}, response, nil
	}
}
