package model

import "time"

type Comment struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	Content   string     `gorm:"size:500;not null" json:"content"`
	ArticleID uint       `gorm:"index;not null" json:"articleId"`
	Article   Article    `json:"article"`
	UserID    uint       `gorm:"index;not null" json:"userId"`
	User      User       `json:"user"`
	ParentID  *uint      `gorm:"index" json:"parentId"`
	Replies   []Comment  `gorm:"foreignKey:ParentID" json:"replies,omitempty"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
}
