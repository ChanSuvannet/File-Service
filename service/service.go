package service

import (
	"errors"
	"mime/multipart"
	"my-project/models"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ReadFile reads a file from the database and serves it
func ReadFile(filename string, download bool, c *gin.Context) error {
	var file models.File

	// Fetch file details from the database using the filename
	if err := models.DB.Where("filename = ?", filename).First(&file).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("file not found in database")
		}
		return err
	}

	// Check if the file exists in the specified path
	if _, err := os.Stat(file.Path); os.IsNotExist(err) {
		return errors.New("file not found on disk")
	}

	// Set the correct content type (MIME type) for the response based on the file
	c.Header("Content-Type", file.MimeType)

	// Serve the file as an attachment if requested
	if download {
		c.Header("Content-Disposition", "attachment; filename="+file.OriginalName)
	} else {
		c.Header("Content-Disposition", "inline; filename="+file.OriginalName)
	}

	// Send the file to the client
	c.File(file.Path)
	return nil
}

// UpdateFile updates file information in the database and replaces the file if a new one is provided.
func UploadFile(file *multipart.FileHeader) (map[string]interface{}, error) {
	// Create upload folder if it doesn't exist
	uploadFolder := "public/uploads"
	if err := os.MkdirAll(uploadFolder, os.ModePerm); err != nil {
		return nil, errors.New("failed to create upload directory")
	}

	// Generate a unique file name without extension
	fileName := uuid.New().String() // No file extension added here
	filePath := filepath.Join(uploadFolder, fileName)

	// Open the uploaded file
	src, err := file.Open()
	if err != nil {
		return nil, errors.New("failed to open uploaded file")
	}
	defer src.Close()

	// Create destination file
	dst, err := os.Create(filePath)
	if err != nil {
		return nil, errors.New("failed to create destination file")
	}
	defer dst.Close()

	// Copy file content
	if _, err := dst.ReadFrom(src); err != nil {
		return nil, errors.New("failed to save file")
	}

	// Save metadata in the database, excluding the extension
	fileRecord := models.File{
		Filename:     fileName, 
		OriginalName: file.Filename,
		MimeType:     file.Header.Get("Content-Type"),
		Path:         filePath,
		Size:         file.Size,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := models.DB.Create(&fileRecord).Error; err != nil {
		return nil, err
	}

	// Prepare the response with detailed metadata
	response := map[string]interface{}{
		"uri":          fileName,
		"originalname": file.Filename,
		"mimetype":     file.Header.Get("Content-Type"),
		"size":         file.Size,
	}

	return response, nil
}

// UploadProductImage handles saving an image specifically for products
func UploadProductImage(file *multipart.FileHeader) (map[string]interface{}, error) {
	uploadFolder := "public/uploads/products"
	if err := os.MkdirAll(uploadFolder, os.ModePerm); err != nil {
		return nil, errors.New("failed to create upload directory")
	}

	// Generate a unique file name without extension
	fileName := uuid.New().String() // Include extension
	filePath := filepath.Join(uploadFolder, fileName)

	// Open the uploaded file
	src, err := file.Open()
	if err != nil {
		return nil, errors.New("failed to open uploaded file")
	}
	defer src.Close()

	// Create destination file
	dst, err := os.Create(filePath)
	if err != nil {
		return nil, errors.New("failed to create destination file")
	}
	defer dst.Close()

	// Copy file content
	if _, err := dst.ReadFrom(src); err != nil {
		return nil, errors.New("failed to save file")
	}

	// Save metadata in the database, excluding the extension
	fileRecord := models.FileProduct{
		Filename:     fileName,
		OriginalName: file.Filename,
		MimeType:     file.Header.Get("Content-Type"),
		Path:         filePath,
		Size:         file.Size,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := models.DB.Create(&fileRecord).Error; err != nil {
		return nil, err
	}

	// Prepare the response with detailed metadata
	response := map[string]interface{}{
		"uri":          fileName,
		"originalname": file.Filename,
		"mimetype":     file.Header.Get("Content-Type"),
		"size":         file.Size,
	}

	return response, nil
}


