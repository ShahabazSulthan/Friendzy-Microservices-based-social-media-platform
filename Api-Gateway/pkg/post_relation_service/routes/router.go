package router_post

import (
	"github.com/ShahabazSulthan/Friendzy_apiGateway/pkg/Auth_Service/middleware"
	handler_post "github.com/ShahabazSulthan/Friendzy_apiGateway/pkg/post_relation_service/handler"
	"github.com/gofiber/fiber/v2"
)

func PostNrelUserRoutes(app *fiber.App,
	postHandler *handler_post.PostHandler,
	middleware *middleware.Middleware,
	relationHandler *handler_post.RelationHandler,
	commentHandler *handler_post.CommentHandler) {

	app.Use(middleware.UserAuthorizationMiddleWare)
	{
		postManagement := app.Group("/post")
		{
			postManagement.Post("/", postHandler.AddNewPost)
			postManagement.Get("/", postHandler.GetAllPostByUser)
			postManagement.Delete("/:postid", postHandler.DeletePost)
			postManagement.Patch("/", postHandler.EditPost)

			likemanagement := postManagement.Group("/like")
			{
				likemanagement.Post("/:postid", postHandler.LikePost)
				likemanagement.Delete("/:postid", postHandler.UnLikePost)
			}

			commentManagement := postManagement.Group("/comment")
			{
				commentManagement.Get("/:postid", commentHandler.FetchPostComments)
				commentManagement.Post("/", commentHandler.AddComment)
				commentManagement.Delete("/:commentid", commentHandler.DeleteComment)
				commentManagement.Patch("/", commentHandler.EditComment)
			}

		}
		followRelationshipManagement := app.Group("/relation")
		{
			followRelationshipManagement.Post("/follow/:followingid", relationHandler.Follow)
			followRelationshipManagement.Delete("/unfollow/:unfollowingid", relationHandler.UnFollow)
		}

		exploremanagement := app.Group("/explore")
		{
			exploremanagement.Get("/", postHandler.GetMostLovedPostsFromGlobalUser)
			exploremanagement.Get("/home", postHandler.GetAllRelatedPostsForHomeScreen)
			exploremanagement.Get("/random", postHandler.GetAllRandomPosts)
		}

	}
}
