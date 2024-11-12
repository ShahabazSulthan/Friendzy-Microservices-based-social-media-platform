package server

import (
	"context"
	"log"
	"time"

	"github.com/ShahabazSulthan/friendzy_post/pkg/model/requestmodel"
	"github.com/ShahabazSulthan/friendzy_post/pkg/pb"
	interface_usecase "github.com/ShahabazSulthan/friendzy_post/pkg/usecase/interface"
)

type ChatNCallSvc struct {
	ChatUsecase interface_usecase.IChatUseCase
	pb.ChatNCallServiceServer
}

func NewChatNCallServiceServer(chatUseCase interface_usecase.IChatUseCase) *ChatNCallSvc {
	return &ChatNCallSvc{
		ChatUsecase: chatUseCase,
	}
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

func (c *ChatNCallSvc) CreateNewGroup(ctx context.Context, req *pb.RequestNewGroup) (*pb.ResponseNewGroup, error) {
	// Map request data to NewGroupInfo model
	groupDataInput := requestmodel.NewGroupInfo{
		GroupName:    req.GroupName,
		GroupMembers: req.GroupMembers,
		CreatorID:    req.CreatorID,
		CreateAt:     time.Now(),
	}

	// Call use case to create the group
	if err := c.ChatUsecase.CreateNewGroup(&groupDataInput); err != nil {
		log.Printf("Error creating new group in usecase: %v", err)
		return &pb.ResponseNewGroup{ErrorMessage: err.Error()}, nil
	}

	// Return successful response
	return &pb.ResponseNewGroup{}, nil
}

func (c *ChatNCallSvc) GetGroupMembersInfo(ctx context.Context, req *pb.RequestGroupMembersInfo) (*pb.ResponseGroupMembersInfo, error) {

	groupMembers, err := c.ChatUsecase.GroupMembersList(&req.GroupID)
	if err != nil {
		return &pb.ResponseGroupMembersInfo{ErrorMessage: err.Error()}, nil
	}

	return &pb.ResponseGroupMembersInfo{GroupMembers: *groupMembers}, nil
}

func (c *ChatNCallSvc) GetUserGroupChatSummary(ctx context.Context, req *pb.RequestGroupChatSummary) (*pb.ResponseGroupChatSummary, error) {

	chatSummary, err := c.ChatUsecase.GetUserGroupChatSummary(&req.UserID, &req.Limit, &req.Offset)
	if err != nil {
		return &pb.ResponseGroupChatSummary{ErrorMessage: err.Error()}, nil
	}

	var singleSummarySlice []*pb.SingleGroupChatDetails

	for i := range *chatSummary {
		singleSummarySlice = append(singleSummarySlice, &pb.SingleGroupChatDetails{
			GroupID:              (*chatSummary)[i].GroupID,
			GroupName:            (*chatSummary)[i].GroupName,
			GroupProfileImageURL: (*chatSummary)[i].GroupProfileImgURL,
			LastMessageContent:   (*chatSummary)[i].LastMessage,
			TimeStamp:            (*chatSummary)[i].StringTime,
			SenderID:             (*chatSummary)[i].SenderID,
			SenderUserName:       (*chatSummary)[i].SenderUserName,
		})

	}

	return &pb.ResponseGroupChatSummary{SingleEntity: singleSummarySlice}, nil
}

func (c *ChatNCallSvc) GetOneToManyChats(ctx context.Context, req *pb.RequestGetOneToManyChats) (*pb.ResponseGetOneToManyChats, error) {

	chatData, err := c.ChatUsecase.GetOneToManyChats(&req.UserID, &req.GroupID, &req.Limit, &req.Offset)
	if err != nil {
		return &pb.ResponseGetOneToManyChats{ErrorMessage: err.Error()}, nil
	}

	var repeatedData []*pb.SingleOneToManyChat
	for i := range *chatData {
		repeatedData = append(repeatedData, &pb.SingleOneToManyChat{
			MessageID:             (*chatData)[i].MessageID,
			SenderID:              (*chatData)[i].SenderID,
			SenderUserName:        (*chatData)[i].SenderUserName,
			SenderProfileImageURL: (*chatData)[i].SenderProfileImgURL,
			GroupID:               req.GroupID,
			Content:               (*chatData)[i].Content,
			TimeStamp:             (*chatData)[i].StringTime,
		})
	}

	return &pb.ResponseGetOneToManyChats{Chat: repeatedData}, nil
}

func (c *ChatNCallSvc) AddMembersToGroup(ctx context.Context, req *pb.RequestAddGroupMembers) (*pb.ResponseAddGroupMembers, error) {
	var inputData requestmodel.AddNewMemberToGroup

	inputData.UserID = req.UserID
	inputData.GroupID = req.GroupID
	inputData.GroupMembers = req.MemberIDs

	err := c.ChatUsecase.AddNewMembersToGroup(&inputData)
	if err != nil {
		return &pb.ResponseAddGroupMembers{ErrorMessage: err.Error()}, nil
	}

	return &pb.ResponseAddGroupMembers{}, nil
}

func (c *ChatNCallSvc) RemoveMemberFromGroup(ctx context.Context, req *pb.RequestRemoveGroupMember) (*pb.ResponseRemoveGroupMember, error) {
	var inputData requestmodel.RemoveMemberFromGroup

	inputData.UserID = req.UserID
	inputData.GroupID = req.GroupID
	inputData.MemberID = req.MemberID

	err := c.ChatUsecase.RemoveMemberFromGroup(&inputData)
	if err != nil {
		return &pb.ResponseRemoveGroupMember{ErrorMessage: err.Error()}, nil
	}
	return &pb.ResponseRemoveGroupMember{}, nil
}
