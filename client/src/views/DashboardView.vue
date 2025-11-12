<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import AppHeader from '@/components/AppHeader.vue'
import FilterBar from '@/components/FilterBar.vue'
import PostCard from '@/components/PostCard.vue'
import DashboardSection from '@/components/DashboardSection.vue'
import AppFooter from '@/components/AppFooter.vue'
import { fetchRedditPosts, getCurrentUser, type RedditPost as APIRedditPost } from '@/lib/api'
import type { RedditPost, FilterOptions } from '@/types'

const router = useRouter()
const posts = ref<RedditPost[]>([])
const fetchedPosts = ref<RedditPost[]>([])
const isLoading = ref(false)
const error = ref<string | null>(null)
const isAuthenticated = ref(false)

// Check authentication on mount
onMounted(async () => {
  try {
    const user = await getCurrentUser()
    if (user) {
      isAuthenticated.value = true
    } else {
      // Not authenticated, redirect to home
      router.push('/')
    }
  } catch (err) {
    console.error('Failed to get user:', err)
    router.push('/')
  }
})

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
    comments: apiPost.num_comments,
    saved: false,
    thumbnail: apiPost.thumbnail && !apiPost.thumbnail.startsWith('self')
      ? apiPost.thumbnail
      : undefined,
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

const handleSave = (id: string) => {
  const post = posts.value.find(p => p.id === id)
  if (post) {
    post.saved = true
  }

  const fetchedPost = fetchedPosts.value.find(p => p.id === id)
  if (fetchedPost) {
    fetchedPost.saved = true
  }
}

const handleOpen = (id: string) => {
  console.log('Opening post in Notion:', id)
  // TODO: Implement Notion API to save post
}
</script>

<template>
  <div class="relative min-h-screen bg-black">
    <div class="absolute inset-0 bg-[url(/image.png)] bg-repeat opacity-10 pointer-events-none"></div>

    <div class="relative z-10">
      <AppHeader />

      <main>
        <section class="py-8 px-6 mt-8">
          <div class="container mx-auto">
            <h1 class="text-4xl font-bold text-white mb-8">Your Dashboard</h1>

            <!-- Error Message -->
            <div v-if="error" class="mb-6 p-4 rounded-xl bg-red-500/20 border border-red-500/30 text-red-400">
              <p class="font-semibold">{{ error }}</p>
            </div>

            <FilterBar @fetch="handleFetch" />

            <!-- Loading State -->
            <div v-if="isLoading" class="mt-8 text-center">
              <div class="inline-block animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-cyan-500"></div>
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
                @save="handleSave"
                @open="handleOpen"
              />
            </div>
          </div>
        </section>

        <DashboardSection
          v-if="posts.length > 0 && !isLoading"
          :posts="posts"
          @save="handleSave"
          @open="handleOpen"
        />
      </main>

      <AppFooter />
    </div>
  </div>
</template>
