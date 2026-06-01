<template>
  <div class="flex items-center justify-center min-h-screen bg-gradient-to-br from-gray-50 to-gray-100 px-4">
    <div class="text-center max-w-md w-full">
      <!-- 404 Icon/Illustration -->
      <div class="mb-8">
        <div class="text-8xl font-bold text-gray-300 mb-2">404</div>
        <div class="text-2xl font-semibold text-gray-800 mb-2">Page Not Found</div>
        <p class="text-gray-600 mb-6">
          Sorry, the page you're looking for doesn't exist or has been moved.
        </p>
      </div>

      <!-- Error Description -->
      <div class="bg-white rounded-lg shadow-md p-6 mb-6">
        <p class="text-gray-700 text-sm mb-4">
          The URL you're trying to access is not valid. Please check the address and try again.
        </p>
        <div class="bg-gray-50 rounded p-3 text-left text-xs text-gray-600 break-all font-mono">
          {{ currentPath }}
        </div>
      </div>

      <!-- Action Buttons -->
      <div class="flex flex-col gap-3">
        <button
          @click="goHome"
          class="w-full bg-blue-600 hover:bg-blue-700 text-white font-semibold py-3 px-4 rounded-lg transition-colors duration-200"
        >
          Go to Home
        </button>
        <button
          @click="goBack"
          class="w-full bg-white border-2 border-gray-300 hover:bg-gray-50 text-gray-700 font-semibold py-3 px-4 rounded-lg transition-colors duration-200"
        >
          Go Back
        </button>
      </div>

      <!-- Quick Links -->
      <div class="mt-8 pt-6 border-t border-gray-200">
        <p class="text-sm text-gray-600 mb-4">Quick Links:</p>
        <div class="flex flex-wrap justify-center gap-2">
          <router-link
            to="/home"
            class="text-xs bg-blue-50 hover:bg-blue-100 text-blue-600 px-3 py-2 rounded transition-colors"
          >
            Home
          </router-link>
          <router-link
            to="/news"
            class="text-xs bg-blue-50 hover:bg-blue-100 text-blue-600 px-3 py-2 rounded transition-colors"
          >
            News
          </router-link>
          <router-link
            to="/history"
            class="text-xs bg-blue-50 hover:bg-blue-100 text-blue-600 px-3 py-2 rounded transition-colors"
          >
            History
          </router-link>
          <router-link
            to="/account"
            class="text-xs bg-blue-50 hover:bg-blue-100 text-blue-600 px-3 py-2 rounded transition-colors"
          >
            Account
          </router-link>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'

const router = useRouter()
const route = useRoute()

const currentPath = computed(() => {
  return window.location.pathname || route.path
})

const goHome = () => {
  router.push('/')
}

const goBack = () => {
  window.history.back()
}
</script>

<style scoped>
/* Smooth transitions for buttons */
button {
  transition:
    background-color 0.2s ease,
    border-color 0.2s ease;
}
</style>
