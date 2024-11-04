package responsemodel_post

type CommentResponse struct {
	PostID          string `json:"PostId,omitempty"`
	UserID          string `json:"UserId,omitempty"`
	CommentText     string `json:"CommentText,omitempty" `
	ParentCommentID string `json:"ParentCommentId,omitempty" `
}

type CommentEditResponse struct {
	UserId      string `json:"UserId,omitempty"`
	CommentId   string `json:"CommentId,omitempty" `
	CommentText string `json:"CommentText,omitempty" `
}
