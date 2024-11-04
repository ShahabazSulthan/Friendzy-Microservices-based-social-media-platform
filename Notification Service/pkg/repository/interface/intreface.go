package interface_notification

import (
	"time"

	"github.com/ShahabazSulthan/Friendzy_Notification/pkg/models/requestmodel"
	"github.com/ShahabazSulthan/Friendzy_Notification/pkg/models/responsemodel"
)

type INotifRepo interface {
	CreateNewNotification(msg *requestmodel.KafkaNotification) error
	GetNotificationsForUser(userId, limit, offset *string) (*[]responsemodel.NotificationModel, error)
	CalculatePostAge(createdAt time.Time) (string, error)
}
