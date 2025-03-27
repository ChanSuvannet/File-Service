package database

import (
	"fmt"
	"log"
	"my-project/models"

	"gorm.io/gorm"
)

// Migrate will perform the database migration
func Migrate(DB *gorm.DB) {
	// Auto migrate the File model (will create the table if it doesn't exist)
	if err := DB.AutoMigrate(&models.File{}); err != nil {
		log.Fatalf("Error migrating database: %v", err)
	}
	fmt.Println("Database migrated successfully")
}

// Seed will insert initial data into the database
func Seed(DB *gorm.DB) {
	// Insert a sample file record
	file := models.File{
		Filename:     "example.txt",
		OriginalName: "Example File",
		MimeType:     "text/plain",
		Path:         "/files/example.txt",
		Size:         1024,
	}

	// Check if data exists before seeding
	var count int64
	DB.Model(&models.File{}).Count(&count)
	if count == 0 {
		if err := DB.Create(&file).Error; err != nil {
			log.Fatalf("Error seeding data: %v", err)
		}
		fmt.Println("Database seeded with initial data")
	} else {
		fmt.Println("Database already has data, skipping seeding.")
	}
}
