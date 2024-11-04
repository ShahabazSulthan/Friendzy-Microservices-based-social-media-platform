package handler_notif

import (
	"context"
	"fmt"
	"time"

	config "github.com/ShahabazSulthan/Friendzy_apiGateway/pkg/Config"
	"github.com/ShahabazSulthan/Friendzy_apiGateway/pkg/Notification_Service/pb"
	responsemodel_notif "github.com/ShahabazSulthan/Friendzy_apiGateway/pkg/Notification_Service/responsemodel"
	"github.com/gofiber/fiber/v2"
)

type NotifHandler struct {
	Client pb.NotificationServiceClient
	Config *config.Config
}

func NewNotificationHandler(
	client *pb.NotificationServiceClient,
	config *config.Config) *NotifHandler {
	return &NotifHandler{
		Client: *client,
		Config: config,
	}
}

func (svc *NotifHandler) GetNotificationsForUser(ctx *fiber.Ctx) error {

	userId := ctx.Locals("userId")
	limit, offset := ctx.Query("limit", "10"), ctx.Query("offset", "0")

	context, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	resp, err := svc.Client.GetNotificationsForUser(context, &pb.RequestGetNotifications{
		UserId: fmt.Sprint(userId),
		Limit:  limit,
		OffSet: offset,
	})

	if err != nil {
		fmt.Println("----------notification service down--------")

		return ctx.Status(fiber.StatusServiceUnavailable).
			JSON(responsemodel_notif.CommonResponse{
				StatusCode: fiber.StatusServiceUnavailable,
				Message:    "can't fetch Notifications",
				Error:      err.Error(),
			})
	}

	if resp.ErrorMessage != "" {
		return ctx.Status(fiber.StatusBadRequest).
			JSON(responsemodel_notif.CommonResponse{
				StatusCode: fiber.StatusBadRequest,
				Message:    "can't fetch Notifications",
				Data:       resp,
				Error:      resp.ErrorMessage,
			})
	}

	return ctx.Status(fiber.StatusOK).
		JSON(responsemodel_notif.CommonResponse{
			StatusCode: fiber.StatusOK,
			Message:    "Notifications fetched succesfully",
			Data:       resp,
			Error:      nil,
		})
}
