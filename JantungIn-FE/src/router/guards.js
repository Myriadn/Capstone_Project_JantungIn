import authService from '@/services/AuthService'

/**
 * Auth Guard Setup
 * Handles route protection, token validation, and role-based access control
 */

// Define protected routes by role
const routesByRole = {
  user: ['home', 'news', 'history', 'account'],
  dokter: [
    'homeAdmin',
    'diagnoseAdmin',
    'newsAdmin',
    'historyAdmin',
    'resultAdmin',
    'accountAdmin',
  ],
}

// Define all protected routes (regardless of role)
const allProtectedRoutes = [...routesByRole.user, ...routesByRole.dokter, 'diagnosisPrint']

// Define auth-only routes (need logout before access)
const authOnlyRoutes = ['login', 'register', 'admin']

/**
 * Check if route requires authentication
 * @param {string} routeName - Name of the route
 * @returns {boolean} True if route is protected
 */
function isProtectedRoute(routeName) {
  return allProtectedRoutes.includes(routeName)
}

/**
 * Check if route is auth-only (login/register/admin)
 * @param {string} routeName - Name of the route
 * @returns {boolean} True if route is auth-only
 */
function isAuthOnlyRoute(routeName) {
  return authOnlyRoutes.includes(routeName)
}

/**
 * Check if route is for specific role
 * @param {string} routeName - Name of the route
 * @param {string} role - User role ('user' or 'dokter')
 * @returns {boolean} True if user with this role can access route
 */
function canAccessRoute(routeName, role) {
  if (!role) return false

  // Check if route requires specific role
  if (routesByRole[role]?.includes(routeName)) {
    return true
  }

  return false
}

/**
 * Setup authentication guards for the router
 * @param {Router} router - Vue Router instance
 */
export function setupAuthGuards(router) {
  router.beforeEach((to, from, next) => {
    // Get user from localStorage (which includes token)
    const user = authService.getCurrentUser()
    const userRole = user?.role
    const token = user?.token

    // Debug logging
    // console.log('[Auth Guard]', {
    //   to: to.name,
    //   from: from.name,
    //   hasToken: !!token,
    //   userRole,
    //   isProtected: isProtectedRoute(to.name),
    //   isAuthOnly: isAuthOnlyRoute(to.name),
    // })

    // RULE 1: Accessing protected route without token
    if (isProtectedRoute(to.name) && !token) {
      console.warn(`[Auth Guard] No token, redirecting to /login from ${to.name}`)
      next({ name: 'login' })
      return
    }

    // RULE 2: Accessing protected route without proper role
    if (isProtectedRoute(to.name) && token && !canAccessRoute(to.name, userRole)) {
      console.warn(
        `[Auth Guard] User role "${userRole}" cannot access ${to.name}, redirecting to default dashboard`,
      )

      // Redirect to appropriate dashboard based on role
      if (userRole === 'dokter') {
        next({ name: 'homeAdmin' })
      } else {
        next({ name: 'home' })
      }
      return
    }

    // RULE 3: Already logged in, trying to access login/register/admin
    if (isAuthOnlyRoute(to.name) && token) {
      console.log('[Auth Guard] Already logged in, redirecting to dashboard')

      // If trying to access /admin (dokter login page) but already logged in as user
      if (to.name === 'admin' && userRole !== 'dokter') {
        next({ name: 'home' })
        return
      }

      // If trying to access /login or /register but already logged in
      if (to.name === 'login' || to.name === 'register') {
        // Redirect to appropriate dashboard
        next({
          name: userRole === 'dokter' ? 'homeAdmin' : 'home',
        })
        return
      }

      // Allow access to /admin if user is dokter
      if (to.name === 'admin' && userRole === 'dokter') {
        next()
        return
      }
    }

    // RULE 4: Allow 404 and print routes without restriction
    if (to.name === 'notFound' || to.meta.layout === 'print') {
      next()
      return
    }

    // RULE 5: Allow normal navigation
    next()
  })

  // Optional: Handle route errors
  router.afterEach((to, from, failure) => {
    if (failure) {
      console.error('[Router] Navigation failed:', failure)
    }
  })
}

/**
 * Get default redirect path based on user role
 * @param {Object} user - User object
 * @returns {string} Route name to redirect to
 */
export function getDefaultRedirectPath(user) {
  if (!user) return 'login'

  if (user.role === 'dokter') {
    return 'homeAdmin'
  }

  return 'home'
}

/**
 * Check if user can access specific route
 * Useful for UI logic (show/hide menu items)
 * @param {string} routeName - Route name
 * @param {Object} user - User object
 * @returns {boolean} True if user can access
 */
export function hasAccessToRoute(routeName, user) {
  if (!user || !user.role) return false

  return canAccessRoute(routeName, user.role)
}
