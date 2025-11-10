package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"re2no/auth"
	"re2no/database"
	"re2no/models"
	"time"

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
	log.Println("=== OAuth Callback Started ===")

	state := c.Query("state")
	code := c.Query("code")
	errorParam := c.Query("error")

	log.Printf("State: %s, Code: %s", state, code[:20]+"...")

	// Check for OAuth errors
	if errorParam != "" {
		log.Printf("ERROR: OAuth error parameter: %s", errorParam)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errorParam,
		})
		return
	}

	// Validate state
	if !stateStore[state] {
		log.Printf("ERROR: Invalid state parameter")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid state parameter",
		})
		return
	}
	delete(stateStore, state)
	log.Println("State validated successfully")

	// Exchange code for token and get user info
	log.Println("Exchanging code for token...")
	notionUser, err := auth.GetNotionUser(c.Request.Context(), code)
	if err != nil {
		log.Printf("ERROR: Failed to exchange token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to exchange token: " + err.Error(),
		})
		return
	}
	log.Println("Token exchange successful")

	// Debug: Log the entire response
	log.Printf("BotID: %s, WorkspaceID: %s, WorkspaceName: %s",
		notionUser.BotID, notionUser.WorkspaceID, notionUser.WorkspaceName)

	ownerJSON, _ := json.MarshalIndent(notionUser.Owner, "", "  ")
	log.Printf("Notion Owner Object:\n%s", string(ownerJSON))

	// Extract user info from Owner map
	var notionUserID, userName, avatarURL, email string

	// The Owner object structure is: { "workspace": true } or { "type": "user", "user": {...} }
	// Let's check both possibilities
	if userObj, ok := notionUser.Owner["user"].(map[string]interface{}); ok {
		// Case 1: Owner has a "user" key
		if id, ok := userObj["id"].(string); ok {
			notionUserID = id
		}
		if name, ok := userObj["name"].(string); ok {
			userName = name
		}
		if avatar, ok := userObj["avatar_url"].(string); ok {
			avatarURL = avatar
		}
		if person, ok := userObj["person"].(map[string]interface{}); ok {
			if personEmail, ok := person["email"].(string); ok {
				email = personEmail
			}
		}
	} else if workspace, ok := notionUser.Owner["workspace"].(bool); ok && workspace {
		// Case 2: Owner is workspace - use workspace info
		log.Printf("OAuth granted by workspace, using bot_id as identifier")
		notionUserID = notionUser.BotID // Use bot_id as unique identifier
		userName = notionUser.WorkspaceName
	}

	// If still empty, log the full response for debugging
	if notionUserID == "" {
		log.Printf("ERROR: Could not extract user ID from Notion response")
		log.Printf("Full NotionUser: AccessToken=%s, BotID=%s, WorkspaceID=%s, WorkspaceName=%s",
			notionUser.AccessToken[:20]+"...", notionUser.BotID, notionUser.WorkspaceID, notionUser.WorkspaceName)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to extract user information from Notion response",
		})
		return
	}

	log.Printf("Extracted user info - ID: %s, Name: %s, Email: %s", notionUserID, userName, email)

	// Find or create user in database
	log.Println("Looking up user in database...")
	var user models.User
	result := database.DB.Where("notion_user_id = ?", notionUserID).First(&user)

	if result.Error != nil {
		log.Println("User not found, creating new user...")
		// Create new user
		user = models.User{
			NotionUserID:  notionUserID,
			WorkspaceID:   notionUser.WorkspaceID,
			WorkspaceName: notionUser.WorkspaceName,
			BotID:         notionUser.BotID,
			Name:          userName,
			AvatarURL:     avatarURL,
			Email:         email,
		}

		if err := database.DB.Create(&user).Error; err != nil {
			log.Printf("ERROR: Failed to create user: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to create user: " + err.Error(),
			})
			return
		}
		log.Printf("User created successfully with ID: %d", user.ID)
	} else {
		log.Printf("User found with ID: %d", user.ID)
	}

	// Create or update session
	log.Println("Managing session...")
	var session models.Session
	sessionResult := database.DB.Where("user_id = ?", user.ID).First(&session)

	expiresAt := time.Now().Add(30 * 24 * time.Hour) // 30 days

	if sessionResult.Error != nil {
		log.Println("Session not found, creating new session...")
		// Create new session
		session = models.Session{
			UserID:      user.ID,
			AccessToken: notionUser.AccessToken,
			TokenType:   "Bearer",
			ExpiresAt:   expiresAt,
		}

		if err := database.DB.Create(&session).Error; err != nil {
			log.Printf("ERROR: Failed to create session: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to create session: " + err.Error(),
			})
			return
		}
		log.Println("Session created successfully")
	} else {
		log.Println("Session found, updating...")
		// Update existing session
		session.AccessToken = notionUser.AccessToken
		session.ExpiresAt = expiresAt
		database.DB.Save(&session)
		log.Println("Session updated successfully")
	}

	// Generate JWT token
	log.Println("Generating JWT token...")
	token, err := auth.GenerateToken(user.ID, user.Email)
	if err != nil {
		log.Printf("ERROR: Failed to generate token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to generate token: " + err.Error(),
		})
		return
	}
	log.Println("JWT token generated successfully")

	// Set HTTP-only cookie
	c.SetCookie(
		"auth_token",
		token,
		int(24*time.Hour.Seconds()), // 24 hours
		"/",
		"",
		false, // Set to true in production with HTTPS
		true,  // HTTP-only
	)
	log.Println("Auth cookie set")

	// Redirect to frontend dashboard
	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:5173"
	}

	log.Printf("Redirecting to: %s/dashboard?auth=success", frontendURL)
	log.Println("=== OAuth Callback Completed Successfully ===")

	c.Redirect(http.StatusTemporaryRedirect, frontendURL+"/dashboard?auth=success")
	c.Redirect(http.StatusTemporaryRedirect, frontendURL+"/dashboard?auth=success")
}

// HandleGetUser returns the current authenticated user
func HandleGetUser(c *gin.Context) {
	// Get token from cookie
	tokenString, err := c.Cookie("auth_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "not authenticated",
		})
		return
	}

	// Validate token
	claims, err := auth.ValidateToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid token",
		})
		return
	}

	// Get user from database
	var user models.User
	if err := database.DB.First(&user, claims.UserID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "user not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

// HandleLogout handles user logout
func HandleLogout(c *gin.Context) {
	// Get token from cookie
	tokenString, err := c.Cookie("auth_token")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "already logged out",
		})
		return
	}

	// Validate token to get user ID
	claims, err := auth.ValidateToken(tokenString)
	if err == nil {
		// Delete session from database
		database.DB.Where("user_id = ?", claims.UserID).Delete(&models.Session{})
	}

	// Clear cookie
	c.SetCookie(
		"auth_token",
		"",
		-1,
		"/",
		"",
		false,
		true,
	)

	c.JSON(http.StatusOK, gin.H{
		"message": "logged out successfully",
	})
}
