package domain

import "time"

type Notification struct {
	NotificaitonID uint   `gorm:"primarykey"`
	UserID         uint   `gorm:"not null"`
	ActorID        uint   `gorm:"not null"`
	ActionType     string `gorm:"not null"`
	TargetID       uint   `gorm:"default:0"`
	TargetType     string
	CommentText    string
	CreatedAt      time.Time
}
