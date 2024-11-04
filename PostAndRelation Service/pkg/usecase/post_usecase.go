package usecase

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
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

type PostUseCase struct {
	PostRepo      interface_repo.IPostRepo
	AuthClient    pb.AuthServiceClient
	KafkaProducer interface_kafka.IKafkaProducer
	Redis         *redis.Client
	Ctx           context.Context
}

func NewPostUseCase(postRepo interface_repo.IPostRepo, authClient pb.AuthServiceClient, kafkaProducer interface_kafka.IKafkaProducer, redis *redis.Client, ctx context.Context) interface_usecase.IPostUseCase {
	return &PostUseCase{
		PostRepo:      postRepo,
		AuthClient:    authClient,
		KafkaProducer: kafkaProducer,
		Redis:         redis,
		Ctx:           ctx,
	}
}

func saveMediaLocally(media *[]byte, filePath string) error {
	// Write the file to the local directory
	return os.WriteFile(filePath, *media, 0644)
}

func (p *PostUseCase) AddNewPost(data *[]*pb.SingleMedia, caption *string, userId *string) error {

	localFolder := "PostAndRelationService/posts/" // Local directory path

	// Create the folder if it doesn't exist
	err := os.MkdirAll(localFolder, os.ModePerm)
	if err != nil {
		fmt.Printf("Error creating directory: %v\n", err)
		return err
	}

	var postData requestmodel.AddPostData

	for i, file := range *data {
		
		extension := strings.ReplaceAll(file.ContentType, "/", ".")  
		fileName := fmt.Sprintf("%s_%d.%s", *userId, i+1, extension) 

		// Create the file path in the local folder
		filePath := filepath.Join(localFolder, fileName)

		// Save the media content locally
		err := saveMediaLocally(&file.Media, filePath)
		if err != nil {
			fmt.Printf("Error saving file %d: %v\n", i+1, err)
			return err
		}

		// Add the file path as the media URL (local path)
		postData.MediaURLs = append(postData.MediaURLs, filePath)
	}

	postData.Caption = caption
	postData.UserId = userId

	err = p.PostRepo.AddNewPost(&postData)
	if err != nil {
		return err
	}

	cacheKey2 := "userFeed"
	cacheKey3 := "Homefeed"

	cache.DeleteFeedEntry(p.Ctx, p.Redis, cacheKey2)
	cache.DeleteFeedEntry(p.Ctx, p.Redis, cacheKey3)
	return nil
}

func (p *PostUseCase) GetAllPosts(userId, limit, offset *string) (*[]responsemodel.PostData, error) {
	// Step 1: Fetch all active posts from the repository
	postData, err := p.PostRepo.GetAllActivePostByUser(userId, limit, offset)
	if err != nil {
		return nil, err
	}

	// Step 2: Get user details via the AuthClient service
	context, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	userData, err := p.AuthClient.GetUserDetailsLiteForPostView(context, &pb.RequestUserId{UserId: *userId})
	if err != nil {
		log.Fatal(err) // It's generally better to return the error instead of log.Fatal in production
		return nil, err
	}

	// Handle error from user data
	if userData.ErrorMessage != "" {
		return nil, errors.New(userData.ErrorMessage)
	}

	// Step 3: Loop through each post and enrich data with media, like/comment counts, and user details
	for i, post := range *postData {
		// Set user details (username and profile image)
		(*postData)[i].UserName = userData.UserName
		(*postData)[i].UserProfileImgURL = userData.UserProfileImgURL

		// Get post ID as string
		postIdString := fmt.Sprint(post.PostId)

		// Fetch associated media URLs for the post
		postMedias, err := p.PostRepo.GetPostMediaById(&postIdString)
		if err != nil {
			return nil, err
		}
		(*postData)[i].MediaUrl = *postMedias
		LikeCommentCount, err := p.PostRepo.GetPostLikeAndCommentsCount(&postIdString)
		if err != nil {
			return nil, err
		}

		(*postData)[i].LikesCount = LikeCommentCount.LikesCount
		(*postData)[i].CommentsCount = LikeCommentCount.CommentsCount

		// Calculate the age of the post (time since creation)
		(*postData)[i].PostAge, err = p.PostRepo.CalculatePostAge(int((*postData)[i].PostId))
		if err != nil {
			return nil, err
		}
	}

	// Return the enriched post data
	return postData, nil
}

func (p *PostUseCase) DeletePost(postId, userId *string) error {

	err := p.PostRepo.DeletePostMedias(postId)
	if err != nil {
		return err
	}
	err = p.PostRepo.DeletePostById(postId, userId)
	if err != nil {
		return err
	}

	cacheKey2 := "userFeed"
	cacheKey3 := "Homefeed"

	cache.DeleteFeedEntry(p.Ctx, p.Redis, cacheKey2)
	cache.DeleteFeedEntry(p.Ctx, p.Redis, cacheKey3)

	return nil
}

func (p *PostUseCase) EditPost(request *requestmodel.EditPost) error {

	err := p.PostRepo.EditPost(request)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (p *PostUseCase) LikePost(postId, userId *string) *error {

	postCrreatorId, err := p.PostRepo.GetPostCreatorId(postId)
	if err != nil {
		return &err
	}

	inserted, err := p.PostRepo.LikePost(postId, userId)
	if err != nil {
		fmt.Println("err", err)
		return &err
	}

	if *postCrreatorId != *postId {
		//var message requestmodel.KafkaNotification
		if inserted {
			// message.UserID = *postCrreatorId
			// message.ActorID = *userId
			// message.ActionType = "like"
			// message.TargetID = *postId
			// message.TargetType = "post"
			// message.CreatedAt = time.Now()

			// err := p.kafkaProducer.KafkaNotificationProducer(&message)
			// if err != nil {
			// 	return &err
			// }
			fmt.Println("Post liked successfully by user:", *userId)
		}
	}

	cacheKey2 := "userFeed"
	cacheKey3 := "Homefeed"

	cache.DeleteFeedEntry(p.Ctx, p.Redis, cacheKey2)
	cache.DeleteFeedEntry(p.Ctx, p.Redis, cacheKey3)
	return nil
}

func (p *PostUseCase) UnLikePost(postId, userId *string) *error {
	err := p.PostRepo.UnLikePost(postId, userId)
	if err != nil {
		fmt.Println("Err: ", err)
		return &err
	}

	cacheKey2 := "userFeed"
	cacheKey3 := "Homefeed"

	cache.DeleteFeedEntry(p.Ctx, p.Redis, cacheKey2)
	cache.DeleteFeedEntry(p.Ctx, p.Redis, cacheKey3)
	return nil
}

func (p *PostUseCase) GetAllRelatedPostsForHomeScreen(userId, limit, offset *string) (*[]responsemodel.PostData, error) {
	cacheKey := "Homefeed"
	// Attempt to retrieve from cache
	var cachedPosts []responsemodel.PostData
	if err := cache.CacheGet(p.Ctx, p.Redis, cacheKey, &cachedPosts); err == nil {
		// Cache hit
		return &cachedPosts, nil
	} else if err != redis.Nil {
		// Redis error (not cache miss)
		return nil, fmt.Errorf("redis error: %w", err)
	}

	postData, err := p.PostRepo.GetAllActiveRelatedPostsForHomeScreen(userId, limit, offset)
	if err != nil {
		return nil, err
	}

	for i, split := range *postData {
		context, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		userData, err := p.AuthClient.GetUserDetailsLiteForPostView(context, &pb.RequestUserId{UserId: fmt.Sprint((*postData)[i].UserId)})
		if err != nil || userData.ErrorMessage != "" {
			return nil, errors.New(fmt.Sprint(err) + userData.ErrorMessage)
		}

		(*postData)[i].UserName = userData.UserName
		(*postData)[i].UserProfileImgURL = userData.UserProfileImgURL

		postIdString := fmt.Sprint(split.PostId)
		postMedias, err := p.PostRepo.GetPostMediaById(&postIdString)
		if err != nil {
			return nil, err
		}
		(*postData)[i].MediaUrl = *postMedias
		LikeCommentCount, err := p.PostRepo.GetPostLikeAndCommentsCount(&postIdString)
		if err != nil {
			return nil, err
		}

		(*postData)[i].LikesCount = LikeCommentCount.LikesCount
		(*postData)[i].CommentsCount = LikeCommentCount.CommentsCount

		(*postData)[i].PostAge, _ = p.PostRepo.CalculatePostAge(int((*postData)[i].PostId))
	}

	// Cache enriched data
	if err := cache.CacheSet(p.Ctx, p.Redis, cacheKey, postData, 10*time.Minute); err != nil {
		fmt.Printf("Failed to cache user feed data: %v\n", err)
	}
	return postData, nil
}

func (p *PostUseCase) GetMostLovedPostsFromGlobalUser(userId, limit, offset *string) (*[]responsemodel.PostData, error) {
	cacheKey := "userFeed"

	// Attempt to retrieve from cache
	var cachedPosts []responsemodel.PostData
	if err := cache.CacheGet(p.Ctx, p.Redis, cacheKey, &cachedPosts); err == nil {
		// Cache hit
		return &cachedPosts, nil
	} else if err != redis.Nil {
		// Redis error (not cache miss)
		return nil, fmt.Errorf("redis error: %w", err)
	}

	// Cache miss: Retrieve from PostRepo
	postData, err := p.PostRepo.GetMostLovedPostsFromGlobalUser(*userId, *limit, *offset)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}

	// Enrich post data
	for i := range *postData {
		postIdStr := fmt.Sprint((*postData)[i].PostId)
		ctx, cancel := context.WithTimeout(p.Ctx, 10*time.Second)
		defer cancel()

		// Fetch and add user details
		userData, err := p.AuthClient.GetUserDetailsLiteForPostView(ctx, &pb.RequestUserId{UserId: fmt.Sprint((*postData)[i].UserId)})
		if err != nil || userData.ErrorMessage != "" {
			return nil, fmt.Errorf("error getting user details for post %s: %v %s", postIdStr, err, userData.ErrorMessage)
		}
		(*postData)[i].UserName = userData.UserName
		(*postData)[i].UserProfileImgURL = userData.UserProfileImgURL

		// Fetch media URLs for each post
		postMedias, err := p.PostRepo.GetPostMediaById(&postIdStr)
		if err == nil {
			(*postData)[i].MediaUrl = *postMedias
		}

		// Fetch like and comment counts
		likeCommentCount, err := p.PostRepo.GetPostLikeAndCommentsCount(&postIdStr)
		if err == nil {
			(*postData)[i].LikesCount = likeCommentCount.LikesCount
			(*postData)[i].CommentsCount = likeCommentCount.CommentsCount
		}

		// Calculate post age
		postAge, err := p.PostRepo.CalculatePostAge(int((*postData)[i].PostId))
		if err == nil {
			(*postData)[i].PostAge = postAge
		}
	}

	// Cache enriched data
	if err := cache.CacheSet(p.Ctx, p.Redis, cacheKey, postData, 10*time.Minute); err != nil {
		fmt.Printf("Failed to cache user feed data: %v\n", err)
	}

	return postData, nil
}

func (p *PostUseCase) GetRandomPosts(limit, offset *string) (*[]responsemodel.PostData, error) {

	postData, err := p.PostRepo.GetRandomPosts(*limit, *offset)
	if err != nil {
		return nil, err
	}

	for i, split := range *postData {
		context, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		userData, err := p.AuthClient.GetUserDetailsLiteForPostView(context, &pb.RequestUserId{UserId: fmt.Sprint((*postData)[i].UserId)})
		if err != nil || userData.ErrorMessage != "" {
			return nil, errors.New(fmt.Sprint(err) + userData.ErrorMessage)
		}

		(*postData)[i].UserName = userData.UserName
		(*postData)[i].UserProfileImgURL = userData.UserProfileImgURL

		postIdString := fmt.Sprint(split.PostId)
		postMedias, err := p.PostRepo.GetPostMediaById(&postIdString)
		if err != nil {
			return nil, err
		}
		(*postData)[i].MediaUrl = *postMedias
		LikeCommentCount, err := p.PostRepo.GetPostLikeAndCommentsCount(&postIdString)
		if err != nil {
			return nil, err
		}

		(*postData)[i].LikesCount = LikeCommentCount.LikesCount
		(*postData)[i].CommentsCount = LikeCommentCount.CommentsCount

		(*postData)[i].PostAge, _ = p.PostRepo.CalculatePostAge(int((*postData)[i].PostId))
	}

	return postData, nil
}
