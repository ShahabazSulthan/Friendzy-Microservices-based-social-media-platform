package domain

import "time"

type OTP struct {
	ID         uint `gorm:"primary key"`
	Email      string
	OTP        int
	Expiration time.Time
}
