package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	requestmodel "github.com/ShahabazSulthan/friendzy_post/pkg/models/requestModel"
	responsemodel "github.com/ShahabazSulthan/friendzy_post/pkg/models/responseModel"
	"github.com/ShahabazSulthan/friendzy_post/pkg/pb"
	interface_repo "github.com/ShahabazSulthan/friendzy_post/pkg/repository/interface"
	interface_usecase "github.com/ShahabazSulthan/friendzy_post/pkg/usecase/interface"
	cache "github.com/ShahabazSulthan/friendzy_post/pkg/utils/Cache"
	interface_kafka "github.com/ShahabazSulthan/friendzy_post/pkg/utils/Kakfa/interface"
	"github.com/go-redis/redis/v8"
)

type CommentUseCase struct {
	CommentRepo interface_repo.ICommentRepo
	AuthClient  pb.AuthServiceClient
	Kafka       interface_kafka.IKafkaProducer
	PostRepo    interface_repo.IPostRepo
	Redis       *redis.Client
	Ctx         context.Context
}

func NewCommentUsecase(commentRepo interface_repo.ICommentRepo,
	authClient *pb.AuthServiceClient,
	kafka interface_kafka.IKafkaProducer,
	postRepo interface_repo.IPostRepo, redis *redis.Client, ctx context.Context) interface_usecase.ICommentUseCase {
	return &CommentUseCase{
		CommentRepo: commentRepo,
		AuthClient:  *authClient,
		Kafka:       kafka,
		PostRepo:    postRepo,
		Redis:       redis,
		Ctx:         ctx,
	}
}

func (c *CommentUseCase) AddNewComment(input *requestmodel.CommentRequest) error {
	// Check if the comment is a reply to another reply
	if input.ParentCommentId != 0 {
		isReplyToReply, err := c.CommentRepo.CheckingCommentHierarchy(&input.ParentCommentId)
		if err != nil {
			fmt.Println("Error checking comment hierarchy:", err)
			return err
		}
		if isReplyToReply {
			return errors.New("you can't reply to a comment reply")
		}
	}

	// Add the comment to the repository
	err := c.CommentRepo.AddComment(input)
	if err != nil {
		return err
	}

	// Temporary log message to indicate successful addition
	if input.ParentCommentId == 0 {
		fmt.Println("New comment added on post:", input.PostId)
	} else {
		fmt.Println("Reply added to comment:", input.ParentCommentId)
	}

	// Kafka notification logic is commented out
	
		var message requestmodel.KafkaNotification

		if input.ParentCommentId == 0 {
			strPostId := fmt.Sprint(input.PostId)
			PostCreatorId, err := c.PostRepo.GetPostCreatorId(&strPostId)
			if err != nil {
				return err
			}
			message.UserID = *PostCreatorId
			message.ActorID = input.UserId
			message.ActionType = "comment"
			message.TargetID = fmt.Sprint(input.PostId)
			message.TargetType = "post"
			message.CommentText = input.CommentText
			message.CreatedAt = time.Now()
		} else {
			ParentCommentCreatorId, err := c.CommentRepo.FindCommentCreatorId(&input.ParentCommentId)
			if err != nil {
				return err
			}
			message.UserID = *ParentCommentCreatorId
			message.ActorID = input.UserId
			message.ActionType = "reply"
			message.TargetID = fmt.Sprint(input.ParentCommentId)
			message.TargetType = "comment"
			message.CommentText = input.CommentText
			message.CreatedAt = time.Now()
		}

		if message.UserID != message.ActorID {
			err = c.Kafka.KafkaNotificationProducer(&message)
			if err != nil {
				return err
			}
		}


	cacheKey2 := "userFeed"
	cacheKey3 := "Homefeed"

	cache.DeleteFeedEntry(c.Ctx, c.Redis, cacheKey2)
	cache.DeleteFeedEntry(c.Ctx, c.Redis, cacheKey3)
	return nil
}

func (c *CommentUseCase) DeleteComment(userId, commentId *string) error {

	isParent, err := c.CommentRepo.DeleteCommentAndReturnIsParentStat(userId, commentId)
	if err != nil {
		return err
	}

	if isParent {
		err = c.CommentRepo.DeleteChildComments(commentId)
		if err != nil {
			return err
		}
	}
	cacheKey2 := "userFeed"
	cacheKey3 := "Homefeed"

	cache.DeleteFeedEntry(c.Ctx, c.Redis, cacheKey2)
	cache.DeleteFeedEntry(c.Ctx, c.Redis, cacheKey3)
	return nil
}

func (c *CommentUseCase) EditComment(userId, commentText *string, commentId *uint64) error {

	err := c.CommentRepo.EditComment(userId, commentText, commentId)
	if err != nil {
		return err
	}
	
	return nil
}

func (c *CommentUseCase) FetchPostComments(userId, postId, limit, offset *string) (*[]responsemodel.ParentComments, error) {
	parentComments, err := c.CommentRepo.FetchParentCommentsOfPost(userId, postId, limit, offset)
	if err != nil {
		return nil, err
	}

	for i := range *parentComments {
		// Fetch child comments for each parent comment
		childComments, err := c.CommentRepo.FetchChildCommentsOfComment(&(*parentComments)[i].CommentId)
		if err != nil {
			return nil, err
		}

		// Handle child comments' user data and comment age
		for j := range *childComments {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
			defer cancel()

			userData, err := c.AuthClient.GetUserDetailsLiteForPostView(ctx, &pb.RequestUserId{UserId: fmt.Sprint((*childComments)[j].UserID)})
			if err != nil || userData.ErrorMessage != "" {
				return nil, fmt.Errorf("error fetching user data: %v %s", err, userData.ErrorMessage)
			}

			(*childComments)[j].UseName = userData.UserName
			(*childComments)[j].UserProfileImgURL = userData.UserProfileImgURL
			(*childComments)[j].CommentAge, _ = c.CommentRepo.CalculateCommntAge(int((*childComments)[j].CommentId))
		}

		// Fetch parent comment's user data and comment age
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		userData, err := c.AuthClient.GetUserDetailsLiteForPostView(ctx, &pb.RequestUserId{UserId: fmt.Sprint((*parentComments)[i].UserID)})
		if err != nil || userData.ErrorMessage != "" {
			return nil, fmt.Errorf("error fetching user data: %v %s", err, userData.ErrorMessage)
		}

		(*parentComments)[i].UseName = userData.UserName
		(*parentComments)[i].UserProfileImgURL = userData.UserProfileImgURL
		(*parentComments)[i].CommentAge, _ = c.CommentRepo.CalculateCommntAge(int((*parentComments)[i].CommentId))

		// Assign child comments to the parent comment
		(*parentComments)[i].ChildComments = *childComments
	}

	return parentComments, nil
}
