package handler_chat

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/IBM/sarama"
	requestmodel_chat "github.com/ShahabazSulthan/Friendzy_apiGateway/pkg/Chat_Call_Service/model/requestmodel"
	responsemodel_chat "github.com/ShahabazSulthan/Friendzy_apiGateway/pkg/Chat_Call_Service/model/responsemodel"
	"github.com/ShahabazSulthan/Friendzy_apiGateway/pkg/Chat_Call_Service/pb"
	config "github.com/ShahabazSulthan/Friendzy_apiGateway/pkg/Config"
	responsemodel_post "github.com/ShahabazSulthan/Friendzy_apiGateway/pkg/post_relation_service/model/responsemodel"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

type ChatWebsockethandler struct {
	Client      pb.ChatNCallServiceClient
	LocationInd *time.Location
	Config      *config.Config
}

func NewChatWebsocketHandler(client *pb.ChatNCallServiceClient, config *config.Config) *ChatWebsockethandler {
	if client == nil || config == nil {
		log.Println("client or config is nil, ensure dependencies are initialized")
		return nil
	}
	locationInd, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		log.Println("Error fetching location:", err)
	}
	return &ChatWebsockethandler{
		Client:      *client,
		LocationInd: locationInd,
		Config:      config,
	}
}

var UserSocketMap = make(map[string]*websocket.Conn)

func (svc *ChatWebsockethandler) WsConnection(ctx *fiber.Ctx) error {
	websocket.New(func(conn *websocket.Conn) {
		var messageModel requestmodel_chat.MessageRequest
		var (
			msg []byte
			err error
		)

		userId := conn.Locals("userId")
		userIdStr := fmt.Sprint(userId)

		UserSocketMap[userIdStr] = conn

		defer conn.Close()
		defer delete(UserSocketMap, userIdStr)

		for {
			if _, msg, err = conn.ReadMessage(); err != nil {
				log.Println("read:", err)
				sendErrMessageWS(userIdStr, err)
				break
			}
			err = json.Unmarshal(msg, &messageModel)
			if err != nil {
				log.Println("read:", err)
				sendErrMessageWS(userIdStr, err)
				break
			}
			messageModel.SenderID = userIdStr
			messageModel.TimeStamp = time.Now()

			validate := validator.New(validator.WithRequiredStructEnabled())
			err = validate.Struct(messageModel)
			if err != nil {
				if ve, ok := err.(validator.ValidationErrors); ok {
					for _, e := range ve {
						switch e.Field() {
						case "Type":
							sendErrMessageWS(userIdStr, errors.New("no Type found in input"))
						}
					}
				}
				break
			}
			switch messageModel.Type {
			case "OneToOne":
				svc.OnetoOneMessage(&messageModel)
			case "OneToMany":
				svc.OnetoMany(&messageModel)
			case "TypingStatus":
				svc.TypingStatus(&messageModel)
			default:
				sendErrMessageWS(userIdStr, errors.New("message Type should be OneToOne,OneToMany or TypingStatus ,no other types allowed"))
			}
		}
	})(ctx)

	return nil
}

func (svc *ChatWebsockethandler) TypingStatus(msgModel *requestmodel_chat.MessageRequest) {
	if msgModel.RecipientID == "" {
		sendErrMessageWS(msgModel.SenderID, errors.New("no RecipientID found in input"))
		return // Return immediately after sending error
	}

	var MsgModel requestmodel_chat.TypingStatusRequest
	MsgModel.SenderID = msgModel.SenderID
	MsgModel.RecipientID = msgModel.RecipientID
	MsgModel.Type = msgModel.Type
	MsgModel.TypingStat = msgModel.TypingStat

	conn, ok := UserSocketMap[MsgModel.RecipientID]
	if ok {
		data, err := MarshalStructJson(MsgModel)
		if err != nil {
			sendErrMessageWS(MsgModel.SenderID, err)
			return
		}
		err = conn.WriteMessage(websocket.TextMessage, *data)
		if err != nil {
			fmt.Println("error sending to recipient:", err)
			sendErrMessageWS(MsgModel.SenderID, err)
		}
	}
}

func (svc *ChatWebsockethandler) OnetoOneMessage(msgModel *requestmodel_chat.MessageRequest) {
	if msgModel.RecipientID == "" {
		sendErrMessageWS(msgModel.SenderID, errors.New("no RecipientID found in input"))
		return // Return immediately after sending error
	}
	if msgModel.Content == "" || len(msgModel.Content) > 100 {
		sendErrMessageWS(msgModel.SenderID, errors.New("message content should be less than 100 characters"))
		return
	}

	var OneToOneMsgModel requestmodel_chat.OnetoOneMessageRequest
	OneToOneMsgModel.SenderID = msgModel.SenderID
	OneToOneMsgModel.RecipientID = msgModel.RecipientID
	OneToOneMsgModel.Content = msgModel.Content
	OneToOneMsgModel.TimeStamp = msgModel.TimeStamp
	OneToOneMsgModel.Status = "pending"
	OneToOneMsgModel.Type = msgModel.Type

	conn, ok := UserSocketMap[OneToOneMsgModel.RecipientID]
	if ok {
		OneToOneMsgModel.TimeStamp = OneToOneMsgModel.TimeStamp.In(svc.LocationInd)
		data, err := MarshalStructJson(OneToOneMsgModel)
		if err != nil {
			sendErrMessageWS(OneToOneMsgModel.SenderID, err)
			return
		}
		err = conn.WriteMessage(websocket.TextMessage, *data)
		if err != nil {
			fmt.Println("error sending to recipient:", err)
			sendErrMessageWS(OneToOneMsgModel.SenderID, err)
			return
		}
		OneToOneMsgModel.Status = "send"
	}

	fmt.Println("------check status is pending or send--------", OneToOneMsgModel)
	fmt.Println("Adding to Kafka producer for transporting to service and storing")
	svc.KafkaProducerUpdateOneToOneMessage(&OneToOneMsgModel)
}

func (svc *ChatWebsockethandler) KafkaProducerUpdateOneToOneMessage(message *requestmodel_chat.OnetoOneMessageRequest) error {
	fmt.Println("---------------to KafkaProducerUpdateOneToOneMessage:", *message)

	configs := sarama.NewConfig()
	configs.Producer.Return.Successes = true
	configs.Producer.Retry.Max = 5

	producer, err := sarama.NewSyncProducer([]string{svc.Config.KafkaPort}, configs)
	if err != nil {
		return fmt.Errorf("failed to create Kafka producer: %w", err)
	}
	defer producer.Close()

	msgJson, err := MarshalStructJson(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: svc.Config.KafkaTopicOneToOne,
		Key:   sarama.StringEncoder(message.RecipientID),
		Value: sarama.StringEncoder(*msgJson),
	}
	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		fmt.Println("error sending message to Kafka producer:", err)
		return fmt.Errorf("failed to send message to Kafka: %w", err)
	}

	log.Printf("[producer] partition id: %d; offset: %d; value: %v\n", partition, offset, msg)
	return nil
}

func (svc *ChatWebsockethandler) KafkaProducerUpdateOneToManyMessage(message *requestmodel_chat.OnetoManyMessageRequest) error {
	fmt.Println("---------------to KafkaProducerUpdateOneToManyMessage:", *message)

	configs := sarama.NewConfig()
	configs.Producer.Return.Successes = true
	configs.Producer.Retry.Max = 5

	producer, err := sarama.NewSyncProducer([]string{svc.Config.KafkaPort}, configs)
	if err != nil {
		return err
	}

	msgJson, _ := MarshalStructJson(message)

	msg := &sarama.ProducerMessage{Topic: svc.Config.KafkaTopicOneToMany,
		Key:   sarama.StringEncoder(message.GroupID),
		Value: sarama.StringEncoder(*msgJson)}
	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		fmt.Println("err sending message to kafkaproducer ", err)
	}
	log.Printf("[producer] partition id: %d; offset:%d, value: %v\n", partition, offset, msg)
	return nil
}

func (svc *ChatWebsockethandler) GetOneToOneChats(ctx *fiber.Ctx) error {
	userId := ctx.Locals("userId")

	limit := ctx.Query("limit", "12")
	offset := ctx.Query("offset", "0")

	recipientId := ctx.Params("recipientid")
	if recipientId == "" {
		return ctx.Status(fiber.ErrBadRequest.Code).
			JSON(responsemodel_post.CommonResponse{
				StatusCode: fiber.ErrBadRequest.Code,
				Message:    "can't get chat (possible reason: no input)",
				Error:      "no recipientId found in request",
			})
	}

	c, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := svc.Client.GetOneToOneChats(c, &pb.RequestUserOneToOneChat{
		SenderID:   fmt.Sprint(userId),
		RecieverID: recipientId,
		Limit:      limit,
		Offset:     offset,
	})

	if err != nil {
		fmt.Println("----------chat service down--------")
		return ctx.Status(fiber.StatusServiceUnavailable).
			JSON(responsemodel_post.CommonResponse{
				StatusCode: fiber.StatusServiceUnavailable,
				Message:    "can't get chat",
				Error:      err.Error(),
			})
	}

	if resp.ErrorMessage != "" {
		return ctx.Status(fiber.StatusBadRequest).
			JSON(responsemodel_post.CommonResponse{
				StatusCode: fiber.StatusBadRequest,
				Message:    "can't get chat",
				Data:       resp,
				Error:      resp.ErrorMessage,
			})
	}

	return ctx.Status(fiber.StatusOK).
		JSON(responsemodel_post.CommonResponse{
			StatusCode: fiber.StatusOK,
			Data:       resp,
			Message:    "chat fetched successfully",
		})
}

func (svc *ChatWebsockethandler) GetRecentChatProfileDetails(ctx *fiber.Ctx) error {
	userId := ctx.Locals("userId")
	limit := ctx.Query("limit", "12")
	offset := ctx.Query("offset", "0")

	c, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := svc.Client.GetRecentChatProfiles(c, &pb.RequestRecentChatProfiles{
		SenderID: fmt.Sprint(userId),
		Limit:    limit,
		Offset:   offset,
	})

	if err != nil {
		fmt.Println("----------chat service down--------, err:", err)
		return ctx.Status(fiber.StatusServiceUnavailable).
			JSON(responsemodel_post.CommonResponse{
				StatusCode: fiber.StatusServiceUnavailable,
				Message:    "can't get recent chat profiles",
				Error:      err.Error(),
			})
	}

	if resp.ErrorMessage != "" {
		fmt.Println("-----------------------", resp.ErrorMessage)
		return ctx.Status(fiber.StatusBadRequest).
			JSON(responsemodel_post.CommonResponse{
				StatusCode: fiber.StatusBadRequest,
				Message:    "can't get recent chat profiles",
				Data:       resp,
				Error:      resp.ErrorMessage,
			})
	}

	return ctx.Status(fiber.StatusOK).
		JSON(responsemodel_post.CommonResponse{
			StatusCode: fiber.StatusOK,
			Data:       resp,
			Message:    "recent chat profiles fetched successfully",
		})
}

func sendErrMessageWS(userid string, err error) {
	conn, ok := UserSocketMap[userid]
	if ok {
		conn.WriteMessage(websocket.TextMessage, []byte(err.Error()))
	}
}

func MarshalStructJson(msgModel interface{}) (*[]byte, error) {
	data, err := json.Marshal(msgModel)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (svc *ChatWebsockethandler) OnetoMany(msgModel *requestmodel_chat.MessageRequest) {
	if msgModel.GroupID == "" {
		sendErrMessageWS(msgModel.SenderID, errors.New("no GroupID found in input"))
		return
	}
	if msgModel.Content == "" || len(msgModel.Content) > 100 {
		sendErrMessageWS(msgModel.SenderID, errors.New("message content should be less than 100 words "))
		return
	}

	context, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	resp, err := svc.Client.GetGroupMembersInfo(context, &pb.RequestGroupMembersInfo{GroupID: msgModel.GroupID})
	if err != nil {
		log.Println("-----from handler:OnetoMany() chatNcall service down while calling GetGroupMembersInfo()---")
		sendErrMessageWS(msgModel.SenderID, err)
		return
	}
	if resp.ErrorMessage != "" {
		sendErrMessageWS(msgModel.SenderID, errors.New(resp.ErrorMessage))
		return
	}

	var OneToManyMsgModel requestmodel_chat.OnetoManyMessageRequest
	OneToManyMsgModel.SenderID = msgModel.SenderID
	OneToManyMsgModel.GroupID = msgModel.GroupID
	OneToManyMsgModel.Content = msgModel.Content
	OneToManyMsgModel.TimeStamp = msgModel.TimeStamp
	OneToManyMsgModel.Status = "pending"
	OneToManyMsgModel.Type = msgModel.Type

	for i := range resp.GroupMembers {
		if (resp.GroupMembers)[i] == msgModel.SenderID {
			continue
		}
		conn, ok := UserSocketMap[(resp.GroupMembers[i])]
		if ok {
			OneToManyMsgModel.TimeStamp = OneToManyMsgModel.TimeStamp.In(svc.LocationInd)
			data, err := MarshalStructJson(OneToManyMsgModel)
			if err != nil {
				sendErrMessageWS(OneToManyMsgModel.SenderID, err)
				return
			}
			err = conn.WriteMessage(websocket.TextMessage, *data)
			if err != nil {
				fmt.Println("error sending to recipient", err)
				return
			}
			OneToManyMsgModel.Status = "send"
		}
	}

	fmt.Println("------check status is pending or send--------", OneToManyMsgModel)
	fmt.Println("Adding to kafkaproducer for transporting to service and storing")
	svc.KafkaProducerUpdateOneToManyMessage(&OneToManyMsgModel)
}

func (svc *ChatWebsockethandler) CreateNewGroup(ctx *fiber.Ctx) error {
	var newGroupData requestmodel_chat.NewGroupInfo
	userId := ctx.Locals("userId")

	// Convert userId to integer
	userIdInt, err := strconv.Atoi(fmt.Sprint(userId))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(responsemodel_post.CommonResponse{
			StatusCode: fiber.StatusBadRequest,
			Message:    "Invalid user ID",
			Error:      err.Error(),
		})
	}

	// Parse and validate the request body
	if err := ctx.BodyParser(&newGroupData); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(responsemodel_post.CommonResponse{
			StatusCode: fiber.StatusBadRequest,
			Message:    "Invalid or missing JSON input",
			Error:      err.Error(),
		})
	}

	// Ensure the creator is added as a group member
	newGroupData.CreatorID = fmt.Sprint(userIdInt)
	isMember := false
	for _, member := range newGroupData.GroupMembers {
		if member == uint64(userIdInt) {
			isMember = true
			break
		}
	}
	if !isMember {
		newGroupData.GroupMembers = append(newGroupData.GroupMembers, uint64(userIdInt))
	}

	// Validate new group data structure
	validate := validator.New()
	if err := validate.Struct(newGroupData); err != nil {
		if ve, ok := err.(validator.ValidationErrors); ok {
			validationResponse := responsemodel_chat.NewGroupInfo{}
			for _, e := range ve {
				switch e.Field() {
				case "GroupName":
					validationResponse.GroupName = "should contain less than 20 characters"
				case "GroupMembers":
					validationResponse.GroupMembers = "should be unique, maximum 12 members, and IDs should be valid numbers"
				}
			}
			return ctx.Status(fiber.StatusBadRequest).JSON(responsemodel_post.CommonResponse{
				StatusCode: fiber.StatusBadRequest,
				Message:    "Invalid group data",
				Data:       validationResponse,
				Error:      err.Error(),
			})
		}
	}

	// Call the gRPC method to create a new group
	resp, err := svc.Client.CreateNewGroup(ctx.Context(), &pb.RequestNewGroup{
		GroupName:    newGroupData.GroupName,
		GroupMembers: newGroupData.GroupMembers,
		CreatorID:    newGroupData.CreatorID,
		CreatedAt:    time.Now().Format(time.RFC3339),
	})

	// Handle gRPC errors and response error message
	if err != nil {
		log.Printf("Service unavailable: error creating group - %v", err)
		return ctx.Status(fiber.StatusServiceUnavailable).JSON(responsemodel_post.CommonResponse{
			StatusCode: fiber.StatusServiceUnavailable,
			Message:    "Unable to create group, service unavailable",
			Error:      err.Error(),
		})
	}
	if resp.ErrorMessage != "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(responsemodel_post.CommonResponse{
			StatusCode: fiber.StatusBadRequest,
			Message:    "Unable to create group",
			Error:      resp.ErrorMessage,
		})
	}

	// Return successful response
	return ctx.Status(fiber.StatusOK).JSON(responsemodel_post.CommonResponse{
		StatusCode: fiber.StatusOK,
		Message:    "Group created successfully",
	})
}

func (svc *ChatWebsockethandler) GetUserGroupsAndLastMessage(ctx *fiber.Ctx) error {
	userId := ctx.Locals("userId")
	limit, offset := ctx.Query("limit", "12"), ctx.Query("offset", "0")

	context, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	resp, err := svc.Client.GetUserGroupChatSummary(context, &pb.RequestGroupChatSummary{
		UserID: fmt.Sprint(userId),
		Limit:  limit,
		Offset: offset,
	})

	if err != nil {
		log.Println("-----error: from handler:GetUserGroupsAndLastMessage(),chatNcall service down while calling GetUserGroupChatSummary()")

		return ctx.Status(fiber.StatusServiceUnavailable).
			JSON(responsemodel_post.CommonResponse{
				StatusCode: fiber.StatusServiceUnavailable,
				Message:    "failed to fetch groupchat summary",
				Error:      err.Error(),
			})
	}

	if resp.ErrorMessage != "" {
		return ctx.Status(fiber.StatusBadRequest).
			JSON(responsemodel_post.CommonResponse{
				StatusCode: fiber.StatusBadRequest,
				Message:    "failed to fetch groupchat summary",
				Data:       resp,
				Error:      resp.ErrorMessage,
			})
	}

	return ctx.Status(fiber.StatusOK).
		JSON(responsemodel_post.CommonResponse{
			StatusCode: fiber.StatusOK,
			Data:       resp.SingleEntity,
			Message:    "groupchat summary fetched succesfully",
		})

}

func (svc *ChatWebsockethandler) GetGroupChats(ctx *fiber.Ctx) error {
	userId := ctx.Locals("userId")

	groupId := ctx.Params("groupid")
	limit, offset := ctx.Query("limit", "12"), ctx.Query("offset", "0")

	context, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	resp, err := svc.Client.GetOneToManyChats(context, &pb.RequestGetOneToManyChats{
		UserID:  fmt.Sprint(userId),
		GroupID: groupId,
		Limit:   limit,
		Offset:  offset,
	})

	if err != nil {
		log.Println("-----error: from handler:GetGroupChats(),chatNcall service down while calling GetOneToManyChats()")

		return ctx.Status(fiber.StatusServiceUnavailable).
			JSON(responsemodel_post.CommonResponse{
				StatusCode: fiber.StatusServiceUnavailable,
				Message:    "failed to fetch groupchat",
				Error:      err.Error(),
			})
	}

	if resp.ErrorMessage != "" {
		return ctx.Status(fiber.StatusBadRequest).
			JSON(responsemodel_post.CommonResponse{
				StatusCode: fiber.StatusBadRequest,
				Message:    "failed to fetch groupchat",
				Data:       resp,
				Error:      resp.ErrorMessage,
			})
	}

	fmt.Println("Response ",resp)
	return ctx.Status(fiber.StatusOK).
		JSON(responsemodel_post.CommonResponse{
			StatusCode: fiber.StatusOK,
			Data:       resp.Chat,
			Message:    "groupchat fetched succesfully",
		})
}

func (svc *ChatWebsockethandler) AddMembersToGroup(ctx *fiber.Ctx) error {
	userId := ctx.Locals("userId")

	var addMembersInput requestmodel_chat.AddNewMembersToGroup

	if err := ctx.BodyParser(&addMembersInput); err != nil {
		return ctx.Status(fiber.ErrBadRequest.Code).
			JSON(responsemodel_post.CommonResponse{
				StatusCode: fiber.ErrBadRequest.Code,
				Message:    "can't add new GroupMembers(possible-reason:invalid/no json input)",
				Error:      err.Error(),
			})
	}

	var validationReponse responsemodel_chat.AddNewMembersToGroupResponse
	validate := validator.New(validator.WithRequiredStructEnabled())
	err := validate.Struct(addMembersInput)
	if err != nil {
		if ve, ok := err.(validator.ValidationErrors); ok {
			for _, e := range ve {
				switch e.Field() {
				case "GroupID":
					validationReponse.GroupID = "should contain less than 35 characters"
				case "GroupMembers":
					validationReponse.GroupMembers = "Should be unique,maximum 12 members and id should be a number"
				}
			}
			return ctx.Status(fiber.ErrBadRequest.Code).
				JSON(responsemodel_post.CommonResponse{
					StatusCode: fiber.ErrBadRequest.Code,
					Message:    "can't add new GroupMembers",
					Data:       validationReponse,
					Error:      err.Error(),
				})
		}
	}

	context, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	resp, err := svc.Client.AddMembersToGroup(context, &pb.RequestAddGroupMembers{
		UserID:    fmt.Sprint(userId),
		GroupID:   addMembersInput.GroupID,
		MemberIDs: addMembersInput.GroupMembers,
	})

	if err != nil {
		log.Println("-----error: from handler:AddMembersToGroup(),chatNcall service down while calling AddMembersToGroup()")

		return ctx.Status(fiber.StatusServiceUnavailable).
			JSON(responsemodel_post.CommonResponse{
				StatusCode: fiber.StatusServiceUnavailable,
				Message:    "can't add new GroupMembers",
				Error:      err.Error(),
			})
	}

	if resp.ErrorMessage != "" {
		return ctx.Status(fiber.StatusBadRequest).
			JSON(responsemodel_post.CommonResponse{
				StatusCode: fiber.StatusBadRequest,
				Message:    "can't add new GroupMembers",
				Data:       resp,
				Error:      resp.ErrorMessage,
			})
	}

	return ctx.Status(fiber.StatusOK).
		JSON(responsemodel_post.CommonResponse{
			StatusCode: fiber.StatusOK,
			Message:    "new groupmembers added succesfully",
		})

}

func (svc *ChatWebsockethandler) RemoveAMemberFromGroup(ctx *fiber.Ctx) error {
	// Get the user ID from the context (assuming it's set during authentication)
	userId := ctx.Locals("userId")

	// Define a struct to capture the input data
	var inputData requestmodel_chat.RemoveMemberFromGroup

	// Parse the request body into the inputData struct
	if err := ctx.BodyParser(&inputData); err != nil {
		return ctx.Status(fiber.ErrBadRequest.Code).
			JSON(responsemodel_post.CommonResponse{
				StatusCode: fiber.ErrBadRequest.Code,
				Message:    "can't remove GroupMember (possible reason: invalid or no JSON input)",
				Error:      err.Error(),
			})
	}

	// Validate the input data
	var validationResponse responsemodel_chat.RemoveMemberFromGroup
	validate := validator.New(validator.WithRequiredStructEnabled())
	err := validate.Struct(inputData) // Validate the inputData instead of validationResponse
	if err != nil {
		if ve, ok := err.(validator.ValidationErrors); ok {
			for _, e := range ve {
				switch e.Field() {
				case "GroupID":
					validationResponse.GroupID = "should contain less than 35 characters"
				case "MemberID":
					validationResponse.MemberID = "should be a valid id"
				}
			}
			return ctx.Status(fiber.ErrBadRequest.Code).
				JSON(responsemodel_post.CommonResponse{
					StatusCode: fiber.ErrBadRequest.Code,
					Message:    "can't remove GroupMember",
					Data:       validationResponse,
					Error:      err.Error(),
				})
		}
	}

	// Call the RemoveMemberFromGroup method in the Client (e.g., chat service)
	context, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := svc.Client.RemoveMemberFromGroup(context, &pb.RequestRemoveGroupMember{
		UserID:   fmt.Sprint(userId),
		GroupID:  inputData.GroupID,
		MemberID: inputData.MemberID,
	})

	// Error handling for the call to RemoveMemberFromGroup
	if err != nil {
		log.Println("-----error: from handler:RemoveAMemberFromGroup(), chatNcall service down while calling RemoveMemberFromGroup()")
		return ctx.Status(fiber.StatusServiceUnavailable).
			JSON(responsemodel_post.CommonResponse{
				StatusCode: fiber.StatusServiceUnavailable,
				Message:    "can't remove GroupMember",
				Error:      err.Error(),
			})
	}

	// Check for errors in the response from the chat service
	if resp.ErrorMessage != "" {
		return ctx.Status(fiber.StatusBadRequest).
			JSON(responsemodel_post.CommonResponse{
				StatusCode: fiber.StatusBadRequest,
				Message:    "can't remove GroupMember",
				Data:       resp,
				Error:      resp.ErrorMessage,
			})
	}

	// Return a successful response if everything goes well
	return ctx.Status(fiber.StatusOK).
		JSON(responsemodel_post.CommonResponse{
			StatusCode: fiber.StatusOK,
			Message:    "group member removed successfully",
		})
}

