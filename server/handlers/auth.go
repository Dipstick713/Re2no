package handlers

import (
	"net/http"
	"re2no/auth"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
)

// Store for OAuth state tokens (in production, use Redis or similar)
var stateStore = make(map[string]bool)

// HandleNotionLogin initiates the Notion OAuth flow
func HandleNotionLogin(c *gin.Context) {
	state := uuid.New().String()
	stateStore[state] = true

	url := auth.NotionOAuthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
	c.JSON(http.StatusOK, gin.H{
		"url": url,
	})
}

// HandleNotionCallback handles the OAuth callback from Notion
func HandleNotionCallback(c *gin.Context) {
	state := c.Query("state")
	code := c.Query("code")
	errorParam := c.Query("error")

	// Check for OAuth errors
	if errorParam != "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errorParam,
		})
		return
	}

	// Validate state
	if !stateStore[state] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid state parameter",
		})
		return
	}
	delete(stateStore, state)

	// Exchange code for token and get user info
	user, err := auth.GetNotionUser(c.Request.Context(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Here you would typically:
	// 1. Store user info in your database
	// 2. Create a session
	// 3. Generate a JWT token
	// 4. Set a cookie

	// Log the successful authentication
	_ = user // TODO: Store user info when database is implemented

	// For now, we'll redirect to frontend with success
	frontendURL := "http://localhost:5173/dashboard?auth=success"
	c.Redirect(http.StatusTemporaryRedirect, frontendURL)
}

// HandleGetUser returns the current authenticated user
func HandleGetUser(c *gin.Context) {
	// This will be implemented later with JWT/session validation
	c.JSON(http.StatusOK, gin.H{
		"message": "User endpoint - to be implemented with session management",
	})
}

// HandleLogout handles user logout
func HandleLogout(c *gin.Context) {
	// This will be implemented later with JWT/session invalidation
	c.JSON(http.StatusOK, gin.H{
		"message": "Logged out successfully",
	})
}
