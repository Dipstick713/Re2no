package reddit

import (
"encoding/json"
"fmt"
"io"
"net/http"
"net/url"
"time"
)

// RedditClient handles Reddit API requests
type RedditClient struct {
	HTTPClient *http.Client
	UserAgent  string
}

// NewRedditClient creates a new Reddit API client
func NewRedditClient() *RedditClient {
	return &RedditClient{
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		UserAgent: "Re2no:v1.0.0 (by /u/your_username)",
	}
}

// RedditPost represents a single Reddit post
type RedditPost struct {
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	Author      string  `json:"author"`
	Subreddit   string  `json:"subreddit"`
	Score       int     `json:"score"`
	URL         string  `json:"url"`
	Permalink   string  `json:"permalink"`
	CreatedUTC  float64 `json:"created_utc"`
	NumComments int     `json:"num_comments"`
	Thumbnail   string  `json:"thumbnail"`
	SelfText    string  `json:"selftext"`
	IsVideo     bool    `json:"is_video"`
}

// RedditResponse represents the Reddit API response structure
type RedditResponse struct {
	Data struct {
		Children []struct {
			Data RedditPost `json:"data"`
		} `json:"children"`
		After  string `json:"after"`
		Before string `json:"before"`
	} `json:"data"`
}

// FetchPostsParams holds the parameters for fetching posts
type FetchPostsParams struct {
	Subreddit string
	Sort      string
	TimeRange string
	Limit     int
	After     string
}

// FetchPosts fetches posts from Reddit based on the given parameters
func (c *RedditClient) FetchPosts(params FetchPostsParams) ([]RedditPost, error) {
	if params.Subreddit == "" {
		params.Subreddit = "all"
	}
	if params.Sort == "" {
		params.Sort = "hot"
	}
	if params.Limit <= 0 || params.Limit > 100 {
		params.Limit = 25
	}

	baseURL := fmt.Sprintf("https://www.reddit.com/r/%s/%s.json", params.Subreddit, params.Sort)
	
	urlParams := url.Values{}
	urlParams.Add("limit", fmt.Sprintf("%d", params.Limit))
	
	if params.TimeRange != "" && (params.Sort == "top" || params.Sort == "controversial") {
		urlParams.Add("t", params.TimeRange)
	}
	
	if params.After != "" {
		urlParams.Add("after", params.After)
	}

	fullURL := baseURL + "?" + urlParams.Encode()

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", c.UserAgent)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch posts: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("reddit API returned status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var redditResp RedditResponse
	if err := json.Unmarshal(body, &redditResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	posts := make([]RedditPost, 0, len(redditResp.Data.Children))
	for _, child := range redditResp.Data.Children {
		posts = append(posts, child.Data)
	}

	return posts, nil
}

// SearchPosts searches for posts containing a keyword
func (c *RedditClient) SearchPosts(subreddit, keyword string, sort string, limit int) ([]RedditPost, error) {
	if subreddit == "" {
		subreddit = "all"
	}
	if sort == "" {
		sort = "relevance"
	}
	if limit <= 0 || limit > 100 {
		limit = 25
	}

	baseURL := fmt.Sprintf("https://www.reddit.com/r/%s/search.json", subreddit)
	
	urlParams := url.Values{}
	urlParams.Add("q", keyword)
	urlParams.Add("restrict_sr", "true")
	urlParams.Add("sort", sort)
	urlParams.Add("limit", fmt.Sprintf("%d", limit))

	fullURL := baseURL + "?" + urlParams.Encode()

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", c.UserAgent)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to search posts: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("reddit API returned status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var redditResp RedditResponse
	if err := json.Unmarshal(body, &redditResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	posts := make([]RedditPost, 0, len(redditResp.Data.Children))
	for _, child := range redditResp.Data.Children {
		posts = append(posts, child.Data)
	}

	return posts, nil
}
