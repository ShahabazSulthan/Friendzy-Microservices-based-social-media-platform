package client_auth

import (
	"fmt"

	"github.com/ShahabazSulthan/Friendzy_apiGateway/pkg/Auth_Service/pb"
	config "github.com/ShahabazSulthan/Friendzy_apiGateway/pkg/Config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitAuthClient(Config *config.Config) (*pb.AuthServiceClient, error) {
	cc, err := grpc.Dial(Config.AuthSvcUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println("Error Connect auth Clinet : ", err)
		return nil, err
	}

	client := pb.NewAuthServiceClient(cc)

	return &client, nil
}
