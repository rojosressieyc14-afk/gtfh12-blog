package model

import (
	"time"
)

const (
	TargetTypeArticle = "article"
	TargetTypeProject = "project"

	RiskLevelLow  = "low"
	RiskLevelMid  = "mid"
	RiskLevelHigh = "high"
)

type AIReviewRecord struct {
	ID                 uint      `gorm:"primaryKey" json:"id"`
	TargetType         string    `gorm:"size:20;index:idx_ai_review_target;not null" json:"targetType"`
	TargetID           uint      `gorm:"index:idx_ai_review_target;not null" json:"targetId"`
	OperatorID         uint      `gorm:"index;not null" json:"operatorId"`
	Operator           User      `json:"operator"`
	RiskLevel          string    `gorm:"size:10;not null;default:low" json:"riskLevel"`
	RiskLabels         []string  `gorm:"serializer:json" json:"riskLabels"`
	Summary            string    `gorm:"size:500" json:"summary"`
	SuspiciousSegments []string  `gorm:"serializer:json" json:"suspiciousSegments"`
	Suggestion         string    `gorm:"size:500" json:"suggestion"`
	ModelName          string    `gorm:"size:60" json:"modelName"`
	CreatedAt          time.Time `json:"createdAt"`
	UpdatedAt          time.Time `json:"updatedAt"`
}
