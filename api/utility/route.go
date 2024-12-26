package utility

import (
	"boilerplate-api/lib/router"
)

// SetupRoutes Setup sets up route for util entities
func SetupRoutes(
	router router.Router,
	utilityController Controller,
) {
	utils := router.V1.Group("/utils")
	{
		utils.POST("/files/upload", utilityController.FileUploadHandler)
		utils.GET("/images/signed_url", utilityController.GetSignedUrl)
		utils.POST("/s3-file-upload", utilityController.FileUploadS3Handler)
	}
}
