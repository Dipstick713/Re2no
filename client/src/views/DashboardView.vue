<script setup lang="ts">
import { ref } from 'vue'
import AppHeader from '@/components/AppHeader.vue'
import FilterBar from '@/components/FilterBar.vue'
import PostCard from '@/components/PostCard.vue'
import DashboardSection from '@/components/DashboardSection.vue'
import AppFooter from '@/components/AppFooter.vue'
import { mockPosts } from '@/data/mockData'
import type { RedditPost } from '@/types'

const posts = ref<RedditPost[]>(mockPosts)
const fetchedPosts = ref<RedditPost[]>([])

const handleFetch = () => {
  // Simulate fetching posts - get first 3 unsaved posts
  fetchedPosts.value = posts.value.filter(p => !p.saved).slice(0, 3)
}

const handleSave = (id: string) => {
  const post = posts.value.find(p => p.id === id)
  if (post) {
    post.saved = true
  }
}

const handleOpen = (id: string) => {
  console.log('Opening post in Notion:', id)
  // This would open the post in Notion
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
            <FilterBar @fetch="handleFetch" />
          </div>
        </section>

        <section v-if="fetchedPosts.length > 0" class="py-8 px-6">
          <div class="container mx-auto">
            <h2 class="text-3xl font-bold text-white mb-6">Fetched Reddit Posts</h2>
            <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
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

        <DashboardSection :posts="posts" @save="handleSave" @open="handleOpen" />
      </main>

      <AppFooter />
    </div>
  </div>
</template>
