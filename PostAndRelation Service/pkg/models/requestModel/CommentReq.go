package requestmodel

type CommentRequest struct {
	PostId          uint64
	UserId          string
	CommentText     string
	ParentCommentId uint64
}
