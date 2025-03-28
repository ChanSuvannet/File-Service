package middleware

import (
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SingleFileMulter is a middleware that handles single file upload validation
func SingleFileMulter() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the file from the request
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
			c.Abort()
			return
		}

		// Validate file type (allow only images, adjust as needed)
		if !isAllowedFileType(file) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type"})
			c.Abort()
			return
		}

		// Set file in context for later use
		c.Set("uploadedFile", file)

		// Proceed to the next handler
		c.Next()
	}
}

// isAllowedFileType checks if the file type is allowed
func isAllowedFileType(file *multipart.FileHeader) bool {
	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/gif":  true,
	}

	// Open the file to read its MIME type
	src, err := file.Open()
	if err != nil {
		return false
	}
	defer src.Close()

	// Get the first 512 bytes to detect file type
	buffer := make([]byte, 512)
	_, err = src.Read(buffer)
	if err != nil {
		return false
	}

	// Detect the content type
	contentType := http.DetectContentType(buffer)
	return allowedTypes[contentType]
}
