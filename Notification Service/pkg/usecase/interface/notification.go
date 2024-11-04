package interface_usecase

import "github.com/ShahabazSulthan/Friendzy_Notification/pkg/models/responsemodel"

type INotifUseCase interface {
	KafkaMessageConsumer()
	GetNotificationsForUser(userId, limit, offset *string) (*[]responsemodel.NotificationModel, error)
}
