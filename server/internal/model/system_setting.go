package model

import "time"

type SystemSetting struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Key       string    `gorm:"size:80;uniqueIndex;not null" json:"key"`
	Value     string    `gorm:"size:255;not null" json:"value"`
	Note      string    `gorm:"size:255" json:"note"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
