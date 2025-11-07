package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"golang.org/x/oauth2"
)

var (
	NotionOAuthConfig *oauth2.Config
)

func InitNotionOAuth() {
	NotionOAuthConfig = &oauth2.Config{
		ClientID:     os.Getenv("NOTION_CLIENT_ID"),
		ClientSecret: os.Getenv("NOTION_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("NOTION_REDIRECT_URI"),
		Scopes:       []string{},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://api.notion.com/v1/oauth/authorize",
			TokenURL: "https://api.notion.com/v1/oauth/token",
		},
	}
}

type NotionUser struct {
	AccessToken          string                 `json:"access_token"`
	TokenType            string                 `json:"token_type"`
	BotID                string                 `json:"bot_id"`
	WorkspaceID          string                 `json:"workspace_id"`
	WorkspaceName        string                 `json:"workspace_name,omitempty"`
	WorkspaceIcon        string                 `json:"workspace_icon,omitempty"`
	Owner                map[string]interface{} `json:"owner"`
	DuplicatedTemplateID string                 `json:"duplicated_template_id,omitempty"`
}

func GetNotionUser(ctx context.Context, code string) (*NotionUser, error) {
	token, err := NotionOAuthConfig.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange token: %w", err)
	}

	// Get user info from Notion
	client := NotionOAuthConfig.Client(ctx, token)
	resp, err := client.Get("https://api.notion.com/v1/users/me")
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var user NotionUser
	user.AccessToken = token.AccessToken
	user.TokenType = token.TokenType

	if err := json.Unmarshal(body, &user); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user info: %w", err)
	}

	return &user, nil
}
