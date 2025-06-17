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

var DB *gorm.DB

func loadEnv() {
	if os.Getenv("ENV") != "production" {
		if err := godotenv.Load(); err != nil {
			log.Println("⚠️ Could not load .env file, using system environment")
		}
	}
}

func connectDatabase() {
	var err error
	DB, err = config.ConnectDB()
	if err != nil {
		log.Fatal("❌ Error connecting to database:", err)
	}
	database.Migrate(DB)
	fmt.Println("✅ Database connected successfully")
}

func notFoundHandler(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Cannot %s %s", c.Request.Method, c.Request.URL)})
}

func errorHandler(c *gin.Context) {
	c.Next()
	if len(c.Errors) > 0 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": c.Errors.Last().Error()})
	}
}

func staticFileHandler(c *gin.Context) {
	path := c.Param("filepath")

	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		c.JSON(http.StatusOK, gin.H{"url": path, "isExternal": true})
		return
	}

	if strings.Contains(path, "../") {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Directory traversal not allowed"})
		return
	}

	fullPath := filepath.Join("public", path)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	c.File(fullPath)
}

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(corsMiddleware())
	r.Use(errorHandler)

	r.LoadHTMLFiles(filepath.Join("view", "index.html"))

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{"title": "File Upload"})
	})

	api := r.Group("/api", errorHandler)
	routes.SetupRoutes(api)

	r.Static("/public", "./public")
	r.NoRoute(notFoundHandler)

	return r
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3001")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

func main() {
	loadEnv()
	connectDatabase()

	r := setupRouter()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("\n🚀 Application running on: http://localhost:%s\n", port)
	r.Run(":" + port)
}
