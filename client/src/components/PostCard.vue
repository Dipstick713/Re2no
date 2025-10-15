<script setup lang="ts">
import { GripVertical } from 'lucide-vue-next'
import type { RedditPost } from '@/types'

defineProps<{
  post: RedditPost
}>()

const emit = defineEmits<{
  save: [id: string]
  open: [id: string]
}>()
</script>

<template>
  <div class="bg-black/40 border border-white/10 rounded-2xl p-6 hover:border-white/20 hover:bg-black/60 transition-all backdrop-blur-xl group">
    <div class="flex items-start justify-between mb-3">
      <div class="flex-1">
        <h3 class="text-lg font-semibold text-white mb-2 group-hover:text-cyan-400 transition-colors">{{ post.title }}</h3>
        <div class="flex items-center gap-2 text-sm text-gray-400 mb-3">
          <span class="px-2 py-1 rounded-lg bg-white/5 border border-white/10">{{ post.subreddit }}</span>
        </div>
      </div>
      <button class="text-gray-500 hover:text-gray-300 transition-colors">
        <GripVertical :size="20" />
      </button>
    </div>
    <p class="text-gray-400 text-sm mb-4 line-clamp-2">{{ post.content }}</p>
    <button
      v-if="!post.saved"
      @click="emit('save', post.id)"
      class="w-full px-4 py-2.5 rounded-xl bg-gradient-to-r from-blue-500 via-cyan-400 to-emerald-400 hover:from-blue-400 hover:via-cyan-300 hover:to-emerald-300 text-white font-semibold transition-all shadow-lg shadow-cyan-500/25"
    >
      Save to Notion
    </button>
    <button
      v-else
      @click="emit('open', post.id)"
      class="w-full border border-cyan-500/50 text-cyan-400 hover:bg-cyan-500/10 font-semibold px-4 py-2.5 rounded-xl transition-all"
    >
      Open in Notion
    </button>
  </div>
</template>
