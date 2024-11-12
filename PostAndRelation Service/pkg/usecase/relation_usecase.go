package usecase

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	requestmodel "github.com/ShahabazSulthan/friendzy_post/pkg/models/requestModel"
	"github.com/ShahabazSulthan/friendzy_post/pkg/pb"
	interface_repo "github.com/ShahabazSulthan/friendzy_post/pkg/repository/interface"
	interface_usecase "github.com/ShahabazSulthan/friendzy_post/pkg/usecase/interface"
	cache "github.com/ShahabazSulthan/friendzy_post/pkg/utils/Cache"
	interface_kafka "github.com/ShahabazSulthan/friendzy_post/pkg/utils/Kakfa/interface"
	"github.com/go-redis/redis/v8"
)

type RelationUsecase struct {
	RelationRepo interface_repo.IRelationRepo
	PostRepo     interface_repo.IPostRepo
	AuthClient   pb.AuthServiceClient
	Kafka        interface_kafka.IKafkaProducer
	Redis        *redis.Client
	Ctx          context.Context
}

func NewRelationUseCase(relationRepo interface_repo.IRelationRepo,
	postRepo interface_repo.IPostRepo,
	authClient pb.AuthServiceClient,
	kafka interface_kafka.IKafkaProducer, redis *redis.Client, ctx context.Context) interface_usecase.IRelationUseCase {
	return &RelationUsecase{
		RelationRepo: relationRepo,
		PostRepo:     postRepo,
		AuthClient:   authClient,
		Kafka:        kafka,
		Redis:        redis,
		Ctx:          ctx,
	}
}

func (r *RelationUsecase) GetCountsForUserProfile(userId *string) (*uint, *uint, *uint, *error) {
	a, b, err := r.RelationRepo.GetFollowerAndFollowingCountofUser(userId)
	if err != nil {
		return nil, nil, nil, err
	}
	c, err := r.PostRepo.GetPostCountOfUser(userId)
	if err != nil {
		return nil, nil, nil, err
	}
	return a, b, c, nil
}

func (r *RelationUsecase) GetFollowersIds(userId *string) (*[]uint64, error) {

	userIdSlice, err := r.RelationRepo.GetFollowersIdsOfUser(userId)
	if err != nil {
		return nil, err
	}
	return userIdSlice, nil
}

func (r *RelationUsecase) GetFollowingsIds(userId *string) (*[]uint64, error) {

	userIdSlice, err := r.RelationRepo.GetFollowingsIdsOfUser(userId)
	if err != nil {
		return nil, err
	}
	return userIdSlice, nil
}

func (r *RelationUsecase) UserAFollowingUserBorNot(userId, userBId *string) (bool, error) {

	resp, err := r.RelationRepo.UserAFollowingUserBorNot(userId, userBId)
	if err != nil {
		return false, err
	}
	cacheKey2 := "userFeed"
	cacheKey3 := "Homefeed"

	cache.DeleteFeedEntry(r.Ctx, r.Redis, cacheKey2)
	cache.DeleteFeedEntry(r.Ctx, r.Redis, cacheKey3)
	return resp, err
}

func (r *RelationUsecase) Follow(userId, userBId *string) *error {

	if *userId == *userBId {
		err := errors.New("user can't follow themselves")
		return &err
	}
	var message requestmodel.KafkaNotification
	// Step 1: Create a context with a timeout for the AuthClient call
	context, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Step 2: Check if the userBId exists using the AuthClient
	resp, err := r.AuthClient.CheckUserExist(context, &pb.RequestUserId{
		UserId: *userBId,
	})
	if err != nil {
		log.Fatal(err)
	}
	if resp.ErrorMessage != "" {
		err = errors.New(resp.ErrorMessage)
		return &err
	}
	if !resp.ExistStatus {
		err = errors.New("no user exists with this ID, enter a valid user ID")
		return &err
	}

	// Step 3: Initiate the follow relationship between userId and userBId
	inserted, err := r.RelationRepo.InitiateFollowRelationship(userId, userBId)
	if err != nil {
		return &err
	}

	// Step 4: If follow was successfully initiated, return nil (success)
	if inserted {
		message.UserID = *userBId
		message.ActorID = *userId
		message.ActionType = "follow"
		message.TargetID = "0"
		message.CreatedAt = time.Now()

		err = r.Kafka.KafkaNotificationProducer(&message)
		if err != nil {
			return &err
		}
		fmt.Printf("User %s successfully followed User %s\n", *userId, *userBId)
	}
	cacheKey2 := "userFeed"
	cacheKey3 := "Homefeed"

	cache.DeleteFeedEntry(r.Ctx, r.Redis, cacheKey2)
	cache.DeleteFeedEntry(r.Ctx, r.Redis, cacheKey3)

	return nil
}

func (r *RelationUsecase) UnFollow(userId, userBId *string) error {

	if *userId == *userBId {
		err := errors.New("dont use same userid")
		return err
	}
	// Step 1: Create a context with a timeout for the AuthClient call
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Check if AuthClient is initialized
	if r.AuthClient == nil {
		return errors.New("AuthClient is not initialized")
	}

	// Check if RelationRepo is initialized
	if r.RelationRepo == nil {
		return errors.New("RelationRepo is not initialized")
	}

	// Step 2: Check if the userBId exists using the AuthClient
	resp, err := r.AuthClient.CheckUserExist(ctx, &pb.RequestUserId{
		UserId: *userBId,
	})
	if err != nil {
		return fmt.Errorf("error checking if user exists: %v", err)
	}

	if resp.ErrorMessage != "" {
		return errors.New(resp.ErrorMessage)
	}

	if !resp.ExistStatus {
		return errors.New("no user exists with this ID, enter a valid user ID")
	}

	// Step 3: Initiate the unfollow relationship between userId and userBId
	err = r.RelationRepo.InitiateUnFollowRelationship(userId, userBId)
	if err != nil {
		return fmt.Errorf("error initiating unfollow: %v", err)
	}

	// Step 4: Log or confirm the successful unfollow
	fmt.Printf("User %s successfully unfollowed User %s\n", *userId, *userBId)

	cacheKey2 := "userFeed"
	cacheKey3 := "Homefeed"

	cache.DeleteFeedEntry(r.Ctx, r.Redis, cacheKey2)
	cache.DeleteFeedEntry(r.Ctx, r.Redis, cacheKey3)

	return nil
}
