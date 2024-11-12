package router_chat

import (
	"fmt"

	"github.com/ShahabazSulthan/Friendzy_apiGateway/pkg/Auth_Service/middleware"
	handler_chat "github.com/ShahabazSulthan/Friendzy_apiGateway/pkg/Chat_Call_Service/handler"
	responsemodel_post "github.com/ShahabazSulthan/Friendzy_apiGateway/pkg/post_relation_service/model/responsemodel"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func ChatNcallRoutes(app *fiber.App,
	webSocHandler *handler_chat.ChatWebsockethandler,
	middleware *middleware.Middleware) {

	app.Use(middleware.UserAuthorizationMiddleWare)
	{
		chatManagement := app.Group("/chat")
		{
			chatManagement.Get("/onetoonechats/:recipientid", webSocHandler.GetOneToOneChats)
			chatManagement.Get("/recentchatprofiles/", webSocHandler.GetRecentChatProfileDetails)

			chatManagement.Use(HttptoWsConnectionUpgrader)
			{
				chatManagement.Get("/ws", webSocHandler.WsConnection)

			}

		}

		groupChatManagemanent := app.Group("/groupchat")
		{
			groupChatManagemanent.Post("/", webSocHandler.CreateNewGroup)
			groupChatManagemanent.Post("/add", webSocHandler.AddMembersToGroup)
			groupChatManagemanent.Post("/remove", webSocHandler.RemoveAMemberFromGroup)
			groupChatManagemanent.Get("/summary", webSocHandler.GetUserGroupsAndLastMessage)
			groupChatManagemanent.Get("/:groupid", webSocHandler.GetGroupChats)
		}

	}
}

func HttptoWsConnectionUpgrader(ctx *fiber.Ctx) error {
	if websocket.IsWebSocketUpgrade(ctx) {
		ctx.Locals("allowed", true)
		return ctx.Next()
	}

	fmt.Println("-------------websocket.IsWebSocketUpgrade(ctx)------------,returned false:::::::::::::")
	return ctx.Status(fiber.ErrUpgradeRequired.Code).
		JSON(responsemodel_post.CommonResponse{
			StatusCode: fiber.ErrUpgradeRequired.Code,
			Message:    "requires websocket connection",
		})
}
