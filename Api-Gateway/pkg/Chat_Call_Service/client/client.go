package client_chat

import (
	"fmt"

	"github.com/ShahabazSulthan/Friendzy_apiGateway/pkg/Chat_Call_Service/pb"
	config "github.com/ShahabazSulthan/Friendzy_apiGateway/pkg/Config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitChatNcallClient(config *config.Config) (*pb.ChatNCallServiceClient, error) {
	cc, err := grpc.Dial(config.ChatSvcUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println("-------", err)
		return nil, err
	}

	Client := pb.NewChatNCallServiceClient(cc)

	return &Client, nil
}
