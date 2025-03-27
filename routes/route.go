package routes

import (
	"my-project/controller"

	"github.com/gin-gonic/gin"
)

// SetupRoutes sets up the main routes of the application
func SetupRoutes(router *gin.Engine) {
	// File Route Group
	fileRouter := router.Group("/file")
	{
		// File-related routes
		controller.SetupFileRoutes(fileRouter)
	}
}
