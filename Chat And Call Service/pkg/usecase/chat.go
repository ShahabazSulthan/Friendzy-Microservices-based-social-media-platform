package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
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

func (c *ChatUsecase) KafkaOneToManyMessageConsumer() {
	fmt.Println("========= Kafka One-to-Many Message Consumer Initiated ================")

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
	partitionConsumer, err := consumer.ConsumePartition(c.Kafka.KafkaTopicOneToMany, 0, sarama.OffsetNewest)
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

	// Consume messages in a loop with panic recovery
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from panic in message consumer loop: %v", r)
		}
	}()

	for {
		select {
		case message := <-partitionConsumer.Messages():
			// Process the message
			msg, err := unmarshalOneToManyChatMessage(message.Value)
			if err != nil {
				log.Printf("Failed to unmarshal message: %v", err)
				continue
			}
			fmt.Printf("Processed message: %v\n", msg)

			// Store the message in the database with error handling
			if err := c.ChatRepo.StoreOneToManyChatToDB(msg); err != nil {
				log.Printf("Failed to store message to DB: %v", err)
				// Add retry or other error handling if necessary
			}

		case err := <-partitionConsumer.Errors():
			log.Printf("Error received from Kafka partition consumer: %v", err)
		}
	}
}



func unmarshalOneToManyChatMessage(data []byte) (*requestmodel.OneToManyMessageRequest, error) {
	var store requestmodel.OneToManyMessageRequest

	err := json.Unmarshal(data, &store)
	if err != nil {
		return nil, err
	}
	return &store, nil
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

func (c *ChatUsecase) CreateNewGroup(groupInfo *requestmodel.NewGroupInfo) error {
	for _, member := range groupInfo.GroupMembers {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		resp, err := c.Client.CheckUserExist(ctx, &pb.RequestUserId{UserId: fmt.Sprint(member)})
		if err != nil {
			log.Printf("Error calling CheckUserExist in authSvc: %v", err)
			// Retry mechanism or custom error handling logic here
			return fmt.Errorf("auth service unavailable: %w", err)
		}

		if resp.ErrorMessage != "" {
			return errors.New(resp.ErrorMessage)
		}
		if !resp.ExistStatus {
			return fmt.Errorf("no user found with id %d", member)
		}
	}

	err := c.ChatRepo.CreateNewGroup(groupInfo)
	if err != nil {
		log.Printf("Error creating new group: %v", err)
		return fmt.Errorf("database error: %w", err)
	}
	return nil
}


func (c *ChatUsecase) GroupMembersList(groupId *string) (*[]string, error) {
	groupMembers, err := c.ChatRepo.GetGroupMembersList(groupId)
	if err != nil {
		return nil, err
	}

	var memberIds []string
	for _, member := range *groupMembers {
		memberIds = append(memberIds, strconv.Itoa(int(member)))
	}

	return &memberIds, nil
}

func (c *ChatUsecase) GetUserGroupChatSummary(userId, limit, offset *string) (*[]responsemodel.GroupChatSummaryResponse, error) {
	var groupChatSummary []responsemodel.GroupChatSummaryResponse
	var singlegroupChatSummary responsemodel.GroupChatSummaryResponse

	recentGroupProfiles, err := c.ChatRepo.GetRecentGroupProfilesOfUser(userId, limit, offset)
	if err != nil {
		return nil, err
	}

	fmt.Println("----------", recentGroupProfiles)

	for i := range *recentGroupProfiles {
		lastMessageDetails, err := c.ChatRepo.GetGroupLastMessageDetailsByGroupId(&(*recentGroupProfiles)[i].GroupID)
		if err != nil {
			return nil, err
		}

		if lastMessageDetails != nil {

			ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
			defer cancel()

			resp, err := c.Client.GetUserDetailsLiteForPostView(ctx, &pb.RequestUserId{
				UserId: lastMessageDetails.SenderID,
			})
			if err != nil {
				log.Println("-----error: from usecase:GetUserGroupChatSummary() authSvc down while calling GetUserDetailsLiteForPostView(),error:", err)
				return nil, err
			}
			if resp.ErrorMessage != "" {
				return nil, errors.New(resp.ErrorMessage)
			}

			singlegroupChatSummary.LastMessage = lastMessageDetails.LastMessage
			singlegroupChatSummary.SenderID = lastMessageDetails.SenderID
			singlegroupChatSummary.SenderUserName = resp.UserName
			singlegroupChatSummary.StringTime = lastMessageDetails.StringTime
		}
		singlegroupChatSummary.GroupID = ((*recentGroupProfiles)[i].GroupID)
		singlegroupChatSummary.GroupName = (*recentGroupProfiles)[i].GroupName
		singlegroupChatSummary.GroupProfileImgURL = (*recentGroupProfiles)[i].GroupProfileImgURL

		groupChatSummary = append(groupChatSummary, singlegroupChatSummary)
	}

	fmt.Println("Group summury ",groupChatSummary)
	return &groupChatSummary, nil
}

func (c *ChatUsecase) GetOneToManyChats(userid, groupid, limit, offset *string) (*[]responsemodel.OneToManyChatResponse, error) {
	belongs, err := c.ChatRepo.CheckUserIsGroupMember(userid, groupid)
	if err != nil {
		return nil, err
	}
	if !belongs {
		return nil, fmt.Errorf("can't access chat data,user with id %s does not belongs to group with id %s", *userid, *groupid)
	}

	userChats, err := c.ChatRepo.GetOneToManyChats(groupid, limit, offset)
	if err != nil {
		return nil, err
	}

	fmt.Println("userchat ",userChats)

	for i := range *userChats {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		resp, err := c.Client.GetUserDetailsLiteForPostView(ctx, &pb.RequestUserId{
			UserId: (*userChats)[i].SenderID,
		})
		if err != nil {
			log.Println("-----error: from usecase:GetOneToManyChats() authSvc down while calling GetUserDetailsLiteForPostView(),error:", err)
			return nil, err
		}
		if resp.ErrorMessage != "" {
			return nil, errors.New(resp.ErrorMessage)
		}

		(*userChats)[i].SenderUserName = resp.UserName
		(*userChats)[i].SenderProfileImgURL = resp.UserProfileImgURL
	}

	fmt.Println("UserGroup Chat ",userChats)
	return userChats, nil
}

func (c *ChatUsecase) AddNewMembersToGroup(inputData *requestmodel.AddNewMemberToGroup) error {
	isMember, err := c.ChatRepo.CheckUserIsGroupMember(&inputData.UserID, &inputData.GroupID)
	if err != nil {
		return err
	}
	if !isMember {
		return errors.New("you can't add members to this group,cause you are not a member of this group")
	}

	for _, member := range inputData.GroupMembers {
		context, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		resp, err := c.Client.CheckUserExist(context, &pb.RequestUserId{UserId: fmt.Sprint(member)})
		if err != nil {
			log.Println("-------error: from usecase:AddNewMembersToGroup() authSvc down while calling CheckUserExist()-----")
			return err
		}
		if resp.ErrorMessage != "" {
			return errors.New(resp.ErrorMessage)
		}
		if !resp.ExistStatus {
			newErr := fmt.Sprintf("no user found with id %d,please enter valid userId", member)
			return errors.New(newErr)
		}
	}

	err = c.ChatRepo.AddNewMembersToGroupByGroupId(inputData)
	if err != nil {
		return err
	}

	return nil
}

func (c *ChatUsecase) RemoveMemberFromGroup(inputData *requestmodel.RemoveMemberFromGroup) error {
	isMember, err := c.ChatRepo.CheckUserIsGroupMember(&inputData.UserID, &inputData.GroupID)
	if err != nil {
		return err
	}
	if !isMember {
		return errors.New("you can't remove members from this group,cause you are not a member of this group")
	}

	err = c.ChatRepo.RemoveGroupMember(inputData)
	if err != nil {
		return err
	}

	memberCount, err := c.ChatRepo.CountMembersInGroup(inputData.GroupID)
	if err != nil {
		return err
	}
	if memberCount == 0 {
		err := c.ChatRepo.DeleteOneToManyChatsByGroupId(inputData.GroupID)
		if err != nil {
			return err
		}
		err = c.ChatRepo.DeleteGroupDataByGroupId(inputData.GroupID)
		if err != nil {
			return err
		}
	}

	return nil
}
