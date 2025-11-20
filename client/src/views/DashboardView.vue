<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ChevronDown, Loader2, Info } from 'lucide-vue-next'
import AppHeader from '@/components/AppHeader.vue'
import FilterBar from '@/components/FilterBar.vue'
import PostCard from '@/components/PostCard.vue'
import DashboardSection from '@/components/DashboardSection.vue'
import AppFooter from '@/components/AppFooter.vue'
import ToastContainer from '@/components/ToastContainer.vue'
import { useToast } from '@/composables/useToast'
import {
  fetchRedditPosts,
  getCurrentUser,
  storeAuthToken,
  getNotionDatabases,
  saveToNotion,
  getSavedPosts,
  deleteSavedPost,
  type RedditPost as APIRedditPost,
  type NotionDatabase
} from '@/lib/api'
import type { RedditPost, FilterOptions } from '@/types'

const router = useRouter()
const toast = useToast()
const posts = ref<RedditPost[]>([])
const fetchedPosts = ref<RedditPost[]>([])
const isLoading = ref(false)
const error = ref<string | null>(null)
const isAuthenticated = ref(false)
const databases = ref<NotionDatabase[]>([])
const selectedDatabase = ref<string>('')
const savingPostId = ref<string | null>(null)
const deletingPostId = ref<string | null>(null)
const showDatabaseDropdown = ref(false)
const loadingDatabases = ref(false)
const showInstructions = ref(false)

// Check authentication on mount and load databases
onMounted(async () => {
  try {
    // Check if we have a token in the URL (from OAuth redirect)
    const urlParams = new URLSearchParams(window.location.search)
    const token = urlParams.get('token')

    if (token) {
      // Store token in localStorage and validate
      await storeAuthToken(token)
      // Remove token from URL for security
      window.history.replaceState({}, document.title, '/dashboard')
      toast.success('Successfully connected to Notion!')
    }

    const user = await getCurrentUser()
    if (user) {
      isAuthenticated.value = true
      // Load user's Notion databases
      await loadDatabases()
      // Load saved posts from database
      await loadSavedPosts()
    } else {
      // Not authenticated, redirect to home
      router.push('/')
    }
  } catch (err) {
    const errorMessage = err instanceof Error ? err.message : 'Authentication failed'
    toast.error(errorMessage)
    router.push('/')
  }
})

// Load Notion databases
const loadDatabases = async () => {
  loadingDatabases.value = true
  error.value = null
  try {
    // Add 1 second minimum delay for loader visibility
    const [databasesResult] = await Promise.all([
      getNotionDatabases(),
      new Promise(resolve => setTimeout(resolve, 1000))
    ])
    databases.value = databasesResult
    // Auto-select first database if available
    if (databases.value.length > 0 && databases.value[0]) {
      selectedDatabase.value = databases.value[0].id
    }
  } catch (err) {
    const errorMessage = err instanceof Error ? err.message : 'Failed to load Notion databases. Please check your connection.'
    error.value = errorMessage
    toast.error(errorMessage)
  } finally {
    loadingDatabases.value = false
  }
}

// Load saved posts from database
const loadSavedPosts = async () => {
  try {
    const savedPostsData = await getSavedPosts()

    // Process saved posts to add displayContent
    const processedSavedPosts = savedPostsData.map((savedPost: RedditPost) => {
      const content = savedPost.content || 'No content available'
      let displayContent = content

      // Check if content is a URL (image/video/link)
      if (content.startsWith('http://') || content.startsWith('https://')) {
        const isVideo = content.includes('v.redd.it') || content.includes('youtube.com') || content.includes('youtu.be')
        const isImage = content.match(/\.(jpg|jpeg|png|gif|webp)$/i) || content.includes('i.redd.it') || content.includes('imgur.com')

        if (isVideo) {
          displayContent = `Video: <a href="${content}" target="_blank" rel="noopener noreferrer" class="text-cyan-400 hover:text-cyan-300 underline">${content}</a>`
        } else if (isImage) {
          displayContent = `Image: <a href="${content}" target="_blank" rel="noopener noreferrer" class="text-cyan-400 hover:text-cyan-300 underline">${content}</a>`
        } else {
          displayContent = `Link: <a href="${content}" target="_blank" rel="noopener noreferrer" class="text-cyan-400 hover:text-cyan-300 underline">${content}</a>`
        }
      }

      return {
        ...savedPost,
        displayContent,
        saved: true
      }
    })

    // Merge saved posts with current posts
    // Mark posts as saved if they exist in savedPostsData
    posts.value = posts.value.map(post => {
      const savedPost = processedSavedPosts.find((sp: RedditPost) => sp.id === post.id)
      if (savedPost) {
        return {
          ...post,
          saved: true,
          notionPageUrl: savedPost.notionPageUrl,
          displayContent: post.displayContent || savedPost.displayContent
        }
      }
      return post
    })

    // Add any saved posts that aren't in the current posts
    const currentPostIds = new Set(posts.value.map(p => p.id))
    const additionalSavedPosts = processedSavedPosts.filter((sp: RedditPost) => !currentPostIds.has(sp.id))
    posts.value = [...additionalSavedPosts, ...posts.value]

  } catch {
    // Silently fail - saved posts will load on next refresh
  }
}

const getSelectedDatabaseTitle = () => {
  const db = databases.value.find(d => d.id === selectedDatabase.value)
  return db ? db.title : 'Select a database...'
}

const selectDatabase = (dbId: string) => {
  selectedDatabase.value = dbId
  showDatabaseDropdown.value = false
}

// Convert API post to internal post format
const convertPost = (apiPost: APIRedditPost): RedditPost => {
  // If no selftext, check if it's an image/video post and show the URL
  let content = apiPost.selftext
  let displayContent = content

  if (!content || content.trim() === '') {
    // Store the original URL for saving to Notion
    content = apiPost.url || 'No content available'

    // Create a clickable link for display in the UI
    if (apiPost.url && (apiPost.url.includes('i.redd.it') || apiPost.url.includes('imgur') || apiPost.is_video)) {
      displayContent = `Image/Video: <a href="${apiPost.url}" target="_blank" rel="noopener noreferrer" class="text-cyan-400 hover:text-cyan-300 underline">${apiPost.url}</a>`
    } else if (apiPost.url) {
      displayContent = `Link: <a href="${apiPost.url}" target="_blank" rel="noopener noreferrer" class="text-cyan-400 hover:text-cyan-300 underline">${apiPost.url}</a>`
    } else {
      displayContent = 'No content available'
    }
  } else {
    displayContent = content
  }

  return {
    id: apiPost.id,
    title: apiPost.title,
    content: content, // Plain content for saving to Notion
    displayContent: displayContent, // HTML content for display
    url: `https://reddit.com${apiPost.permalink}`,
    subreddit: apiPost.subreddit,
    author: apiPost.author,
    score: apiPost.score,
    created: new Date(apiPost.created_utc * 1000).toISOString(),
    saved: false,
  }
}

const handleFetch = async (filters: FilterOptions) => {
  isLoading.value = true
  error.value = null

  try {
    // First, get all saved posts from database to check against
    let savedPostsData: RedditPost[] = []
    try {
      savedPostsData = await getSavedPosts()
      console.log('📦 Saved posts from DB:', savedPostsData.length, savedPostsData.map(p => p.id))
    } catch (err) {
      console.error('Failed to load saved posts:', err)
      // Continue even if we can't load saved posts
    }

    const apiPosts = await fetchRedditPosts({
      subreddits: filters.subreddits,
      keyword: filters.keyword,
      sort: filters.sortBy,
      dateRange: filters.dateRange === 'all' ? undefined : filters.dateRange,
      limit: filters.numberOfPosts,
    })

    // Convert API posts to internal format
    const convertedPosts = apiPosts.map(convertPost)
    console.log('🔄 Converted posts:', convertedPosts.length, convertedPosts.map(p => p.id))

    // Update fetched posts
    fetchedPosts.value = convertedPosts

    // Create a map of saved posts for quick lookup
    const savedPostsMap = new Map(savedPostsData.map(sp => [sp.id, sp]))
    console.log('🗺️ Saved posts map:', Array.from(savedPostsMap.keys()))

    // Merge with existing saved posts (don't replace them)
    // Keep saved posts that aren't in the new fetch
    const newPostIds = new Set(convertedPosts.map(p => p.id))
    const savedPosts = posts.value.filter(p => p.saved && !newPostIds.has(p.id))

    // Check if any newly fetched posts were previously saved
    const mergedNewPosts = convertedPosts.map(newPost => {
      const savedPost = savedPostsMap.get(newPost.id)
      if (savedPost) {
        console.log('✅ Found saved post:', newPost.id, newPost.title.substring(0, 30))
        // Mark as saved and preserve Notion URL
        return {
          ...newPost,
          saved: true,
          notionPageUrl: savedPost.notionPageUrl
        }
      }
      return newPost
    })

    console.log('🎯 Merged posts with saved status:', mergedNewPosts.filter(p => p.saved).length, 'saved')

    // Combine saved posts + newly fetched posts
    posts.value = [...savedPosts, ...mergedNewPosts]

  } catch (err) {
    const errorMessage = err instanceof Error ? err.message : 'Failed to fetch posts'
    error.value = errorMessage

    if (errorMessage.includes('Not authenticated')) {
      toast.error('Session expired. Please login again.')
      // Redirect to login
      router.push('/')
    } else {
      toast.error(errorMessage)
    }
  } finally {
    isLoading.value = false
  }
}

const handleSave = async (id: string) => {
  const post = posts.value.find(p => p.id === id) || fetchedPosts.value.find(p => p.id === id)

  if (!post) {
    return
  }

  if (!selectedDatabase.value) {
    error.value = 'Please select a Notion database first. Go to your Notion workspace and create a database.'
    toast.error('Please select a Notion database first')
    return
  }

  savingPostId.value = id
  error.value = null

  try {
    const response = await saveToNotion({
      title: post.title,
      subreddit: post.subreddit,
      content: post.content,
      author: post.author,
      score: post.score,
      url: post.url,
      reddit_id: post.id,
      database_id: selectedDatabase.value,
    })

    // Update the post with Notion page URL and mark as saved
    const mainPost = posts.value.find(p => p.id === id)
    if (mainPost) {
      mainPost.saved = true
      mainPost.notionPageUrl = response.notion_page_url
    }

    const fetchedPost = fetchedPosts.value.find(p => p.id === id)
    if (fetchedPost) {
      fetchedPost.saved = true
      fetchedPost.notionPageUrl = response.notion_page_url
    }

    // Show success message
    toast.success('Post saved to Notion successfully!')

  } catch (err) {
    const errorMessage = err instanceof Error ? err.message : 'Failed to save post to Notion'
    error.value = errorMessage
    toast.error(errorMessage)
  } finally {
    savingPostId.value = null
  }
}

const handleOpen = (id: string) => {
  const post = posts.value.find(p => p.id === id) || fetchedPosts.value.find(p => p.id === id)

  if (!post) {
    error.value = 'Post not found'
    return
  }

  if (post.notionPageUrl) {
    window.open(post.notionPageUrl, '_blank')
  } else {
    error.value = 'Notion page URL not available. Try saving the post again.'
    toast.error('Notion page URL not available. Try saving the post again.')
  }
}

const handleDelete = async (id: string) => {
  const post = posts.value.find(p => p.id === id)

  if (!post) {
    return
  }

  deletingPostId.value = id

  try {
    await deleteSavedPost(id)

    // Remove the post from both arrays
    posts.value = posts.value.filter(p => p.id !== id)
    fetchedPosts.value = fetchedPosts.value.map(p => {
      if (p.id === id) {
        return { ...p, saved: false, notionPageUrl: undefined }
      }
      return p
    })

    toast.success('Post deleted successfully!')
  } catch (err) {
    const errorMessage = err instanceof Error ? err.message : 'Failed to delete post'
    toast.error(errorMessage)
  } finally {
    deletingPostId.value = null
  }
}
</script>

<template>
  <div class="relative min-h-screen bg-black flex flex-col">
    <div class="absolute inset-0 bg-[url(/image.png)] bg-repeat opacity-10 pointer-events-none"></div>

    <ToastContainer />

    <div class="relative z-10 flex flex-col min-h-screen">
      <AppHeader />

      <main class="flex-1">
        <section class="py-8 px-6 mt-8">
          <div class="container mx-auto">
            <h1 class="text-4xl font-bold text-white mb-6">Your Dashboard</h1>

            <!-- Database Selector -->
            <div v-if="loadingDatabases" class="mb-6">
              <label class="block text-sm font-medium text-gray-300 mb-2">
                Select Notion Database
              </label>
              <div class="w-full md:w-96 px-4 py-2.5 rounded-xl bg-black/40 border border-white/20 flex items-center justify-center gap-2">
                <Loader2 :size="20" class="animate-spin text-cyan-500" />
                <span class="text-gray-400">Loading databases...</span>
              </div>
            </div>

            <div v-else-if="databases.length > 0" class="mb-6">
              <label class="block text-sm font-medium text-gray-300 mb-2">
                Select Notion Database
              </label>
              <div class="relative w-full md:w-96">
                <button
                  @click="showDatabaseDropdown = !showDatabaseDropdown"
                  class="w-full px-4 py-2.5 rounded-xl bg-black/40 border border-white/20 text-white focus:border-cyan-500 focus:outline-none transition-colors flex items-center justify-between"
                >
                  <span>{{ getSelectedDatabaseTitle() }}</span>
                  <ChevronDown :size="20" :class="{ 'rotate-180': showDatabaseDropdown }" class="transition-transform" />
                </button>

                <!-- Custom Dropdown -->
                <div
                  v-if="showDatabaseDropdown"
                  class="absolute z-50 w-full mt-2 rounded-xl bg-neutral-900 border border-white/20 shadow-xl max-h-60 overflow-y-auto"
                >
                  <button
                    v-for="db in databases"
                    :key="db.id"
                    @click="selectDatabase(db.id)"
                    class="w-full px-4 py-3 text-left text-white hover:bg-cyan-500/20 transition-colors first:rounded-t-xl last:rounded-b-xl"
                    :class="{ 'bg-black/40': selectedDatabase === db.id }"
                  >
                    {{ db.title }}
                  </button>
                </div>
              </div>
              <p class="mt-2 text-sm text-gray-400">
                Posts will be saved to this database
              </p>
            </div>            <!-- Warning if no databases -->
            <div v-else-if="isAuthenticated && !isLoading" class="mb-6 p-4 rounded-xl bg-red-500/20 border border-red-500/30 text-red-400">
              <div class="flex items-start gap-3">
                <div class="flex-1">
                  <p class="font-semibold mb-2">No Notion databases found</p>
                  <p class="text-sm">Please create a database in your Notion workspace to save posts.</p>
                </div>
                <button
                  @click="showInstructions = true"
                  class="flex-shrink-0 p-2 rounded-lg hover:bg-red-500/30 transition-colors"
                  title="How to create a database"
                >
                  <Info :size="20" />
                </button>
              </div>
            </div>

            <!-- Error Message -->
            <div v-if="error" class="mb-6 p-4 rounded-xl bg-red-500/20 border border-red-500/30 text-red-400">
              <p class="font-semibold">{{ error }}</p>
            </div>

            <FilterBar @fetch="handleFetch" />

            <!-- Loading State -->
            <div v-if="isLoading" class="mt-8 text-center">
              <Loader2 :size="48" class="inline-block animate-spin text-cyan-500" />
              <p class="text-gray-400 mt-4">Fetching posts from Reddit...</p>
            </div>
          </div>
        </section>

        <section v-if="fetchedPosts.length > 0 && !isLoading" class="py-8 px-6">
          <div class="container mx-auto">
            <h2 class="text-3xl font-bold text-white mb-6">Fetched Reddit Posts ({{ fetchedPosts.length }})</h2>
            <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
              <PostCard
                v-for="post in fetchedPosts"
                :key="post.id"
                :post="post"
                :is-saving="savingPostId === post.id"
                :is-deleting="deletingPostId === post.id"
                @save="handleSave"
                @open="handleOpen"
                @delete="handleDelete"
              />
            </div>
          </div>
        </section>

        <DashboardSection
          :posts="posts"
          :saving-post-id="savingPostId"
          :deleting-post-id="deletingPostId"
          @save="handleSave"
          @open="handleOpen"
          @delete="handleDelete"
        />
      </main>

      <AppFooter />
    </div>

    <!-- Instructions Modal -->
    <div
      v-if="showInstructions"
      class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/80 backdrop-blur-sm"
      @click.self="showInstructions = false"
    >
      <div class="bg-black/90 border border-white/20 rounded-2xl p-6 max-w-2xl w-full max-h-[80vh] overflow-y-auto">
        <div class="flex items-start justify-between mb-4">
          <h3 class="text-2xl font-bold text-white">How to Create a Notion Database</h3>
          <button
            @click="showInstructions = false"
            class="text-gray-400 hover:text-white transition-colors p-1"
          >
            <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <line x1="18" y1="6" x2="6" y2="18"></line>
              <line x1="6" y1="6" x2="18" y2="18"></line>
            </svg>
          </button>
        </div>

        <div class="space-y-4 text-gray-300">
          <div>
            <h4 class="text-lg font-semibold text-white mb-2">Step 1: Open Notion</h4>
            <p class="text-sm">Go to your Notion workspace in your browser or desktop app.</p>
          </div>

          <div>
            <h4 class="text-lg font-semibold text-white mb-2">Step 2: Create a New Page</h4>
            <p class="text-sm">Click on "+ New page" in the sidebar or press <code class="px-1.5 py-0.5 bg-white/10 rounded">Cmd/Ctrl + N</code></p>
          </div>

          <div>
            <h4 class="text-lg font-semibold text-white mb-2">Step 3: Add a Database</h4>
            <p class="text-sm mb-2">In your new page, type <code class="px-1.5 py-0.5 bg-white/10 rounded">/table</code> or <code class="px-1.5 py-0.5 bg-white/10 rounded">/database</code> and select one of these options:</p>
            <ul class="list-disc list-inside space-y-1 text-sm ml-4">
              <li><strong>Table - Inline</strong>: Creates a database in the current page</li>
              <li><strong>Table - Full page</strong>: Creates a new page with a database</li>
            </ul>
          </div>

          <div>
            <h4 class="text-lg font-semibold text-white mb-2">Step 4: Name Your Database</h4>
            <p class="text-sm">Give it a name like "Reddit Posts" or "Saved Articles"</p>
          </div>

          <div>
            <h4 class="text-lg font-semibold text-white mb-2">Step 5: Refresh This Page</h4>
            <p class="text-sm">Once created, refresh this page and your new database should appear in the dropdown above.</p>
          </div>

          <div class="pt-4 border-t border-white/10">
            <p class="text-sm text-gray-400">
              <strong class="text-cyan-400">Tip:</strong> The app will automatically map Reddit post data to your database properties. You can customize the database structure in Notion after creation.
            </p>
          </div>
        </div>

        <button
          @click="showInstructions = false"
          class="mt-6 w-full px-4 py-2.5 rounded-xl bg-gradient-to-r from-blue-500 via-cyan-400 to-emerald-400 hover:from-blue-400 hover:via-cyan-300 hover:to-emerald-300 text-white font-semibold transition-all"
        >
          Got it!
        </button>
      </div>
    </div>
  </div>
</template>
