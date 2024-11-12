package interface_chat

import (
	"github.com/ShahabazSulthan/friendzy_post/pkg/model/requestmodel"
	"github.com/ShahabazSulthan/friendzy_post/pkg/model/responsemodel"
)

type IChatRepo interface {
	StoreOneToOneChatToDB(chatData *requestmodel.OneToOneChatRequest) (*string, error)
	UpdateChatStatus(senderId, recipientId *string) error
	GetOneToOneChats(senderId, recipientId, limit, offset *string) (*[]responsemodel.OneToOneChatResponse, error)
	RecentChatProfileData(senderid, limit, offset *string) (*[]responsemodel.RecentChatProfileResponse, error)

	CreateNewGroup(groupInfo *requestmodel.NewGroupInfo) error
	GetGroupMembersList(groupId *string) (*[]uint64, error)
	StoreOneToManyChatToDB(msg *requestmodel.OneToManyMessageRequest) error
	GetRecentGroupProfilesOfUser(userId, limit, offset *string) (*[]responsemodel.GroupInfoLite, error)

	GetGroupLastMessageDetailsByGroupId(groupid *string) (*responsemodel.OneToManyMessageLite, error)
	CheckUserIsGroupMember(userid, groupid *string) (bool, error)
	GetOneToManyChats(groupId, limit, offset *string) (*[]responsemodel.OneToManyChatResponse, error)
	AddNewMembersToGroupByGroupId(inputData *requestmodel.AddNewMemberToGroup) error
	
	RemoveGroupMember(inputData *requestmodel.RemoveMemberFromGroup) error
	CountMembersInGroup(groupId string) (int, error)
	DeleteOneToManyChatsByGroupId(groupId string) error
	DeleteGroupDataByGroupId(groupId string) error
}
