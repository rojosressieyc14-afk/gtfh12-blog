package model

import "time"

type ModerationHit struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `gorm:"index;not null" json:"userId"`
	User        User      `json:"user"`
	Scene       string    `gorm:"size:40;index;not null" json:"scene"`
	Field       string    `gorm:"size:80;not null" json:"field"`
	MatchedWord string    `gorm:"size:120;not null" json:"matchedWord"`
	Snippet     string    `gorm:"size:255" json:"snippet"`
	CreatedAt   time.Time `json:"createdAt"`
}
