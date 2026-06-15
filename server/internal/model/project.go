package model

import "time"

const (
	ProjectDraft     = "draft"
	ProjectPending   = "pending"
	ProjectPublished = "published"
	ProjectRejected  = "rejected"
)

type Project struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	Title        string     `gorm:"size:120;not null" json:"title"`
	Summary      string     `gorm:"size:300" json:"summary"`
	RoleLabel    string     `gorm:"size:120" json:"roleLabel"`
	Duration     string     `gorm:"size:120" json:"duration"`
	TeamLabel    string     `gorm:"size:120" json:"teamLabel"`
	Content      string     `gorm:"type:longtext;not null" json:"content"`
	CoverImage   string     `gorm:"size:255" json:"coverImage"`
	TechStacks   []string   `gorm:"serializer:json" json:"techStacks"`
	Highlights   []string   `gorm:"serializer:json" json:"highlights"`
	Process      []string   `gorm:"serializer:json" json:"process"`
	Challenges   []string   `gorm:"serializer:json" json:"challenges"`
	Solutions    []string   `gorm:"serializer:json" json:"solutions"`
	Results      []string   `gorm:"serializer:json" json:"results"`
	DemoURL      string     `gorm:"size:255" json:"demoUrl"`
	RepoURL      string     `gorm:"size:255" json:"repoUrl"`
	Status       string     `gorm:"size:20;index;not null;default:draft" json:"status"`
	RejectReason string     `gorm:"size:255" json:"rejectReason"`
	IsPrivate    bool       `gorm:"not null;default:false" json:"isPrivate"`
	IsFeatured   bool       `gorm:"not null;default:false" json:"isFeatured"`
	SortOrder    int        `gorm:"not null;default:0" json:"sortOrder"`
	AuthorID     uint       `gorm:"index;not null" json:"authorId"`
	Author       User       `json:"author"`
	PublishedAt  *time.Time `json:"publishedAt"`
	CreatedAt    time.Time  `json:"createdAt"`
	UpdatedAt    time.Time  `json:"updatedAt"`
}
