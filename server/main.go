package main

import (
	"log"
	"net/http"
	"os"
	"re2no/auth"
	"re2no/database"
	"re2no/handlers"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Initialize database connection
	if err := database.Connect(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Run database migrations
	if err := database.Migrate(); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Initialize Notion OAuth
	auth.InitNotionOAuth()

	router := gin.Default()

	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})

	// Health check
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Re2no API is running",
		})
	})

	// Auth routes
	router.GET("/api/auth/notion/login", handlers.HandleNotionLogin)
	router.GET("/api/auth/notion/callback", handlers.HandleNotionCallback)
	router.GET("/api/auth/user", handlers.HandleGetUser)
	router.POST("/api/auth/logout", handlers.HandleLogout)

	// Reddit routes
	router.GET("/api/reddit", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"data":   "This will fetch Reddit posts soon.",
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Server running on http://localhost:%s", port)
	log.Printf("📝 Notion OAuth callback: %s", os.Getenv("NOTION_REDIRECT_URI"))

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
