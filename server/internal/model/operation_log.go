package model

import "time"

type OperationLog struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	OperatorID  uint      `gorm:"index;not null" json:"operatorId"`
	Operator    User      `json:"operator"`
	Action      string    `gorm:"size:80;not null" json:"action"`
	TargetType  string    `gorm:"size:40;not null" json:"targetType"`
	TargetID    uint      `gorm:"not null" json:"targetId"`
	Description string    `gorm:"size:255;not null" json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
}
