package client_notif

import (
	"fmt"

	config "github.com/ShahabazSulthan/Friendzy_apiGateway/pkg/Config"
	"github.com/ShahabazSulthan/Friendzy_apiGateway/pkg/Notification_Service/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitNotificationClient(config *config.Config) (*pb.NotificationServiceClient, error) {
	cc, err := grpc.Dial(config.NotifSvcUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println("------err-", err)
		return nil, err
	}

	Client := pb.NewNotificationServiceClient(cc)

	return &Client, nil
}
