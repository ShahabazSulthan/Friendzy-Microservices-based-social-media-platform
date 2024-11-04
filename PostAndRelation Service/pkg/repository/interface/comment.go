package interface_repo

import (
	requestmodel "github.com/ShahabazSulthan/friendzy_post/pkg/models/requestModel"
	responsemodel "github.com/ShahabazSulthan/friendzy_post/pkg/models/responseModel"
)

type ICommentRepo interface {
	CheckingCommentHierarchy(input *uint64) (bool, error)
	AddComment(input *requestmodel.CommentRequest) error
	DeleteCommentAndReturnIsParentStat(userId, commentId *string) (bool, error)
	DeleteChildComments(parentCommentId *string) error
	EditComment(userId, commentText *string, commentId *uint64) error
	FetchParentCommentsOfPost(userId, postId, limit, offset *string) (*[]responsemodel.ParentComments, error)
	FetchChildCommentsOfComment(parentCommentId *uint) (*[]responsemodel.ChildComments, error)
	FindCommentCreatorId(CommentId *uint64) (*string, error)
	CalculateCommntAge(postID int) (string, error)
}
