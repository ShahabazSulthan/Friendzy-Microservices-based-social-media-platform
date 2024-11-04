package server

import (
	"context"

	"github.com/ShahabazSulthan/Friendzy_Notification/pkg/pb"
	interface_usecase "github.com/ShahabazSulthan/Friendzy_Notification/pkg/usecase/interface"
)

type NotifService struct {
	NotifUseCase interface_usecase.INotifUseCase
	pb.NotificationServiceServer
}

func NewNotifServiceServer(notifUseCase interface_usecase.INotifUseCase) *NotifService {
	return &NotifService{NotifUseCase: notifUseCase}
}

func (u *NotifService) GetNotificationsForUser(ctx context.Context, req *pb.RequestGetNotifications) (*pb.ResponseGetNotifications, error) {

	notificationsData, err := u.NotifUseCase.GetNotificationsForUser(&req.UserId, &req.Limit, &req.OffSet)
	if err != nil {
		return &pb.ResponseGetNotifications{
			ErrorMessage: err.Error(),
		}, nil
	}

	var SingleNotification []*pb.SingleNotification

	for i := range *notificationsData {
		SingleNotification = append(SingleNotification, &pb.SingleNotification{
			NotificationID:     (*notificationsData)[i].NotificaitonID,
			UserID:             (*notificationsData)[i].UserID,
			ActorID:            (*notificationsData)[i].ActorID,
			ActorUserName:      (*notificationsData)[i].ActorUserName,
			ActorProfileImgURL: (*notificationsData)[i].ActorProfileImgURL,
			ActionType:         (*notificationsData)[i].ActionType,
			TargetID:           (*notificationsData)[i].TargetID,
			TargetType:         (*notificationsData)[i].TargetType,
			CommentText:        (*notificationsData)[i].CommentText,
			NotificationAge:    (*notificationsData)[i].NotificationAge,
		})
	}

	return &pb.ResponseGetNotifications{
		Notifications: SingleNotification,
	}, nil

}
