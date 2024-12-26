package utility

import (
	"net/http"
	"path/filepath"

	"boilerplate-api/lib/config"
	"boilerplate-api/lib/json_response"
	"boilerplate-api/lib/utils"
	"boilerplate-api/services/aws"
	"github.com/gin-gonic/gin"
)

type Controller struct {
	logger   config.Logger
	s3Bucket aws.S3BucketService
	service  Service
}

func NewController(
	logger config.Logger,
	s3Bucket aws.S3BucketService,
	service Service,
) Controller {
	return Controller{
		logger:   logger,
		s3Bucket: s3Bucket,
		service:  service,
	}
}

//	@Tags			UtilityApi
//	@Summary		handles file upload
//	@Description	handles file upload
//	@Security		Bearer
//	@Produce		application/json
//	@Param			file	formData	file		true	"Upload File"
//	@Success		200		{object}	Response	"File Uploaded Successfully"
//	@Failure		400		{object}	json_response.Error[string]
//	@Router			/api/v1/utils/file-upload [post]
//	@Id				FileUpload
func (uc Controller) FileUploadHandler(ctx *gin.Context) {
	file, uploadFile, err := ctx.Request.FormFile("file")
	if err != nil {
		uc.logger.Error("Error Get File from request :: ", err.Error())
		ctx.JSON(
			http.StatusBadRequest, json_response.Error[string]{
				Error:   err.Error(),
				Message: "Failed to get file from request",
			},
		)
		return
	}

	message, response, err := uc.service.UploadImage(file, uploadFile)
	if err != nil {
		uc.logger.Error("Error Upload File from request :: ", err.Error())
		ctx.JSON(
			message.StatusCode, json_response.Error[string]{
				Error:   err.Error(),
				Message: message.Message,
			},
		)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

//	@Tags			UtilityApi
//	@Summary		GetSignedUrl
//	@Description	generate signed url
//	@Security		Bearer
//	@Produce		application/json
//	@Param			image_url	query		string	false	"Image Url"
//	@Success		200			{object}	json_response.Data[string]
//	@Failure		400			{object}	json_response.Error[string]
//	@Router			/api/v1/utils/images/signed_url [get]
//	@Id				GetSignedUrl
func (uc Controller) GetSignedUrl(ctx *gin.Context) {
	imageUrl := ctx.Query("image_url")
	if imageUrl == "" {
		ctx.JSON(
			http.StatusBadRequest, json_response.Error[string]{
				Message: "Image Url is invalid",
			},
		)
	}

	signedUrl, err := uc.service.GetSignedUrl(imageUrl)
	if err != nil {
		uc.logger.Error("Error Failed to convert signed url:", err.Error())
		ctx.JSON(
			http.StatusOK, json_response.Error[string]{
				Message: "Error Failed to convert signed url",
			},
		)
		return
	}

	ctx.Redirect(http.StatusFound, signedUrl)
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
		ctx.JSON(
			http.StatusBadRequest, json_response.Error[string]{
				Error:   err.Error(),
				Message: "Failed to get file from request",
			},
		)
		return
	}
	var input Input
	err = ctx.ShouldBind(&input)
	if err != nil {
		uc.logger.Error("Error Failed to bind input:: ", err.Error())
		ctx.JSON(
			http.StatusBadRequest, json_response.Error[string]{
				Error:   err.Error(),
				Message: "Failed to bind",
			},
		)
		return
	}

	fileExtension := filepath.Ext(fileHeader.Filename)
	fileName := utils.GenerateRandomFileName() + fileExtension
	originalFileNamePath := *input.Path + "/" + fileName

	uploadedFileURL, err := uc.s3Bucket.UploadToS3(file, fileHeader, originalFileNamePath)
	if err != nil {
		uc.logger.Error("Error Failed to upload File:: ", err.Error())
		ctx.JSON(
			http.StatusBadRequest, json_response.Error[string]{
				Error:   err.Error(),
				Message: "Failed to upload file to s3 bucket",
			},
		)
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
