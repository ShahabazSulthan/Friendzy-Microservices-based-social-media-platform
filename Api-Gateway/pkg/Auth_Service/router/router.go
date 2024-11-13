package router_auth

import (
	handler_auth "github.com/ShahabazSulthan/Friendzy_apiGateway/pkg/Auth_Service/handler"
	"github.com/ShahabazSulthan/Friendzy_apiGateway/pkg/Auth_Service/middleware"
	"github.com/gofiber/fiber/v2"
)

func AuthUserRoutes(app *fiber.App, userHandler *handler_auth.UserHandler, middleware *middleware.Middleware) {

	app.Get("/verify/:userID/:verificationID", userHandler.OnlinePayment)

	// Admin routes (protected with AdminAuthMiddleware)
	adminRoutes := app.Group("/admin")
	{
		// Admin login is public
		adminRoutes.Post("/login", userHandler.AdminLogin)

		// Apply admin middleware to protected admin routes
		adminRoutes.Use(middleware.AdminAuthMiddleware)

		// Protected admin routes
		adminRoutes.Patch("/blockuser", userHandler.BlockUser)
		adminRoutes.Patch("/unblockuser", userHandler.UnblockUser)
		adminRoutes.Get("/users", userHandler.GetAllUsers)
	}

	// Public routes
	app.Post("/signup", userHandler.UserSignUp)
	app.Post("/verify", userHandler.UserOTPVerication)
	app.Post("/login", userHandler.UserLogin)
	app.Post("/forgotpassword", userHandler.ForgotPasswordRequest)
	app.Patch("/resetpassword", userHandler.ResetPassword) // No middleware here
	app.Get("/accessgenerator", middleware.AccessRegenerator)
	app.Get("/log", userHandler.GetLogFile)
	

	// Apply user authentication middleware for user-protected routes
	app.Use(middleware.UserAuthorizationMiddleWare)

	app.Get("/bluetickusers",userHandler.GetAllVerifiedUsers)

	// User profile management routes (protected with UserAuthorizationMiddleWare)
	profileManagement := app.Group("/profile")
	{
		profileManagement.Get("/", userHandler.GetUserProfile)
		profileManagement.Patch("/edit", userHandler.EditUserProfile)
		profileManagement.Post("/setprofileimage", userHandler.SetProfileImage)
		profileManagement.Get("/followers", userHandler.GetFollowersDetails)
		profileManagement.Get("/following", userHandler.GetFollowingsDetails)
		profileManagement.Post("/blueTick",userHandler.CreateBlueTickPaymentHandler)
		profileManagement.Post("/verifypayment",userHandler.VerifyBlueTickPaymentHandler)
	}

	// Explore management routes (protected with UserAuthorizationMiddleWare)
	exploremanagement := app.Group("/explore")
	{
		exploremanagement.Get("/profile/:userbid", userHandler.GetAnotherUserProfile)

		// Search routes (protected with UserAuthorizationMiddleWare)
		searchmanagement := exploremanagement.Group("/search")
		{
			searchmanagement.Get("/user/:searchtext", userHandler.SearchUser)
		}
	}
}
