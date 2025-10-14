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
  <section class="py-12 px-6">
    <div class="container mx-auto">
      <div class="flex items-center justify-between mb-8">
        <h2 class="text-3xl font-bold text-white">Dashboard / <span class="text-gray-400">Saved Ideas</span></h2>
        <div class="flex gap-2 bg-black/40 backdrop-blur-xl p-1.5 rounded-xl border border-white/10">
          <button
            @click="activeTab = 'all'"
            :class="[
              'px-5 py-2 rounded-lg font-medium transition-all',
              activeTab === 'all'
                ? 'bg-gradient-to-r from-teal-500 to-cyan-500 text-white shadow-lg shadow-teal-500/25'
                : 'text-gray-400 hover:text-white hover:bg-white/5'
            ]"
          >
            All
          </button>
          <button
            @click="activeTab = 'saved'"
            :class="[
              'px-5 py-2 rounded-lg font-medium transition-all',
              activeTab === 'saved'
                ? 'bg-gradient-to-r from-teal-500 to-cyan-500 text-white shadow-lg shadow-teal-500/25'
                : 'text-gray-400 hover:text-white hover:bg-white/5'
            ]"
          >
            Saved
          </button>
          <button
            @click="activeTab = 'recent'"
            :class="[
              'px-5 py-2 rounded-lg font-medium transition-all',
              activeTab === 'recent'
                ? 'bg-gradient-to-r from-teal-500 to-cyan-500 text-white shadow-lg shadow-teal-500/25'
                : 'text-gray-400 hover:text-white hover:bg-white/5'
            ]"
          >
            Recent
          </button>
        </div>
      </div>
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
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
