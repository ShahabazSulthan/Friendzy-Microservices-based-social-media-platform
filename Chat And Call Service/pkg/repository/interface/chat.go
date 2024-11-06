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
}
