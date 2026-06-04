package model

import "time"

type SensitiveWord struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Word      string    `gorm:"size:120;uniqueIndex;not null" json:"word"`
	Category  string    `gorm:"size:40;index;not null;default:custom" json:"category"`
	Note      string    `gorm:"size:255" json:"note"`
	IsEnabled bool      `gorm:"not null;default:true" json:"isEnabled"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
