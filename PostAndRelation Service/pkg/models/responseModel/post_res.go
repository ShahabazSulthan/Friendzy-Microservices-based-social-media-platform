package responsemodel

import "time"

type PostData struct {
	UserId            uint
	BlueTick          string
	UserName          string
	UserProfileImgURL string
	PostId            uint
	Caption           string
	CreatedAt         time.Time
	PostAge           string
	MediaUrl          []string `gorm:"-"`

	IsLiked       bool
	LikesCount    uint
	CommentsCount uint
}

type PostAge struct {
	AgeMinutes float64 `json:"age_minutes"`
	AgeHours   float64 `json:"age_hours"`
	AgeDays    float64 `json:"age_days"`
	AgeWeeks   float64 `json:"age_weeks"`
	AgeMonths  float64 `json:"age_months"`
	AgeYears   float64 `json:"age_years"`
}
