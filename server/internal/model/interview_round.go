package model

import "time"

type InterviewRound struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	SessionID   uint      `gorm:"index;not null" json:"sessionId"`
	RoundNumber int       `gorm:"not null" json:"roundNumber"`
	Question    string    `gorm:"type:text;not null" json:"question"`
	Answer      string    `gorm:"type:text" json:"answer"`
	Score       int       `json:"score"`
	Feedback    string    `gorm:"type:text" json:"feedback"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
