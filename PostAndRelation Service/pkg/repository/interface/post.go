package interface_repo

import (
	requestmodel "github.com/ShahabazSulthan/friendzy_post/pkg/models/requestModel"
	responsemodel "github.com/ShahabazSulthan/friendzy_post/pkg/models/responseModel"
)

type IPostRepo interface {
	AddNewPost(postData *requestmodel.AddPostData) error
	DeletePostById(postId, userId *string) error
	DeletePostMedias(postId *string) error
	EditPost(inputData *requestmodel.EditPost) error

	GetAllActivePostByUser(userId, limit, offset *string) (*[]responsemodel.PostData, error)
	GetPostMediaById(postId *string) (*[]string, error)
	GetPostCountOfUser(userId *string) (*uint, *error)
	CalculatePostAge(postID int) (string, error)

	LikePost(postId, userId *string) (bool, error)
	UnLikePost(postId, userId *string) error
	RemovePostLikesByPostId(postId *string) error
	GetPostCreatorId(postId *string) (*string, error)

	GetPostLikeAndCommentsCount(postId *string) (*responsemodel.LikeCommentCounts, error)
	GetAllActiveRelatedPostsForHomeScreen(userId, limit, offset *string) (*[]responsemodel.PostData, error)
	GetMostLovedPostsFromGlobalUser(userId, limit, offset string) (*[]responsemodel.PostData, error)
	GetRandomPosts(limit, offset string) (*[]responsemodel.PostData, error)
}
