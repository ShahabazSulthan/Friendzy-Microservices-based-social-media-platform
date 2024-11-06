package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/IBM/sarama"
	"github.com/ShahabazSulthan/friendzy_post/pkg/config"
	"github.com/ShahabazSulthan/friendzy_post/pkg/model/requestmodel"
	"github.com/ShahabazSulthan/friendzy_post/pkg/model/responsemodel"
	"github.com/ShahabazSulthan/friendzy_post/pkg/pb"
	interface_chat "github.com/ShahabazSulthan/friendzy_post/pkg/repository/interface"
	interface_usecase "github.com/ShahabazSulthan/friendzy_post/pkg/usecase/interface"
)

type ChatUsecase struct {
	ChatRepo interface_chat.IChatRepo
	Client   pb.AuthServiceClient
	Kafka    *config.ApacheKafka
}

func NewChatUseCase(
	chatRepo interface_chat.IChatRepo,
	client *pb.AuthServiceClient,
	config *config.ApacheKafka) interface_usecase.IChatUseCase {
	return &ChatUsecase{
		ChatRepo: chatRepo,
		Client:   *client,
		Kafka:    config,
	}
}

func (c *ChatUsecase) KafkaOneToOneMessageConsumer() {
	fmt.Println("========= Kafka One-to-One Message Consumer Initiated ================")

	// Initialize Kafka consumer configuration
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true // Enable error reporting on the consumer

	// Create a new Kafka consumer
	consumer, err := sarama.NewConsumer([]string{c.Kafka.KafkaPort}, config)
	if err != nil {
		log.Printf("Failed to create Kafka consumer: %v", err)
		return
	}
	defer func() {
		if err := consumer.Close(); err != nil {
			log.Printf("Failed to close Kafka consumer: %v", err)
		}
	}()
	fmt.Println("Kafka consumer created successfully")

	// Create a partition consumer for the specified topic and partition
	partitionConsumer, err := consumer.ConsumePartition(c.Kafka.KafkaTopicOneToOne, 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatalf("Failed to start partition consumer: %v", err)
		return
	}
	defer func() {
		if err := partitionConsumer.Close(); err != nil {
			log.Printf("Failed to close partition consumer: %v", err)
		}
	}()
	fmt.Println("Partition consumer started successfully")

	// Consume messages in a loop
	for {
		select {
		case message := <-partitionConsumer.Messages():
			fmt.Println("Received message from Kafka topic")
			msg, err := unmarshalOneToOneChatMessage(message.Value)
			if err != nil {
				log.Printf("Failed to unmarshal message: %v", err)
				continue
			}
			fmt.Printf("Processed message: %v\n", msg)

			// Store the message in the database
			if _, err := c.ChatRepo.StoreOneToOneChatToDB(msg); err != nil {
				log.Printf("Failed to store message to DB: %v", err)
			}

		case err := <-partitionConsumer.Errors():
			log.Printf("Error received from Kafka partition consumer: %v", err)
		}
	}
}

func unmarshalOneToOneChatMessage(data []byte) (*requestmodel.OneToOneChatRequest, error) {
	var store requestmodel.OneToOneChatRequest

	err := json.Unmarshal(data, &store)
	if err != nil {
		return nil, err
	}
	return &store, nil
}

func (c *ChatUsecase) GetOneToOneChats(senderId, recipientId, limit, offset *string) (*[]responsemodel.OneToOneChatResponse, error) {
	// Update the chat status for messages between the sender and recipient
	if err := c.ChatRepo.UpdateChatStatus(senderId, recipientId); err != nil {
		log.Printf("Failed to update chat status: %v", err)
		return nil, fmt.Errorf("could not update chat status: %w", err)
	}

	// Retrieve the one-to-one chat messages with pagination
	userChats, err := c.ChatRepo.GetOneToOneChats(senderId, recipientId, limit, offset)
	if err != nil {
		log.Printf("Failed to retrieve one-to-one chats: %v", err)
		return nil, fmt.Errorf("could not fetch chats: %w", err)
	}

	return userChats, nil
}

func (c *ChatUsecase) GetRecentChatProfilesPlusChatData(senderID, limit, offset *string) (*[]responsemodel.RecentChatProfileResponse, error) {
	// Fetch recent chat profile data
	recentChatData, err := c.ChatRepo.RecentChatProfileData(senderID, limit, offset)
	if err != nil {
		log.Printf("Error retrieving recent chat profile data: %v", err)
		return nil, fmt.Errorf("failed to fetch recent chat profiles: %w", err)
	}

	// Iterate over recent chat data to fetch additional user details
	for i, chatProfile := range *recentChatData {
		// Use a context with a timeout for the RPC call to fetch user details
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Call to external service to get user details
		resp, err := c.Client.GetUserDetailsLiteForPostView(ctx, &pb.RequestUserId{
			UserId: chatProfile.UserId,
		})
		if err != nil {
			log.Printf("Error calling GetUserDetailsLiteForPostView: %v", err)
			return nil, fmt.Errorf("auth service is unavailable: %w", err)
		}
		if resp.ErrorMessage != "" {
			log.Printf("Error from auth service: %s", resp.ErrorMessage)
			return nil, errors.New(resp.ErrorMessage)
		}

		// Update profile data with fetched user details
		(*recentChatData)[i].UserName = resp.UserName
		(*recentChatData)[i].UserProfileImgURL = resp.UserProfileImgURL
	}

	return recentChatData, nil
}
