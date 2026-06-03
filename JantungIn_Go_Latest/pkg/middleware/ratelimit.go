package middleware

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"jantungin-api-server/pkg/utils"

	"github.com/gin-gonic/gin"
)

// RateLimitStore tracks rate limits per identifier (IP or user ID)
type RateLimitStore struct {
	mu       sync.RWMutex
	attempts map[string][]time.Time
}

func NewRateLimitStore() *RateLimitStore {
	return &RateLimitStore{
		attempts: make(map[string][]time.Time),
	}
}

// IsAllowed checks if an identifier is allowed to make a request
// Returns (allowed, retryAfterSeconds)
func (r *RateLimitStore) IsAllowed(identifier string, maxAttempts int, windowSeconds int) (bool, int) {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	windowTime := time.Duration(windowSeconds) * time.Second

	// Get attempts for this identifier
	attempts := r.attempts[identifier]

	// Remove attempts outside the time window
	validAttempts := []time.Time{}
	for _, attempt := range attempts {
		if now.Sub(attempt) < windowTime {
			validAttempts = append(validAttempts, attempt)
		}
	}

	// Check if limit exceeded
	if len(validAttempts) >= maxAttempts {
		// Calculate retry after seconds
		retryAfter := int(windowTime.Seconds()) - int(now.Sub(validAttempts[0]).Seconds())
		if retryAfter < 1 {
			retryAfter = 1
		}
		return false, retryAfter
	}

	// Add new attempt
	validAttempts = append(validAttempts, now)
	r.attempts[identifier] = validAttempts

	return true, 0
}

// Cleanup removes old entries periodically to prevent memory leak
func (r *RateLimitStore) Cleanup() {
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()

		for range ticker.C {
			r.mu.Lock()
			now := time.Now()

			for key, attempts := range r.attempts {
				validAttempts := []time.Time{}
				for _, attempt := range attempts {
					// Keep attempts within last 24 hours
					if now.Sub(attempt) < 24*time.Hour {
						validAttempts = append(validAttempts, attempt)
					}
				}

				if len(validAttempts) == 0 {
					delete(r.attempts, key)
				} else {
					r.attempts[key] = validAttempts
				}
			}

			r.mu.Unlock()
		}
	}()
}

// Global rate limit stores for different endpoints
var (
	loginRateLimitStore     = NewRateLimitStore()
	otpResendRateLimitStore = NewRateLimitStore()
)

func init() {
	loginRateLimitStore.Cleanup()
	otpResendRateLimitStore.Cleanup()
}

// LoginRateLimit middleware - 5 attempts per 15 minutes per IP
// In development environment, the rate limiter is disabled.
func LoginRateLimit(env string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip rate limiting in development
		if env == "development" {
			c.Next()
			return
		}

		ip := c.ClientIP()

		allowed, retryAfter := loginRateLimitStore.IsAllowed(ip, 5, 900) // 5 attempts per 15 minutes

		if !allowed {
			c.Header("Retry-After", strconv.Itoa(retryAfter))
			utils.ErrorResponse(c, http.StatusTooManyRequests, "Too many login attempts. Please try again later.", map[string]interface{}{
				"retry_after_seconds": retryAfter,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// OTPResendRateLimit middleware - 3 attempts per 60 minutes per user ID
// This should be used after AuthRequired middleware
func OTPResendRateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get(AuthUserIDKey)
		if !exists {
			utils.UnauthorizedResponse(c, "User not authenticated for OTP resend")
			c.Abort()
			return
		}

		userIDStr, ok := userID.(string)
		if !ok {
			utils.UnauthorizedResponse(c, "Invalid user identity")
			c.Abort()
			return
		}

		allowed, retryAfter := otpResendRateLimitStore.IsAllowed(userIDStr, 3, 3600) // 3 attempts per 60 minutes

		if !allowed {
			c.Header("Retry-After", strconv.Itoa(retryAfter))
			utils.ErrorResponse(c, http.StatusTooManyRequests, "Too many OTP resend attempts. Please try again later.", map[string]interface{}{
				"retry_after_seconds": retryAfter,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
