package middleware

import (
	"log"
	"net/http"
	"re2no/auth"
	"re2no/database"
	"re2no/models"

	"github.com/gin-gonic/gin"
)

// RequireAuth is a middleware that validates JWT tokens
func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from cookie
		tokenString, err := c.Cookie("auth_token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "authentication required",
			})
			c.Abort()
			return
		}

		// Validate token
		claims, err := auth.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid or expired token",
			})
			c.Abort()
			return
		}

		// Fetch user from database
		var user models.User
		if err := database.DB.Where("id = ?", claims.UserID).First(&user).Error; err != nil {
			log.Printf("[Auth Middleware] Failed to fetch user: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "user not found",
			})
			c.Abort()
			return
		}

		// Set user info in context for handlers to use
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user", &user)

		c.Next()
	}
}
