package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"jantungin-api-server/internal/data/entity"
	"jantungin-api-server/internal/data/repository"
)

const (
	RequestLogTargetUserKey   = "request_log_target_user"
	RequestLogActionStatusKey = "request_log_action_status"
	RequestLogReasonKey       = "request_log_reason"
	RequestLogClientTypeKey   = "request_log_client_type"
)

func RequestTracker(statsRepo repository.StatsRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		latencyMs := time.Since(start).Milliseconds()
		statusCode := c.Writer.Status()

		userAgent := c.Request.UserAgent()
		clientType := userAgent
		if rawClientType, exists := c.Get(RequestLogClientTypeKey); exists {
			if ct, ok := rawClientType.(string); ok && strings.TrimSpace(ct) != "" {
				clientType = ct
			}
		}
		if strings.TrimSpace(clientType) == "" {
			clientType = "-"
		}
		if strings.TrimSpace(userAgent) == "" {
			userAgent = clientType
		}

		targetUser := "-"
		if rawTarget, exists := c.Get(RequestLogTargetUserKey); exists {
			if tu, ok := rawTarget.(string); ok && strings.TrimSpace(tu) != "" {
				targetUser = tu
			}
		}

		actionStatus := ""
		if rawStatus, exists := c.Get(RequestLogActionStatusKey); exists {
			if as, ok := rawStatus.(string); ok && strings.TrimSpace(as) != "" {
				actionStatus = as
			}
		}
		if actionStatus == "" {
			if statusCode >= http.StatusBadRequest {
				actionStatus = "failed"
			} else {
				actionStatus = "success"
			}
		}

		reason := ""
		if rawReason, exists := c.Get(RequestLogReasonKey); exists {
			if rs, ok := rawReason.(string); ok && strings.TrimSpace(rs) != "" {
				reason = rs
			}
		}
		if reason == "" {
			reason = http.StatusText(statusCode)
		}
		if strings.TrimSpace(reason) == "" {
			reason = "unknown"
		}

		var userID *uuid.UUID
		if rawID, exists := c.Get(AuthUserIDKey); exists {
			if idStr, ok := rawID.(string); ok {
				if parsed, err := uuid.Parse(idStr); err == nil {
					userID = &parsed
				}
			}
		}

		log := &entity.RequestLog{
			Method:       c.Request.Method,
			Path:         c.Request.URL.Path,
			StatusCode:   statusCode,
			IP:           c.ClientIP(),
			UserAgent:    userAgent,
			LatencyMs:    latencyMs,
			UserID:       userID,
			ClientType:   clientType,
			TargetUser:   targetUser,
			ActionStatus: actionStatus,
			Reason:       reason,
		}

		// Semua nilai sudah di-capture sebelum goroutine
		// Pakai context.Background() karena c.Request.Context() sudah cancelled
		// saat request selesai dan goroutine baru mulai jalan
		go func(l *entity.RequestLog) {
			_ = statsRepo.InsertRequestLog(context.Background(), l)
		}(log)
	}
}
