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
		log.Printf("[Auth Middleware] %s %s", c.Request.Method, c.Request.URL.Path)
		log.Printf("[Auth Middleware] Origin: %s", c.Request.Header.Get("Origin"))

		// Get token from Authorization header (for cross-domain)
		authHeader := c.GetHeader("Authorization")
		log.Printf("[Auth Middleware] Authorization Header: %s", authHeader)
		var tokenString string

		if authHeader != "" && len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			tokenString = authHeader[7:]
			log.Printf("[Auth Middleware] ✓ Token from Authorization header (length: %d)", len(tokenString))
		} else {
			log.Println("[Auth Middleware] No Authorization header, checking cookie...")
			// Fallback: Try to get token from cookie (for same-domain/local dev)
			var err error
			tokenString, err = c.Cookie("auth_token")
			if err != nil {
				log.Printf("[Auth Middleware] ✗ No authentication found: %v", err)
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": "authentication required",
				})
				c.Abort()
				return
			}
			log.Printf("[Auth Middleware] ✓ Token from cookie (length: %d)", len(tokenString))
		}

		// Validate token
		log.Println("[Auth Middleware] Validating token...")
		claims, err := auth.ValidateToken(tokenString)
		if err != nil {
			log.Printf("[Auth Middleware] ✗ Token validation failed: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid or expired token",
			})
			c.Abort()
			return
		}
		log.Printf("[Auth Middleware] ✓ Token valid for user_id: %d", claims.UserID)

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
