package responsemodel

import "time"

type NotificationModel struct {
	NotificaitonID     uint64
	UserID             uint64
	ActorID            uint64
	ActorUserName      string
	ActorProfileImgURL string
	ActionType         string
	TargetID           uint64
	TargetType         string
	CreatedAt          time.Time
	CommentText        string
	NotificationAge    string
}

type PostAge struct {
	AgeMinutes float64 `json:"age_minutes"`
	AgeHours   float64 `json:"age_hours"`
	AgeDays    float64 `json:"age_days"`
	AgeWeeks   float64 `json:"age_weeks"`
	AgeMonths  float64 `json:"age_months"`
	AgeYears   float64 `json:"age_years"`
}
