package main

import (
	"fmt"
	"log"
	"my-project/config"
	"my-project/database"
	"my-project/routes"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

// Load environment variables from .env file
func loadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}

// Global database variable
var DB *gorm.DB

// Connect to the database
func connectDatabase() {
	var err error
	DB, err = config.ConnectDB()
	if err != nil {
		log.Fatal("Error connecting to database: ", err)
	}
	database.Migrate(DB)

	fmt.Println("Database connected successfully")
}

// NotFound handler
func notFoundHandler(c *gin.Context) {
	c.JSON(404, gin.H{"error": fmt.Sprintf("Cannot %s %s", c.Request.Method, c.Request.URL)})
}

// ErrorHandler middleware
func errorHandler(c *gin.Context) {
	c.Next()
	if len(c.Errors) > 0 {
		c.JSON(422, gin.H{
			"error": c.Errors.Last().Error(),
		})
	}
}

func main() {
	// Load env variables
	loadEnv()

	// Optional: Initialize the database
	connectDatabase()

	// Create a new Gin router
	r := gin.Default()

	// Global error handler
	r.Use(errorHandler)

	// Serve static files from ./public at /public
	r.Static("/public", "./public")

	// Load HTML templates
	r.LoadHTMLFiles(filepath.Join("view", "index.html"))

	// Root route to render HTML
	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{
			"title": "My Page", // Optional dynamic data
		})
	})

	// Simple API endpoint
	r.GET("/api", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to the API",
		})
	})

	// Custom 404 handler
	r.NoRoute(notFoundHandler)
	r.Use(errorHandler)
	// Import other routes
	routes.SetupRoutes(r)

	// Run server on port from env or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
