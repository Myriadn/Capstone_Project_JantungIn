package middleware

import (
	"slices"
	"strings"

	"jantungin-api-server/pkg/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	AuthUserIDKey       = "auth_user_id"
	AuthEmailKey        = "auth_email"
	AuthRoleIDKey       = "auth_role_id"
	AuthRoleCodeKey     = "auth_role_code"
	AuthTokenCookieName = "access_token"
)

func AuthRequired(cfg *utils.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		tokenString := ""

		if authHeader != "" {
			if !strings.HasPrefix(authHeader, "Bearer ") {
				utils.Warn("Invalid authorization header format",
					zap.String("path", c.Request.URL.Path),
				)
				utils.UnauthorizedResponse(c, "Invalid authentication token format")
				c.Abort()
				return
			}

			tokenString = strings.TrimPrefix(authHeader, "Bearer ")
			if tokenString == "" {
				utils.Warn("Empty token after Bearer prefix",
					zap.String("path", c.Request.URL.Path),
				)
				utils.UnauthorizedResponse(c, "Empty authentication token")
				c.Abort()
				return
			}
		} else {
			if cookieToken, err := c.Cookie(AuthTokenCookieName); err == nil {
				tokenString = cookieToken
			}
		}

		if tokenString == "" {
			utils.Warn("Missing authentication token",
				zap.String("path", c.Request.URL.Path),
			)
			utils.UnauthorizedResponse(c, "Missing authentication token")
			c.Abort()
			return
		}

		claims, err := utils.ValidateToken(tokenString, cfg)
		if err != nil {
			utils.Warn("Token validation failed",
				zap.String("path", c.Request.URL.Path),
				zap.Error(err),
			)
			utils.UnauthorizedResponse(c, "Invalid or expired token")
			c.Abort()
			return
		}

		c.Set(AuthUserIDKey, claims.UserID)
		c.Set(AuthEmailKey, claims.Email)
		c.Set(AuthRoleIDKey, claims.RoleID)
		c.Set(AuthRoleCodeKey, claims.RoleCode)

		utils.Debug("Authentication successful",
			zap.String("user_id", claims.UserID),
			zap.String("role_code", claims.RoleCode),
		)

		c.Next()
	}
}

func RoleRequired(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleCode, exists := c.Get(AuthRoleCodeKey)
		if !exists {
			utils.Warn("Role code not found in context",
				zap.String("path", c.Request.URL.Path),
			)
			utils.ForbiddenResponse(c, "Access denied. Insufficient privileges.")
			c.Abort()
			return
		}

		userRole, ok := roleCode.(string)
		if !ok {
			utils.Warn("Invalid role code type in context",
				zap.String("path", c.Request.URL.Path),
			)
			utils.ForbiddenResponse(c, "Access denied. Insufficient privileges.")
			c.Abort()
			return
		}

		allowed := slices.Contains(roles, userRole)

		if !allowed {
			utils.Warn("Insufficient role privileges",
				zap.String("user_role", userRole),
				zap.Strings("required_roles", roles),
				zap.String("path", c.Request.URL.Path),
			)
			utils.ForbiddenResponse(c, "Access denied. Insufficient privileges.")
			c.Abort()
			return
		}

		c.Next()
	}
}
