<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import authService from '@/services/AuthService'
import { useErrorHandler } from '@/utils/errorHandler'

const router = useRouter()
const { t } = useI18n()
const { getErrorMessage } = useErrorHandler()

const username = ref('')
const email = ref('') // Added email field
const password = ref('')
const confirmPassword = ref('')
const dateOfBirth = ref('') // Added date of birth field
const loading = ref(false)
const showPassword = ref(false)
const showConfirmPassword = ref(false)
const errorMessage = ref('')

const navigateToLogin = () => {
  router.push('/')
}

// Function to navigate to doctor login page
const navigateToDoctorLogin = () => {
  router.push('/admin')
}

const togglePasswordVisibility = () => {
  showPassword.value = !showPassword.value
}

const toggleConfirmPasswordVisibility = () => {
  showConfirmPassword.value = !showConfirmPassword.value
}

const validateForm = () => {
  if (password.value !== confirmPassword.value) {
    errorMessage.value = t('errors.validation.passwordsDoNotMatch')
    return false
  }

  if (password.value.length < 6) {
    errorMessage.value = t('errors.validation.passwordLength')
    return false
  }

  // Validate username
  const trimmedUsername = username.value.trim()
  const usernameRegex = /^[A-Za-z0-9]+([._][A-Za-z0-9]+)*$/

  if (trimmedUsername.length < 3 || trimmedUsername.length > 16) {
    errorMessage.value = t('errors.validation.usernameFormat')
    return false
  }

  if (trimmedUsername !== username.value) {
    errorMessage.value = t('errors.validation.usernameFormat')
    return false
  }

  if (!usernameRegex.test(trimmedUsername)) {
    errorMessage.value = t('errors.validation.usernameFormat')
    return false
  }

  // Validate email format if provided
  if (email.value && !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email.value)) {
    errorMessage.value = t('errors.validation.emailFormat')
    return false
  }

  return true
}

const handleRegister = async () => {
  errorMessage.value = ''
  if (!validateForm()) {
    return
  }
  loading.value = true
  try {
    // Panggil backend register
    console.log('Attempting to register with:', {
      username: username.value,
      name: username.value,
      email: email.value || '[not provided]',
      dateOfBirth: dateOfBirth.value || '[not provided]',
    })

    const registerData = await authService.register({
      username: username.value,
      name: username.value,
      email: email.value || undefined,
      password: password.value,
      dateOfBirth: dateOfBirth.value || undefined,
    })

    console.log('Registration successful:', registerData)

    // Store credentials for auto-fill on login page
    sessionStorage.setItem('autoFillUsername', username.value)
    sessionStorage.setItem('autoFillPassword', password.value)
    sessionStorage.setItem('autoFillEmail', email.value || '')

    // Redirect ke login. Backend does not return token on register.
    // User must login with OTP after registration.
    router.push('/')
  } catch (err) {
    console.error('Registration failed:', err)

    // Use error handler to get consistent localized error message
    errorMessage.value = getErrorMessage(err)
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="register-container">
    <div class="register-content">
      <!-- Left side with illustration -->
      <div class="register-illustration">
        <div class="illustration-content">
          <!-- Heart Icon Logo -->
          <div class="logo-container" @click="navigateToDoctorLogin" style="cursor: pointer">
            <svg
              class="heart-logo"
              viewBox="0 0 24 24"
              fill="none"
              xmlns="http://www.w3.org/2000/svg"
            >
              <path
                d="M12 21.35l-1.45-1.32C5.4 15.36 2 12.28 2 8.5 2 5.42 4.42 3 7.5 3c1.74 0 3.41.81 4.5 2.09C13.09 3.81 14.76 3 16.5 3 19.58 3 22 5.42 22 8.5c0 3.78-3.4 6.86-8.55 11.54L12 21.35z"
                fill="currentColor"
              />
            </svg>
          </div>
          <h1 class="illustration-title">JantungIn</h1>
          <p class="illustration-tagline">Your Heart Health Companion</p>

          <!-- Benefits Section -->
          <div class="benefits-card">
            <div class="benefits-header">
              <svg
                xmlns="http://www.w3.org/2000/svg"
                class="benefits-icon"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z"
                />
              </svg>
              <span>Benefits of Registration</span>
            </div>
            <ul class="benefits-list">
              <li>Track your heart health history</li>
              <li>Receive personalized health recommendations</li>
              <li>Consult with heart specialists</li>
              <li>Access your data even when offline</li>
            </ul>
          </div>

          <div class="features-section">
            <div class="feature-item">
              <svg
                xmlns="http://www.w3.org/2000/svg"
                class="feature-icon"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M9 12h6m-6 4h6m-6-8h6M5 8h14a2 2 0 012 2v10a2 2 0 01-2 2H5a2 2 0 01-2-2V10a2 2 0 012-2z"
                />
              </svg>
              <span>Health Records</span>
            </div>
            <div class="feature-item">
              <svg
                xmlns="http://www.w3.org/2000/svg"
                class="feature-icon"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"
                />
              </svg>
              <span>Appointment Scheduling</span>
            </div>
            <div class="feature-item">
              <svg
                xmlns="http://www.w3.org/2000/svg"
                class="feature-icon"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9"
                />
              </svg>
              <span>Health Alerts</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Right side with registration form -->
      <div class="register-form-section">
        <div class="register-form-container">
          <h2 class="register-title">Create Your Account</h2>
          <p class="register-subtitle">Take the first step towards better heart health</p>

          <!-- Error Message -->
          <div v-if="errorMessage" class="error-message">
            <svg
              xmlns="http://www.w3.org/2000/svg"
              class="error-icon"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"
              />
            </svg>
            <span>{{ errorMessage }}</span>
          </div>

          <form @submit.prevent="handleRegister" class="register-form">
            <div class="form-group">
              <label for="username">Username</label>
              <div class="input-container">
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  class="input-icon"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                >
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"
                  />
                </svg>
                <input
                  id="username"
                  v-model="username"
                  type="text"
                  class="form-control"
                  placeholder="Choose a username"
                  required
                />
              </div>
            </div>

            <!-- Optional Email Field -->
            <div class="form-group">
              <label for="email">Email</label>
              <div class="input-container">
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  class="input-icon"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                >
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"
                  />
                </svg>
                <input
                  id="email"
                  v-model="email"
                  type="email"
                  class="form-control"
                  placeholder="Enter your email"
                />
              </div>
            </div>

            <!-- Optional Date of Birth Field -->
            <div class="form-group">
              <label for="dateOfBirth">Date of Birth</label>
              <div class="input-container">
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  class="input-icon"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                >
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"
                  />
                </svg>
                <input
                  id="dateOfBirth"
                  v-model="dateOfBirth"
                  type="date"
                  class="form-control"
                  placeholder="YYYY-MM-DD"
                />
              </div>
            </div>

            <div class="form-group">
              <label for="password">Password</label>
              <div class="input-container">
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  class="input-icon"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                >
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"
                  />
                </svg>
                <input
                  id="password"
                  v-model="password"
                  :type="showPassword ? 'text' : 'password'"
                  class="form-control"
                  placeholder="Create a password"
                  required
                />
                <button type="button" class="toggle-password" @click="togglePasswordVisibility">
                  <svg
                    v-if="showPassword"
                    xmlns="http://www.w3.org/2000/svg"
                    class="h-5 w-5"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                  >
                    <path
                      stroke-linecap="round"
                      stroke-linejoin="round"
                      stroke-width="2"
                      d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21"
                    />
                  </svg>
                  <svg
                    v-else
                    xmlns="http://www.w3.org/2000/svg"
                    class="h-5 w-5"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                  >
                    <path
                      stroke-linecap="round"
                      stroke-linejoin="round"
                      stroke-width="2"
                      d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"
                    />
                    <path
                      stroke-linecap="round"
                      stroke-linejoin="round"
                      stroke-width="2"
                      d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"
                    />
                  </svg>
                </button>
              </div>
              <span class="password-hint">Password must be at least 6 characters</span>
            </div>

            <div class="form-group">
              <label for="confirm-password">Confirm Password</label>
              <div class="input-container">
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  class="input-icon"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                >
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z"
                  />
                </svg>
                <input
                  id="confirm-password"
                  v-model="confirmPassword"
                  :type="showConfirmPassword ? 'text' : 'password'"
                  class="form-control"
                  placeholder="Confirm your password"
                  required
                />
                <button
                  type="button"
                  class="toggle-password"
                  @click="toggleConfirmPasswordVisibility"
                >
                  <svg
                    v-if="showConfirmPassword"
                    xmlns="http://www.w3.org/2000/svg"
                    class="h-5 w-5"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                  >
                    <path
                      stroke-linecap="round"
                      stroke-linejoin="round"
                      stroke-width="2"
                      d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21"
                    />
                  </svg>
                  <svg
                    v-else
                    xmlns="http://www.w3.org/2000/svg"
                    class="h-5 w-5"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                  >
                    <path
                      stroke-linecap="round"
                      stroke-linejoin="round"
                      stroke-width="2"
                      d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"
                    />
                    <path
                      stroke-linecap="round"
                      stroke-linejoin="round"
                      stroke-width="2"
                      d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"
                    />
                  </svg>
                </button>
              </div>
            </div>

            <div class="terms-container">
              <label class="terms-agreement">
                <input type="checkbox" required />
                <span
                  >I agree to the <a href="#" class="terms-link">Terms of Service</a> and
                  <a href="#" class="terms-link">Privacy Policy</a></span
                >
              </label>
            </div>

            <button type="submit" class="btn-primary" :disabled="loading">
              <span v-if="loading" class="loading-spinner"></span>
              <span>Create Account</span>
            </button>

            <div class="register-divider">
              <span>Already have an account?</span>
            </div>

            <button type="button" @click="navigateToLogin" class="btn-secondary">Sign In</button>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.register-container {
  min-height: 100vh;
  width: 100%;
  background: linear-gradient(135deg, #f0f9ff 0%, #e1f5fe 100%);
  font-family: 'Inter', 'Segoe UI', sans-serif;
}

.register-content {
  display: flex;
  min-height: 100vh;
  flex-direction: column;
}

/* Left side - illustration */
.register-illustration {
  display: none; /* Hidden on mobile */
  background: linear-gradient(135deg, #ec4899 0%, #8b5cf6 50%, #3b82f6 100%);
  color: white;
  padding: 2rem;
  position: relative;
  overflow: hidden;
}

.illustration-content {
  position: relative;
  z-index: 10;
  display: flex;
  flex-direction: column;
  justify-content: center;
  height: 100%;
  max-width: 400px;
  margin: 0 auto;
}

.logo-container {
  display: flex;
  align-items: center;
  margin-bottom: 1rem;
  position: relative;
}

.heart-logo {
  width: 3rem;
  height: 3rem;
  color: white;
  filter: drop-shadow(0 0 8px rgba(255, 255, 255, 0.5));
  transition:
    transform 0.3s ease,
    filter 0.3s ease;
}

.logo-container:hover .heart-logo {
  transform: scale(1.1);
  filter: drop-shadow(0 0 12px rgba(255, 255, 255, 0.8));
}

.logo-container::after {
  content: 'Login as Doctor';
  position: absolute;
  top: 100%;
  left: 50%;
  transform: translateX(-50%);
  background-color: rgba(0, 0, 0, 0.8);
  color: white;
  padding: 0.25rem 0.75rem;
  border-radius: 0.25rem;
  font-size: 0.75rem;
  opacity: 0;
  pointer-events: none;
  transition:
    opacity 0.3s ease,
    transform 0.3s ease;
  white-space: nowrap;
}

.logo-container:hover::after {
  opacity: 1;
  transform: translateX(-50%) translateY(5px);
}

.illustration-title {
  font-size: 2.5rem;
  font-weight: 700;
  margin-bottom: 0.5rem;
}

.illustration-tagline {
  font-size: 1.25rem;
  margin-bottom: 2rem;
  opacity: 0.9;
}

/* Benefits Card */
.benefits-card {
  background-color: rgba(255, 255, 255, 0.1);
  border-radius: 1rem;
  padding: 1.5rem;
  margin-bottom: 2rem;
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.benefits-header {
  display: flex;
  align-items: center;
  margin-bottom: 0.75rem;
  font-weight: 600;
}

.benefits-icon {
  width: 1.25rem;
  height: 1.25rem;
  margin-right: 0.5rem;
}

.benefits-list {
  padding-left: 1.5rem;
  margin-top: 0.5rem;
  font-size: 0.95rem;
  line-height: 1.7;
}

.benefits-list li {
  margin-bottom: 0.5rem;
}

.features-section {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 1rem;
  margin-top: auto;
}

.feature-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  font-size: 0.875rem;
}

.feature-icon {
  width: 1.75rem;
  height: 1.75rem;
  margin-bottom: 0.5rem;
}

/* Right side - registration form */
.register-form-section {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 1.5rem;
}

.register-form-container {
  width: 100%;
  max-width: 450px;
  background: white;
  border-radius: 1rem;
  padding: 2rem;
  box-shadow: 0 10px 25px rgba(0, 0, 0, 0.05);
}

.register-title {
  font-size: 1.75rem;
  font-weight: 700;
  color: #1e293b;
  text-align: center;
  margin-bottom: 0.5rem;
}

.register-subtitle {
  color: #64748b;
  font-size: 0.95rem;
  text-align: center;
  margin-bottom: 2rem;
}

.error-message {
  background-color: #f87171;
  color: white;
  padding: 0.75rem;
  border-radius: 0.5rem;
  margin-bottom: 1.5rem;
  text-align: center;
  font-size: 0.875rem;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
}

.error-icon {
  width: 1.25rem;
  height: 1.25rem;
}

.register-form {
  display: flex;
  flex-direction: column;
}

.form-group {
  margin-bottom: 1.5rem;
}

.form-group label {
  display: block;
  font-size: 0.875rem;
  margin-bottom: 0.5rem;
  font-weight: 500;
  color: #334155;
}

.input-container {
  position: relative;
  display: flex;
  align-items: center;
}

.input-icon {
  position: absolute;
  left: 1rem;
  height: 1.25rem;
  width: 1.25rem;
  color: #94a3b8;
}

.form-control {
  width: 100%;
  padding: 0.75rem 1rem 0.75rem 2.75rem;
  border: 1px solid #e2e8f0;
  border-radius: 0.5rem;
  font-size: 1rem;
  transition: all 0.2s;
  background-color: #f8fafc;
}

.form-control:focus {
  outline: none;
  border-color: #3b82f6;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.25);
  background-color: white;
}

.password-hint {
  display: block;
  margin-top: 0.5rem;
  font-size: 0.75rem;
  color: #64748b;
  padding-left: 0.5rem;
}

.toggle-password {
  position: absolute;
  right: 1rem;
  background: transparent;
  border: none;
  color: #94a3b8;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
}

.toggle-password:hover {
  color: #64748b;
}

.terms-container {
  margin-bottom: 1.5rem;
}

.terms-agreement {
  display: flex;
  align-items: flex-start;
  color: #334155;
  font-size: 0.875rem;
}

.terms-agreement input {
  margin-right: 0.5rem;
  margin-top: 0.25rem;
}

.terms-link {
  color: #3b82f6;
  text-decoration: none;
}

.terms-link:hover {
  text-decoration: underline;
}

.btn-primary {
  width: 100%;
  padding: 0.875rem;
  background: linear-gradient(135deg, #3b82f6 0%, #1d4ed8 100%);
  color: white;
  border: none;
  border-radius: 0.5rem;
  font-size: 1rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 1.5rem;
}

.btn-primary:hover:not(:disabled) {
  background: linear-gradient(135deg, #2563eb 0%, #1e40af 100%);
  transform: translateY(-1px);
  box-shadow:
    0 4px 6px -1px rgba(0, 0, 0, 0.1),
    0 2px 4px -1px rgba(0, 0, 0, 0.06);
}

.btn-primary:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}

.loading-spinner {
  width: 1.25rem;
  height: 1.25rem;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-top-color: white;
  border-radius: 50%;
  margin-right: 0.5rem;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

.register-divider {
  position: relative;
  text-align: center;
  margin-bottom: 1.5rem;
}

.register-divider::before {
  content: '';
  position: absolute;
  top: 50%;
  left: 0;
  right: 0;
  height: 1px;
  background-color: #e2e8f0;
}

.register-divider span {
  position: relative;
  background-color: white;
  padding: 0 1rem;
  color: #64748b;
  font-size: 0.875rem;
}

.btn-secondary {
  width: 100%;
  padding: 0.875rem;
  background-color: transparent;
  color: #1e40af;
  border: 1px solid #3b82f6;
  border-radius: 0.5rem;
  font-size: 1rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-secondary:hover {
  background-color: rgba(59, 130, 246, 0.05);
  transform: translateY(-1px);
}

/* Responsive design */
@media (min-width: 1024px) {
  .register-content {
    flex-direction: row;
  }

  .register-illustration {
    display: flex;
    flex: 1;
  }

  .register-form-section {
    flex: 1;
  }
}

/* Animations for illustration background */
.register-illustration::before {
  content: '';
  position: absolute;
  top: -50%;
  left: -50%;
  width: 200%;
  height: 200%;
  background: radial-gradient(circle, rgba(255, 255, 255, 0.1) 0%, transparent 70%);
  animation: rotate 30s linear infinite;
}

@keyframes rotate {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}
</style>
