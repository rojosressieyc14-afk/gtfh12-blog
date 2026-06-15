package model

import "time"

type UserApiKey struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	UserID       uint       `gorm:"index;not null" json:"userId"`
	User         User       `json:"user,omitempty"`
	Provider     string     `gorm:"size:20;not null" json:"provider"`
	EncryptedKey string     `gorm:"type:text;not null" json:"-"`
	KeyPrefix    string     `gorm:"size:20" json:"keyPrefix"`
	BaseURL      string     `gorm:"size:255" json:"baseURL"`
	LastUsedAt   *time.Time `json:"lastUsedAt"`
	CreatedAt    time.Time  `json:"createdAt"`
	UpdatedAt    time.Time  `json:"updatedAt"`
}
