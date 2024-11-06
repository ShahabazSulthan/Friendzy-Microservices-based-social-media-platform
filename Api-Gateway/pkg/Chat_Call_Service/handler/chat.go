package handler_chat

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/IBM/sarama"
	requestmodel_chat "github.com/ShahabazSulthan/Friendzy_apiGateway/pkg/Chat_Call_Service/model/requestmodel"
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
			case "TypingStatus":
				svc.TypingStatus(&messageModel)
			default:
				sendErrMessageWS(userIdStr, errors.New("message Type should be OneToOne,OneToMany,DeleteMessage,UpdateSeenStatus or TypingStatus ,no other types allowed"))
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
