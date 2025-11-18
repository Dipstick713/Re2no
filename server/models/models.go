package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	NotionUserID  string         `gorm:"uniqueIndex;not null" json:"notion_user_id"`
	WorkspaceID   string         `json:"workspace_id"`
	WorkspaceName string         `json:"workspace_name"`
	BotID         string         `json:"bot_id"`
	Email         string         `json:"email"`
	Name          string         `json:"name"`
	AvatarURL     string         `json:"avatar_url"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Sessions    []Session    `gorm:"foreignKey:UserID" json:"-"`
	RedditPosts []RedditPost `gorm:"foreignKey:UserID" json:"-"`
}

type Session struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	UserID       uint      `gorm:"not null;index" json:"user_id"`
	AccessToken  string    `gorm:"type:text;not null" json:"-"` // Encrypted
	TokenType    string    `json:"token_type"`
	ExpiresAt    time.Time `json:"expires_at"`
	RefreshToken string    `gorm:"type:text" json:"-"` // Encrypted
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	// Relations
	User User `gorm:"foreignKey:UserID" json:"-"`
}

type RedditPost struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	UserID        uint      `gorm:"not null;index" json:"user_id"`
	RedditID      string    `gorm:"uniqueIndex;not null" json:"reddit_id"`
	Subreddit     string    `gorm:"index" json:"subreddit"`
	Title         string    `json:"title"`
	Content       string    `gorm:"type:text" json:"content"`
	Author        string    `json:"author"`
	Score         int       `json:"score"`
	URL           string    `json:"url"`
	NotionPageID  string    `json:"notion_page_id"`  // ID of the Notion page created
	NotionPageURL string    `json:"notion_page_url"` // URL to open the Notion page
	SavedAt       time.Time `json:"saved_at"`
	CreatedAt     time.Time `json:"created_at"`

	// Relations
	User User `gorm:"foreignKey:UserID" json:"-"`
}

type OAuthState struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	State     string    `gorm:"uniqueIndex;not null" json:"state"`
	ExpiresAt time.Time `gorm:"index;not null" json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}
