package utility

import (
	"net/http"
	"path/filepath"

	"boilerplate-api/internal/api_errors"
	"boilerplate-api/internal/config"
	"boilerplate-api/internal/json_response"
	"boilerplate-api/internal/utils"
	"boilerplate-api/services/aws"
	"boilerplate-api/services/gcp"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	logger   config.Logger
	env      config.Env
	bucket   gcp.StorageBucketService
	s3Bucket aws.S3BucketService
}

func NewController(logger config.Logger,
	env config.Env,
	bucket gcp.StorageBucketService,
	s3Bucket aws.S3BucketService,
) Controller {
	return Controller{
		logger:   logger,
		env:      env,
		bucket:   bucket,
		s3Bucket: s3Bucket,
	}
}

// Response for the util scope
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    string      `json:"data"`
	Path    string      `json:"path"`
	Value   interface{} `json:"attributes"`
}

const storageURL string = "https://storage.googleapis.com/"

// FileUploadHandler handles file upload
func (uc Controller) FileUploadHandler(ctx *gin.Context) {
	file, uploadFile, err := ctx.Request.FormFile("file")
	if err != nil {
		uc.logger.Error("Error Get File from request :: ", err.Error())
		err := api_errors.BadRequest.Wrap(err, "Failed to get file form request")
		status, errM := api_errors.HandleError(err)
		ctx.JSON(status, json_response.Error{Error: errM})
		return
	}

	fileExtension := filepath.Ext(uploadFile.Filename)
	fileName := utils.GenerateRandomFileName() + fileExtension
	originalFileName := "images/original/" + fileName
	thumbnailFileName := "images/thumbnail/" + fileName

	// File type
	file1, _, _ := ctx.Request.FormFile("file")
	fileHeader := make([]byte, 512)
	if _, err := file1.Read(fileHeader); err != nil {
		uc.logger.Error("Error File Read upload File::", err.Error())
		err := api_errors.BadRequest.Wrap(err, "Failed to read upload  File")
		status, errM := api_errors.HandleError(err)
		ctx.JSON(status, json_response.Error{Error: errM})
		return
	}
	fileType := http.DetectContentType(fileHeader)
	if fileType == "image/png" || fileType == "image/jpg" || fileType == "image/jpeg" || fileType == "image/gif" {
		uploadedOriginalURL, err := uc.bucket.UploadFile(ctx.Request.Context(), file, originalFileName)
		if err != nil {
			uc.logger.Error("Error Failed to upload File::", err.Error())
			err := api_errors.BadRequest.Wrap(err, "Failed to upload File")
			status, errM := api_errors.HandleError(err)
			ctx.JSON(status, json_response.Error{Error: errM})
			return
		}

		//uploadedthumbnail
		thumbnail, err := utils.CreateThumbnail(file, fileType, 200, 0)
		if err != nil {
			uc.logger.Error("Error Failed create thumbnail", err.Error())
			err := api_errors.BadRequest.Wrap(err, "Error Failed create thumbnail")
			status, errM := api_errors.HandleError(err)
			ctx.JSON(status, json_response.Error{Error: errM})
			return
		}
		uploadThumbnailUrl, err := uc.bucket.UploadThumbnailFile(ctx.Request.Context(), thumbnail, thumbnailFileName, fileExtension)
		if err != nil {
			uc.logger.Error("Error Failed to upload File::", err.Error())
			err := api_errors.BadRequest.Wrap(err, "Failed to upload thumbnail File")
			status, errM := api_errors.HandleError(err)
			ctx.JSON(status, json_response.Error{Error: errM})
			return
		}

		response := &Response{
			Success: true,
			Message: "Uploaded Successfully",
			Data:    storageURL + uc.env.StorageBucketName + "/" + uploadedOriginalURL,
			Path:    uploadedOriginalURL,
			Value: map[string]string{
				"original_image_url":   storageURL + uc.env.StorageBucketName + "/" + uploadedOriginalURL,
				"original_image_path":  uploadedOriginalURL,
				"thumbnail_image_url":  storageURL + uc.env.StorageBucketName + "/" + uploadThumbnailUrl,
				"thumbnail_image_path": uploadThumbnailUrl,
			}}
		ctx.JSON(http.StatusOK, response)
		return
	}

	originalFileName = "files/" + fileName
	uploadedFileURL, err := uc.bucket.UploadFile(ctx.Request.Context(), file, originalFileName)
	if err != nil {
		uc.logger.Error("Error Failed to upload File::", err.Error())
		err := api_errors.BadRequest.Wrap(err, "Failed to upload file")
		status, errM := api_errors.HandleError(err)
		ctx.JSON(status, json_response.Error{Error: errM})
		return
	}
	response := &Response{
		Success: true,
		Message: "Uploaded Successfully",
		Data:    storageURL + uc.env.StorageBucketName + "/" + uploadedFileURL,
		Path:    uploadedFileURL,
	}
	ctx.JSON(http.StatusOK, response)
}

// Input model
type Input struct {
	Path *string `form:"path" json:"path" binding:"required"`
}

// FileUploadS3Handler handles aws s3 file upload
func (uc Controller) FileUploadS3Handler(ctx *gin.Context) {
	file, fileHeader, err := ctx.Request.FormFile("file")
	if err != nil {
		uc.logger.Error("Error Get File from request: ", err.Error())
		err := api_errors.BadRequest.Wrap(err, "Failed to get file form request")
		status, errM := api_errors.HandleError(err)
		ctx.JSON(status, json_response.Error{Error: errM})
		return
	}
	var input Input
	err = ctx.ShouldBind(&input)
	if err != nil {
		uc.logger.Error("Error Failed to bind input:: ", err.Error())
		err := api_errors.BadRequest.Wrap(err, "Failed to bind")
		status, errM := api_errors.HandleError(err)
		ctx.JSON(status, json_response.Error{Error: errM})
		return
	}

	fileExtension := filepath.Ext(fileHeader.Filename)
	fileName := utils.GenerateRandomFileName() + fileExtension
	originalFileNamePath := *input.Path + "/" + fileName

	uploadedFileURL, err := uc.s3Bucket.UploadToS3(file, fileHeader, originalFileNamePath)
	if err != nil {
		uc.logger.Error("Error Failed to upload File:: ", err.Error())
		err := api_errors.BadRequest.Wrap(err, "Failed to upload file to s3 bucket")
		status, errM := api_errors.HandleError(err)
		ctx.JSON(status, json_response.Error{Error: errM})
		return
	}

	response := &Response{
		Success: true,
		Message: "Uploaded Successfully",
		Path:    uploadedFileURL,
		Data:    uploadedFileURL,
	}
	ctx.JSON(http.StatusOK, response)
}
