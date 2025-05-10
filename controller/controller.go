package controller

import (
	"encoding/base64"
	"errors"
	"io/ioutil"
	"my-project/service"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type FileController struct{}

// Read handles the GET request for reading a file
func (fc *FileController) Read(c *gin.Context) {
	filename := c.Param("filename")
	download := c.DefaultQuery("download", "false") == "true"

	if err := service.ReadFile(filename, download, c); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}


// Upload handles the POST request for uploading a file
func (fc *FileController) Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File not provided"})
		return
	}

	// Use the service to save file and metadata
	result, err := service.UploadFile(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"file": result,
	})
}

// For Controller Upload and Get Product Image

// Function to handle get all file from database
func (fc *FileController) GetProductImage(c *gin.Context) {
	filename := c.Param("filename")

	// Construct full path
	fullPath := filepath.Join("public", "uploads", "products", filename)

	// Security: Prevent directory traversal
	if strings.Contains(filename, "../") {
		c.JSON(http.StatusForbidden, gin.H{"error": "Invalid filename"})
		return
	}

	// Check if the file exists
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Image not found"})
		return
	}

	// Serve the file
	c.File(fullPath)
}


// UploadProductImage handles product image uploads specifically
func (fc *FileController) UploadProductImage(c *gin.Context) {
	file, err := c.FormFile("product_image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File not provided"})
		return
	}

	filename, err := service.UploadProductImage(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"file": filename,
	})
}


// Base64Upload handles the POST request for uploading a base64 encoded image
func (fc *FileController) Base64Upload(c *gin.Context) {
	var request struct {
		Folder string `json:"folder"`
		Image  string `json:"image"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	buffer, err := base64.StdEncoding.DecodeString(request.Image)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid base64 data"})
		return
	}

	destinationFolder := filepath.Join("public/uploads", sanitize(request.Folder))
	if err := os.MkdirAll(destinationFolder, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create directory"})
		return
	}

	fileName := uuid.New().String()
	filePath := filepath.Join(destinationFolder, fileName)

	if err := ioutil.WriteFile(filePath, buffer, 0644); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write file"})
		return
	}

	if isValid, err := verifyFile(filePath); !isValid {
		os.Remove(filePath)
		c.JSON(http.StatusBadRequest, gin.H{"error": "File verification failed: " + err.Error()})
		return
	}

	// url, err := service.UploadFileByPath(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// c.JSON(http.StatusOK, gin.H{"url": url})
}

func sanitize(text string) string {
	return strings.ToLower(strings.ReplaceAll(text, " ", "_"))
}

func verifyFile(filePath string) (bool, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return false, errors.New("failed to access file")
	}
	if fileInfo.Size() == 0 {
		return false, errors.New("empty file")
	}
	return true, nil
}
