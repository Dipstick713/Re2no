<script setup lang="ts">
import { useRouter } from 'vue-router'

const router = useRouter()

const handleConnectNotion = () => {
  // This would trigger Notion OAuth flow
  router.push('/dashboard')
}

const handleGetStarted = () => {
  const featuresSection = document.getElementById('features')
  if (featuresSection) {
    featuresSection.scrollIntoView({ behavior: 'smooth' })
  }
}
</script>

<template>
  <section class="py-12 sm:py-20 px-4 sm:px-6">
    <div class="container mx-auto max-w-4xl text-center">
      <h1 class="text-4xl sm:text-5xl md:text-6xl font-bold text-white mb-4 sm:mb-6 leading-tight">
        Turn <span class="bg-gradient-to-r from-blue-500 via-cyan-400 to-emerald-400 bg-clip-text text-transparent">Reddit </span>Trends into <span class="bg-gradient-to-r from-blue-500 via-cyan-400 to-emerald-400 bg-clip-text text-transparent">Notion</span> Notes
      </h1>
      <p class="text-lg sm:text-xl text-gray-400 mb-8 sm:mb-10 max-w-2xl mx-auto px-4">
        Seamlessly funnel the best of Reddit directly into Notion. Capture trends, filter noise, and build your knowledge base—on autopilot.
      </p>
      <div class="flex flex-col sm:flex-row gap-4 justify-center px-4">
        <button
          @click="handleConnectNotion"
          class="flex items-center justify-center gap-2 px-6 py-3 rounded-xl bg-gradient-to-r from-blue-500 via-cyan-400 to-emerald-400 hover:from-blue-400 hover:via-cyan-300 hover:to-emerald-300 text-white font-semibold transition-all shadow-lg shadow-cyan-500/25"
        >
          Connect Notion
        </button>
        <button
          @click="handleGetStarted"
          class="px-6 py-3 rounded-xl border border-white/20 bg-white/5 hover:bg-white/10 text-white font-medium transition-all backdrop-blur-sm"
        >
          Learn More
        </button>
      </div>
    </div>
  </section>
</template>
