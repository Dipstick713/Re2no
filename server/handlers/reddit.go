package handlers

import (
"log"
"net/http"
"re2no/reddit"
"strconv"
"strings"

"github.com/gin-gonic/gin"
)

var redditClient = reddit.NewRedditClient()

// HandleFetchPosts fetches Reddit posts based on query parameters
func HandleFetchPosts(c *gin.Context) {
	log.Println("=== Fetching Reddit Posts ===")

	// Get query parameters
	subreddits := c.Query("subreddits")      // Comma-separated list
	keyword := c.Query("keyword")
	sortBy := c.Query("sort")
	dateRange := c.Query("date_range")
	limitStr := c.Query("limit")

	// Parse limit
	limit := 25
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil {
			limit = parsedLimit
		}
	}

	// Default to hot if not specified
	if sortBy == "" {
		sortBy = "hot"
	}

	log.Printf("Query params - Subreddits: %s, Keyword: %s, Sort: %s, DateRange: %s, Limit: %d",
subreddits, keyword, sortBy, dateRange, limit)

	// Split subreddits
	subredditList := []string{"all"}
	if subreddits != "" {
		subredditList = strings.Split(subreddits, ",")
		// Clean up subreddit names
		for i, sub := range subredditList {
			sub = strings.TrimSpace(sub)
			sub = strings.TrimPrefix(sub, "r/")
			subredditList[i] = sub
		}
	}

	// Fetch posts from each subreddit
	allPosts := []reddit.RedditPost{}

	for _, subreddit := range subredditList {
		log.Printf("Fetching from r/%s...", subreddit)

		var posts []reddit.RedditPost
		var err error

		if keyword != "" {
			// Search with keyword
			posts, err = redditClient.SearchPosts(subreddit, keyword, sortBy, limit)
		} else {
			// Fetch without keyword
			posts, err = redditClient.FetchPosts(reddit.FetchPostsParams{
Subreddit: subreddit,
Sort:      sortBy,
TimeRange: dateRange,
Limit:     limit,
})
		}

		if err != nil {
			log.Printf("ERROR: Failed to fetch from r/%s: %v", subreddit, err)
			continue
		}

		log.Printf("Fetched %d posts from r/%s", len(posts), subreddit)
		allPosts = append(allPosts, posts...)
	}

	log.Printf("=== Total Posts Fetched: %d ===", len(allPosts))

	c.JSON(http.StatusOK, gin.H{
"posts": allPosts,
"count": len(allPosts),
})
}
