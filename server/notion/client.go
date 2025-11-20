package notion

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/jomei/notionapi"
)

type NotionClient struct {
	client *notionapi.Client
}

type SavePostRequest struct {
	Title      string `json:"title" binding:"required"`
	Subreddit  string `json:"subreddit" binding:"required"`
	Content    string `json:"content"`
	Author     string `json:"author"`
	Score      int    `json:"score"`
	URL        string `json:"url" binding:"required"`
	RedditID   string `json:"reddit_id" binding:"required"`
	DatabaseID string `json:"database_id" binding:"required"`
}

type SavePostResponse struct {
	NotionPageID  string `json:"notion_page_id"`
	NotionPageURL string `json:"notion_page_url"`
}

// NewNotionClient creates a new Notion API client with the given access token
func NewNotionClient(accessToken string) *NotionClient {
	return &NotionClient{
		client: notionapi.NewClient(notionapi.Token(accessToken)),
	}
}

// SaveRedditPost saves a Reddit post to a Notion database (flexible properties)
func (nc *NotionClient) SaveRedditPost(req SavePostRequest) (*SavePostResponse, error) {
	log.Printf("[Notion] Saving Reddit post to Notion database: %s", req.DatabaseID)

	// Parse database ID
	dbID := notionapi.DatabaseID(req.DatabaseID)

	// First, retrieve the database to see what properties it has
	ctx := context.Background()
	database, err := nc.client.Database.Get(ctx, dbID)
	if err != nil {
		log.Printf("[Notion] Error fetching database schema: %v", err)
		return nil, fmt.Errorf("failed to fetch database: %w", err)
	}

	log.Printf("[Notion] Database schema loaded, building properties dynamically")

	// Build properties dynamically based on what exists in the database
	properties := nc.buildPropertiesFromSchema(database.Properties, req)

	// Create content blocks
	children := nc.createContentBlocks(req.Content, req.URL)

	// Create the page request
	createPageReq := &notionapi.PageCreateRequest{
		Parent: notionapi.Parent{
			Type:       notionapi.ParentTypeDatabaseID,
			DatabaseID: dbID,
		},
		Properties: properties,
		Children:   children,
	}

	// Create the page in Notion
	page, err := nc.client.Page.Create(ctx, createPageReq)
	if err != nil {
		log.Printf("[Notion] Error creating page: %v", err)
		return nil, fmt.Errorf("failed to create Notion page: %w", err)
	}

	log.Printf("[Notion] Successfully created page: %s", page.ID)

	return &SavePostResponse{
		NotionPageID:  string(page.ID),
		NotionPageURL: page.URL,
	}, nil
}

// buildPropertiesFromSchema builds properties dynamically based on the database schema
func (nc *NotionClient) buildPropertiesFromSchema(schema notionapi.PropertyConfigs, req SavePostRequest) notionapi.Properties {
	properties := notionapi.Properties{}
	now := notionapi.Date(time.Now())

	log.Printf("[Notion] Building properties from schema with %d properties", len(schema))

	// Iterate through database properties and match with our data
	for propName, propConfig := range schema {
		propType := propConfig.GetType()
		propNameLower := strings.ToLower(strings.ReplaceAll(propName, " ", "_"))

		log.Printf("[Notion] Processing property: %s (type: %s)", propName, propType)

		switch propType {
		case notionapi.PropertyConfigTypeTitle:
			// Always set the title property
			properties[propName] = notionapi.TitleProperty{
				Title: []notionapi.RichText{
					{Text: &notionapi.Text{Content: req.Title}},
				},
			}
			log.Printf("[Notion] Set title property: %s", propName)

		case notionapi.PropertyConfigTypeRichText:
			// Match rich text properties by name
			var textValue string
			switch {
			case strings.Contains(propNameLower, "subreddit"):
				textValue = req.Subreddit
			case strings.Contains(propNameLower, "author"):
				textValue = req.Author
			case strings.Contains(propNameLower, "reddit") && strings.Contains(propNameLower, "id"):
				textValue = req.RedditID
			case strings.Contains(propNameLower, "content"):
				textValue = req.Content
			}

			if textValue != "" {
				properties[propName] = notionapi.RichTextProperty{
					RichText: []notionapi.RichText{
						{Text: &notionapi.Text{Content: textValue}},
					},
				}
				log.Printf("[Notion] Set rich text property: %s = %s", propName, textValue[:min(50, len(textValue))])
			}

		case notionapi.PropertyConfigTypeNumber:
			// Match number properties
			if strings.Contains(propNameLower, "score") || strings.Contains(propNameLower, "upvote") {
				properties[propName] = notionapi.NumberProperty{
					Number: float64(req.Score),
				}
				log.Printf("[Notion] Set number property: %s = %d", propName, req.Score)
			}

		case notionapi.PropertyConfigTypeURL:
			// Match URL properties
			if strings.Contains(propNameLower, "url") || strings.Contains(propNameLower, "link") || strings.Contains(propNameLower, "reddit") {
				properties[propName] = notionapi.URLProperty{
					URL: req.URL,
				}
				log.Printf("[Notion] Set URL property: %s = %s", propName, req.URL)
			}

		case notionapi.PropertyConfigTypeDate:
			// Match date properties
			if strings.Contains(propNameLower, "saved") || strings.Contains(propNameLower, "created") || strings.Contains(propNameLower, "date") {
				properties[propName] = notionapi.DateProperty{
					Date: &notionapi.DateObject{Start: &now},
				}
				log.Printf("[Notion] Set date property: %s", propName)
			}
		}
	}

	log.Printf("[Notion] Built %d properties for page creation", len(properties))
	return properties
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// createContentBlocks creates Notion blocks from the Reddit post content
func (nc *NotionClient) createContentBlocks(content, url string) []notionapi.Block {
	blocks := []notionapi.Block{}

	// Add the Reddit link as a bookmark
	blocks = append(blocks, notionapi.BookmarkBlock{
		BasicBlock: notionapi.BasicBlock{
			Object: notionapi.ObjectTypeBlock,
			Type:   notionapi.BlockTypeBookmark,
		},
		Bookmark: notionapi.Bookmark{
			URL: url,
		},
	})

	// Add a divider
	blocks = append(blocks, notionapi.DividerBlock{
		BasicBlock: notionapi.BasicBlock{
			Object: notionapi.ObjectTypeBlock,
			Type:   notionapi.BlockTypeDivider,
		},
	})

	// If there's content, add it as paragraphs
	if content != "" {
		// Check if content is a URL (image/video link)
		if strings.HasPrefix(content, "http://") || strings.HasPrefix(content, "https://") {
			// Create a clickable link block for media URLs
			blocks = append(blocks, notionapi.ParagraphBlock{
				BasicBlock: notionapi.BasicBlock{
					Object: notionapi.ObjectTypeBlock,
					Type:   notionapi.BlockTypeParagraph,
				},
				Paragraph: notionapi.Paragraph{
					RichText: []notionapi.RichText{
						{
							Type: notionapi.ObjectTypeText,
							Text: &notionapi.Text{Content: "Media Link: "},
						},
						{
							Type: notionapi.ObjectTypeText,
							Text: &notionapi.Text{
								Content: content,
								Link:    &notionapi.Link{Url: content},
							},
							Annotations: &notionapi.Annotations{
								Color: notionapi.ColorBlue,
							},
						},
					},
				},
			})
		} else {
			// Split content into paragraphs (by double newline or single newline)
			paragraphs := strings.Split(content, "\n\n")
			if len(paragraphs) == 1 {
				paragraphs = strings.Split(content, "\n")
			}

			for _, para := range paragraphs {
				para = strings.TrimSpace(para)
				if para == "" {
					continue
				}

				// Notion has a 2000 character limit per rich text block
				// Split long paragraphs if needed
				if len(para) > 2000 {
					chunks := splitTextIntoChunks(para, 2000)
					for _, chunk := range chunks {
						blocks = append(blocks, nc.createParagraphBlock(chunk))
					}
				} else {
					blocks = append(blocks, nc.createParagraphBlock(para))
				}
			}
		}
	} else {
		// Add a placeholder if there's no content
		blocks = append(blocks, nc.createParagraphBlock("No content available for this post."))
	}

	return blocks
}

// createParagraphBlock creates a paragraph block from text
func (nc *NotionClient) createParagraphBlock(text string) notionapi.ParagraphBlock {
	return notionapi.ParagraphBlock{
		BasicBlock: notionapi.BasicBlock{
			Object: notionapi.ObjectTypeBlock,
			Type:   notionapi.BlockTypeParagraph,
		},
		Paragraph: notionapi.Paragraph{
			RichText: []notionapi.RichText{
				{
					Text: &notionapi.Text{
						Content: text,
					},
				},
			},
		},
	}
}

// splitTextIntoChunks splits text into chunks of specified size
func splitTextIntoChunks(text string, chunkSize int) []string {
	var chunks []string
	for len(text) > 0 {
		if len(text) <= chunkSize {
			chunks = append(chunks, text)
			break
		}
		// Try to split at a word boundary
		splitIndex := chunkSize
		for splitIndex > 0 && text[splitIndex] != ' ' && text[splitIndex] != '\n' {
			splitIndex--
		}
		if splitIndex == 0 {
			splitIndex = chunkSize
		}
		chunks = append(chunks, text[:splitIndex])
		text = text[splitIndex:]
		text = strings.TrimSpace(text)
	}
	return chunks
}

// GetDatabases retrieves all databases accessible to the integration
func (nc *NotionClient) GetDatabases() ([]notionapi.Database, error) {
	log.Printf("[Notion] Fetching accessible databases")

	// Search for databases
	searchReq := &notionapi.SearchRequest{
		Filter: notionapi.SearchFilter{
			Property: "object",
			Value:    "database",
		},
	}

	ctx := context.Background()
	searchResp, err := nc.client.Search.Do(ctx, searchReq)
	if err != nil {
		log.Printf("[Notion] Error searching databases: %v", err)
		return nil, fmt.Errorf("failed to search databases: %w", err)
	}

	databases := make([]notionapi.Database, 0)
	for _, result := range searchResp.Results {
		if result.GetObject() == "database" {
			// Type assert to database
			if db, ok := result.(*notionapi.Database); ok {
				databases = append(databases, *db)
			}
		}
	}

	log.Printf("[Notion] Found %d databases", len(databases))
	return databases, nil
}

// CreateRedditPostsDatabase creates a new database with the required schema for Reddit posts
func (nc *NotionClient) CreateRedditPostsDatabase(parentPageID string) (*notionapi.Database, error) {
	log.Printf("[Notion] Creating Reddit Posts database in page: %s", parentPageID)

	ctx := context.Background()

	createDBReq := &notionapi.DatabaseCreateRequest{
		Parent: notionapi.Parent{
			Type:   notionapi.ParentTypePageID,
			PageID: notionapi.PageID(parentPageID),
		},
		Title: []notionapi.RichText{
			{Text: &notionapi.Text{Content: "Reddit Posts"}},
		},
		Properties: notionapi.PropertyConfigs{
			"Title": notionapi.TitlePropertyConfig{
				Type: notionapi.PropertyConfigTypeTitle,
			},
			"Subreddit": notionapi.RichTextPropertyConfig{
				Type: notionapi.PropertyConfigTypeRichText,
			},
			"Author": notionapi.RichTextPropertyConfig{
				Type: notionapi.PropertyConfigTypeRichText,
			},
			"Score": notionapi.NumberPropertyConfig{
				Type: notionapi.PropertyConfigTypeNumber,
			},
			"Reddit URL": notionapi.URLPropertyConfig{
				Type: notionapi.PropertyConfigTypeURL,
			},
			"Reddit ID": notionapi.RichTextPropertyConfig{
				Type: notionapi.PropertyConfigTypeRichText,
			},
			"Saved At": notionapi.DatePropertyConfig{
				Type: notionapi.PropertyConfigTypeDate,
			},
		},
	}

	database, err := nc.client.Database.Create(ctx, createDBReq)
	if err != nil {
		log.Printf("[Notion] Error creating database: %v", err)
		return nil, fmt.Errorf("failed to create database: %w", err)
	}

	log.Printf("[Notion] Successfully created database: %s", database.ID)
	return database, nil
}

// DeletePage deletes (archives) a Notion page by its ID
func (nc *NotionClient) DeletePage(pageID string) error {
	log.Printf("[Notion] Deleting page: %s", pageID)

	ctx := context.Background()

	// Archive the page (Notion API doesn't have direct delete, uses archive)
	_, err := nc.client.Block.Delete(ctx, notionapi.BlockID(pageID))
	if err != nil {
		log.Printf("[Notion] Error deleting page: %v", err)
		return fmt.Errorf("failed to delete Notion page: %w", err)
	}

	log.Printf("[Notion] Successfully deleted page: %s", pageID)
	return nil
}
