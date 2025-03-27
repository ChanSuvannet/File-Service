package controller

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// Maximum file size (512MB)
const maxFileSize = 512 * 1024 * 1024

// FileUploadHandler handles the file upload process
func FileUploadHandler(c *gin.Context) {
	// Validate form
	if err := c.Request.ParseMultipartForm(maxFileSize); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "File size limit exceeded or form parse error.",
		})
		return
	}

	// Get the file from the form data
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "No file uploaded! Please select a file.",
		})
		return
	}
	defer file.Close()

	// Custom validation check (you can implement your own validation here)
	if err := validateFile(file); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Extract the folder name from the request body
	folder := sanitizeFolderName(c.DefaultPostForm("folder", "unknown"))

	// Create destination folder
	destinationFolder := fmt.Sprintf("public/uploads/%s/", folder)
	if err := os.MkdirAll(destinationFolder, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create destination folder.",
		})
		return
	}

	// Create the file in the destination folder
	fileName := generateFileName() // Function to generate the new file name (e.g., based on timestamp or UUID)
	filePath := filepath.Join(destinationFolder, fileName)
	dstFile, err := os.Create(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create file in destination folder.",
		})
		return
	}
	defer dstFile.Close()

	// Copy the file data to the new file
	if _, err := io.Copy(dstFile, file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to save the uploaded file.",
		})
		return
	}

	// Successfully uploaded
	c.JSON(http.StatusOK, gin.H{
		"message": "File uploaded successfully",
		"path":    filePath,
	})
}

// sanitizeFolderName sanitizes the folder name (removes unwanted characters)
func sanitizeFolderName(text string) string {
	return strings.ToLower(strings.ReplaceAll(text, " ", "_"))
}

// validateFile is a placeholder function for custom file validation
func validateFile(file multipart.File) error {
	// Implement your custom file validation here (e.g., check file type, size, etc.)
	return nil
}

// generateFileName generates a unique file name (e.g., timestamp or UUID)
func generateFileName() string {
	// You can use a timestamp or any other mechanism to generate a unique name
	// This is a placeholder example
	return fmt.Sprintf("%d.txt", time.Now().Unix())
}
