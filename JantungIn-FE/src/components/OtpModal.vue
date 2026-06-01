<script setup>
import { ref, computed, nextTick, onMounted, onUnmounted } from 'vue'

const props = defineProps({
  isOpen: {
    type: Boolean,
    default: false,
  },
  otpDigits: {
    type: Array,
    default: () => Array(6).fill(''),
  },
  otpStatus: {
    type: String,
    default: 'idle', // idle, verifying, success, error
  },
  otpErrorMessage: {
    type: String,
    default: '',
  },
  isOtpResending: {
    type: Boolean,
    default: false,
  },
  maskedEmail: {
    type: String,
    default: '',
  },
})

const emit = defineEmits(['close', 'verify', 'resend'])

const otpInputRefs = ref([])
const resendCountdown = ref(0)
const resendTimerInterval = ref(null)

const isOtpLocked = computed(() => props.otpStatus === 'verifying' || props.otpStatus === 'success')
const otpInputStateClass = computed(() => ({
  'otp-inputs--error': props.otpStatus === 'error',
  'otp-inputs--success': props.otpStatus === 'success',
}))
const isOtpComplete = computed(() => props.otpDigits.every((digit) => digit !== ''))

const focusOtpInput = (index) => {
  const input = otpInputRefs.value[index]
  if (input) {
    input.focus()
  }
}

const handleOtpInput = (index, event) => {
  const rawValue = event.target.value
  const digits = rawValue.replace(/[^0-9]/g, '').split('')

  if (digits.length > 1) {
    const newOtpDigits = [...props.otpDigits]
    newOtpDigits.splice(index, 6 - index, ...digits.slice(0, 6 - index))
    emit('update:otpDigits', newOtpDigits)

    const nextIndex = Math.min(index + digits.length, 5)
    nextTick(() => {
      focusOtpInput(nextIndex)
    })
  } else {
    const newOtpDigits = [...props.otpDigits]
    newOtpDigits[index] = digits[0] || ''
    emit('update:otpDigits', newOtpDigits)

    if (digits.length === 1 && index < 5) {
      const nextFocusIndex = index + 1
      nextTick(() => {
        focusOtpInput(nextFocusIndex)
      })
    }
  }
}

const handleOtpKeydown = (index, event) => {
  if (event.key === 'Backspace' && !props.otpDigits[index] && index > 0) {
    const newOtpDigits = [...props.otpDigits]
    newOtpDigits[index - 1] = ''
    emit('update:otpDigits', newOtpDigits)
    nextTick(() => {
      focusOtpInput(index - 1)
    })
  } else if (event.key === 'ArrowLeft' && index > 0) {
    focusOtpInput(index - 1)
  } else if (event.key === 'ArrowRight' && index < 5) {
    focusOtpInput(index + 1)
  }
}

const handleOtpPaste = (event) => {
  event.preventDefault()
  const pasted = (event.clipboardData || window.clipboardData).getData('text')
  const digits = pasted.replace(/[^0-9]/g, '').split('')

  if (digits.length > 0) {
    const newOtpDigits = digits.slice(0, 6).concat(Array(6).fill('')).slice(0, 6)
    emit('update:otpDigits', newOtpDigits)

    const focusIndex = Math.min(digits.length, 5)
    nextTick(() => {
      focusOtpInput(focusIndex)
    })
  }
}

const handleClose = () => {
  emit('close')
}

const handleVerify = () => {
  emit('verify')
}

const handleResend = () => {
  // Prevent double-click and abuse by checking resend conditions
  // Use props.isOtpResending instead of isOtpResending.value since it's a prop
  if (props.isOtpResending || isOtpLocked.value || resendCountdown.value > 0) {
    console.warn('Resend OTP blocked - already resending, locked, or countdown active')
    return
  }

  console.log('Emitting resend event')
  emit('resend')
  startResendCountdown()
}

const startResendCountdown = () => {
  resendCountdown.value = 30
  if (resendTimerInterval.value) {
    clearInterval(resendTimerInterval.value)
  }
  resendTimerInterval.value = setInterval(() => {
    resendCountdown.value -= 1
    if (resendCountdown.value <= 0) {
      clearInterval(resendTimerInterval.value)
      resendTimerInterval.value = null
    }
  }, 1000)
}

onMounted(() => {
  return () => {
    if (resendTimerInterval.value) {
      clearInterval(resendTimerInterval.value)
    }
  }
})

onUnmounted(() => {
  if (resendTimerInterval.value) {
    clearInterval(resendTimerInterval.value)
  }
})
</script>

<template>
  <div v-if="isOpen" class="otp-modal-backdrop">
    <div class="otp-modal-card" role="dialog" aria-modal="true">
      <button class="otp-close-btn" type="button" @click="handleClose" :disabled="isOtpLocked">
        <img src="/images/otp/otp-close.png" alt="Close" />
      </button>

      <img class="otp-shield" src="/images/otp/otp-shield.png" alt="Security" />

      <p class="otp-description">
        We've sent a 6-digit verification code to your email
        <span class="otp-email">({{ maskedEmail || 'your email' }})</span>
        Please enter it below.
      </p>

      <p v-if="otpStatus === 'error'" class="otp-error-message">
        {{ otpErrorMessage || 'OTP is Invalid or Expired' }}
      </p>

      <div :class="['otp-inputs', otpInputStateClass]">
        <input
          v-for="(_, index) in otpDigits"
          :key="`otp-${index}`"
          class="otp-input"
          type="text"
          inputmode="numeric"
          maxlength="1"
          :value="otpDigits[index]"
          :disabled="isOtpLocked"
          ref="otpInputRefs"
          :autocomplete="index === 0 ? 'one-time-code' : 'off'"
          @input="(event) => handleOtpInput(index, event)"
          @keydown="(event) => handleOtpKeydown(index, event)"
          @paste="handleOtpPaste"
        />
      </div>

      <div v-if="otpStatus === 'success'" class="otp-success">
        <p class="otp-success-text">Please wait, Redirecting to Dashboard . . .</p>
        <div class="otp-wait-spinner"></div>
      </div>

      <button
        v-if="otpStatus !== 'success'"
        type="button"
        class="otp-verify-btn"
        @click="handleVerify"
        :disabled="!isOtpComplete || isOtpLocked"
      >
        <span v-if="otpStatus === 'verifying'">Verifying...</span>
        <span v-else>Verify &amp; Proceed →</span>
      </button>

      <p v-if="otpStatus !== 'success'" class="otp-resend">
        Didn't receive the code?
        <button
          type="button"
          class="otp-resend-link"
          @click="handleResend"
          :disabled="isOtpResending || isOtpLocked || resendCountdown > 0"
        >
          <span v-if="resendCountdown > 0">Resend in {{ resendCountdown }}s</span>
          <span v-else>Resend OTP</span>
        </button>
      </p>
    </div>
  </div>
</template>

<style scoped>
.otp-modal-backdrop {
  position: fixed;
  inset: 0;
  background: rgba(168, 196, 255, 0.79);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 2000;
  padding: 1.5rem;
}

.otp-modal-card {
  width: 559px;
  max-width: 100%;
  background: #ffffff;
  border-radius: 10px;
  padding: 2.5rem 2rem 2rem;
  box-shadow: 0 20px 50px rgba(15, 23, 42, 0.15);
  position: relative;
  text-align: center;
}

.otp-close-btn {
  position: absolute;
  top: 1.5rem;
  right: 1.5rem;
  width: 30px;
  height: 30px;
  border: none;
  background: transparent;
  cursor: pointer;
  padding: 0;
}

.otp-close-btn:disabled {
  cursor: not-allowed;
  opacity: 0.6;
}

.otp-close-btn img {
  width: 30px;
  height: 30px;
}

.otp-shield {
  width: 90px;
  height: 90px;
  margin: 0 auto 1rem;
}

.otp-description {
  color: #424754;
  font-size: 18px;
  line-height: 1.4;
  max-width: 360px;
  margin: 0 auto 0.75rem;
}

.otp-email {
  color: #2563eb;
  font-weight: 600;
}

.otp-error-message {
  color: #c32c2f;
  font-size: 16px;
  font-weight: 700;
  margin-bottom: 0.75rem;
}

.otp-inputs {
  display: flex;
  justify-content: center;
  gap: 12px;
  margin: 1rem 0 1.5rem;
}

.otp-input {
  width: 56px;
  height: 66px;
  border-radius: 10px;
  border: 1px solid transparent;
  background: rgba(37, 99, 235, 0.15);
  text-align: center;
  font-size: 36px;
  font-weight: 700;
  font-family: 'Afacad', 'Inter', 'Segoe UI', sans-serif;
  color: #000000;
}

.otp-input:focus {
  outline: none;
  border-color: #2563eb;
  box-shadow: 0 0 0 3px rgba(37, 99, 235, 0.25);
}

.otp-input:disabled {
  cursor: not-allowed;
  opacity: 0.8;
}

.otp-inputs--success .otp-input {
  background: rgba(52, 199, 89, 0.5);
}

.otp-inputs--error .otp-input {
  background: rgba(199, 52, 52, 0.5);
  border-color: rgba(199, 52, 52, 0.7);
}

.otp-verify-btn {
  width: 417px;
  max-width: 100%;
  height: 57px;
  background: #2563eb;
  border-radius: 10px;
  border: none;
  color: #ffffff;
  font-size: 24px;
  font-weight: 700;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto 1rem;
}

.otp-verify-btn:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}

.otp-resend {
  font-size: 18px;
  font-weight: 700;
  color: #000000;
}

.otp-resend-link {
  background: none;
  border: none;
  color: #2563eb;
  cursor: pointer;
  font-weight: 700;
  padding: 0;
}

.otp-resend-link:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.otp-success {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1rem;
  margin-top: 0.5rem;
}

.otp-success-text {
  color: #048217;
  font-size: 18px;
}

.otp-wait-spinner {
  width: 50px;
  height: 50px;
  border: 3px solid rgba(52, 199, 89, 0.2);
  border-top-color: #34c759;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin: 0 auto;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

@media (max-width: 640px) {
  .otp-modal-card {
    padding: 2rem 1.5rem;
  }

  .otp-description {
    font-size: 16px;
  }

  .otp-inputs {
    gap: 8px;
  }

  .otp-input {
    width: 42px;
    height: 52px;
    font-size: 28px;
  }

  .otp-verify-btn {
    font-size: 20px;
    height: 52px;
  }
}
</style>
