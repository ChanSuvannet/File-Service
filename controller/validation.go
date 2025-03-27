package controller

import (
	"errors"
	"my-project/exceptions"
	"regexp"

	"github.com/gin-gonic/gin"
)

// ValidationRequest validates the folder field
func ValidationRequest(c *gin.Context) error {
	var request struct {
		Folder string `json:"folder" validate:"required,min=2,alphanum"`
		Image  string `json:"image" validate:"required"`
	}

	// Bind JSON input to the struct
	if err := c.ShouldBindJSON(&request); err != nil {
		return err
	}

	// Validate folder with custom logic (for regex and length checks)
	if request.Folder == "" {
		return errors.New("Field folder is required")
	} else if !isValidFolder(request.Folder) {
		return errors.New("Field folder is invalid")
	}

	// Validate image with base64 regex
	if !isValidBase64Image(request.Image) {
		return errors.New("Field image must be a valid base64")
	}

	return nil
}

// isValidFolder checks if the folder field is valid (only alphanumeric characters and hyphens allowed)
func isValidFolder(folder string) bool {
	re := regexp.MustCompile("^[A-Za-z0-9-]+$")
	return re.MatchString(folder)
}

// isValidBase64Image checks if the image field contains valid base64 data for an image
func isValidBase64Image(image string) bool {
	// Check if it's a valid base64-encoded image (PNG, JPG, JPEG, GIF)
	re := regexp.MustCompile(`^data:image\/(png|jpg|jpeg|gif);base64,[A-Za-z0-9+/=]+$`)
	return re.MatchString(image)
}

// UploadValidationMiddleware is a middleware that validates the request before processing
func UploadValidationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Validate request data
		if err := ValidationRequest(c); err != nil {
			// Return Unprocessable Entity status with validation error
			exception := exceptions.NewUnprocessableEntityException("Invalid Entity", []string{err.Error()})
			c.JSON(exception.StatusCode, exception)
			c.Abort()
			return
		}

		// Continue with the request if validation passes
		c.Next()
	}
}
