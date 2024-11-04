package main

import (
	"fmt"
	"log"

	di_auth "github.com/ShahabazSulthan/Friendzy_apiGateway/pkg/Auth_Service/di"
	config "github.com/ShahabazSulthan/Friendzy_apiGateway/pkg/Config"
	di_post "github.com/ShahabazSulthan/Friendzy_apiGateway/pkg/post_relation_service/di"
	"github.com/gofiber/fiber/v2"
)

func main() {

	config, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	app := fiber.New()

	middleware, err := di_auth.InitAuthClient(app, config)
	if err != nil {
		log.Fatal(err)
	}

	err = di_post.InitPostNrelClient(app, config, middleware)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Port Running  ", config.Port)

	err = app.Listen(config.Port)
	if err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}
