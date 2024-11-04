package di

import (
	"fmt"

	"github.com/ShahabazSulthan/Friendzy_Notification/pkg/client"
	"github.com/ShahabazSulthan/Friendzy_Notification/pkg/config"
	"github.com/ShahabazSulthan/Friendzy_Notification/pkg/db"
	"github.com/ShahabazSulthan/Friendzy_Notification/pkg/repository"
	"github.com/ShahabazSulthan/Friendzy_Notification/pkg/server"
	"github.com/ShahabazSulthan/Friendzy_Notification/pkg/usecase"
	"github.com/ShahabazSulthan/Friendzy_Notification/pkg/utils/hash"
)

func InitializeNotificationServer(config *config.Config) (*server.NotifService, error) {

	hashUtil := hash.NewHashUtil()

	DB, err := db.ConnectDatabase(&config.DB, hashUtil)
	if err != nil {
		fmt.Println("ERROR CONNECTING DB FROM DI.GO")
		return nil, err
	}

	authClient, err := client.InitAuthServiceClient(config)
	if err != nil {
		fmt.Println("--------err--------", err)
		return nil, err
	}

	notifRepo := repository.NewNotifRepo(DB)
	notifUseCase := usecase.NewNotifUseCase(notifRepo, config.KafkaConfig, authClient)

	go notifUseCase.KafkaMessageConsumer()

	return server.NewNotifServiceServer(notifUseCase), nil
}
