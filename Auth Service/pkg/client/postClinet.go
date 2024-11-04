package client

import (
	"fmt"

	"github.com/ShahabazSulthan/Friendzy_Auth/pkg/config"
	"github.com/ShahabazSulthan/Friendzy_Auth/pkg/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitPostClientService(config *config.Config) (*pb.PostNrelServiceClient, error) {
	cc, err := grpc.Dial(config.PortMngr.PostNrelSvcUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println("Error in connect post clinet")
		return nil, err
	}

	client := pb.NewPostNrelServiceClient(cc)

	return &client, nil
}
