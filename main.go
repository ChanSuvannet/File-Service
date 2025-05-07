package main

import (
	"fmt"
	"log"
	"my-project/config"
	"my-project/database"
	"my-project/routes"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

// Load environment variables from .env file
func loadEnv() {
	// Only load .env file if not running inside Docker
	if os.Getenv("ENV") != "production" {
		if err := godotenv.Load(); err != nil {
			log.Println("⚠️ Could not load .env file, using system environment")
		}
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

func staticFileHandler(c *gin.Context) {
    path := c.Param("filepath")
    
    // If path is an external URL, return it directly
    if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
        c.JSON(http.StatusOK, gin.H{
            "url": path,
            "isExternal": true,
        })
        return
    }

    // Security: Prevent directory traversal
    if strings.Contains(path, "../") {
        c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
            "error": "Directory traversal not allowed",
        })
        return
    }

    fullPath := filepath.Join("public", path)
    
    if _, err := os.Stat(fullPath); os.IsNotExist(err) {
        c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
            "error": "File not found",
        })
        return
    }

    c.File(fullPath)
}

func main() {
	// Load env variables
	loadEnv()

	// Optional: Initialize the database
	connectDatabase()
	

	// Create a new Gin router
	r := gin.Default()

	// Add CORS middleware
    r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3001")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	})

	// Global error handler
	r.Use(errorHandler)

	// Load HTML templates
	r.LoadHTMLFiles(filepath.Join("view", "index.html"))

	// Root route to render HTML
	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{
			"title": "File Upload",
		})
	})

	// Group API routes under /api
	api := r.Group("/api", errorHandler)
	{
		routes.SetupRoutes(api)
	}

	r.Static("/public", "./public")

	// Custom 404 handler
	r.NoRoute(notFoundHandler)

	// Run server on port from env or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("🚀 Application running on: http://localhost:%s\n", port)
	r.Run(":" + port)
}