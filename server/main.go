package main

import (
	"log"
	"net/http"
	"os"
	"re2no/auth"
	"re2no/database"
	"re2no/handlers"
	"re2no/middleware"

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

	// Initialize JWT
	auth.InitJWT()

	// Initialize Notion OAuth
	auth.InitNotionOAuth()

	router := gin.Default()

	// CORS middleware - Allow credentials for authentication
	router.Use(func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if origin == "" {
			origin = "http://localhost:5173" // Default for development
		}

		c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
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

	// Auth routes (public)
	router.GET("/api/auth/notion/login", handlers.HandleNotionLogin)
	router.GET("/api/auth/notion/callback", handlers.HandleNotionCallback)

	// Auth routes (protected)
	authRoutes := router.Group("/api/auth")
	authRoutes.Use(middleware.RequireAuth())
	{
		authRoutes.GET("/user", handlers.HandleGetUser)
		authRoutes.POST("/logout", handlers.HandleLogout)
	}

	// Reddit routes (protected)
	redditRoutes := router.Group("/api/reddit")
	redditRoutes.Use(middleware.RequireAuth())
	{
		redditRoutes.GET("/posts", handlers.HandleFetchPosts)
	}

	// Notion routes (protected)
	notionRoutes := router.Group("/api/notion")
	notionRoutes.Use(middleware.RequireAuth())
	{
		notionRoutes.POST("/save", handlers.HandleSaveToNotion)
		notionRoutes.GET("/databases", handlers.HandleGetDatabases)
		notionRoutes.GET("/saved-posts", handlers.HandleGetSavedPosts)
		notionRoutes.POST("/create-database", handlers.HandleCreateRedditDatabase)
	}

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
