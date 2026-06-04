package model

import "time"

const (
	RoleAdmin  = "admin"
	RoleUser   = "user"
	UserActive = "active"
	UserBanned = "banned"
)

type User struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Username    string    `gorm:"size:32;uniqueIndex;not null" json:"username"`
	Password    string    `gorm:"size:255;not null" json:"-"`
	Role        string    `gorm:"size:16;not null;default:user" json:"role"`
	Status      string    `gorm:"size:16;not null;default:active" json:"status"`
	Avatar      string    `gorm:"size:255" json:"avatar"`
	Headline    string    `gorm:"size:120" json:"headline"`
	CurrentRole string    `gorm:"size:120" json:"currentRole"`
	YearsLabel  string    `gorm:"size:80" json:"yearsLabel"`
	Motto       string    `gorm:"size:255" json:"motto"`
	Location    string    `gorm:"size:120" json:"location"`
	Email       string    `gorm:"size:120" json:"email"`
	ResumeURL   string    `gorm:"size:255" json:"resumeUrl"`
	WebsiteURL  string    `gorm:"size:255" json:"websiteUrl"`
	GithubURL   string    `gorm:"size:255" json:"githubUrl"`
	GiteeURL    string    `gorm:"size:255" json:"giteeUrl"`
	JuejinURL   string    `gorm:"size:255" json:"juejinUrl"`
	CSDNURL     string    `gorm:"size:255" json:"csdnUrl"`
	Skills      []string  `gorm:"serializer:json" json:"skills"`
	FocusAreas  []string  `gorm:"serializer:json" json:"focusAreas"`
	Bio         string    `gorm:"size:2000" json:"bio"`
	BanReason   string    `gorm:"size:255" json:"banReason"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
