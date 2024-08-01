package utility

import (
	"boilerplate-api/internal/router"
)

// SetupRoutes Setup sets up route for util entities
func SetupRoutes(
	router router.Router,
	utilityController Controller,
) {
	utils := router.V1.Group("/utils")
	{
		utils.POST("/file-upload", utilityController.FileUploadHandler)
		utils.POST("/s3-file-upload", utilityController.FileUploadS3Handler)
	}
}
