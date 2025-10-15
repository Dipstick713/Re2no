<script setup lang="ts">
import { ref, computed } from 'vue'
import PostCard from './PostCard.vue'
import type { RedditPost } from '@/types'

const props = defineProps<{
  posts: RedditPost[]
}>()

const activeTab = ref<'all' | 'saved' | 'recent'>('all')

const emit = defineEmits<{
  save: [id: string]
  open: [id: string]
}>()

const filteredPosts = computed(() => {
  if (activeTab.value === 'all') {
    return props.posts
  } else if (activeTab.value === 'saved') {
    return props.posts.filter(p => p.saved)
  } else {
    // Recent - sort by date
    return [...props.posts].sort((a, b) =>
      new Date(b.created).getTime() - new Date(a.created).getTime()
    )
  }
})
</script>

<template>
  <section class="py-8 sm:py-12 px-4 sm:px-6">
    <div class="container mx-auto">
      <div class="flex flex-col sm:flex-row items-start sm:items-center justify-between mb-6 sm:mb-8 gap-4">
        <h2 class="text-2xl sm:text-3xl font-bold text-white">Dashboard / <span class="text-gray-400">Saved Ideas</span></h2>
        <div class="flex gap-2 bg-black/40 backdrop-blur-xl p-1.5 rounded-xl border border-white/10 w-full sm:w-auto overflow-x-auto">
          <button
            @click="activeTab = 'all'"
            :class="[
              'px-4 sm:px-5 py-2 rounded-lg font-medium transition-all whitespace-nowrap text-sm sm:text-base',
              activeTab === 'all'
                ? 'bg-gradient-to-r from-blue-500 via-cyan-400 to-emerald-400 text-white shadow-lg shadow-cyan-500/25'
                : 'text-gray-400 hover:text-white hover:bg-white/5'
            ]"
          >
            All
          </button>
          <button
            @click="activeTab = 'saved'"
            :class="[
              'px-4 sm:px-5 py-2 rounded-lg font-medium transition-all whitespace-nowrap text-sm sm:text-base',
              activeTab === 'saved'
                ? 'bg-gradient-to-r from-blue-500 via-cyan-400 to-emerald-400 text-white shadow-lg shadow-cyan-500/25'
                : 'text-gray-400 hover:text-white hover:bg-white/5'
            ]"
          >
            Saved
          </button>
          <button
            @click="activeTab = 'recent'"
            :class="[
              'px-4 sm:px-5 py-2 rounded-lg font-medium transition-all whitespace-nowrap text-sm sm:text-base',
              activeTab === 'recent'
                ? 'bg-gradient-to-r from-blue-500 via-cyan-400 to-emerald-400 text-white shadow-lg shadow-cyan-500/25'
                : 'text-gray-400 hover:text-white hover:bg-white/5'
            ]"
          >
            Recent
          </button>
        </div>
      </div>
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4 sm:gap-6">
        <PostCard
          v-for="post in filteredPosts"
          :key="post.id"
          :post="post"
          @save="emit('save', $event)"
          @open="emit('open', $event)"
        />
      </div>
    </div>
  </section>
</template>
