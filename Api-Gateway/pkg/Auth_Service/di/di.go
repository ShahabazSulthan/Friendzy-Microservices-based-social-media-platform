package di_auth

import (
	"log"

	client_auth "github.com/ShahabazSulthan/Friendzy_apiGateway/pkg/Auth_Service/client"
	handler_auth "github.com/ShahabazSulthan/Friendzy_apiGateway/pkg/Auth_Service/handler"
	"github.com/ShahabazSulthan/Friendzy_apiGateway/pkg/Auth_Service/middleware"
	router_auth "github.com/ShahabazSulthan/Friendzy_apiGateway/pkg/Auth_Service/router"
	config "github.com/ShahabazSulthan/Friendzy_apiGateway/pkg/Config"
	"github.com/gofiber/fiber/v2"
)

func InitAuthClient(app *fiber.App, config *config.Config) (*middleware.Middleware, error) {

	client, err := client_auth.InitAuthClient(config)
	if err != nil {
		log.Fatal(err)
	}

	middleware := middleware.NewAuthMiddleware(client)

	userHandler := handler_auth.NewAuthUserHandler(client)

	router_auth.AuthUserRoutes(app, userHandler, middleware)

	return middleware, nil
}
