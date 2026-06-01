import { createRouter, createWebHistory } from 'vue-router'

// Using lazy loading for all components
// This will improve initial loading performance by loading routes only when needed

const routes = [
  {
    path: '/print/diagnosis/:id',
    name: 'diagnosisPrint',
    component: () => import('@/Page/Print/DiagnosisPrint.vue'),
    meta: { layout: 'print' },
  },
  {
    path: '/',
    name: 'login',
    component: () => import('@/Page/User/auth/Login-user.vue'),
    meta: { transition: 'fade' },
  },
  {
    path: '/register',
    name: 'register',
    component: () => import('@/Page/User/auth/Register-user.vue'),
    meta: { transition: 'fade' },
  },
  {
    path: '/home',
    name: 'home',
    component: () => import('@/Page/User/HomeUserPage.vue'),
    meta: { transition: 'fade' },
  },
  {
    path: '/news',
    name: 'news',
    component: () => import('@/Page/User/NewsPage.vue'),
    meta: { transition: 'fade' },
  },
  {
    path: '/history',
    name: 'history',
    component: () => import('@/Page/User/HistoryPage.vue'),
    meta: { transition: 'fade' },
  },
  {
    path: '/account',
    name: 'account',
    component: () => import('@/Page/User/AccountPage.vue'),
    meta: { transition: 'fade' },
  },
  {
    path: '/admin',
    name: 'admin',
    component: () => import('@/Page/Admin/auth/Login-doctor.vue'),
    meta: { transition: 'fade' },
  },
  {
    path: '/home-admin',
    name: 'homeAdmin',
    component: () => import('@/Page/Admin/HomeAdminPage.vue'),
    meta: { transition: 'fade' },
  },
  {
    path: '/diagnose-admin',
    name: 'diagnoseAdmin',
    component: () => import('@/Page/Admin/DiagnoseAdminPage.vue'),
    meta: { transition: 'fade' },
  },
  {
    path: '/news-admin',
    name: 'newsAdmin',
    component: () => import('@/Page/Admin/NewsAdminPage.vue'),
    meta: { transition: 'fade' },
  },
  {
    path: '/history-admin',
    name: 'historyAdmin',
    component: () => import('@/Page/Admin/HistoryAdminPage.vue'),
    meta: { transition: 'fade' },
  },
  {
    path: '/result-admin',
    name: 'resultAdmin',
    component: () => import('@/Page/Admin/ResultAdminPage.vue'),
    meta: { transition: 'fade' },
  },
  {
    path: '/account-admin',
    name: 'accountAdmin',
    component: () => import('@/Page/Admin/AccountAdminPage.vue'),
    meta: { transition: 'fade' },
  },
  // Catch-all route for 404 - must be at the end
  {
    path: '/:pathMatch(.*)*',
    name: 'notFound',
    component: () => import('@/Page/NotFoundPage.vue'),
    meta: { layout: 'notfound', transition: 'fade' },
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior(to, from, savedPosition) {
    // Return to saved position if available (when using browser back/forward buttons)
    if (savedPosition) {
      return savedPosition
    }
    // Scroll to top for all other navigation
    return { top: 0 }
  },
})

export default router
