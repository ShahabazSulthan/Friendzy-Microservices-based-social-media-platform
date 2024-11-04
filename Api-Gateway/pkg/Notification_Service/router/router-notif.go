package router_notif

import (
	"github.com/ShahabazSulthan/Friendzy_apiGateway/pkg/Auth_Service/middleware"
	handler_notif "github.com/ShahabazSulthan/Friendzy_apiGateway/pkg/Notification_Service/handler"
	"github.com/gofiber/fiber/v2"
)

func NotificationRoutes(app *fiber.App,
	NotifHandler *handler_notif.NotifHandler,
	middleware *middleware.Middleware) {

	app.Use(middleware.UserAuthorizationMiddleWare)
	{
		notifManger := app.Group("/notification")
		{
			notifManger.Get("/", NotifHandler.GetNotificationsForUser)
		}

	}

}
