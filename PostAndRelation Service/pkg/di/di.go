package di

import (
	"context"
	"fmt"
	"log"

	"github.com/ShahabazSulthan/friendzy_post/pkg/client"
	"github.com/ShahabazSulthan/friendzy_post/pkg/config"
	"github.com/ShahabazSulthan/friendzy_post/pkg/db"
	"github.com/ShahabazSulthan/friendzy_post/pkg/repository"
	"github.com/ShahabazSulthan/friendzy_post/pkg/server"
	"github.com/ShahabazSulthan/friendzy_post/pkg/usecase"
	cache "github.com/ShahabazSulthan/friendzy_post/pkg/utils/Cache"
	hashpassword "github.com/ShahabazSulthan/friendzy_post/pkg/utils/Hash_Password"
	kafka "github.com/ShahabazSulthan/friendzy_post/pkg/utils/Kakfa"
)

// InitializePostAndRelationServer initializes the server for posts and relations
func InitializePostAndRelationServer(config *config.Config) (*server.PostService, error) {
	// Initialize hash utility
	hashUtil := hashpassword.NewHashUtil()

	// Connect to the database
	DB, err := db.ConnectDatabase(&config.DB, hashUtil)
	if err != nil {
		fmt.Println("Error connecting to the database:", err)
		return nil, err
	}

	// Initialize Redis client
	redisClient, err := cache.NewRedisHelper(config.Redis)
	if err != nil {
		fmt.Println("Error connecting to Redis:", err)
		return nil, err
	}

	// Initialize the Kafka producer and check connection
	kafkaProducer := kafka.NewKafkaProducer(config.Kafka)
	if err := kafka.CheckKafkaConnection([]string{config.Kafka.KafkaPort}, nil); err != nil {
		log.Println("Failed to connect to Kafka:", err)
		return nil, err
	}
	fmt.Println("Kafka connection established successfully.")

	// Initialize the authentication client
	authClient, err := client.InitAuthServiceClient(config)
	if err != nil {
		fmt.Println("Error connecting to auth service client:", err)
		return nil, err
	}

	// Create context for repositories and use cases
	ctx := context.Background()

	// Create repositories
	postRepo := repository.NewPostRepo(DB, redisClient.Client, ctx)
	relationRepo := repository.NewRelationRepo(DB, redisClient.Client, ctx)
	commentRepo := repository.NewCommentRepo(DB, redisClient.Client, ctx)

	// Create use cases
	postUsecase := usecase.NewPostUseCase(postRepo, *authClient, kafkaProducer, redisClient.Client, ctx)
	relationUsecase := usecase.NewRelationUseCase(relationRepo, postRepo, *authClient, kafkaProducer, redisClient.Client, ctx)
	commentUseCase := usecase.NewCommentUsecase(commentRepo, authClient, kafkaProducer, postRepo, redisClient.Client, ctx)

	// Return a new instance of PostAndRelation service
	return server.NewPostAndRelation(postUsecase, relationUsecase, commentUseCase), nil
}
