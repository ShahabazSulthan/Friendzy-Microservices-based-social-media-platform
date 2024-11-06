package client

import (
	"fmt"

	"github.com/ShahabazSulthan/friendzy_post/pkg/config"
	"github.com/ShahabazSulthan/friendzy_post/pkg/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitAuthServiceClient(config *config.PortManager) (*pb.AuthServiceClient, error) {
	cc, err := grpc.NewClient(config.AuthSvcUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println("-------", err)
		return nil, err
	}

	Client := pb.NewAuthServiceClient(cc)

	return &Client, nil
}

