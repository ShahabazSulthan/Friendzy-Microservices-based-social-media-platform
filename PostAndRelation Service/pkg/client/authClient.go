package client

import (
	"fmt"

	"github.com/ShahabazSulthan/friendzy_post/pkg/config"
	"github.com/ShahabazSulthan/friendzy_post/pkg/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitAuthServiceClient(config *config.Config) (*pb.AuthServiceClient, error) {
	cc, err := grpc.Dial(config.PortMngr.AuthSvcUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println("Error in Init Auth client : ", err)
		return nil, err
	}

	client := pb.NewAuthServiceClient(cc)

	return &client, nil
}
