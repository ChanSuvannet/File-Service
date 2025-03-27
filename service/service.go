package service

import (
	"errors"
	"fmt"
	"my-project/models"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ReadFile reads a file from the database
func ReadFile(filename string, download bool, c *gin.Context) error {
	var file models.File

	// Fetch the file details from the database using models.DB
	if err := models.DB.Where("filename = ?", filename).First(&file).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("invalid file name")
		}
		return err
	}

	// Check if the file exists in the specified path
	if _, err := os.Stat(file.Path); os.IsNotExist(err) {
		return errors.New("file not found")
	}

	// Serve the file as an attachment
	c.FileAttachment(file.Path, file.Filename)
	return nil
}

func UploadFile(filePath, originalName, mimeType string, size int64) (map[string]interface{}, error) {
	file := models.File{
		Filename:     filepath.Base(filePath),
		OriginalName: originalName,
		MimeType:     mimeType,
		Path:         filePath,
		Size:         size,
	}

	if err := models.DB.Create(&file).Error; err != nil {
		return nil, err
	}

	uri := fmt.Sprintf("api/file/%s", file.Filename)
	data := map[string]interface{}{
		"uri":          uri,
		"originalname": file.OriginalName,
		"mimetype":     file.MimeType,
		"size":         file.Size,
	}

	return map[string]interface{}{
		"data":    data,
		"message": "File has been uploaded successfully.",
	}, nil
}
