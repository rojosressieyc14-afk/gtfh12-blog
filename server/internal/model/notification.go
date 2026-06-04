package model

import "time"

const (
	NotificationTypeArticleReview  = "article_review"
	NotificationTypeProjectReview  = "project_review"
	NotificationTypeArticleComment = "article_comment"
	NotificationTypeCommentReply   = "comment_reply"
	NotificationTypeModeration     = "moderation"
	NotificationTypeSystem         = "system"
)

type Notification struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uint           `gorm:"index;not null" json:"userId"`
	Title     string         `gorm:"size:120;not null" json:"title"`
	Content   string         `gorm:"size:255;not null" json:"content"`
	Type      string         `gorm:"size:40;index;not null;default:system" json:"type"`
	ActionURL string         `gorm:"size:255" json:"actionUrl"`
	Payload   map[string]any `gorm:"serializer:json" json:"payload"`
	IsRead    bool           `gorm:"not null;default:false" json:"isRead"`
	CreatedAt time.Time      `json:"createdAt"`
}
