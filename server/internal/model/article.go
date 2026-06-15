package model

import "time"

const (
	ArticleDraft     = "draft"
	ArticlePending   = "pending"
	ArticlePublished = "published"
	ArticleRejected  = "rejected"
)

type Article struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	Title        string     `gorm:"size:120;not null" json:"title"`
	Summary      string     `gorm:"size:300" json:"summary"`
	Content      string     `gorm:"type:longtext;not null" json:"content"`
	CoverImage   string     `gorm:"size:255" json:"coverImage"`
	Status       string     `gorm:"size:20;index;not null;default:draft" json:"status"`
	RejectReason string     `gorm:"size:255" json:"rejectReason"`
	AuthorID     uint       `gorm:"index;not null" json:"authorId"`
	Author       User       `json:"author"`
	CategoryID   *uint      `gorm:"index" json:"categoryId"`
	Category     *Category  `json:"category"`
	Tags         []Tag      `gorm:"many2many:article_tags;" json:"tags"`
	Comments     []Comment  `json:"comments,omitempty"`
	IsPrivate    bool       `gorm:"not null;default:false" json:"isPrivate"`
	ViewCount    int64      `gorm:"not null;default:0" json:"viewCount"`
	LikesCount   int64      `gorm:"-" json:"likesCount"`
	FavoritesCount int64    `gorm:"-" json:"favoritesCount"`
	IsLiked      bool       `gorm:"-" json:"isLiked"`
	IsFavorited  bool       `gorm:"-" json:"isFavorited"`
	PublishedAt  *time.Time `json:"publishedAt"`
	CreatedAt    time.Time  `json:"createdAt"`
	UpdatedAt    time.Time  `json:"updatedAt"`
}

type ArticleReview struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	ArticleID   uint      `gorm:"index;not null" json:"articleId"`
	ReviewerID  uint      `gorm:"index;not null" json:"reviewerId"`
	Action      string    `gorm:"size:16;not null" json:"action"`
	Reason      string    `gorm:"size:255" json:"reason"`
	CreatedAt   time.Time `json:"createdAt"`
}

type ArticleLike struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ArticleID uint      `gorm:"uniqueIndex:idx_article_like;not null" json:"articleId"`
	UserID    uint      `gorm:"uniqueIndex:idx_article_like;not null" json:"userId"`
	CreatedAt time.Time `json:"createdAt"`
}

type ArticleFavorite struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ArticleID uint      `gorm:"uniqueIndex:idx_article_favorite;not null" json:"articleId"`
	UserID    uint      `gorm:"uniqueIndex:idx_article_favorite;not null" json:"userId"`
	CreatedAt time.Time `json:"createdAt"`
}
