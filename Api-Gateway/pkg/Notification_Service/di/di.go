package di_notif

import (
	"github.com/ShahabazSulthan/Friendzy_apiGateway/pkg/Auth_Service/middleware"
	config "github.com/ShahabazSulthan/Friendzy_apiGateway/pkg/Config"
	client_notif "github.com/ShahabazSulthan/Friendzy_apiGateway/pkg/Notification_Service/client"
	handler_notif "github.com/ShahabazSulthan/Friendzy_apiGateway/pkg/Notification_Service/handler"
	router_notif "github.com/ShahabazSulthan/Friendzy_apiGateway/pkg/Notification_Service/router"
	"github.com/gofiber/fiber/v2"
)

func InitNotificationClient(app *fiber.App, config *config.Config, middleware *middleware.Middleware) error {

	client, err := client_notif.InitNotificationClient(config)
	if err != nil {
		return err
	}

	notifhandler := handler_notif.NewNotificationHandler(client, config)

	router_notif.NotificationRoutes(app, notifhandler, middleware)

	return nil
}
