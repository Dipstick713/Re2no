<script setup lang="ts">
import { Zap, Shield, Layers, Clock } from 'lucide-vue-next'

const features = [
  {
    icon: Zap,
    title: 'Fast & Efficient',
    description: 'Quickly fetch and filter Reddit posts with advanced search capabilities'
  },
  {
    icon: Shield,
    title: 'Secure Integration',
    description: 'Connect securely to your Notion workspace with OAuth authentication'
  },
  {
    icon: Layers,
    title: 'Smart Organization',
    description: 'Automatically organize posts into structured Notion databases'
  },
  {
    icon: Clock,
    title: 'Real-time Sync',
    description: 'Keep your ideas synchronized across Reddit and Notion seamlessly'
  }
]
</script>

<template>
  <section id="features" class="py-12 sm:py-20 px-4 sm:px-6">
    <div class="container mx-auto">
      <div class="text-center mb-12 sm:mb-16">
        <h2 class="text-3xl sm:text-4xl font-bold text-white mb-4">
          Master Your <span class="bg-gradient-to-r from-teal-400 to-cyan-400 bg-clip-text text-transparent">Research</span> Workflow
        </h2>
        <p class="text-lg sm:text-xl text-gray-400 max-w-2xl mx-auto px-4">
          Save time and stay organized with powerful features designed for content creators and researchers
        </p>
      </div>

      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 sm:gap-6">
        <div
          v-for="feature in features"
          :key="feature.title"
          class="bg-black/40 border border-white/10 rounded-2xl p-6 hover:border-white/20 hover:bg-black/60 transition-all backdrop-blur-xl group"
        >
          <div class="w-12 h-12 rounded-xl bg-gradient-to-br from-teal-500 to-cyan-500 flex items-center justify-center mb-4 group-hover:scale-110 transition-transform">
            <component :is="feature.icon" :size="24" class="text-white" />
          </div>
          <h3 class="text-xl font-semibold text-white mb-2">{{ feature.title }}</h3>
          <p class="text-gray-400">{{ feature.description }}</p>
        </div>
      </div>
    </div>
  </section>
</template>
