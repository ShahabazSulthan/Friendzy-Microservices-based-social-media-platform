package di_chat

import (
	"github.com/ShahabazSulthan/Friendzy_apiGateway/pkg/Auth_Service/middleware"
	client_chat "github.com/ShahabazSulthan/Friendzy_apiGateway/pkg/Chat_Call_Service/client"
	handler_chat "github.com/ShahabazSulthan/Friendzy_apiGateway/pkg/Chat_Call_Service/handler"
	router_chat "github.com/ShahabazSulthan/Friendzy_apiGateway/pkg/Chat_Call_Service/router"
	config "github.com/ShahabazSulthan/Friendzy_apiGateway/pkg/Config"
	"github.com/gofiber/fiber/v2"
)

func InitChatNcallClient(app *fiber.App, config *config.Config, middleware *middleware.Middleware) error {

	client, err := client_chat.InitChatNcallClient(config)
	if err != nil {
		return err
	}

	webSocHandler := handler_chat.NewChatWebsocketHandler(client, config)

	router_chat.ChatNcallRoutes(app, webSocHandler, middleware)

	return nil
}
