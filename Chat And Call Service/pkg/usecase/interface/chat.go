package interface_usecase

import (
	"github.com/ShahabazSulthan/friendzy_post/pkg/model/requestmodel"
	"github.com/ShahabazSulthan/friendzy_post/pkg/model/responsemodel"
)

type IChatUseCase interface {
	KafkaOneToOneMessageConsumer()
	KafkaOneToManyMessageConsumer()
	
	GetOneToOneChats(senderId, recipientId, limit, offset *string) (*[]responsemodel.OneToOneChatResponse, error)
	GetRecentChatProfilesPlusChatData(senderid, limit, offset *string) (*[]responsemodel.RecentChatProfileResponse, error)

	CreateNewGroup(groupInfo *requestmodel.NewGroupInfo) error
	GroupMembersList(groupId *string) (*[]string, error)

	GetUserGroupChatSummary(userId, limit, offset *string) (*[]responsemodel.GroupChatSummaryResponse, error)
	GetOneToManyChats(userid, groupid, limit, offset *string) (*[]responsemodel.OneToManyChatResponse, error)

	AddNewMembersToGroup(inputData *requestmodel.AddNewMemberToGroup) error
	RemoveMemberFromGroup(inputData *requestmodel.RemoveMemberFromGroup) error
}
