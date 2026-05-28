import './assets/main.css'

import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import i18n from './i18n'
import { registerSW } from './registerSW'
import { setupServiceProvider } from './services/serviceProvider'
import { setupAuthGuards } from './router/guards'
import VueLazyload from 'vue-lazyload'
import { lazyImageDirective } from './utils/lazyLoadUtils'

// Register service worker for PWA
registerSW()

// Entry point for the application
const app = createApp(App)

// Setup MVVM architecture with services
setupServiceProvider(app)

// Setup authentication guards for route protection
setupAuthGuards(router)

// Use router
app.use(router)

// Use i18n for internationalization
app.use(i18n)

// Use Vue Lazyload for images with v-lazy directive
app.use(VueLazyload, {
  preLoad: 1.3,
  attempt: 1,
  loading: '/images/loading-placeholder.svg',
  error: '/images/error-placeholder.svg',
})

// Register custom lazy image directive for regular img tags
app.directive('lazy-img', lazyImageDirective)

// Mount app
app.mount('#app')

/*
  Halaman diatur menggunakan vue-router di src/router/index.js
*/
