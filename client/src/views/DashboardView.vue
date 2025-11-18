<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ChevronDown, Loader2 } from 'lucide-vue-next'
import AppHeader from '@/components/AppHeader.vue'
import FilterBar from '@/components/FilterBar.vue'
import PostCard from '@/components/PostCard.vue'
import DashboardSection from '@/components/DashboardSection.vue'
import AppFooter from '@/components/AppFooter.vue'
import {
  fetchRedditPosts,
  getCurrentUser,
  getNotionDatabases,
  saveToNotion,
  type RedditPost as APIRedditPost,
  type NotionDatabase
} from '@/lib/api'
import type { RedditPost, FilterOptions } from '@/types'

const router = useRouter()
const posts = ref<RedditPost[]>([])
const fetchedPosts = ref<RedditPost[]>([])
const isLoading = ref(false)
const error = ref<string | null>(null)
const isAuthenticated = ref(false)
const databases = ref<NotionDatabase[]>([])
const selectedDatabase = ref<string>('')
const savingPostId = ref<string | null>(null)
const showDatabaseDropdown = ref(false)
const loadingDatabases = ref(false)

// Check authentication on mount and load databases
onMounted(async () => {
  try {
    const user = await getCurrentUser()
    if (user) {
      isAuthenticated.value = true
      // Load user's Notion databases
      await loadDatabases()
    } else {
      // Not authenticated, redirect to home
      router.push('/')
    }
  } catch (err) {
    console.error('Failed to get user:', err)
    router.push('/')
  }
})

// Load Notion databases
const loadDatabases = async () => {
  loadingDatabases.value = true
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
    console.error('Failed to load databases:', err)
    // Don't fail the whole app if databases can't be loaded
  } finally {
    loadingDatabases.value = false
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
  return {
    id: apiPost.id,
    title: apiPost.title,
    content: apiPost.selftext || 'No content available',
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
    console.log('Fetching with filters:', filters)

    const apiPosts = await fetchRedditPosts({
      subreddits: filters.subreddits,
      keyword: filters.keyword,
      sort: filters.sortBy,
      dateRange: filters.dateRange === 'all' ? undefined : filters.dateRange,
      limit: filters.numberOfPosts,
    })

    console.log(`Fetched ${apiPosts.length} posts`)

    // Convert API posts to internal format
    const convertedPosts = apiPosts.map(convertPost)

    // Update both fetched and main posts
    fetchedPosts.value = convertedPosts
    posts.value = convertedPosts

  } catch (err) {
    console.error('Failed to fetch posts:', err)
    error.value = err instanceof Error ? err.message : 'Failed to fetch posts'

    if (error.value.includes('Not authenticated')) {
      // Redirect to login
      router.push('/')
    }
  } finally {
    isLoading.value = false
  }
}

const handleSave = async (id: string) => {
  const post = posts.value.find(p => p.id === id) || fetchedPosts.value.find(p => p.id === id)

  if (!post) {
    console.error('Post not found:', id)
    return
  }

  if (!selectedDatabase.value) {
    error.value = 'Please select a Notion database first. Go to your Notion workspace and create a database.'
    return
  }

  savingPostId.value = id
  error.value = null

  try {
    console.log('Saving post to Notion:', post.title)

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

    console.log('Successfully saved to Notion:', response)

    // Mark post as saved in both lists
    const mainPost = posts.value.find(p => p.id === id)
    if (mainPost) {
      mainPost.saved = true
    }

    const fetchedPost = fetchedPosts.value.find(p => p.id === id)
    if (fetchedPost) {
      fetchedPost.saved = true
    }

    // Show success message (you can add a toast notification here)
    console.log('Post saved successfully! View it at:', response.notion_page_url)

  } catch (err) {
    console.error('Failed to save post:', err)
    error.value = err instanceof Error ? err.message : 'Failed to save post to Notion'
  } finally {
    savingPostId.value = null
  }
}

const handleOpen = (id: string) => {
  console.log('Opening post in Notion:', id)
  // TODO: Store Notion page URL and open it
  error.value = 'Notion page URL will be available after implementing storage'
}
</script>

<template>
  <div class="relative min-h-screen bg-black flex flex-col">
    <div class="absolute inset-0 bg-[url(/image.png)] bg-repeat opacity-10 pointer-events-none"></div>

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
            <div v-else-if="isAuthenticated && !isLoading" class="mb-6 p-4 rounded-xl bg-yellow-500/20 border border-yellow-500/30 text-yellow-400">
              <p class="font-semibold mb-2">No Notion databases found</p>
              <p class="text-sm">Please create a database in your Notion workspace to save posts.</p>
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
                @save="handleSave"
                @open="handleOpen"
              />
            </div>
          </div>
        </section>

        <DashboardSection
          v-if="posts.length > 0 && !isLoading"
          :posts="posts"
          :saving-post-id="savingPostId"
          @save="handleSave"
          @open="handleOpen"
        />
      </main>

      <AppFooter />
    </div>
  </div>
</template>
