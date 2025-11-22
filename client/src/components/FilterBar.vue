<script setup lang="ts">
import { ref } from 'vue'
import { Plus, X, ChevronDown } from 'lucide-vue-next'
import { useToast } from '@/composables/useToast'
import type { FilterOptions } from '@/types'

const toast = useToast()

const filters = ref<FilterOptions>({
  subreddits: ['r/learncode'],
  keyword: '',
  dateRange: 'week',
  sortBy: 'hot',
  numberOfPosts: 50,
  filterType: 'all'
})

const showFilterDropdown = ref(false)
const showDateRangeDropdown = ref(false)
const showSortByDropdown = ref(false)

const filterOptions = [
  { value: 'all', label: 'All' },
  { value: 'unsaved', label: 'Unsaved' }
] as const

const dateRangeOptions = [
  { value: 'week', label: 'This Week' },
  { value: 'month', label: 'This Month' },
  { value: 'year', label: 'This Year' },
  { value: 'all', label: 'All Time' }
] as const

const sortByOptions = [
  { value: 'hot', label: 'Hot' },
  { value: 'new', label: 'New' },
  { value: 'top', label: 'Top' },
  { value: 'rising', label: 'Rising' }
] as const

const getFilterLabel = () => {
  return filterOptions.find(opt => opt.value === filters.value.filterType)?.label || 'All'
}

const getDateRangeLabel = () => {
  return dateRangeOptions.find(opt => opt.value === filters.value.dateRange)?.label || 'This Week'
}

const getSortByLabel = () => {
  return sortByOptions.find(opt => opt.value === filters.value.sortBy)?.label || 'Hot'
}

const selectFilter = (value: 'all' | 'unsaved') => {
  filters.value.filterType = value
  showFilterDropdown.value = false
}

const selectDateRange = (value: 'week' | 'month' | 'year' | 'all') => {
  filters.value.dateRange = value
  showDateRangeDropdown.value = false
}

const selectSortBy = (value: 'hot' | 'new' | 'top' | 'rising') => {
  filters.value.sortBy = value
  showSortByDropdown.value = false
}

const newSubreddit = ref('')

const emit = defineEmits<{
  filter: [filters: FilterOptions]
  fetch: [filters: FilterOptions]
}>()

const checkSubredditExists = async (subreddit: string): Promise<boolean> => {
  try {
    // Remove 'r/' prefix for the API call
    const subredditName = subreddit.replace('r/', '')
    const response = await fetch(`https://www.reddit.com/r/${subredditName}/about.json`)

    if (!response.ok) {
      return false
    }

    const data = await response.json()
    // Check if the subreddit data exists and is not private
    return data && data.data && !data.error
  } catch (error) {
    console.error('Error checking subreddit:', error)
    return false
  }
}

const handleFetch = () => {
  emit('fetch', filters.value)
}

const addSubreddit = async () => {
  if (newSubreddit.value.trim()) {
    const subreddit = newSubreddit.value.trim().toLowerCase()
    const formatted = subreddit.startsWith('r/') ? subreddit : `r/${subreddit}`

    if (filters.value.subreddits.includes(formatted)) {
      toast.info('Subreddit already added')
      newSubreddit.value = ''
      return
    }

    // Show loading state
    const originalValue = newSubreddit.value
    newSubreddit.value = 'Checking...'

    // Check if subreddit exists
    const exists = await checkSubredditExists(formatted)

    if (!exists) {
      toast.error(`Subreddit ${formatted} does not exist or is private`)
      newSubreddit.value = originalValue
      return
    }

    filters.value.subreddits.push(formatted)
    toast.success(`Added ${formatted}`)
    newSubreddit.value = ''
  }
}

const removeSubreddit = (subreddit: string) => {
  filters.value.subreddits = filters.value.subreddits.filter(s => s !== subreddit)
}

const handleKeyPress = (event: KeyboardEvent) => {
  if (event.key === 'Enter') {
    event.preventDefault()
    addSubreddit()
  }
}

const postsInputValue = ref(filters.value.numberOfPosts.toString())

const handlePostsInput = (event: Event) => {
  const target = event.target as HTMLInputElement
  postsInputValue.value = target.value

  const value = parseInt(target.value)
  if (!isNaN(value) && value >= 10 && value <= 100) {
    filters.value.numberOfPosts = value
  }
}

const handlePostsBlur = () => {
  const value = parseInt(postsInputValue.value)

  // On blur, ensure we have a valid value
  if (isNaN(value) || value < 10) {
    filters.value.numberOfPosts = 10
  } else if (value > 100) {
    filters.value.numberOfPosts = 100
  }

  postsInputValue.value = filters.value.numberOfPosts.toString()
}
</script>

<template>
  <div class="border border-white/10 rounded-2xl p-4 sm:p-6 bg-black/40 backdrop-blur-xl">
    <!-- Subreddit Tags Section -->
    <div class="mb-4">
      <label class="block text-sm text-gray-400 mb-2 font-medium">Subreddits</label>
      <div class="flex flex-wrap gap-2 mb-3">
        <div
          v-for="subreddit in filters.subreddits"
          :key="subreddit"
          class="flex items-center gap-2 px-3 py-1.5 rounded-lg bg-cyan-500/20 border border-cyan-500/30 text-cyan-400 group"
        >
          <span class="text-sm font-medium">{{ subreddit }}</span>
          <button
            @click="removeSubreddit(subreddit)"
            class="hover:bg-cyan-500/30 rounded p-0.5 transition-colors"
          >
            <X class="w-3.5 h-3.5" />
          </button>
        </div>
      </div>
      <div class="flex flex-col sm:flex-row gap-2">
        <input
          v-model="newSubreddit"
          type="text"
          placeholder="Add subreddit (e.g., programming)"
          @keypress="handleKeyPress"
          class="flex-1 bg-white/5 border border-white/10 rounded-xl px-4 py-2.5 text-white placeholder-gray-500 focus:outline-none focus:border-cyan-500/50 focus:bg-white/10 transition-all text-sm sm:text-base"
        />
        <button
          @click="addSubreddit"
          class="px-4 py-2.5 rounded-xl bg-cyan-500/20 border border-cyan-500/30 hover:bg-cyan-500/30 text-cyan-400 font-medium transition-all flex items-center justify-center gap-2"
        >
          <Plus class="w-4 h-4" />
          Add
        </button>
      </div>
    </div>

    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 mb-6">
      <div>
        <label class="block text-sm text-gray-400 mb-2 font-medium">Keyword</label>
        <input
          v-model="filters.keyword"
          type="text"
          placeholder="programming"
          class="w-full bg-white/5 border border-white/10 rounded-xl px-4 py-2.5 text-white placeholder-gray-500 focus:outline-none focus:border-cyan-500/50 focus:bg-white/10 transition-all text-sm sm:text-base"
        />
      </div>
      <div class="relative">
        <label class="block text-sm text-gray-400 mb-2 font-medium">Date Range</label>
        <button
          @click="showDateRangeDropdown = !showDateRangeDropdown"
          class="w-full bg-black/40 border border-white/20 rounded-xl px-4 py-2.5 text-white focus:outline-none focus:border-cyan-500 transition-all text-sm sm:text-base flex items-center justify-between"
        >
          <span>{{ getDateRangeLabel() }}</span>
          <ChevronDown :size="20" :class="{ 'rotate-180': showDateRangeDropdown }" class="transition-transform" />
        </button>
        <div
          v-if="showDateRangeDropdown"
          class="absolute z-50 w-full mt-2 rounded-xl bg-neutral-900 border border-white/20 shadow-xl max-h-60 overflow-y-auto"
        >
          <button
            v-for="option in dateRangeOptions"
            :key="option.value"
            @click="selectDateRange(option.value)"
            class="w-full px-4 py-3 text-left text-white hover:bg-cyan-500/20 transition-colors first:rounded-t-xl last:rounded-b-xl"
            :class="{ 'bg-black/40': filters.dateRange === option.value }"
          >
            {{ option.label }}
          </button>
        </div>
      </div>
      <div class="relative">
        <label class="block text-sm text-gray-400 mb-2 font-medium">Sort By</label>
        <button
          @click="showSortByDropdown = !showSortByDropdown"
          class="w-full bg-black/40 border border-white/20 rounded-xl px-4 py-2.5 text-white focus:outline-none focus:border-cyan-500 transition-all text-sm sm:text-base flex items-center justify-between"
        >
          <span>{{ getSortByLabel() }}</span>
          <ChevronDown :size="20" :class="{ 'rotate-180': showSortByDropdown }" class="transition-transform" />
        </button>
        <div
          v-if="showSortByDropdown"
          class="absolute z-50 w-full mt-2 rounded-xl bg-neutral-900 border border-white/20 shadow-xl max-h-60 overflow-y-auto"
        >
          <button
            v-for="option in sortByOptions"
            :key="option.value"
            @click="selectSortBy(option.value)"
            class="w-full px-4 py-3 text-left text-white hover:bg-cyan-500/20 transition-colors first:rounded-t-xl last:rounded-b-xl"
            :class="{ 'bg-black/40': filters.sortBy === option.value }"
          >
            {{ option.label }}
          </button>
        </div>
      </div>
      <div class="relative">
        <label class="block text-sm text-gray-400 mb-2 font-medium">Filter</label>
        <button
          @click="showFilterDropdown = !showFilterDropdown"
          class="w-full bg-black/40 border border-white/20 rounded-xl px-4 py-2.5 text-white focus:outline-none focus:border-cyan-500 transition-all text-sm sm:text-base flex items-center justify-between"
        >
          <span>{{ getFilterLabel() }}</span>
          <ChevronDown :size="20" :class="{ 'rotate-180': showFilterDropdown }" class="transition-transform" />
        </button>
        <div
          v-if="showFilterDropdown"
          class="absolute z-50 w-full mt-2 rounded-xl bg-neutral-900 border border-white/20 shadow-xl max-h-60 overflow-y-auto"
        >
          <button
            v-for="option in filterOptions"
            :key="option.value"
            @click="selectFilter(option.value)"
            class="w-full px-4 py-3 text-left text-white hover:bg-cyan-500/20 transition-colors first:rounded-t-xl last:rounded-b-xl"
            :class="{ 'bg-black/40': filters.filterType === option.value }"
          >
            {{ option.label }}
          </button>
        </div>
      </div>
    </div>
    <div class="flex flex-col sm:flex-row items-stretch sm:items-center justify-between gap-4">
      <div class="flex flex-col sm:flex-row items-stretch sm:items-center gap-4 w-full sm:w-auto">
        <span class="text-gray-400 font-medium text-sm sm:text-base">Number of Posts</span>
        <div class="flex items-center gap-3">
          <input
            v-model.number="filters.numberOfPosts"
            type="range"
            min="10"
            max="100"
            @input="postsInputValue = filters.numberOfPosts.toString()"
            class="flex-1 sm:w-32 accent-cyan-500"
          />
          <input
            :value="postsInputValue"
            type="text"
            inputmode="numeric"
            pattern="[0-9]*"
            @input="handlePostsInput"
            @blur="handlePostsBlur"
            class="w-16 bg-white/5 border border-white/10 rounded-lg px-2 py-1 text-white text-center focus:outline-none focus:border-cyan-500/50 focus:bg-white/10 transition-all text-sm font-semibold [appearance:textfield] [&::-webkit-outer-spin-button]:appearance-none [&::-webkit-inner-spin-button]:appearance-none"
          />
        </div>
      </div>
      <div class="flex flex-col sm:flex-row items-stretch sm:items-center gap-3 sm:gap-4">
        <button
          @click="handleFetch"
          class="px-6 py-2.5 rounded-xl bg-gradient-to-r from-blue-500 via-cyan-400 to-emerald-400 hover:from-blue-400 hover:via-cyan-300 hover:to-emerald-300 text-white font-semibold transition-all shadow-lg shadow-cyan-500/25 text-sm sm:text-base"
        >
          Fetch Posts
        </button>
      </div>
    </div>
  </div>
</template>
