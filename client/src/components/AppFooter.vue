<script setup lang="ts">
import { Github } from 'lucide-vue-next'

const currentYear = new Date().getFullYear()
</script>

<template>
  <footer class="border-t border-white/10 bg-black/40 backdrop-blur-xl mt-20">
    <div class="container mx-auto px-6 py-8">
      <div class="flex items-center justify-between">
        <div class="flex items-center gap-2 text-gray-400">
          <span>© {{ currentYear }} Re2no. All rights reserved.</span>
        </div>

        <div class="flex items-center gap-4">
          <a
            href="https://twitter.com"
            target="_blank"
            rel="noopener noreferrer"
            class="w-10 h-10 rounded-lg border border-white/10 bg-white/5 hover:bg-white/10 flex items-center justify-center transition-all group"
          >
            <svg
              class="w-5 h-5 text-gray-400 group-hover:text-white transition-colors"
              fill="currentColor"
              viewBox="0 0 24 24"
            >
              <path d="M18.244 2.25h3.308l-7.227 8.26 8.502 11.24H16.17l-5.214-6.817L4.99 21.75H1.68l7.73-8.835L1.254 2.25H8.08l4.713 6.231zm-1.161 17.52h1.833L7.084 4.126H5.117z"/>
            </svg>
          </a>

          <a
            href="https://github.com"
            target="_blank"
            rel="noopener noreferrer"
            class="w-10 h-10 rounded-lg border border-white/10 bg-white/5 hover:bg-white/10 flex items-center justify-center transition-all group"
          >
            <Github :size="20" class="text-gray-400 group-hover:text-white transition-colors" />
          </a>
        </div>
      </div>
    </div>
  </footer>
</template>
