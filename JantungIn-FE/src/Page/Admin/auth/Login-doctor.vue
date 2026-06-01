<script setup>
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import authService from '@/services/AuthService'
import OtpModal from '@/components/OtpModal.vue'

const router = useRouter()
const email = ref('')
const password = ref('')
const loading = ref(false)
const showPassword = ref(false)
const errorMessage = ref('')

// OTP Modal state
const isOtpModalOpen = ref(false)
const otpDigits = ref(Array.from({ length: 6 }, () => ''))
const otpStatus = ref('idle')
const otpErrorMessage = ref('')
const isOtpResending = ref(false)
const otpContext = ref({
  userId: '',
  email: '',
  name: '',
  role: '',
  otpExpiresIn: 0,
})

const maskEmail = (emailAddress) => {
  if (!emailAddress) return ''
  const [name, domain] = emailAddress.split('@')
  if (!domain) return emailAddress
  const visibleCount = Math.min(3, name.length)
  const visible = name.slice(0, visibleCount)
  return `${visible}***@${domain}`
}

const maskedOtpEmail = computed(() => maskEmail(otpContext.value.email))

const resetOtpState = () => {
  otpDigits.value = Array.from({ length: 6 }, () => '')
  otpStatus.value = 'idle'
  otpErrorMessage.value = ''
}

const openOtpModal = (otpData) => {
  otpContext.value = {
    userId: otpData.id || '',
    email: otpData.email || '',
    name: otpData.name || '',
    role: otpData.role || 'dokter',
    otpExpiresIn: otpData.otpExpiresIn || 0,
  }
  resetOtpState()
  isOtpModalOpen.value = true
}

const closeOtpModal = () => {
  isOtpModalOpen.value = false
  resetOtpState()
  otpContext.value = {
    userId: '',
    email: '',
    name: '',
    role: '',
    otpExpiresIn: 0,
  }
}

// Function to navigate to user login page
const navigateToUserLogin = () => {
  router.push('/')
}

const handleLogin = async () => {
  loading.value = true
  errorMessage.value = ''
  try {
    const loginResult = await authService.loginWithEmail(email.value, password.value)

    if (loginResult?.otpRequired) {
      openOtpModal(loginResult)
      return
    }

    router.push('/home-admin')
  } catch (err) {
    errorMessage.value = err.message || 'Login failed. Please try again.'
  } finally {
    loading.value = false
  }
}

const togglePasswordVisibility = () => {
  showPassword.value = !showPassword.value
}

const getOtpErrorMessage = (error) => {
  const message = error?.message?.toLowerCase?.() || ''
  if (message.includes('otp') || error?.status === 401) {
    return 'OTP tidak valid atau sudah kadaluarsa'
  }
  return error?.message || 'Verifikasi OTP gagal'
}

const verifyOtp = async () => {
  const otpCode = otpDigits.value.join('')

  if (!otpContext.value.userId) {
    otpStatus.value = 'error'
    otpErrorMessage.value = 'OTP verification failed'
    return
  }

  try {
    otpStatus.value = 'verifying'
    otpErrorMessage.value = ''

    await authService.verifyOtp(otpContext.value.userId, otpCode, {
      email: email.value,
      password: password.value,
    })

    otpStatus.value = 'success'

    setTimeout(() => {
      closeOtpModal()
      router.push('/home-admin')
    }, 1200)
  } catch (error) {
    otpStatus.value = 'error'
    otpErrorMessage.value = getOtpErrorMessage(error)
  }
}

const resendOtp = async () => {
  if (isOtpResending.value) return

  try {
    isOtpResending.value = true
    otpErrorMessage.value = ''

    const loginResult = await authService.loginWithEmail(email.value, password.value)

    if (loginResult?.otpRequired) {
      openOtpModal(loginResult)
      return
    }

    router.push('/home-admin')
  } catch (error) {
    otpStatus.value = 'error'
    otpErrorMessage.value = getOtpErrorMessage(error)
  } finally {
    isOtpResending.value = false
  }
}
</script>

<template>
  <div class="login-container">
    <div class="login-content">
      <!-- Left side - Medical illustration/branding -->
      <div class="login-branding">
        <div class="brand-content">
          <h1 class="brand-title">JantungIn</h1>
          <p class="brand-subtitle">Advanced Cardiovascular Healthcare Platform</p>
          <div class="brand-features">
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
                  d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
                />
              </svg>
              <span>AI-Powered Diagnostics</span>
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
                  d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2"
                />
              </svg>
              <span>Patient Records Management</span>
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
                  d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"
                />
              </svg>
              <span>Risk Analysis & Reporting</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Right side - Login form -->
      <div class="login-form-container">
        <div class="login-card">
          <div class="login-header">
            <div class="logo-container">
              <svg
                class="heart-logo"
                viewBox="0 0 24 24"
                fill="none"
                xmlns="http://www.w3.org/2000/svg"
                @click="navigateToUserLogin"
                style="cursor: pointer"
              >
                <path
                  d="M12 21.35l-1.45-1.32C5.4 15.36 2 12.28 2 8.5 2 5.42 4.42 3 7.5 3c1.74 0 3.41.81 4.5 2.09C13.09 3.81 14.76 3 16.5 3 19.58 3 22 5.42 22 8.5c0 3.78-3.4 6.86-8.55 11.54L12 21.35z"
                  fill="currentColor"
                />
              </svg>
            </div>
            <h2 class="login-title">Welcome Back</h2>
            <p class="login-subtitle">Sign in to access your dashboard</p>
          </div>

          <form @submit.prevent="handleLogin">
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
                    d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"
                  />
                </svg>
                <input
                  id="email"
                  v-model="email"
                  type="email"
                  class="form-control"
                  placeholder="Enter your email"
                  required
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
                  placeholder="Enter your password"
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
            </div>

            <!-- Error message - moved here for better positioning -->
            <div class="error-message" v-if="errorMessage">
              {{ errorMessage }}
            </div>

            <button type="submit" class="btn-primary" :disabled="loading">
              <span v-if="loading" class="loading-spinner"></span>
              <span>{{ loading ? 'Signing in...' : 'Sign In' }}</span>
            </button>

            <div class="login-divider">
              <span>Healthcare Professional Access</span>
            </div>

            <p class="secure-note">
              <svg
                xmlns="http://www.w3.org/2000/svg"
                class="secure-icon"
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
              Secure, encrypted connection
            </p>
          </form>
        </div>
      </div>
    </div>

    <!-- OTP Modal Component -->
    <OtpModal
      :is-open="isOtpModalOpen"
      :otp-digits="otpDigits"
      :otp-status="otpStatus"
      :otp-error-message="otpErrorMessage"
      :is-otp-resending="isOtpResending"
      :masked-email="maskedOtpEmail"
      @close="closeOtpModal"
      @verify="verifyOtp"
      @resend="resendOtp"
      @update:otp-digits="otpDigits = $event"
    />
  </div>
</template>

<style scoped>
.login-container {
  min-height: 100vh;
  width: 100%;
  background: linear-gradient(135deg, #f0f9ff 0%, #e1f5fe 100%);
  font-family: 'Inter', 'Segoe UI', sans-serif;
  overflow: hidden;
}

.login-content {
  display: flex;
  min-height: 100vh;
}

/* Left side styling */
.login-branding {
  display: none;
  background: linear-gradient(135deg, #2563eb 0%, #1e40af 100%);
  color: white;
  padding: 2rem;
  position: relative;
  overflow: hidden;
}

.brand-content {
  position: relative;
  z-index: 10;
  display: flex;
  flex-direction: column;
  justify-content: center;
  height: 100%;
  max-width: 420px;
  margin: 0 auto;
}

.logo-container {
  display: flex;
  align-items: center;
  margin-bottom: 1.5rem;
  position: relative;
}

.heart-logo {
  position: relative;
  width: 3rem;
  height: 3rem;
  color: #f43f5e;
  filter: drop-shadow(0 0 8px rgba(244, 63, 94, 0.5));
  transition:
    transform 0.3s ease,
    filter 0.3s ease;
}

.heart-logo:hover {
  transform: scale(1.1);
  filter: drop-shadow(0 0 12px rgba(244, 63, 94, 0.8));
}

.brand-title {
  font-size: 2.5rem;
  font-weight: 700;
  margin-bottom: 0.5rem;
  background: linear-gradient(to right, #ffffff, #e0f2fe);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.brand-subtitle {
  font-size: 1.25rem;
  margin-bottom: 3rem;
  opacity: 0.9;
}

.brand-features {
  margin-top: 3rem;
}

.feature-item {
  display: flex;
  align-items: center;
  margin-bottom: 1.5rem;
}

.feature-icon {
  width: 1.5rem;
  height: 1.5rem;
  margin-right: 1rem;
  color: #60a5fa;
}

/* Right side - Login form */
.login-form-container {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 1.5rem;
}

.login-card {
  width: 100%;
  max-width: 450px;
  background: white;
  border-radius: 1rem;
  padding: 2.5rem 2rem;
  box-shadow: 0 10px 25px rgba(0, 0, 0, 0.05);
}

.login-header {
  text-align: center;
  margin-bottom: 2.5rem;
}

.login-title {
  font-size: 1.75rem;
  font-weight: 700;
  color: #1e40af;
  margin-bottom: 0.5rem;
}

.login-subtitle {
  color: #64748b;
  font-size: 0.95rem;
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

.remember-forgot {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
  font-size: 0.875rem;
}

.remember-me {
  display: flex;
  align-items: center;
  color: #334155;
  cursor: pointer;
}

.remember-me input {
  margin-right: 0.5rem;
  cursor: pointer;
}

.forgot-password {
  color: #3b82f6;
  text-decoration: none;
}

.forgot-password:hover {
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

.error-message {
  color: #ef4444;
  background-color: #fee2e2;
  border: 1px solid #f87171;
  border-radius: 0.5rem;
  padding: 0.75rem;
  margin-bottom: 1.5rem;
  font-size: 0.875rem;
  text-align: center;
}

.login-divider {
  position: relative;
  text-align: center;
  margin: 1.5rem 0;
}

.login-divider::before {
  content: '';
  position: absolute;
  top: 50%;
  left: 0;
  right: 0;
  height: 1px;
  background-color: #e2e8f0;
}

.login-divider span {
  position: relative;
  background-color: white;
  padding: 0 1rem;
  color: #64748b;
  font-size: 0.875rem;
}

.secure-note {
  display: flex;
  align-items: center;
  justify-content: center;
  color: #64748b;
  font-size: 0.875rem;
  margin-top: 2rem;
}

.secure-icon {
  width: 1rem;
  height: 1rem;
  margin-right: 0.5rem;
  color: #64748b;
}

/* Responsive adjustments */
@media (min-width: 1024px) {
  .login-branding {
    display: block;
    flex: 1;
  }

  .login-form-container {
    flex: 1;
  }
}

/* Animated background for branding section */
.login-branding::before {
  content: '';
  position: absolute;
  top: -50%;
  left: -50%;
  width: 200%;
  height: 200%;
  background: radial-gradient(circle, rgba(59, 130, 246, 0.1) 0%, rgba(29, 78, 216, 0.1) 100%);
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
