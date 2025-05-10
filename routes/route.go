package routes

import (
	"my-project/controller"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(api *gin.RouterGroup) {
	fileController := new(controller.FileController)
    
    api.GET("/file/:filename", fileController.Read)
    api.POST("/file/upload-single", fileController.Upload) 

	api.GET("/file/product/image/:filename", fileController.GetProductImage)
	api.POST("/file/product/upload-image", fileController.UploadProductImage)
}
