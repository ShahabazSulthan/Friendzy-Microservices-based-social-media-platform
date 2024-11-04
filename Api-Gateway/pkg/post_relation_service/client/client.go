package client_post

import (
	"fmt"

	config "github.com/ShahabazSulthan/Friendzy_apiGateway/pkg/Config"
	"github.com/ShahabazSulthan/Friendzy_apiGateway/pkg/post_relation_service/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitPostNrelClient(config *config.Config) (*pb.PostNrelServiceClient, error) {
	cc, err := grpc.Dial(config.PostNrelSvcUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println("-------", err)
		return nil, err
	}

	Client := pb.NewPostNrelServiceClient(cc)

	return &Client, nil
}
