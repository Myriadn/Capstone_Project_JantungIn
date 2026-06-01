import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import authService from '@/services/AuthService'
import { useErrorHandler } from '@/utils/errorHandler'

/**
 * Login ViewModel
 * Handles login form state and authentication logic
 */
export function useLoginViewModel() {
  const router = useRouter()
  const { t } = useI18n()
  const { getErrorMessage } = useErrorHandler()

  // State
  const username = ref('')
  const password = ref('')
  const rememberMe = ref(false)
  const showPassword = ref(false)

  // UI state
  const isLoading = ref(false)
  const errorMessage = ref('')
  const isOfflineMode = ref(!navigator.onLine)

  // OTP state
  const isOtpModalOpen = ref(false)
  const otpDigits = ref(Array.from({ length: 6 }, () => ''))
  const otpStatus = ref('idle')
  const otpErrorMessage = ref('')
  const isOtpResending = ref(false)
  const lastOtpResendTime = ref(0)
  const otpResendAttempts = ref(0)
  const otpContext = ref({
    userId: '',
    email: '',
    name: '',
    role: '',
    otpExpiresIn: 0,
  })

  /**
   * Computed property to check if form is valid
   */
  const isFormValid = computed(() => {
    return username.value.trim() !== '' && password.value.trim() !== ''
  })

  const maskEmail = (email) => {
    if (!email) return ''
    const [name, domain] = email.split('@')
    if (!domain) return email

    const visibleCount = Math.min(3, name.length)
    const visible = name.slice(0, visibleCount)
    return `${visible}***@${domain}`
  }

  const otpCode = computed(() => otpDigits.value.join(''))
  const isOtpComplete = computed(() => otpDigits.value.every((digit) => digit !== ''))
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
      role: otpData.role || 'user',
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

  const getOtpErrorMessage = (error) => {
    const message = error?.message?.toLowerCase?.() || ''

    if (message.includes('otp') || error?.status === 401) {
      return 'OTP is Invalid or Expired'
    }

    return error?.message || 'OTP verification failed'
  }

  /**
   * Handle login form submission
   */
  const handleLogin = async () => {
    if (!isFormValid.value) {
      if (username.value.trim() === '') {
        errorMessage.value = t('errors.validation.requiredField')
      } else {
        errorMessage.value = t('errors.auth.missingFields')
      }
      return
    }

    try {
      isLoading.value = true
      errorMessage.value = ''

      const user = await authService.login(username.value, password.value)

      console.log('Login successful:', user)

      if (user?.otpRequired) {
        openOtpModal(user)
        return
      }

      // Redirect to news page after successful login
      router.push('/home')
    } catch (error) {
      console.error('Login error:', error)

      // Use error handler to get localized error message
      errorMessage.value = getErrorMessage(error)

      throw error
    } finally {
      isLoading.value = false
    }
  }

  const verifyOtp = async () => {
    if (!isOtpComplete.value) {
      otpStatus.value = 'error'
      otpErrorMessage.value = 'OTP harus 6 digit'
      return
    }

    if (!otpContext.value.userId) {
      otpStatus.value = 'error'
      otpErrorMessage.value = 'OTP verification failed'
      return
    }

    try {
      otpStatus.value = 'verifying'
      otpErrorMessage.value = ''

      await authService.verifyOtp(otpContext.value.userId, otpCode.value, {
        username: username.value,
        password: password.value,
      })

      otpStatus.value = 'success'

      setTimeout(() => {
        closeOtpModal()
        router.push('/home')
      }, 1200)
    } catch (error) {
      otpStatus.value = 'error'
      otpErrorMessage.value = getOtpErrorMessage(error)
    }
  }

  const resendOtp = async () => {
    // Rate limiting checks to prevent abuse
    if (isOtpResending.value) {
      console.warn('OTP resend already in progress')
      return
    }

    const now = Date.now()
    const timeSinceLastResend = now - lastOtpResendTime.value

    // Check cooldown (30 seconds minimum between resends - matches countdown timer)
    if (timeSinceLastResend < 30000) {
      const waitSeconds = Math.ceil((30000 - timeSinceLastResend) / 1000)
      otpStatus.value = 'error'
      otpErrorMessage.value = `Please wait ${waitSeconds}s before resending OTP`
      console.warn('OTP resend cooldown active, wait time:', waitSeconds, 's')
      return
    }

    // Check max attempts (3 attempts max)
    if (otpResendAttempts.value >= 3) {
      otpStatus.value = 'error'
      otpErrorMessage.value = 'Maximum resend attempts reached. Please try again later.'
      console.warn('Max OTP resend attempts exceeded')
      return
    }

    try {
      isOtpResending.value = true
      otpErrorMessage.value = ''

      const user = await authService.login(username.value, password.value)

      // Successfully resent OTP - update tracking
      otpResendAttempts.value += 1
      lastOtpResendTime.value = now
      console.info('OTP resend successful, attempt:', otpResendAttempts.value)

      if (user?.otpRequired) {
        openOtpModal(user)
        return
      }

      router.push('/home')
    } catch (error) {
      otpStatus.value = 'error'
      otpErrorMessage.value = getOtpErrorMessage(error)
      console.error('OTP resend error:', error)
    } finally {
      isOtpResending.value = false
    }
  }

  /**
   * Navigate to register page
   */
  const goToRegister = () => {
    router.push('/register')
  }

  /**
   * Reset login form
   */
  const resetForm = () => {
    username.value = ''
    password.value = ''
    errorMessage.value = ''
  }

  /**
   * Toggle password visibility
   */
  const togglePasswordVisibility = () => {
    showPassword.value = !showPassword.value
  }

  // Monitor online status
  window.addEventListener('online', () => {
    isOfflineMode.value = false
  })
  window.addEventListener('offline', () => {
    isOfflineMode.value = true
  })

  return {
    // State
    username,
    password,
    rememberMe,
    showPassword,
    isLoading,
    errorMessage,
    isOfflineMode,
    isFormValid,
    isOtpModalOpen,
    otpDigits,
    otpStatus,
    otpErrorMessage,
    isOtpResending,
    maskedOtpEmail,
    isOtpComplete,

    // Methods
    handleLogin,
    verifyOtp,
    resendOtp,
    closeOtpModal,
    goToRegister,
    resetForm,
    togglePasswordVisibility,
  }
}
