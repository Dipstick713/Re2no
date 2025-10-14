<script setup lang="ts">
import { BookOpen, Github } from 'lucide-vue-next'
import { useRouter } from 'vue-router'

const router = useRouter()

const handleGetStarted = () => {
  router.push('/dashboard')
}
</script>

<template>
  <header class="sticky top-4 z-50 px-6">
    <div class="container mx-auto border border-white/10 bg-black/10 backdrop-blur-xl rounded-2xl">
      <nav class="flex items-center justify-between px-6 py-4">
        <router-link to="/" class="flex items-center gap-2">
          <div class="w-8 h-8 rounded-lg bg-gradient-to-br from-teal-400 to-cyan-500 flex items-center justify-center">
            <BookOpen :size="18" class="text-white" />
          </div>
          <span class="text-xl font-semibold text-white">Re2no</span>
        </router-link>
        <div class="flex items-center gap-8">
          <a href="/" class="text-sm text-gray-300 hover:text-white transition-colors">Home</a>
          <router-link to="/dashboard" class="text-sm text-gray-300 hover:text-white transition-colors">Dashboard</router-link>
        </div>
        <div class="flex items-center gap-3">
          <a
            href="https://github.com/Dipstick713/Re2no"
            target="_blank"
            rel="noopener noreferrer"
            class="flex items-center gap-2 px-4 py-2 rounded-lg border border-white/20 bg-white/5 hover:bg-white/10 text-white text-sm font-medium transition-all"
          >
            <Github :size="16" />
            Contribute
          </a>
          <button
            @click="handleGetStarted"
            class="flex items-center gap-2 px-5 py-2 rounded-lg bg-gradient-to-r from-teal-500 to-cyan-500 hover:from-teal-400 hover:to-cyan-400 text-white text-sm font-semibold transition-all shadow-lg shadow-teal-500/25"
          >
            <span>→</span>
            Get Started
          </button>
        </div>
      </nav>
    </div>
  </header>
</template>
