/**
 * Error handling utilities for consistent error messages
 */
import { useI18n } from 'vue-i18n'

/**
 * Get error message from error object with i18n support
 * @param {Error} error - Error object
 * @returns {string} Localized error message
 */
export function useErrorHandler() {
  const { t } = useI18n()

  /**
   * Map API error to localized message
   * @param {Error} error - Error object
   * @returns {string} Localized error message
   */
  const getErrorMessage = (error) => {
    // Default message in case we can't determine a more specific one
    let messageKey = 'errors.api.unknownError'

    // Network related errors
    if (error?.isNetworkError || !navigator.onLine) {
      messageKey = 'errors.network.offline'
    } else if (error?.code === 'ECONNABORTED' || error?.isTimeout) {
      messageKey = 'errors.network.timeout'
    } else if (error?.status === 0 || error?.message?.includes('NetworkError')) {
      messageKey = 'errors.network.connectionRefused'
    }
    // HTTP status errors
    else if (error?.status) {
      switch (error.status) {
        case 400:
          messageKey = 'errors.api.badRequest'
          break
        case 401:
          messageKey = 'errors.auth.invalidCredentials'
          break
        case 403:
          messageKey = 'errors.auth.accessDenied'
          break
        case 404:
          if (error.endpoint?.includes('/auth/') || error.endpoint?.includes('/login')) {
            messageKey = 'errors.auth.userNotFound'
          } else {
            messageKey = 'errors.api.resourceNotFound'
          }
          break
        case 409:
          if (error.message?.toLowerCase().includes('username')) {
            messageKey = 'errors.auth.usernameAlreadyRegistered'
          } else if (error.message?.includes('Email')) {
            messageKey = 'errors.auth.emailAlreadyRegistered'
          }
          break
        case 422:
          // Validation errors, try to find a specific message
          if (error.message?.toLowerCase().includes('username')) {
            messageKey = 'errors.validation.usernameFormat'
          } else if (error.message?.includes('Email')) {
            messageKey = 'errors.validation.emailFormat'
          } else if (error.message?.includes('Password')) {
            messageKey = 'errors.validation.passwordLength'
          } else {
            messageKey = 'errors.validation.requiredField'
          }
          break
        case 500:
        case 502:
        case 503:
          messageKey = 'errors.network.serverError'
          break
      }
    }
    // Specific error types
    else if (error?.message) {
      // Handle specific error messages
      if (error.message.toLowerCase().includes('username')) {
        if (error.message.includes('terdaftar')) {
          messageKey = 'errors.auth.usernameAlreadyRegistered'
        } else {
          messageKey = 'errors.validation.usernameFormat'
        }
      } else if (error.message.includes('email')) {
        if (error.message.includes('format')) {
          messageKey = 'errors.validation.emailFormat'
        } else if (error.message.includes('terdaftar')) {
          messageKey = 'errors.auth.emailAlreadyRegistered'
        }
      } else if (error.message.includes('token')) {
        if (error.message.includes('expired')) {
          messageKey = 'errors.auth.tokenExpired'
        } else {
          messageKey = 'errors.auth.invalidToken'
        }
      } else if (error.message.includes('password')) {
        messageKey = 'errors.validation.passwordLength'
      } else if (error.message.includes('required fields')) {
        messageKey = 'errors.auth.missingFields'
      } else if (error.message.includes('Invalid response format')) {
        messageKey = 'errors.api.invalidResponse'
      } else if (error.message.includes('missing token')) {
        messageKey = 'errors.api.missingToken'
      }
    }

    // Get the translated message using the determined message key
    let errorMessage = t(messageKey)

    // If we have a specific message from the server, include it for developers
    if (import.meta.env.DEV && error?.message) {
      errorMessage += ` (${error.message})`
    }

    return errorMessage
  }

  /**
   * Get error code and category from error
   * @param {Error} error - Error object
   * @returns {Object} Error category and code
   */
  const getErrorInfo = (error) => {
    let category = 'unknown'
    let code = 'unknown'

    if (error?.isNetworkError || !navigator.onLine) {
      category = 'network'
      code = 'offline'
    } else if (error?.status) {
      category = error.status >= 400 && error.status < 500 ? 'client' : 'server'
      code = `http${error.status}`
    } else if (
      error?.message?.toLowerCase().includes('username') ||
      error?.message?.includes('password')
    ) {
      category = 'validation'
      code = error?.message?.toLowerCase().includes('username') ? 'username' : 'password'
    }

    return { category, code }
  }

  return {
    getErrorMessage,
    getErrorInfo,
  }
}

/**
 * Format validation error for form fields
 * @param {string} field - Field name
 * @param {string} type - Validation error type
 * @returns {string} Validation error message key
 */
export function getValidationErrorKey(field, type) {
  const errorTypeMap = {
    required: 'requiredField',
    minLength: 'minLength',
    maxLength: 'maxLength',
    pattern: 'format',
    email: 'emailFormat',
  }

  const fieldMap = {
    username: 'username',
    email: 'email',
    password: 'password',
    name: 'name',
    dateOfBirth: 'date',
  }

  const errorType = errorTypeMap[type] || 'format'
  const fieldName = fieldMap[field] || field

  return `errors.validation.${fieldName}${errorType.charAt(0).toUpperCase() + errorType.slice(1)}`
}
