package interface_usecase

import "github.com/ShahabazSulthan/friendzy_post/pkg/model/responsemodel"

type IChatUseCase interface {
	KafkaOneToOneMessageConsumer()
	GetOneToOneChats(senderId, recipientId, limit, offset *string) (*[]responsemodel.OneToOneChatResponse, error)
	GetRecentChatProfilesPlusChatData(senderid, limit, offset *string) (*[]responsemodel.RecentChatProfileResponse, error)
}
