package interface_usecase

import (
	requestmodel "github.com/ShahabazSulthan/friendzy_post/pkg/models/requestModel"
	responsemodel "github.com/ShahabazSulthan/friendzy_post/pkg/models/responseModel"
)

type ICommentUseCase interface {
	AddNewComment(input *requestmodel.CommentRequest) error
	DeleteComment(userId, commentId *string) error
	EditComment(userId, commentText *string, commentId *uint64) error
	FetchPostComments(userId, postId, limit, offset *string) (*[]responsemodel.ParentComments, error)
	
}
