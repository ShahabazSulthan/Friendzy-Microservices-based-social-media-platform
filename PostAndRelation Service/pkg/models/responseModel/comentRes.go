package responsemodel

import "time"

type ChildComments struct {
	CommentId         uint
	PostID            uint
	UserID            uint
	UseName           string `gorm:"-"`
	UserProfileImgURL string `gorm:"-"`
	CommentText       string
	ParentCommentID   uint
	CreatedAt         time.Time
	CommentAge        string `gorm:"-"`
}

type ParentComments struct {
	CommentId          uint
	PostID             uint
	UserID             uint
	UseName            string `gorm:"-"`
	UserProfileImgURL  string `gorm:"-"`
	CommentText        string
	ParentCommentID    uint
	CreatedAt          time.Time
	CommentAge         string          `gorm:"-"`
	ChildCommentsCount uint            `gorm:"-"`
	ChildComments      []ChildComments `gorm:"-"`
}

type LikeCommentCounts struct {
	LikesCount    uint `gorm:"column:likes_count"`
	CommentsCount uint `gorm:"column:comments_count"`
}
