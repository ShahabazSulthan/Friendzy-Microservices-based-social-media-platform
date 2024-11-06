package di

import (
	"fmt"

	"github.com/ShahabazSulthan/friendzy_post/pkg/client"
	"github.com/ShahabazSulthan/friendzy_post/pkg/config"
	"github.com/ShahabazSulthan/friendzy_post/pkg/db"
	"github.com/ShahabazSulthan/friendzy_post/pkg/repository"
	"github.com/ShahabazSulthan/friendzy_post/pkg/server"
	"github.com/ShahabazSulthan/friendzy_post/pkg/usecase"
)

func InitializeChatNCallSvc(config *config.Config) (*server.ChatNCallSvc, error) {
	// Connect to MongoDB
	DB, err := db.ConnectDatabaseMongo(&config.MongoDB)
	if err != nil {
		fmt.Println("Error connecting to MongoDB in InitializeChatNCallSvc")
		return nil, fmt.Errorf("database connection error: %w", err)
	}

	// Initialize Auth Service Client
	authClient, err := client.InitAuthServiceClient(&config.PortMngr)
	if err != nil {
		fmt.Println("Error initializing Auth Service client in InitializeChatNCallSvc")
		return nil, fmt.Errorf("auth service client setup error: %w", err)
	}

	// Set up repository and use case
	chatRepo := repository.NewChatRepo(*DB)
	chatUseCase := usecase.NewChatUseCase(chatRepo, authClient, &config.Kafka)

	// Start Kafka message consumer in a separate goroutine
	go chatUseCase.KafkaOneToOneMessageConsumer()

	// Initialize and return the ChatNCall Service server
	return server.NewChatNCallServiceServer(chatUseCase), nil
}
