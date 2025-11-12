<script setup lang="ts">
import { BookOpen, Github, Menu, X } from 'lucide-vue-next'
import { useRouter, useRoute } from 'vue-router'
import { ref, onMounted } from 'vue'
import { getNotionAuthUrl, getCurrentUser, logout } from '@/lib/api'

const router = useRouter()
const route = useRoute()
const mobileMenuOpen = ref(false)
const isAuthenticated = ref(false)
const isLoading = ref(false)

onMounted(async () => {
  try {
    const user = await getCurrentUser()
    isAuthenticated.value = !!user
  } catch {
    isAuthenticated.value = false
  }

  if (route.query.auth === 'success') {
    isAuthenticated.value = true
  }
})

const handleGetStarted = async () => {
  if (isAuthenticated.value) {
    router.push('/dashboard')
  } else {
    await handleLogin()
  }
  mobileMenuOpen.value = false
}

const handleLogin = async () => {
  try {
    isLoading.value = true
    const authUrl = await getNotionAuthUrl()
    window.location.href = authUrl
  } catch (err) {
    console.error('Failed to get auth URL:', err)
    alert('Failed to initiate login. Please try again.')
  } finally {
    isLoading.value = false
  }
}

const handleLogout = async () => {
  try {
    await logout()
    isAuthenticated.value = false
    router.push('/')
  } catch (err) {
    console.error('Failed to logout:', err)
  }
  mobileMenuOpen.value = false
}

const toggleMobileMenu = () => {
  mobileMenuOpen.value = !mobileMenuOpen.value
}
</script>

<template>
  <header class="sticky top-4 z-50 px-4 sm:px-6">
    <div class="max-w-6xl mx-auto border border-white/10 bg-black/10 backdrop-blur-xl rounded-2xl">
      <nav class="flex items-center justify-between px-4 sm:px-6 py-4 relative">
        <router-link to="/" class="flex items-center gap-2 flex-shrink-0">
          <div class="w-8 h-8 rounded-lg bg-gradient-to-br from-blue-500 via-cyan-400 to-emerald-400 flex items-center justify-center">
            <BookOpen :size="18" class="text-white" />
          </div>
          <span class="text-xl font-semibold text-white">Re2no</span>
        </router-link>

        <div class="hidden md:flex items-center gap-8 absolute left-1/2 transform -translate-x-1/2">
          <router-link to="/" class="text-sm text-gray-300 hover:text-white transition-colors">Home</router-link>
          <router-link to="/dashboard" class="text-sm text-gray-300 hover:text-white transition-colors">Dashboard</router-link>
        </div>

        <div class="hidden md:flex items-center gap-3 flex-shrink-0">
          <a
            href="https://github.com/Dipstick713/Re2no"
            target="_blank"
            rel="noopener noreferrer"
            class="flex items-center gap-2 px-4 py-2 rounded-lg border border-white/20 bg-white/5 hover:bg-white/10 text-white text-sm font-medium transition-all"
          >
            <Github :size="16" />
            <span>GitHub</span>
          </a>

          <button
            v-if="!isAuthenticated"
            @click="handleGetStarted"
            :disabled="isLoading"
            class="px-6 py-2 rounded-lg bg-gradient-to-r from-blue-500 via-cyan-400 to-emerald-400 hover:from-blue-400 hover:via-cyan-300 hover:to-emerald-300 text-white font-semibold transition-all shadow-lg shadow-cyan-500/25 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {{ isLoading ? 'Loading...' : 'Get Started' }}
          </button>

          <button
            v-else
            @click="handleLogout"
            class="px-6 py-2 rounded-lg border border-white/20 bg-white/5 hover:bg-white/10 text-white font-semibold transition-all"
          >
            Logout
          </button>
        </div>

        <button
          @click="toggleMobileMenu"
          class="md:hidden p-2 rounded-lg border border-white/20 bg-white/5 hover:bg-white/10 text-white transition-all"
        >
          <Menu v-if="!mobileMenuOpen" :size="20" />
          <X v-else :size="20" />
        </button>
      </nav>

      <div
        v-if="mobileMenuOpen"
        class="md:hidden border-t border-white/10 px-4 py-4 space-y-3"
      >
        <router-link to="/" @click="mobileMenuOpen = false" class="block text-sm text-gray-300 hover:text-white transition-colors py-2">Home</router-link>
        <router-link to="/dashboard" @click="mobileMenuOpen = false" class="block text-sm text-gray-300 hover:text-white transition-colors py-2">Dashboard</router-link>
        <a
          href="https://github.com/Dipstick713/Re2no"
          target="_blank"
          rel="noopener noreferrer"
          class="flex items-center gap-2 px-4 py-2 rounded-lg border border-white/20 bg-white/5 hover:bg-white/10 text-white text-sm font-medium transition-all"
        >
          <Github :size="16" />
          GitHub
        </a>
        <button
          v-if="!isAuthenticated"
          @click="handleGetStarted"
          :disabled="isLoading"
          class="w-full px-6 py-2 rounded-lg bg-gradient-to-r from-blue-500 via-cyan-400 to-emerald-400 hover:from-blue-400 hover:via-cyan-300 hover:to-emerald-300 text-white font-semibold transition-all shadow-lg shadow-cyan-500/25 disabled:opacity-50 disabled:cursor-not-allowed"
        >
          {{ isLoading ? 'Loading...' : 'Get Started' }}
        </button>
        <button
          v-else
          @click="handleLogout"
          class="w-full px-6 py-2 rounded-lg border border-white/20 bg-white/5 hover:bg-white/10 text-white font-semibold transition-all"
        >
          Logout
        </button>
      </div>
    </div>
  </header>
</template>
