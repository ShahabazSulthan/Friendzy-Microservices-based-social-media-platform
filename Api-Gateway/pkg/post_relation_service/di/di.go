package di_post

import (
	"github.com/ShahabazSulthan/Friendzy_apiGateway/pkg/Auth_Service/middleware"
	config "github.com/ShahabazSulthan/Friendzy_apiGateway/pkg/Config"
	client_post "github.com/ShahabazSulthan/Friendzy_apiGateway/pkg/post_relation_service/client"
	handler_post "github.com/ShahabazSulthan/Friendzy_apiGateway/pkg/post_relation_service/handler"
	router_post "github.com/ShahabazSulthan/Friendzy_apiGateway/pkg/post_relation_service/routes"
	"github.com/gofiber/fiber/v2"
)

func InitPostNrelClient(app *fiber.App, config *config.Config, middleware *middleware.Middleware) error {

	client, err := client_post.InitPostNrelClient(config)
	if err != nil {
		return err
	}

	postHandler := handler_post.NewPostHandler(client)
	relationHandler := handler_post.NewRelationHandler(client)
	commentHandler := handler_post.NewCommentHandler(client)

	router_post.PostNrelUserRoutes(app, postHandler, middleware, relationHandler, commentHandler)

	return nil
}
