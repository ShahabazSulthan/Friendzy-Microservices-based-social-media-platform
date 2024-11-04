package client

import (
	"fmt"

	"github.com/ShahabazSulthan/Friendzy_Notification/pkg/config"
	"github.com/ShahabazSulthan/Friendzy_Notification/pkg/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitAuthServiceClient(config *config.Config) (*pb.AuthServiceClient, error) {
	cc, err := grpc.Dial(config.PortMngr.AuthSvcUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println("-------", err)
		return nil, err
	}

	Client := pb.NewAuthServiceClient(cc)

	return &Client, nil
}
