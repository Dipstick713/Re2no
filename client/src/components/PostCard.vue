<script setup lang="ts">
import { Loader2, Trash2 } from 'lucide-vue-next'
import type { RedditPost } from '@/types'

defineProps<{
  post: RedditPost
  isSaving?: boolean
  isDeleting?: boolean
}>()

const emit = defineEmits<{
  save: [id: string]
  open: [id: string]
  delete: [id: string]
}>()
</script>

<template>
  <div class="bg-black/40 border border-white/10 rounded-2xl p-6 hover:border-white/20 hover:bg-black/60 transition-all backdrop-blur-xl group flex flex-col h-full">
    <div class="flex items-start justify-between mb-3">
      <div class="flex-1">
        <a
          :href="post.url"
          target="_blank"
          rel="noopener noreferrer"
          class="text-lg font-semibold text-white mb-2 group-hover:text-cyan-400 transition-colors block hover:underline"
        >
          {{ post.title }}
        </a>
        <div class="flex items-center gap-2 text-sm text-gray-400 mb-3">
          <span class="px-2 py-1 rounded-lg bg-white/5 border border-white/10">{{ post.subreddit }}</span>
        </div>
      </div>
    </div>
    <p class="text-gray-400 text-sm mb-4 line-clamp-2 flex-grow" v-html="post.displayContent || post.content"></p>
    <button
      v-if="!post.saved"
      @click="emit('save', post.id)"
      :disabled="isSaving"
      class="w-full px-4 py-2.5 rounded-xl bg-gradient-to-r from-blue-500 via-cyan-400 to-emerald-400 hover:from-blue-400 hover:via-cyan-300 hover:to-emerald-300 text-white font-semibold transition-all shadow-lg shadow-cyan-500/25 disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2 mt-auto"
    >
      <Loader2 v-if="isSaving" :size="16" class="animate-spin" />
      <span>{{ isSaving ? 'Saving...' : 'Save to Notion' }}</span>
    </button>
    <div v-else class="flex gap-2 mt-auto">
      <button
        @click="emit('open', post.id)"
        class="flex-1 border border-cyan-500/50 bg-cyan-500/20 text-cyan-400 hover:bg-cyan-500/10 font-semibold px-4 py-2.5 rounded-xl transition-all"
      >
        Open in Notion
      </button>
      <button
        @click="emit('delete', post.id)"
        :disabled="isDeleting"
        class="px-4 py-2.5 rounded-xl border border-red-500/50 bg-red-500/20 text-red-400 hover:bg-red-500/30 transition-all disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center"
        title="Delete from saved posts"
      >
        <Trash2 v-if="!isDeleting" :size="18" />
        <Loader2 v-else :size="18" class="animate-spin" />
      </button>
    </div>
  </div>
</template>
