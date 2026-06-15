package model

import "time"

const (
	InterviewStatusInProgress = "in_progress"
	InterviewStatusCompleted  = "completed"
)

type InterviewSession struct {
	ID             uint             `gorm:"primaryKey" json:"id"`
	UserID         uint             `gorm:"index;not null" json:"userId"`
	User           User             `json:"user"`
	Position       string           `gorm:"size:200;not null" json:"position"`
	ResumeText     string           `gorm:"type:text" json:"resumeText"`
	TotalQuestions int              `gorm:"not null;default:5" json:"totalQuestions"`
	Status         string           `gorm:"size:20;not null;default:in_progress" json:"status"`
	ApiKeyID       *uint            `gorm:"index" json:"apiKeyId"`
	CreatedAt      time.Time        `json:"createdAt"`
	UpdatedAt      time.Time        `json:"updatedAt"`
	Rounds         []InterviewRound `gorm:"foreignKey:SessionID" json:"rounds,omitempty"`
}
