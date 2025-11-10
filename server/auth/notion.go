package auth

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
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
	// Notion requires manual token exchange because their OAuth doesn't follow standard
	tokenURL := "https://api.notion.com/v1/oauth/token"

	// Create the request body
	authHeader := NotionOAuthConfig.ClientID + ":" + NotionOAuthConfig.ClientSecret
	encodedAuth := "Basic " + encodeBase64(authHeader)

	payload := map[string]string{
		"grant_type":   "authorization_code",
		"code":         code,
		"redirect_uri": NotionOAuthConfig.RedirectURL,
	}

	jsonPayload, _ := json.Marshal(payload)

	req, err := http.NewRequestWithContext(ctx, "POST", tokenURL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, fmt.Errorf("failed to create token request: %w", err)
	}

	req.Header.Set("Authorization", encodedAuth)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange token: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	log.Printf("Notion OAuth Response Status: %d", resp.StatusCode)
	log.Printf("Notion OAuth Response Body: %s", string(body))

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("notion API returned status %d: %s", resp.StatusCode, string(body))
	}

	var user NotionUser
	if err := json.Unmarshal(body, &user); err != nil {
		return nil, fmt.Errorf("failed to unmarshal token response: %w", err)
	}

	return &user, nil
}

func encodeBase64(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}
