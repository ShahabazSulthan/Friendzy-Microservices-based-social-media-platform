package domain

import "gorm.io/gorm"

type is_verified string

const (
	Blocked  is_verified = "blocked"
	Deleted  is_verified = "deleted"
	Pending  is_verified = "pending"
	Active   is_verified = "active"
	Verified is_verified = "verified"
	Rejected is_verified = "rejected"
)

type User struct {
	gorm.Model
	Name          string
	UserName      string
	Email         string
	Password      string
	Bio           string
	ProfileImgUrl string
	Links         string
	Status        is_verified `gorm:"default:pending"`
}
