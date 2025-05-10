package config

import (
	"fmt"
	"my-project/models"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DatabaseConfig holds the database connection details
type DatabaseConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
	Dialect  string
}

// LoadConfig initializes database configuration from environment variables
func LoadConfig() DatabaseConfig {
	return DatabaseConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Username: os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		Database: os.Getenv("DB_DATABASE"),
		Dialect:  os.Getenv("DB_CONNECTION"),
	}
}

// ConnectDB establishes a database connection using GORM
func ConnectDB() (*gorm.DB, error) {
	config := LoadConfig()
	var dialector gorm.Dialector

	switch config.Dialect {
	case "mysql":
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			config.Username, config.Password, config.Host, config.Port, config.Database)
		dialector = mysql.Open(dsn)
	case "postgres":
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
			config.Host, config.Username, config.Password, config.Database, config.Port)
		dialector = postgres.Open(dsn)
	default:
		return nil, fmt.Errorf("unsupported database dialect: %s", config.Dialect)
	}

	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Initialize the DB global variable in models
	models.DB = db

	if err := db.AutoMigrate(&models.FileProduct{}); err != nil {
		return nil, fmt.Errorf("auto migration failed: %w", err)
	}

	return db, nil
}
