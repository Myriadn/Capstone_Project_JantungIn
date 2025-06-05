import apiService from './ApiService'
import { UserModel } from '@/models/UserModel'
import { useErrorHandler } from '@/utils/errorHandler'

/**
 * Authentication service for login, register, and user management
 */
class AuthService {
  /**
   * User login
   * @param {string} nik - User NIK
   * @param {string} password - User password
   * @returns {Promise<UserModel>} User model with token
   */
  async login(nik, password) {
    try {
      // Try to make the API call if online
      if (navigator.onLine) {
        const response = await apiService.post(
          '/api/v1/auth/login',
          { nik, password },
          {
            showNotifications: true,
            operationName: 'Login',
          },
        )
        console.log('Login API response:', response)

        // The API returns: { statusCode, message, data: { id, name, email, role, token } }
        // Check if the token is in response.data.data.token or response.data.token
        const responseData = response.data
        const userData = responseData.data || responseData

        if (!userData || !userData.token) {
          console.error('Invalid login response format:', responseData)
          throw this.formatError(
            new Error('Invalid response format from server: missing token'),
            'Login failed - invalid server response',
          )
        }

        console.log('Login successful, received token')

        // Set token first before saving to storage
        apiService.setToken(userData.token)

        // Save user data to local storage for offline mode
        this.saveUserToStorage({
          ...userData,
          nik: nik,
          password: password, // Only for offline login, in a real app use more secure methods
        })

        return new UserModel(userData)
      } else {
        // Offline login from cache
        return await this.offlineLogin(nik, password)
      }
    } catch (error) {
      console.error('Login error:', error)
      throw this.formatError(error, 'Login failed')
    }
  }

  /**
   * User login with email
   * @param {string} email - User email
   * @param {string} password - User password
   * @returns {Promise<UserModel>} User model with token
   */
  async loginWithEmail(email, password) {
    try {
      if (navigator.onLine) {
        const response = await apiService.post(
          '/api/v1/auth/login-email',
          { email, password },
          {
            showNotifications: true,
            operationName: 'Login dengan Email',
          },
        )
        console.log('Email login response:', response)

        // Check if the token is in response.data.data.token or response.data.token
        const responseData = response.data
        const userData = responseData.data || responseData

        if (!userData || !userData.token) {
          console.error('Invalid email login response format:', responseData)
          throw this.formatError(
            new Error('Invalid response format from server: missing token'),
            'Login failed - invalid server response',
          )
        }

        this.saveUserToStorage({
          ...userData,
          email: email,
          password: password, // Only for offline login
        })

        apiService.setToken(userData.token)

        return new UserModel(userData)
      } else {
        throw this.formatError(
          new Error('Cannot login with email while offline'),
          'No internet connection',
        )
      }
    } catch (error) {
      console.error('Email login error:', error)
      throw this.formatError(error, 'Email login failed')
    }
  }

  /**
   * Offline login from cached credentials
   * @param {string} nik - User NIK
   * @param {string} password - User password
   * @returns {Promise<UserModel>} User model
   */
  async offlineLogin(nik, password) {
    const userData = localStorage.getItem('jantungin_user')

    if (!userData) {
      throw new Error('No cached user found for offline login')
    }

    const parsedData = JSON.parse(userData)

    if (parsedData.nik === nik && parsedData.password === password) {
      return new UserModel(parsedData)
    } else {
      throw new Error('Invalid credentials for offline login')
    }
  }

  /**
   * User registration
   * @param {Object} userData - User registration data with name, nik, password, and optional email and dateOfBirth
   * @returns {Promise} Registration result
   */
  async register(userData) {
    try {
      // Validate user data
      if (!userData.nik || !userData.password || !userData.name) {
        throw this.formatError(
          new Error('Missing required fields for registration'),
          'Registration failed - missing required fields',
        )
      }

      // Validate NIK format
      if (!/^\d{16}$/.test(userData.nik)) {
        throw this.formatError(new Error('NIK harus 16 digit angka'), 'Invalid NIK format')
      }

      // Validate email if provided
      if (userData.email && !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(userData.email)) {
        throw this.formatError(new Error('Format email tidak valid'), 'Invalid email format')
      }

      console.log('Sending registration request with data:', { ...userData, password: '***' })
      const response = await apiService.post('/api/v1/auth/register', userData, {
        showNotifications: true,
        operationName: 'Registrasi',
      })
      console.log('Registration response:', response)

      // API returns: { statusCode: 201, message: "User registered successfully", data: { id, name, email, role, token } }
      // Check if the token is in response.data.data.token or response.data.token
      const responseData = response.data
      const userData2 = responseData.data || responseData

      if (!userData2 || !userData2.token) {
        console.error('Invalid registration response format:', responseData)
        throw this.formatError(
          new Error('Invalid response format from server during registration'),
          'Registration failed - invalid server response',
        )
      }

      // Set token for future API calls
      apiService.setToken(userData2.token)

      // Save user data for offline usage
      this.saveUserToStorage({
        ...userData2,
        nik: userData.nik,
        password: userData.password, // Only for offline login
      })

      return new UserModel(userData2)
    } catch (error) {
      console.error('Registration error:', error)
      throw this.formatError(error, 'Registration failed')
    }
  }

  /**
   * Doctor/admin login
   * @param {string} email - Admin email
   * @param {string} password - Admin password
   * @returns {Promise} Login result
   */
  async adminLogin(email, password) {
    try {
      // Validate email format
      if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)) {
        throw this.formatError(new Error('Format email tidak valid'), 'Invalid email format')
      }

      // Login request
      console.log('Attempting admin login...')
      const loginResponse = await apiService.post(
        '/api/v1/admin/login',
        { email, password },
        {
          showNotifications: true,
          operationName: 'Admin Login',
        },
      )

      // The login response can either be in data.data or directly in data
      const responseData = loginResponse.data?.data || loginResponse.data

      if (!responseData?.token) {
        console.error('Invalid login response:', loginResponse.data)
        throw this.formatError(
          new Error('Invalid response from server: missing token'),
          'Admin login failed - invalid server response',
        )
      }

      const { token, ...userData } = responseData
      console.log('Login successful, setting token...')

      // Set token for future API calls
      apiService.setToken(token)

      try {
        // Get full admin profile with new token
        console.log('Fetching admin profile...')
        const profileResponse = await apiService.get(
          '/api/v1/admin/profile',
          {},
          {
            showNotifications: true,
            operationName: 'Mengambil Profil Admin',
          },
        )

        // Profile data can be in data.data or directly in data
        const profileData = profileResponse.data?.data || profileResponse.data

        if (!profileData) {
          console.error('Invalid profile response:', profileResponse.data)
          throw this.formatError(
            new Error('Failed to get admin profile'),
            'Failed to retrieve admin profile',
          )
        }

        // Combine login data and profile data
        const fullUserData = {
          ...userData, // Basic data from login (id, name, email, role)
          ...profileResponse.data.data, // Additional data from profile
          token, // Store token for offline/refresh purposes
        }

        // Save complete user data to storage
        this.saveUserToStorage(fullUserData)
        console.log('Admin login complete')

        return fullUserData
      } catch (error) {
        // If profile fetch fails, still return login data
        console.error('Error fetching profile:', error)
        this.saveUserToStorage({ ...userData, token })
        return userData
      }
    } catch (error) {
      console.error('Admin login error:', error)
      throw this.formatError(error, 'Admin login failed')
    }
  }

  /**
   * Format error for consistent error handling
   * @param {Error} error - Original error
   * @param {string} defaultMessage - Default message if no specific handler
   * @returns {Error} Formatted error
   */
  formatError(error, defaultMessage = 'Authentication error occurred') {
    console.error('Auth error:', error)

    // Create a new error with the formatted message
    const formattedError = new Error(error.message || defaultMessage)

    // Copy properties from the original error
    formattedError.originalError = error
    formattedError.status = error.status
    formattedError.endpoint = error.endpoint
    formattedError.data = error.data

    // Add error category for UI handling
    if (error.isNetworkError) {
      formattedError.category = 'network'
    } else if (
      error.message &&
      (error.message.includes('NIK') ||
        error.message.includes('password') ||
        error.message.includes('email') ||
        error.message.includes('required field'))
    ) {
      formattedError.category = 'validation'
    } else if (error.status === 401 || error.status === 403) {
      formattedError.category = 'authentication'
    } else {
      formattedError.category = 'unknown'
    }

    return formattedError
  }

  /**
   * Save user data to local storage
   * @param {Object} userData - User data
   */
  saveUserToStorage(userData) {
    localStorage.setItem('jantungin_user', JSON.stringify(userData))
  }

  /**
   * Get current user from storage
   * @returns {UserModel|null} Current user
   */
  getCurrentUser() {
    try {
      const userData = localStorage.getItem('jantungin_user')
      return userData ? UserModel.fromStorage(JSON.parse(userData)) : null
    } catch (error) {
      console.error('Error getting current user:', error)
      return null
    }
  }

  /**
   * Check if user is authenticated
   * @returns {boolean} Authentication status
   */
  isAuthenticated() {
    return !!this.getCurrentUser()
  }

  /**
   * Check if current user is admin
   * @returns {boolean} Admin status
   */
  isAdmin() {
    const user = this.getCurrentUser()
    return user ? user.role === 'admin' : false
  }

  /**
   * Get user profile from API
   * @returns {Promise<Object>} User profile data
   */
  async getProfile() {
    try {
      // Check if we have auth token
      const token = apiService.getStoredToken()
      if (!token) {
        throw this.formatError(
          new Error('No authentication token found'),
          'Authentication required',
        )
      }

      const response = await apiService.get(
        '/api/v1/auth/profile',
        {},
        {
          showNotifications: true,
          operationName: 'Mengambil Profil',
        },
      )
      console.log('Profile response:', response)

      if (!response.data) {
        throw this.formatError(
          new Error('Invalid profile response from server'),
          'Failed to retrieve profile',
        )
      }

      // Handle nested response data structure
      const responseData = response.data
      const profileData = responseData.data || responseData

      // Update stored user data with latest profile info
      const currentUser = this.getCurrentUser()
      if (currentUser) {
        const updatedUserData = {
          ...JSON.parse(localStorage.getItem('jantungin_user')),
          ...profileData,
        }
        this.saveUserToStorage(updatedUserData)
      }

      return profileData
    } catch (error) {
      console.error('Error getting profile:', error)
      throw this.formatError(error, 'Failed to retrieve profile')
    }
  }

  /**
   * Update user profile
   * @param {Object} profileData - Updated profile data (name, dateOfBirth)
   * @returns {Promise<Object>} Updated profile
   */
  async updateProfile(profileData) {
    try {
      // Validate profile data
      if (profileData.name && profileData.name.trim() === '') {
        throw this.formatError(new Error('Name cannot be empty'), 'Name is required')
      }

      const response = await apiService.put('/api/v1/auth/profile', profileData, {
        showNotifications: true,
        operationName: 'Memperbarui Profil',
      })
      console.log('Profile update response:', response)

      if (!response.data) {
        throw this.formatError(
          new Error('Invalid response from server during profile update'),
          'Failed to update profile',
        )
      }

      // Handle nested response data structure
      const responseData = response.data
      const updatedProfileData = responseData.data || responseData

      // Update stored user data with latest profile info
      const currentUser = this.getCurrentUser()
      if (currentUser) {
        const updatedUserData = {
          ...JSON.parse(localStorage.getItem('jantungin_user')),
          ...updatedProfileData,
        }
        this.saveUserToStorage(updatedUserData)
      }

      return updatedProfileData
    } catch (error) {
      console.error('Error updating profile:', error)
      throw this.formatError(error, 'Failed to update profile')
    }
  }

  /**
   * Logout current user
   */
  logout() {
    localStorage.removeItem('jantungin_user')
    apiService.setToken(null)
  }
}

// Create singleton instance
const authService = new AuthService()

export default authService
