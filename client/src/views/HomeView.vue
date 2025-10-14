<script setup lang="ts">
import AppHeader from '@/components/AppHeader.vue'
import HeroSection from '@/components/HeroSection.vue'
import FeaturesSection from '@/components/FeaturesSection.vue'
import AppFooter from '@/components/AppFooter.vue'
</script>

<template>
  <div class="relative min-h-screen bg-black">
    <div class="absolute inset-0 bg-[url(/image.png)] bg-repeat opacity-10 pointer-events-none"></div>

    <div class="relative z-10">
      <AppHeader />

      <main>
        <HeroSection />
        <FeaturesSection />
      </main>

      <AppFooter />
    </div>
  </div>
</template>
