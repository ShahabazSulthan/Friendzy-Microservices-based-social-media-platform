package interface_usecase

import (
	requestmodel "github.com/ShahabazSulthan/friendzy_post/pkg/models/requestModel"
	responsemodel "github.com/ShahabazSulthan/friendzy_post/pkg/models/responseModel"
	"github.com/ShahabazSulthan/friendzy_post/pkg/pb"
)

type IPostUseCase interface {
	AddNewPost(data *[]*pb.SingleMedia, caption *string, userId *string) error
	GetAllPosts(userId, limit, offset *string) (*[]responsemodel.PostData, error)
	DeletePost(postId, userId *string) error
	EditPost(request *requestmodel.EditPost) error

	LikePost(postId, userId *string) *error
	UnLikePost(postId, userId *string) *error
	GetMostLovedPostsFromGlobalUser(userId, limit, offset *string) (*[]responsemodel.PostData, error)
	GetAllRelatedPostsForHomeScreen(userId, limit, offset *string) (*[]responsemodel.PostData, error)
	GetRandomPosts(limit, offset *string) (*[]responsemodel.PostData, error)
}
