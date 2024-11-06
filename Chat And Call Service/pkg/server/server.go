package server

import (
	"context"

	"github.com/ShahabazSulthan/friendzy_post/pkg/pb"
	interface_usecase "github.com/ShahabazSulthan/friendzy_post/pkg/usecase/interface"
)

type ChatNCallSvc struct {
	ChatUsecase interface_usecase.IChatUseCase
	pb.ChatNCallServiceServer
}

func NewChatNCallServiceServer(chatUseCase interface_usecase.IChatUseCase) *ChatNCallSvc {
	return &ChatNCallSvc{ChatUsecase: chatUseCase}
}

func (c *ChatNCallSvc) GetOneToOneChats(ctx context.Context, req *pb.RequestUserOneToOneChat) (*pb.ResponseUserOneToOneChat, error) {
	// Fetch chat data using the use case method
	respData, err := c.ChatUsecase.GetOneToOneChats(&req.SenderID, &req.RecieverID, &req.Limit, &req.Offset)
	if err != nil {
		return &pb.ResponseUserOneToOneChat{
			ErrorMessage: "Failed to fetch chats: " + err.Error(),
		}, nil
	}

	// Prepare response data
	var chatData []*pb.SingleOneToOneChat
	for _, chat := range *respData {
		chatData = append(chatData, &pb.SingleOneToOneChat{
			MessageID:  chat.MessageID,
			SenderID:   chat.SenderID,
			RecieverID: chat.RecipientID,
			Content:    chat.Content,
			Status:     chat.Status,
			TimeStamp:  chat.StringTime,
		})
	}

	return &pb.ResponseUserOneToOneChat{
		Chat: chatData,
	}, nil
}

func (c *ChatNCallSvc) GetRecentChatProfiles(ctx context.Context, req *pb.RequestRecentChatProfiles) (*pb.ResponseRecentChatProfiles, error) {
	// Fetch recent chat profiles and chat data from the use case
	respData, err := c.ChatUsecase.GetRecentChatProfilesPlusChatData(&req.SenderID, &req.Limit, &req.Offset)
	if err != nil {
		return &pb.ResponseRecentChatProfiles{
			ErrorMessage: "Error retrieving recent chat profiles: " + err.Error(),
		}, nil
	}

	// Prepare response data
	var profileData []*pb.SingelUserAndLastChat
	for _, chatProfile := range *respData {
		profileData = append(profileData, &pb.SingelUserAndLastChat{
			UserID:               chatProfile.UserId,
			UserName:             chatProfile.UserName,
			UserProfileURL:       chatProfile.UserProfileImgURL,
			LastMessageContent:   chatProfile.Content,
			LastMessageTimeStamp: chatProfile.StringTime,
		})
	}

	return &pb.ResponseRecentChatProfiles{
		ActualData: profileData,
	}, nil
}
