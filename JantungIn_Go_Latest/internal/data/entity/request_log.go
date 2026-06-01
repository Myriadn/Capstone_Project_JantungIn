package entity

import (
	"time"

	"github.com/google/uuid"
)

type RequestLog struct {
	ID           int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	Method       string     `gorm:"not null" json:"method"`
	Path         string     `gorm:"not null" json:"path"`
	StatusCode   int        `gorm:"not null" json:"statusCode"`
	IP           string     `gorm:"not null" json:"ip"`
	UserAgent    string     `gorm:"not null" json:"userAgent"`
	LatencyMs    int64      `gorm:"not null;default:0" json:"latencyMs"`
	UserID       *uuid.UUID `gorm:"type:uuid" json:"userId"`
	ClientType   string     `gorm:"not null" json:"clientType"`
	TargetUser   string     `gorm:"not null" json:"targetUser"`
	ActionStatus string     `gorm:"not null" json:"actionStatus"`
	Reason       string     `gorm:"not null" json:"reason"`
	CreatedAt    time.Time  `json:"createdAt"`
}

func (RequestLog) TableName() string {
	return "request_logs"
}
