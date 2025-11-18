<script setup lang="ts">
import { computed } from 'vue'
import PostCard from './PostCard.vue'
import type { RedditPost } from '@/types'

const props = defineProps<{
  posts: RedditPost[]
  savingPostId?: string | null
}>()

const emit = defineEmits<{
  save: [id: string]
  open: [id: string]
}>()

const filteredPosts = computed(() => {
  // Only show saved posts
  return props.posts.filter(p => p.saved)
})
</script>

<template>
  <section class="py-8 sm:py-12 px-4 sm:px-6">
    <div class="container mx-auto">
      <div class="flex flex-col sm:flex-row items-start sm:items-center justify-between mb-6 sm:mb-8 gap-4">
        <h2 class="text-2xl sm:text-3xl font-bold text-white">Saved Posts</h2>
      </div>
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4 sm:gap-6">
        <PostCard
          v-for="post in filteredPosts"
          :key="post.id"
          :post="post"
          :is-saving="savingPostId === post.id"
          @save="emit('save', $event)"
          @open="emit('open', $event)"
        />
      </div>
    </div>
  </section>
</template>
