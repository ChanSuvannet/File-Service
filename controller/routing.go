package controller

import (
	"my-project/middleware"

	"github.com/gin-gonic/gin"
)

// SetupFileRoutes configures file-related routes
func SetupFileRoutes(router *gin.RouterGroup) {
	// Create an instance of FileController
	fileController := &FileController{}

	// File-related routes
	// Route for fetching the file by filename
	router.GET("/:filename", fileController.Read)

	// Route for single file upload (commented out for now)
	router.POST("/upload-single", middleware.SingleFileMulter(), fileController.Upload)

	// Route for get product image (commented out for now)
	router.GET("/product/image/:filename", fileController.GetProductImage)

	// Route for product image upload (commented out for now)
	router.POST("/product/upload-image", fileController.UploadProductImage)

	// Route for base64 upload with validation (commented out for now)
	// router.POST("/upload-base64", validation.UploadValidation(), fileController.Base64Upload)
}
