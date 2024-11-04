package domain

import "time"

type postStatus string

const (
	Normal    postStatus = "normal"
	Archieved postStatus = "archieved"
)

type Post struct {
	PostID uint `gorm:"primarykey"`

	UserID uint `gorm:"not null"`

	Caption string

	CreatedAt time.Time

	PostStatus postStatus `gorm:"default:normal"`
}

type PostMedia struct {
	MediaId uint `gorm:"primarykey"`

	PostID uint `gorm:"not null"`

	Posts Post `gorm:"foreignkey:PostID"`

	MediaUrl string
}

type relationType string

const (
	Follows postStatus = "follows"
	Blocked postStatus = "blocked"
)

type Relationship struct {
	FollowerID uint `gorm:"not null"`

	FollowingID uint `gorm:"not null"`

	RelationType relationType `gorm:"default:follows"`

	UniqueConstraint struct {
		FollowerID  uint `gorm:"uniqueIndex:idx_follower_following"`
		FollowingID uint `gorm:"uniqueIndex:idx_follower_following"`
	} `gorm:"embedded;uniqueIndex:idx_follower_following"`
}

type PostLikes struct {
	UserID uint `gorm:"not null"`

	PostID uint `gorm:"not null"`
	Posts  Post `gorm:"foreignkey:PostID"`

	CreatedAt time.Time `gorm:"autoCreateTime"`

	UniqueConstraint struct {
		UserID uint `gorm:"uniqueIndex:idx_user_post"`
		PostID uint `gorm:"uniqueIndex:idx_user_post"`
	} `gorm:"embedded;uniqueIndex:idx_user_post"`
}

type Comment struct {
	CommentId uint `gorm:"primarykey"`

	PostID uint `gorm:"not null"`
	Posts  Post `gorm:"foreignkey:PostID"`

	UserID uint `gorm:"not null"`

	CommentText string `gorm:"not null"`

	ParentCommentID uint `gorm:"default:0"`

	CreatedAt time.Time `gorm:"autoCreateTime"`
}
