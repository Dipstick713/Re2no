export interface RedditPost {
  id: string
  title: string
  subreddit: string
  content: string
  author: string
  score: number
  created: string
  url: string
  saved: boolean
  notionPageUrl?: string
}

export interface FilterOptions {
  subreddits: string[]
  keyword: string
  dateRange: string
  sortBy: string
  numberOfPosts: number
}
