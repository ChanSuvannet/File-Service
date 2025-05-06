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
		return
	}
}

// Upload handles the POST request for uploading a file
func (fc *FileController) Upload(c *gin.Context) {
    file, err := c.FormFile("file")
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error":   "File upload error",
            "details": err.Error(),
        })
        return
    }

    // Validate file size (5MB limit)
    if file.Size > 5<<20 {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "File size exceeds 5MB limit",
        })
        return
    }

    // Create upload directory
    uploadDir := "public/uploads"
    if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to create upload directory",
        })
        return
    }

    // Generate unique filename
    fileExt := filepath.Ext(file.Filename)
    newFilename := uuid.New().String() + fileExt
    filePath := filepath.Join(uploadDir, newFilename)

    // Save the file
    if err := c.SaveUploadedFile(file, filePath); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to save file",
        })
        return
    }

    // Return relative path
    c.JSON(http.StatusOK, gin.H{
        "message": "File uploaded successfully",
        "path":    "uploads/" + newFilename,
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
