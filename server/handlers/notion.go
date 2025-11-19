package handlers

import (
	"log"
	"net/http"
	"re2no/database"
	"re2no/models"
	"re2no/notion"

	"github.com/gin-gonic/gin"
)

// HandleSaveToNotion saves a Reddit post to the user's Notion workspace
func HandleSaveToNotion(c *gin.Context) {
	log.Println("[Notion Handler] Received save to Notion request")

	// Get user from context (set by auth middleware)
	userInterface, exists := c.Get("user")
	if !exists {
		log.Println("[Notion Handler] User not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	user, ok := userInterface.(*models.User)
	if !ok {
		log.Println("[Notion Handler] Invalid user type in context")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	log.Printf("[Notion Handler] Processing request for user: %s", user.NotionUserID)

	// Parse request body
	var req notion.SavePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[Notion Handler] Invalid request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	log.Printf("[Notion Handler] Saving post: %s to database: %s", req.Title, req.DatabaseID)

	// Get user's latest session for access token
	var session models.Session
	if err := database.DB.Where("user_id = ?", user.ID).Order("expires_at DESC").First(&session).Error; err != nil {
		log.Printf("[Notion Handler] Failed to get user session: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No valid session found. Please login again."})
		return
	}

	// Check if session is expired
	// if time.Now().After(session.ExpiresAt) {
	// 	log.Println("[Notion Handler] Session expired")
	// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "Session expired. Please login again."})
	// 	return
	// }

	// Create Notion client with user's access token
	notionClient := notion.NewNotionClient(session.AccessToken)

	// Save post to Notion
	response, err := notionClient.SaveRedditPost(req)
	if err != nil {
		log.Printf("[Notion Handler] Failed to save post to Notion: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save post to Notion", "details": err.Error()})
		return
	}

	log.Printf("[Notion Handler] Successfully saved post. Page ID: %s", response.NotionPageID)

	// Save to database with all post data including Notion URL
	redditPost := models.RedditPost{
		UserID:        user.ID,
		RedditID:      req.RedditID,
		Subreddit:     req.Subreddit,
		Title:         req.Title,
		Content:       req.Content,
		Author:        req.Author,
		Score:         req.Score,
		URL:           req.URL,
		NotionPageID:  response.NotionPageID,
		NotionPageURL: response.NotionPageURL,
	}

	if err := database.DB.Create(&redditPost).Error; err != nil {
		log.Printf("[Notion Handler] Failed to save post to database: %v", err)
		// Don't fail the request if DB save fails - post is already in Notion
	} else {
		log.Printf("[Notion Handler] Successfully saved to database with URL: %s", response.NotionPageURL)
	}

	c.JSON(http.StatusOK, gin.H{
		"success":         true,
		"notion_page_id":  response.NotionPageID,
		"notion_page_url": response.NotionPageURL,
		"message":         "Post saved to Notion successfully",
	})
}

// HandleGetDatabases retrieves all databases accessible to the user
func HandleGetDatabases(c *gin.Context) {
	log.Println("[Notion Handler] Received get databases request")

	// Get user from context
	userInterface, exists := c.Get("user")
	if !exists {
		log.Println("[Notion Handler] User not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	user, ok := userInterface.(*models.User)
	if !ok {
		log.Println("[Notion Handler] Invalid user type in context")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Get user's latest session
	var session models.Session
	if err := database.DB.Where("user_id = ?", user.ID).Order("expires_at DESC").First(&session).Error; err != nil {
		log.Printf("[Notion Handler] Failed to get user session: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No valid session found. Please login again."})
		return
	}

	// Create Notion client
	notionClient := notion.NewNotionClient(session.AccessToken)

	// Get databases
	databases, err := notionClient.GetDatabases()
	if err != nil {
		log.Printf("[Notion Handler] Failed to get databases: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve databases", "details": err.Error()})
		return
	}

	log.Printf("[Notion Handler] Retrieved %d databases", len(databases))

	// Format response
	dbList := make([]gin.H, len(databases))
	for i, db := range databases {
		title := "Untitled"
		if len(db.Title) > 0 && db.Title[0].PlainText != "" {
			title = db.Title[0].PlainText
		}

		dbList[i] = gin.H{
			"id":    db.ID,
			"title": title,
			"url":   db.URL,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"databases": dbList,
	})
}

// HandleGetSavedPosts retrieves all posts saved by the user
func HandleGetSavedPosts(c *gin.Context) {
	log.Println("[Notion Handler] Received get saved posts request")

	// Get user from context
	userInterface, exists := c.Get("user")
	if !exists {
		log.Println("[Notion Handler] User not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	user, ok := userInterface.(*models.User)
	if !ok {
		log.Println("[Notion Handler] Invalid user type in context")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Get saved posts from database
	var posts []models.RedditPost
	if err := database.DB.Where("user_id = ?", user.ID).Order("created_at DESC").Find(&posts).Error; err != nil {
		log.Printf("[Notion Handler] Failed to get saved posts: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve saved posts"})
		return
	}

	log.Printf("[Notion Handler] Found %d saved posts for user", len(posts))

	c.JSON(http.StatusOK, gin.H{
		"posts": posts,
	})
}

// HandleDeleteSavedPost deletes a saved post from the database and Notion
func HandleDeleteSavedPost(c *gin.Context) {
	log.Println("[Notion Handler] Received delete saved post request")

	// Get user from context
	userInterface, exists := c.Get("user")
	if !exists {
		log.Println("[Notion Handler] User not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	user, ok := userInterface.(*models.User)
	if !ok {
		log.Println("[Notion Handler] Invalid user type in context")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Get Reddit ID from URL parameter
	redditID := c.Param("reddit_id")
	if redditID == "" {
		log.Println("[Notion Handler] Reddit ID not provided")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Reddit ID is required"})
		return
	}

	log.Printf("[Notion Handler] Deleting post with Reddit ID: %s for user: %s", redditID, user.NotionUserID)

	// Find the post in database
	var post models.RedditPost
	if err := database.DB.Where("user_id = ? AND reddit_id = ?", user.ID, redditID).First(&post).Error; err != nil {
		log.Printf("[Notion Handler] Post not found: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	// Get user's session for Notion access token
	var session models.Session
	if err := database.DB.Where("user_id = ?", user.ID).Order("expires_at DESC").First(&session).Error; err != nil {
		log.Printf("[Notion Handler] Failed to get user session: %v", err)
		// Continue with database deletion even if session not found
	} else if post.NotionPageID != "" {
		// Delete from Notion if we have both page ID and session
		notionClient := notion.NewNotionClient(session.AccessToken)
		if err := notionClient.DeletePage(post.NotionPageID); err != nil {
			log.Printf("[Notion Handler] Warning: Failed to delete from Notion (continuing with database deletion): %v", err)
			// Continue with database deletion even if Notion deletion fails
		} else {
			log.Printf("[Notion Handler] Successfully deleted from Notion")
		}
	}

	// Delete from database
	result := database.DB.Delete(&post)
	if result.Error != nil {
		log.Printf("[Notion Handler] Failed to delete post from database: %v", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete post"})
		return
	}

	if result.RowsAffected == 0 {
		log.Printf("[Notion Handler] Post not found or not owned by user")
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	log.Printf("[Notion Handler] Successfully deleted post")

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Post deleted successfully",
	})
}

// HandleCreateRedditDatabase creates a template database for Reddit posts
func HandleCreateRedditDatabase(c *gin.Context) {
	log.Println("[Notion Handler] Received create database request")

	// Get user from context
	userInterface, exists := c.Get("user")
	if !exists {
		log.Println("[Notion Handler] User not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	user, ok := userInterface.(*models.User)
	if !ok {
		log.Println("[Notion Handler] Invalid user type in context")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Parse request - need parent page ID
	var req struct {
		ParentPageID string `json:"parent_page_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[Notion Handler] Invalid request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "parent_page_id is required"})
		return
	}

	// Get user's latest session
	var session models.Session
	if err := database.DB.Where("user_id = ?", user.ID).Order("expires_at DESC").First(&session).Error; err != nil {
		log.Printf("[Notion Handler] Failed to get user session: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No valid session found. Please login again."})
		return
	}

	// Create Notion client
	notionClient := notion.NewNotionClient(session.AccessToken)

	// Create the database
	database, err := notionClient.CreateRedditPostsDatabase(req.ParentPageID)
	if err != nil {
		log.Printf("[Notion Handler] Failed to create database: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create database", "details": err.Error()})
		return
	}

	log.Printf("[Notion Handler] Successfully created database: %s", database.ID)

	c.JSON(http.StatusOK, gin.H{
		"success":      true,
		"database_id":  database.ID,
		"database_url": database.URL,
		"message":      "Reddit Posts database created successfully",
	})
}
