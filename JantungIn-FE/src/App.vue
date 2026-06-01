<script setup>
import Navbar from '@/components/Navbar.vue'
import NavbarAdmin from '@/components/Navbar-admin.vue'
import OfflinePage from '@/components/OfflinePage.vue'
import RefreshApp from '@/components/RefreshApp.vue'
import PWAStatus from '@/components/PWAStatus.vue'
import AdPopupComponent from '@/components/AdPopupComponent.vue'
import { useRoute } from 'vue-router'
import { computed, ref, onMounted, onUnmounted, onErrorCaptured } from 'vue'
import { preloadCriticalImages } from '@/utils/lazyLoadUtils'

const route = useRoute()
const isOnline = ref(navigator.onLine)
const offlineAlert = ref(false)
const hasError = ref(false)
const errorMessage = ref('')

// Preload critical images for performance
preloadCriticalImages([
  // Add paths to critical images here
])

// Check if the current route corresponds to an admin page
const isAdminPage = computed(() => {
  // Admin route names - keep this list in sync with router/index.js
  const adminRouteNames = [
    'homeAdmin', // HomeAdminPage
    'diagnoseAdmin', // DiagnoseAdminPage
    'historyAdmin', // HistoryAdminPage
    'newsAdmin', // NewsAdminPage
    'resultAdmin', // ResultAdminPage
    'accountAdmin', // AccountAdminPage
  ]

  // Check if current route name is in the admin routes list
  return adminRouteNames.includes(route.name)
})

// Check if the current route is a 404 page
const isNotFoundPage = computed(() => {
  return route.name === 'notFound'
})

// Network status handlers
const handleOnline = () => {
  isOnline.value = true
  offlineAlert.value = false
}

const handleOffline = () => {
  isOnline.value = false
  offlineAlert.value = true
  // Hide alert after 5 seconds
  setTimeout(() => {
    offlineAlert.value = false
  }, 5000)
}

// Error handling
onErrorCaptured((err, instance, info) => {
  console.error('App Error:', err)
  console.error('Component:', instance)
  console.error('Info:', info)

  hasError.value = true
  errorMessage.value = err.message || 'An unexpected error occurred.'

  // Prevent error from propagating further
  return false
})

// Function to reset error state
const resetErrorState = () => {
  hasError.value = false
  errorMessage.value = ''
}

onMounted(() => {
  window.addEventListener('online', handleOnline)
  window.addEventListener('offline', handleOffline)

  // Preload critical images for better performance
  preloadCriticalImages([
    '/images/loading-placeholder.svg',
    '/images/error-placeholder.svg',
    '/src/assets/images/diagnose-hero.jpg',
    '/src/assets/images/history-hero.jpg',
    '/images/heart1.jpg',
  ])
})

onUnmounted(() => {
  window.removeEventListener('online', handleOnline)
  window.removeEventListener('offline', handleOffline)
})
</script>

<template>
  <!-- Application Shell Architecture -->
  <div class="flex flex-col min-h-screen">
    <!-- Static Shell Content (Header) -->
    <header v-if="route.meta.layout !== 'print' && !isNotFoundPage">
      <!-- Show admin navbar for admin pages -->
      <NavbarAdmin v-if="isAdminPage && route.name !== 'admin'" />
      <!-- Show regular navbar for user pages -->
      <Navbar
        v-if="
          !isAdminPage &&
          route.name !== 'login' &&
          route.name !== 'register' &&
          route.name !== 'admin'
        "
      />
    </header>

    <!-- Offline Alert -->
    <div
      v-if="offlineAlert"
      class="bg-yellow-500 text-white text-center py-2 px-4 fixed top-16 left-0 right-0 z-50 transition-all duration-300"
    >
      You're currently offline. Some features may be limited.
    </div>

    <!-- Dynamic Content Area -->
    <main class="flex-grow">
      <!-- Display offline page when offline -->
      <OfflinePage v-if="!isOnline" />

      <!-- Display error state when there's an error -->
      <div
        v-else-if="hasError"
        class="error-container flex flex-col items-center justify-center h-full p-6"
      >
        <div class="bg-red-50 border border-red-200 rounded-lg p-6 max-w-lg w-full">
          <h2 class="text-xl font-semibold text-red-700 mb-3">Something went wrong</h2>
          <p class="text-gray-700 mb-4">{{ errorMessage }}</p>
          <p class="text-sm text-gray-600 mb-4">
            This error has been logged and we're working on fixing it.
          </p>
          <div class="flex gap-3">
            <button
              @click="resetErrorState"
              class="bg-red-600 hover:bg-red-700 text-white px-4 py-2 rounded-md transition-colors"
            >
              Try Again
            </button>
            <button
              @click="$router.push('/')"
              class="bg-white border border-gray-300 hover:bg-gray-50 text-gray-700 px-4 py-2 rounded-md transition-colors"
            >
              Go to Home
            </button>
          </div>
        </div>
      </div>

      <!-- Display normal content when online with page transition -->
      <transition :name="route.meta.transition || 'fade'" mode="out-in" v-else>
        <router-view v-slot="{ Component }">
          <component :is="Component" />
        </router-view>
      </transition>
    </main>

    <!-- Static Shell Content (Footer) -->
    <footer v-if="route.meta.layout !== 'print' && !isNotFoundPage">
      <!-- Footer content if needed -->
    </footer>

    <!-- PWA Components -->
    <RefreshApp v-if="route.meta.layout !== 'print' && !isNotFoundPage" />
    <PWAStatus v-if="route.meta.layout !== 'print' && !isNotFoundPage" />

    <!-- Ad Popup Component -->
    <AdPopupComponent v-if="route.meta.layout !== 'print' && !isNotFoundPage" />
  </div>
</template>

<style>
html,
body,
#app {
  height: 100%;
  margin: 0;
  padding: 0;
}

#app {
  display: flex;
  flex-direction: column;
}

/* Ensure pages have proper spacing from navbar */
main {
  min-height: calc(100vh - 4rem); /* 4rem = 64px (height of navbar) */
}

/* Page transition animations */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

/* Slide transitions */
.slide-enter-active,
.slide-leave-active {
  transition:
    transform 0.3s ease,
    opacity 0.3s ease;
}

.slide-enter-from {
  opacity: 0;
  transform: translateX(30px);
}

.slide-leave-to {
  opacity: 0;
  transform: translateX(-30px);
}

/* Lazy loading styles */
.lazy-load {
  opacity: 0;
  transition: opacity 0.5s ease;
}

.lazy-load.loaded {
  opacity: 1;
}

/* Error container */
.error-container {
  min-height: calc(100vh - 8rem);
}
</style>
