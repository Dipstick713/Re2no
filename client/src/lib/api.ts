const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080'

// Token management for cross-domain auth
export function getAuthToken(): string | null {
  return localStorage.getItem('auth_token')
}

export function setAuthToken(token: string): void {
  localStorage.setItem('auth_token', token)
}

export function clearAuthToken(): void {
  localStorage.removeItem('auth_token')
}

function getAuthHeaders(): HeadersInit {
  const token = getAuthToken()
  const headers: HeadersInit = {
    'Content-Type': 'application/json',
  }
  if (token) {
    headers['Authorization'] = `Bearer ${token}`
  }
  return headers
}

interface FetchPostsParams {
  subreddits: string[]
  keyword?: string
  sort: string
  dateRange?: string
  limit: number
}

export interface RedditPost {
  id: string
  title: string
  author: string
  subreddit: string
  score: number
  url: string
  permalink: string
  created_utc: number
  num_comments: number
  thumbnail: string
  selftext: string
  is_video: boolean
}

interface FetchPostsResponse {
  posts: RedditPost[]
  count: number
}

// Fetch Reddit posts with authentication
export async function fetchRedditPosts(params: FetchPostsParams): Promise<RedditPost[]> {
  const queryParams = new URLSearchParams()

  // Add subreddits (comma-separated)
  if (params.subreddits.length > 0) {
    const subreddits = params.subreddits
      .map(s => s.replace('r/', '').trim())
      .join(',')
    queryParams.append('subreddits', subreddits)
  }

  // Add keyword if present
  if (params.keyword && params.keyword.trim()) {
    queryParams.append('keyword', params.keyword.trim())
  }

  // Add sort
  queryParams.append('sort', params.sort)

  // Add date range if present
  if (params.dateRange && params.dateRange !== 'all') {
    queryParams.append('date_range', params.dateRange)
  }

  // Add limit
  queryParams.append('limit', params.limit.toString())

  const url = `${API_BASE_URL}/api/reddit/posts?${queryParams.toString()}`

  const response = await fetch(url, {
    method: 'GET',
    headers: getAuthHeaders(),
  })

  if (!response.ok) {
    if (response.status === 401) {
      throw new Error('Not authenticated. Please log in.')
    }
    throw new Error(`Failed to fetch posts: ${response.statusText}`)
  }

  const data: FetchPostsResponse = await response.json()
  return data.posts
}

// Get current user
export async function getCurrentUser() {
  console.log('🔍 [getCurrentUser] Fetching current user...')
  const token = getAuthToken()
  console.log('🔍 [getCurrentUser] Token in localStorage:', token ? `${token.substring(0, 20)}... (${token.length} chars)` : 'NONE')

  const url = `${API_BASE_URL}/api/auth/user`
  const headers = getAuthHeaders()
  console.log('🔍 [getCurrentUser] Request URL:', url)
  console.log('🔍 [getCurrentUser] Headers:', headers)

  const response = await fetch(url, {
    method: 'GET',
    headers: headers,
  })

  console.log('🔍 [getCurrentUser] Response status:', response.status)

  if (!response.ok) {
    if (response.status === 401) {
      console.log('❌ [getCurrentUser] Not authenticated (401)')
      return null
    }
    const errorText = await response.text()
    console.error('❌ [getCurrentUser] Error:', response.statusText, errorText)
    throw new Error(`Failed to get user: ${response.statusText}`)
  }

  const data = await response.json()
  console.log('✅ [getCurrentUser] Success! User:', data.user?.email)
  return data.user
}

// Store token from OAuth callback
export async function storeAuthToken(token: string) {
  console.log('🔐 [storeAuthToken] Validating and storing token...')
  console.log('🔐 [storeAuthToken] Token preview:', `${token.substring(0, 20)}... (${token.length} chars)`)

  // Validate token by fetching user
  const url = `${API_BASE_URL}/api/auth/user`
  console.log('🔐 [storeAuthToken] Validating with:', url)

  const response = await fetch(url, {
    method: 'GET',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`
    },
  })

  console.log('🔐 [storeAuthToken] Validation response status:', response.status)

  if (!response.ok) {
    const errorText = await response.text()
    console.error('❌ [storeAuthToken] Validation failed:', errorText)
    throw new Error(`Invalid token: ${response.statusText}`)
  }

  const data = await response.json()
  console.log('✅ [storeAuthToken] Token valid! User:', data.user?.email)

  // Store token in localStorage
  setAuthToken(token)
  console.log('✅ [storeAuthToken] Token stored in localStorage')

  return data.user
}

// Logout
export async function logout() {
  const url = `${API_BASE_URL}/api/auth/logout`

  const response = await fetch(url, {
    method: 'POST',
    headers: getAuthHeaders(),
  })

  if (!response.ok) {
    throw new Error(`Failed to logout: ${response.statusText}`)
  }

  // Clear token from localStorage
  clearAuthToken()

  return response.json()
}

// Get Notion OAuth URL
export async function getNotionAuthUrl() {
  const url = `${API_BASE_URL}/api/auth/notion/login`

  const response = await fetch(url, {
    method: 'GET',
    headers: {
      'Content-Type': 'application/json',
    },
  })

  if (!response.ok) {
    throw new Error(`Failed to get auth URL: ${response.statusText}`)
  }

  const data = await response.json()
  return data.url
}

// Notion API functions

export interface SaveToNotionRequest {
  title: string
  subreddit: string
  content: string
  author: string
  score: number
  url: string
  reddit_id: string
  database_id: string
}

export interface SaveToNotionResponse {
  success: boolean
  notion_page_id: string
  notion_page_url: string
  message: string
}

export interface NotionDatabase {
  id: string
  title: string
  url: string
}

// Save a Reddit post to Notion
export async function saveToNotion(post: SaveToNotionRequest): Promise<SaveToNotionResponse> {
  const url = `${API_BASE_URL}/api/notion/save`

  const response = await fetch(url, {
    method: 'POST',
    headers: getAuthHeaders(),
    body: JSON.stringify(post),
  })

  if (!response.ok) {
    if (response.status === 401) {
      throw new Error('Not authenticated. Please log in.')
    }
    const errorData = await response.json()
    throw new Error(errorData.error || `Failed to save to Notion: ${response.statusText}`)
  }

  return response.json()
}

// Get user's Notion databases
export async function getNotionDatabases(): Promise<NotionDatabase[]> {
  const url = `${API_BASE_URL}/api/notion/databases`

  const response = await fetch(url, {
    method: 'GET',
    headers: getAuthHeaders(),
  })

  if (!response.ok) {
    if (response.status === 401) {
      throw new Error('Not authenticated. Please log in.')
    }
    throw new Error(`Failed to get databases: ${response.statusText}`)
  }

  const data = await response.json()
  return data.databases
}

// Get saved posts
export async function getSavedPosts(): Promise<import('@/types').RedditPost[]> {
  const url = `${API_BASE_URL}/api/notion/saved-posts`

  const response = await fetch(url, {
    method: 'GET',
    headers: getAuthHeaders(),
  })

  if (!response.ok) {
    if (response.status === 401) {
      throw new Error('Not authenticated. Please log in.')
    }
    throw new Error(`Failed to get saved posts: ${response.statusText}`)
  }

  interface SavedPostBackend {
    reddit_id: string
    title: string
    subreddit: string
    content: string
    author: string
    score: number
    saved_at: string
    url: string
    notion_page_url: string
  }

  const data: { posts: SavedPostBackend[] } = await response.json()

  // Transform backend response to frontend RedditPost format
  return data.posts.map((post) => ({
    id: post.reddit_id,
    title: post.title,
    subreddit: post.subreddit,
    content: post.content || '',
    author: post.author,
    score: post.score,
    created: post.saved_at,
    url: post.url,
    saved: true,
    notionPageUrl: post.notion_page_url
  }))
}

// Delete a saved post
export async function deleteSavedPost(redditId: string): Promise<void> {
  const url = `${API_BASE_URL}/api/notion/saved-posts/${redditId}`

  const response = await fetch(url, {
    method: 'DELETE',
    headers: getAuthHeaders(),
  })

  if (!response.ok) {
    const errorData = await response.json().catch(() => ({ error: response.statusText }))
    throw new Error(errorData.error || 'Failed to delete post')
  }
}
