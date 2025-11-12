const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080'

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
    credentials: 'include', // Important: Send cookies with request
    headers: {
      'Content-Type': 'application/json',
    },
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
  const url = `${API_BASE_URL}/api/auth/user`

  const response = await fetch(url, {
    method: 'GET',
    credentials: 'include',
    headers: {
      'Content-Type': 'application/json',
    },
  })

  if (!response.ok) {
    if (response.status === 401) {
      return null
    }
    throw new Error(`Failed to get user: ${response.statusText}`)
  }

  const data = await response.json()
  return data.user
}

// Logout
export async function logout() {
  const url = `${API_BASE_URL}/api/auth/logout`

  const response = await fetch(url, {
    method: 'POST',
    credentials: 'include',
    headers: {
      'Content-Type': 'application/json',
    },
  })

  if (!response.ok) {
    throw new Error(`Failed to logout: ${response.statusText}`)
  }

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
