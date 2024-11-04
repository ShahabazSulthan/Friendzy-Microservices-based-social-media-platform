package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/IBM/sarama"
	"github.com/ShahabazSulthan/Friendzy_Notification/pkg/config"
	"github.com/ShahabazSulthan/Friendzy_Notification/pkg/models/requestmodel"
	"github.com/ShahabazSulthan/Friendzy_Notification/pkg/models/responsemodel"
	"github.com/ShahabazSulthan/Friendzy_Notification/pkg/pb"
	interface_notification "github.com/ShahabazSulthan/Friendzy_Notification/pkg/repository/interface"
	interface_usecase "github.com/ShahabazSulthan/Friendzy_Notification/pkg/usecase/interface"
)

type NotifUseCase struct {
	NotifRepo   interface_notification.INotifRepo
	KafkaConfig config.KafkaConfigs
	AuthClient  pb.AuthServiceClient
}

func NewNotifUseCase(notifRepo interface_notification.INotifRepo,
	config config.KafkaConfigs,
	authClient *pb.AuthServiceClient) interface_usecase.INotifUseCase {
	return &NotifUseCase{
		NotifRepo:   notifRepo,
		KafkaConfig: config,
		AuthClient:  *authClient,
	}
}

func (n *NotifUseCase) KafkaMessageConsumer() {
	fmt.Println("--------- Kafka consumer initiated ---------")

	// Configure Sarama settings
	configs := sarama.NewConfig()
	configs.Consumer.Return.Errors = true
	configs.Version = sarama.V2_1_0_0 // Adjust to your Kafka version for compatibility

	// Initialize Kafka consumer
	consumer, err := sarama.NewConsumer([]string{n.KafkaConfig.KafkaPort}, configs)
	if err != nil {
		log.Fatalf("Failed to create Kafka consumer: %v", err)
		return
	}
	defer func() {
		if err := consumer.Close(); err != nil {
			log.Printf("Failed to close Kafka consumer: %v", err)
		}
	}()

	// Set up partition consumer
	partitionConsumer, err := consumer.ConsumePartition(n.KafkaConfig.KafkaTopicNotification, 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatalf("Failed to consume partition: %v", err)
		return
	}
	defer func() {
		if err := partitionConsumer.Close(); err != nil {
			log.Printf("Failed to close partition consumer: %v", err)
		}
	}()

	

	// Listen for messages
	for {
		select {
		case message := <-partitionConsumer.Messages():
			if message == nil {
				log.Println("Received nil message from Kafka, skipping...")
				continue
			}
			msg, err := unmarshalChatMessage(message.Value)
			if err != nil {
				log.Printf("Failed to unmarshal Kafka message: %v", err)
				continue
			}
			fmt.Printf("Received message: %v\n", msg)

			// Store notification in the repository
			err = n.NotifRepo.CreateNewNotification(msg)
			if err != nil {
				log.Printf("Failed to create notification: %v", err)
			}

		case err := <-partitionConsumer.Errors():
			log.Printf("Error consuming messages: %v", err)
		}
	}
}

// unmarshalChatMessage decodes a Kafka message into a KafkaNotification struct.
func unmarshalChatMessage(data []byte) (*requestmodel.KafkaNotification, error) {
	var notification requestmodel.KafkaNotification
	err := json.Unmarshal(data, &notification)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling Kafka message: %v", err)
	}
	return &notification, nil
}

func (n *NotifUseCase) GetNotificationsForUser(userId, limit, offset *string) (*[]responsemodel.NotificationModel, error) {

	notifData, err := n.NotifRepo.GetNotificationsForUser(userId, limit, offset)
	if err != nil {
		return nil, err
	}

	for i := range *notifData {
		context, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		userData, err := n.AuthClient.GetUserDetailsLiteForPostView(context, &pb.RequestUserId{UserId: fmt.Sprint((*notifData)[i].ActorID)})
		if err != nil || userData.ErrorMessage != "" {
			return nil, errors.New(fmt.Sprint(err) + userData.ErrorMessage)
		}
		(*notifData)[i].ActorUserName = userData.UserName
		(*notifData)[i].ActorProfileImgURL = userData.UserProfileImgURL

		// Calculate age based on CreatedAt
		age, err := n.NotifRepo.CalculatePostAge((*notifData)[i].CreatedAt)
		if err != nil {
			return nil, err
		}
		(*notifData)[i].NotificationAge = age
	}

	return notifData, nil
}
