package routes

import (
	"my-project/controller"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(api *gin.RouterGroup) {
	// File Route Group
	fileRouter := api.Group("/file")
	{
		// File-related routes
		controller.SetupFileRoutes(fileRouter)
	}
}
