<script setup lang="ts">
import { CheckCircle, XCircle, Info } from 'lucide-vue-next'
import { useToast } from '@/composables/useToast'

const { toasts, removeToast } = useToast()

const getIconComponent = (type: string) => {
  switch (type) {
    case 'success':
      return CheckCircle
    case 'error':
      return XCircle
    default:
      return Info
  }
}

const getIconColor = (type: string) => {
  switch (type) {
    case 'success':
      return 'text-cyan-400'
    case 'error':
      return 'text-red-400'
    default:
      return 'text-blue-400'
  }
}
</script>

<template>
  <div class="fixed top-4 right-4 z-50 flex flex-col gap-2 pointer-events-none">
    <TransitionGroup name="toast">
      <div
        v-for="toast in toasts"
        :key="toast.id"
        class="bg-black/40 backdrop-blur-xl border border-white/10 rounded-xl px-6 py-4 shadow-lg min-w-[300px] max-w-md pointer-events-auto cursor-pointer hover:bg-black/50 transition-colors"
        @click="removeToast(toast.id)"
      >
        <div class="flex items-center gap-3">
          <component :is="getIconComponent(toast.type)" :class="['w-5 h-5', getIconColor(toast.type)]" />
          <p class="text-white text-sm flex-1">{{ toast.message }}</p>
        </div>
      </div>
    </TransitionGroup>
  </div>
</template>

<style scoped>
.toast-enter-active,
.toast-leave-active {
  transition: all 0.3s ease;
}

.toast-enter-from {
  opacity: 0;
  transform: translateX(100%);
}

.toast-leave-to {
  opacity: 0;
  transform: translateX(100%);
}

.toast-move {
  transition: transform 0.3s ease;
}
</style>
