package model

import "time"

type KnowledgeBase struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `gorm:"index;not null" json:"userId"`
	User        User      `json:"user,omitempty"`
	Name        string    `gorm:"size:200;not null" json:"name"`
	Description string    `gorm:"size:500" json:"description"`
	DocCount    int       `gorm:"not null;default:0" json:"docCount"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type KnowledgeDocument struct {
	ID              uint          `gorm:"primaryKey" json:"id"`
	KnowledgeBaseID uint          `gorm:"index;not null" json:"knowledgeBaseId"`
	KnowledgeBase   KnowledgeBase `json:"knowledgeBase,omitempty"`
	UserID          uint          `gorm:"index;not null" json:"userId"`
	User            User          `json:"user,omitempty"`
	Title           string        `gorm:"size:200" json:"title"`
	Content         string        `gorm:"type:longtext;not null" json:"content"`
	SourceType      string        `gorm:"size:20;not null;default:manual" json:"sourceType"`
	ChunkCount      int           `gorm:"not null;default:0" json:"chunkCount"`
	IsPublic        bool          `gorm:"not null;default:false" json:"isPublic"`
	IsMarkdown      bool          `gorm:"not null;default:true" json:"isMarkdown"`
	CategoryID      *uint         `gorm:"index" json:"categoryId"`
	Category        *Category     `json:"category,omitempty"`
	Tags            []Tag         `gorm:"many2many:kb_document_tags;" json:"tags"`
	ViewCount       int64         `gorm:"not null;default:0" json:"viewCount"`
	PublishedAt     *time.Time    `json:"publishedAt"`
	CreatedAt       time.Time     `json:"createdAt"`
	UpdatedAt       time.Time     `json:"updatedAt"`
}

type KbDocumentTag struct {
	KnowledgeDocumentID uint `gorm:"primaryKey"`
	TagID               uint `gorm:"primaryKey"`
}
