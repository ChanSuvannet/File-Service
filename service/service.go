package service

import (
	"errors"
	"mime"
	"mime/multipart"
	"my-project/models"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ReadFile reads a file from the database or serves it from the public folder if not found in DB
func ReadFile(filename string, download bool, c *gin.Context) error {
	var file models.File

	// Try fetching file details from the database
	err := models.DB.Where("filename = ?", filename).First(&file).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Fallback: Try serving from public folder
			publicPath := filepath.Join("public/static", filename)

			if _, err := os.Stat(publicPath); os.IsNotExist(err) {
				return errors.New("file not found in database or public folder")
			}

			// Optionally detect MIME type based on extension
			mimeType := mime.TypeByExtension(filepath.Ext(publicPath))
			if mimeType == "" {
				mimeType = "application/octet-stream"
			}

			c.Header("Content-Type", mimeType)
			if download {
				c.Header("Content-Disposition", "attachment; filename="+filepath.Base(publicPath))
			} else {
				c.Header("Content-Disposition", "inline; filename="+filepath.Base(publicPath))
			}

			c.File(publicPath)
			return nil
		}
		return err
	}

	// If found in DB, ensure file exists
	if _, err := os.Stat(file.Path); os.IsNotExist(err) {
		return errors.New("file found in database but missing on disk")
	}

	c.Header("Content-Type", file.MimeType)
	if download {
		c.Header("Content-Disposition", "attachment; filename="+file.OriginalName)
	} else {
		c.Header("Content-Disposition", "inline; filename="+file.OriginalName)
	}

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
	uploadFolder := "public/uploads"
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
