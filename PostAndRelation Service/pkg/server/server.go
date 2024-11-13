package server

import (
	"context"

	requestmodel "github.com/ShahabazSulthan/friendzy_post/pkg/models/requestModel"
	"github.com/ShahabazSulthan/friendzy_post/pkg/pb"
	interface_usecase "github.com/ShahabazSulthan/friendzy_post/pkg/usecase/interface"
)

type PostService struct {
	PostUseCase     interface_usecase.IPostUseCase
	RelationUseCase interface_usecase.IRelationUseCase
	CommentUseCase  interface_usecase.ICommentUseCase
	pb.PostNrelServiceServer
}

func NewPostAndRelation(postUsecase interface_usecase.IPostUseCase,
	relationusecase interface_usecase.IRelationUseCase,
	commentUsecase interface_usecase.ICommentUseCase) *PostService {
	return &PostService{
		PostUseCase:     postUsecase,
		RelationUseCase: relationusecase,
		CommentUseCase:  commentUsecase,
	}
}

func (p *PostService) AddNewPost(ctx context.Context, req *pb.RequestAddPost) (*pb.ResponseErrorMessageOnly, error) {
	err := p.PostUseCase.AddNewPost(&req.Media, &req.Caption, &req.UserId)
	if err != nil {
		return &pb.ResponseErrorMessageOnly{
			ErrorMessage: err.Error(),
		}, nil
	}

	return &pb.ResponseErrorMessageOnly{}, nil
}

func (p *PostService) GetAllPostByUser(ctx context.Context, req *pb.RequestGetAllPosts) (*pb.ResponseUserPosts, error) {

	respData, err := p.PostUseCase.GetAllPosts(&req.UserId, &req.Limit, &req.OffSet)
	if err != nil {
		return &pb.ResponseUserPosts{
			ErrorMessage: err.Error(),
		}, nil
	}

	var repeatedData []*pb.PostsDataModel
	for i := range *respData {
		repeatedData = append(repeatedData, &pb.PostsDataModel{
			UserId:            uint64((*respData)[i].UserId),
			UserName:          (*respData)[i].UserName,
			UserProfileImgURL: (*respData)[i].UserProfileImgURL,
			PostId:            uint64((*respData)[i].PostId),
			Caption:           (*respData)[i].Caption,
			PostAge:           (*respData)[i].PostAge,
			LikesCount:        uint64((*respData)[i].LikesCount),
			CommentsCount:     uint64((*respData)[i].CommentsCount),
			MediaUrl:          (*respData)[i].MediaUrl,
		})
	}

	return &pb.ResponseUserPosts{
		PostsData: repeatedData,
	}, nil
}

func (p *PostService) DeletePost(ctx context.Context, req *pb.RequestDeletePost) (*pb.ResponseErrorMessageOnly, error) {
	err := p.PostUseCase.DeletePost(&req.PostId, &req.UserId)
	if err != nil {
		return &pb.ResponseErrorMessageOnly{
			ErrorMessage: err.Error(),
		}, nil
	}

	return &pb.ResponseErrorMessageOnly{}, nil

}

func (p *PostService) EditPost(ctx context.Context, req *pb.RequestEditPost) (*pb.ResponseErrorMessageOnly, error) {

	var request requestmodel.EditPost

	request.UserId = req.UserId
	request.PostId = req.PostId
	request.Caption = req.Caption

	err := p.PostUseCase.EditPost(&request)
	if err != nil {
		return &pb.ResponseErrorMessageOnly{
			ErrorMessage: err.Error(),
		}, nil
	}
	return &pb.ResponseErrorMessageOnly{}, nil
}

func (p *PostService) Follow(ctx context.Context, req *pb.RequestFollowUnFollow) (*pb.ResponseErrorMessageOnly, error) {

	err := p.RelationUseCase.Follow(&req.UserId, &req.UserBId)
	if err != nil {
		return &pb.ResponseErrorMessageOnly{
			ErrorMessage: (*err).Error(),
		}, nil
	}
	return &pb.ResponseErrorMessageOnly{}, nil
}

func (p *PostService) UnFollow(ctx context.Context, req *pb.RequestFollowUnFollow) (*pb.ResponseErrorMessageOnly, error) {

	err := p.RelationUseCase.UnFollow(&req.UserId, &req.UserBId)
	if err != nil {
		return &pb.ResponseErrorMessageOnly{
			ErrorMessage: (err).Error(),
		}, nil
	}
	return &pb.ResponseErrorMessageOnly{}, nil
}

func (p *PostService) GetCountsForUserProfile(ctx context.Context, req *pb.RequestUserIdPnR) (*pb.ResponseGetCounts, error) {

	followersCount, followingCount, postsCount, err := p.RelationUseCase.GetCountsForUserProfile(&req.UserId)
	if err != nil {
		return &pb.ResponseGetCounts{
			ErrorMessage: (*err).Error(),
		}, nil
	}

	return &pb.ResponseGetCounts{
		PostCount:      uint64(*postsCount),
		FollowerCount:  uint64(*followersCount),
		FollowingCount: uint64(*followingCount),
	}, nil

}

func (p *PostService) GetFollowersIds(ctx context.Context, req *pb.RequestUserIdPnR) (*pb.ResposneGetUsersIds, error) {

	userIdSlce, err := p.RelationUseCase.GetFollowersIds(&req.UserId)
	if err != nil {
		return &pb.ResposneGetUsersIds{
			ErrorMessage: err.Error(),
		}, nil
	}

	return &pb.ResposneGetUsersIds{
		UserIds: *userIdSlce,
	}, nil
}

func (p *PostService) GetFollowingsIds(ctx context.Context, req *pb.RequestUserIdPnR) (*pb.ResposneGetUsersIds, error) {

	userIdSlce, err := p.RelationUseCase.GetFollowingsIds(&req.UserId)
	if err != nil {
		return &pb.ResposneGetUsersIds{
			ErrorMessage: err.Error(),
		}, nil
	}

	return &pb.ResposneGetUsersIds{
		UserIds: *userIdSlce,
	}, nil
}

func (u *PostService) UserAFollowingUserBorNot(ctx context.Context, req *pb.RequestFollowUnFollow) (*pb.ResponseUserABrelation, error) {
	resp, err := u.RelationUseCase.UserAFollowingUserBorNot(&req.UserId, &req.UserBId)
	if err != nil {
		return &pb.ResponseUserABrelation{
			ErrorMessage: err.Error(),
		}, nil
	}
	return &pb.ResponseUserABrelation{
		BoolStat: resp,
	}, nil
}

func (u *PostService) LikePost(ctx context.Context, req *pb.RequestLikeUnlikePost) (*pb.ResponseErrorMessageOnly, error) {
	err := u.PostUseCase.LikePost(&req.PostId, &req.UserId)
	if err != nil {
		return &pb.ResponseErrorMessageOnly{
			ErrorMessage: (*err).Error(),
		}, nil
	}
	return &pb.ResponseErrorMessageOnly{}, nil
}

func (u *PostService) UnLikePost(ctx context.Context, req *pb.RequestLikeUnlikePost) (*pb.ResponseErrorMessageOnly, error) {
	err := u.PostUseCase.UnLikePost(&req.PostId, &req.UserId)
	if err != nil {
		return &pb.ResponseErrorMessageOnly{
			ErrorMessage: (*err).Error(),
		}, nil
	}
	return &pb.ResponseErrorMessageOnly{}, nil
}

func (u *PostService) AddComment(ctx context.Context, req *pb.RequestAddComment) (*pb.ResponseErrorMessageOnly, error) {
	var input requestmodel.CommentRequest

	input.UserId = req.UserId
	input.PostId = req.PostId
	input.CommentText = req.CommentText
	input.ParentCommentId = req.ParentCommentId

	err := u.CommentUseCase.AddNewComment(&input)
	if err != nil {
		return &pb.ResponseErrorMessageOnly{
			ErrorMessage: err.Error(),
		}, nil
	}
	return &pb.ResponseErrorMessageOnly{}, nil
}

func (u *PostService) DeleteComment(ctx context.Context, req *pb.RequestCommentDelete) (*pb.ResponseErrorMessageOnly, error) {

	err := u.CommentUseCase.DeleteComment(&req.UserId, &req.CommentId)
	if err != nil {
		return &pb.ResponseErrorMessageOnly{
			ErrorMessage: err.Error(),
		}, nil
	}
	return &pb.ResponseErrorMessageOnly{}, nil
}

func (u *PostService) EditComment(ctx context.Context, req *pb.RequestEditComment) (*pb.ResponseErrorMessageOnly, error) {
	err := u.CommentUseCase.EditComment(&req.UserId, &req.CommentText, &req.CommentId)
	if err != nil {
		return &pb.ResponseErrorMessageOnly{
			ErrorMessage: err.Error(),
		}, nil
	}
	return &pb.ResponseErrorMessageOnly{}, nil
}

func (u *PostService) FetchPostComments(ctx context.Context, req *pb.RequestFetchComments) (*pb.ResponseFetchComments, error) {
	respData, err := u.CommentUseCase.FetchPostComments(&req.UserId, &req.PostId, &req.Limit, &req.OffSet)
	if err != nil {
		return &pb.ResponseFetchComments{
			ErrorMessage: err.Error(),
		}, nil
	}

	var ChildComments []*pb.ChildComments
	var ParentComments []*pb.ParentComments

	for i := range *respData {
		ChildComments = []*pb.ChildComments{}
		for j := range (*respData)[i].ChildComments {
			ChildComments = append(ChildComments, &pb.ChildComments{
				CommentId:         uint64(((*respData)[i].ChildComments)[j].CommentId),
				PostId:            uint64(((*respData)[i].ChildComments)[j].PostID),
				UserId:            uint64(((*respData)[i].ChildComments)[j].UserID),
				UserName:          ((*respData)[i].ChildComments)[j].UseName,
				UserProfileImgURL: ((*respData)[i].ChildComments)[j].UserProfileImgURL,
				CommentText:       ((*respData)[i].ChildComments)[j].CommentText,
				ParentCommentID:   uint64(((*respData)[i].ChildComments)[j].ParentCommentID),
				CommentAge:        ((*respData)[i].ChildComments)[j].CommentAge,
			})
		}
		ParentComments = append(ParentComments, &pb.ParentComments{
			CommentId:         uint64((*respData)[i].CommentId),
			PostId:            uint64((*respData)[i].PostID),
			UserId:            uint64((*respData)[i].UserID),
			UserName:          (*respData)[i].UseName,
			UserProfileImgURL: (*respData)[i].UserProfileImgURL,
			CommentText:       (*respData)[i].CommentText,
			CommentAge:        (*respData)[i].CommentAge,
			ChildCommentCount: uint64(len(ChildComments)),
			ChildComments:     ChildComments,
		})

	}

	return &pb.ResponseFetchComments{
		ParentCommentsCount: uint64(len(ParentComments)),
		ParentComments:      ParentComments,
	}, nil
}

func (u *PostService) GetMostLovedPostsFromGlobalUser(ctx context.Context, req *pb.RequestGetAllPosts) (*pb.ResponseUserPosts, error) {
	respData, err := u.PostUseCase.GetMostLovedPostsFromGlobalUser(&req.UserId, &req.Limit, &req.OffSet)
	if err != nil {
		return &pb.ResponseUserPosts{
			ErrorMessage: err.Error(),
		}, nil
	}

	var repeatedData []*pb.PostsDataModel
	for i := range *respData {
		repeatedData = append(repeatedData, &pb.PostsDataModel{
			UserId:            uint64((*respData)[i].UserId),
			BlueTick:          (*respData)[i].BlueTick,
			UserName:          (*respData)[i].UserName,
			UserProfileImgURL: (*respData)[i].UserProfileImgURL,
			PostId:            uint64((*respData)[i].PostId),
			LikeStatus:        (*respData)[i].IsLiked,
			Caption:           (*respData)[i].Caption,
			LikesCount:        uint64((*respData)[i].LikesCount),
			CommentsCount:     uint64((*respData)[i].CommentsCount),
			PostAge:           (*respData)[i].PostAge,
			MediaUrl:          (*respData)[i].MediaUrl,
		})
	}

	return &pb.ResponseUserPosts{
		PostsData: repeatedData,
	}, nil
}

func (u *PostService) GetAllRelatedPostsForHomeScreen(ctx context.Context, req *pb.RequestGetAllPosts) (*pb.ResponseUserPosts, error) {
	respData, err := u.PostUseCase.GetAllRelatedPostsForHomeScreen(&req.UserId, &req.Limit, &req.OffSet)
	if err != nil {
		return &pb.ResponseUserPosts{
			ErrorMessage: err.Error(),
		}, nil
	}

	var repeatedData []*pb.PostsDataModel
	for i := range *respData {
		repeatedData = append(repeatedData, &pb.PostsDataModel{
			UserId:            uint64((*respData)[i].UserId),
			BlueTick:          (*respData)[i].BlueTick,
			UserName:          (*respData)[i].UserName,
			UserProfileImgURL: (*respData)[i].UserProfileImgURL,
			PostId:            uint64((*respData)[i].PostId),
			LikeStatus:        (*respData)[i].IsLiked,
			Caption:           (*respData)[i].Caption,
			LikesCount:        uint64((*respData)[i].LikesCount),
			CommentsCount:     uint64((*respData)[i].CommentsCount),
			PostAge:           (*respData)[i].PostAge,
			MediaUrl:          (*respData)[i].MediaUrl,
		})
	}

	return &pb.ResponseUserPosts{
		PostsData: repeatedData,
	}, nil
}

func (u *PostService) GetRandomPosts(ctx context.Context, req *pb.RequestGetRandomPosts) (*pb.ResponseUserPosts, error) {
	respData, err := u.PostUseCase.GetRandomPosts(&req.Limit, &req.OffSet)
	if err != nil {
		return &pb.ResponseUserPosts{
			ErrorMessage: err.Error(),
		}, nil
	}

	var repeatedData []*pb.PostsDataModel
	for i := range *respData {
		repeatedData = append(repeatedData, &pb.PostsDataModel{
			UserId:            uint64((*respData)[i].UserId),
			BlueTick:          (*respData)[i].BlueTick,
			UserName:          (*respData)[i].UserName,
			UserProfileImgURL: (*respData)[i].UserProfileImgURL,
			PostId:            uint64((*respData)[i].PostId),
			LikeStatus:        (*respData)[i].IsLiked,
			Caption:           (*respData)[i].Caption,
			LikesCount:        uint64((*respData)[i].LikesCount),
			CommentsCount:     uint64((*respData)[i].CommentsCount),
			PostAge:           (*respData)[i].PostAge,
			MediaUrl:          (*respData)[i].MediaUrl,
		})
	}

	return &pb.ResponseUserPosts{
		PostsData: repeatedData,
	}, nil
}
