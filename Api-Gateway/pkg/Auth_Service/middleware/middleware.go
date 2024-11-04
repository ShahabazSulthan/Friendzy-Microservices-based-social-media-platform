package middleware

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	responsemodel_auth "github.com/ShahabazSulthan/Friendzy_apiGateway/pkg/Auth_Service/models/responseModel"
	"github.com/ShahabazSulthan/Friendzy_apiGateway/pkg/Auth_Service/pb"
	"github.com/gofiber/fiber/v2"
)

type Middleware struct {
	Client pb.AuthServiceClient
}

func NewAuthMiddleware(client *pb.AuthServiceClient) *Middleware {
	return &Middleware{Client: *client}
}

func (m *Middleware) UserAuthorizationMiddleWare(ctx *fiber.Ctx) error {
	accessToken := ctx.Get("x-access-token")

	if accessToken == "" || len(accessToken) <= 20 {
		return ctx.Status(fiber.StatusUnauthorized).
			JSON(responsemodel_auth.CommonResponse{
				StatusCode: fiber.StatusUnauthorized,
				Message:    "error parsing access token from request",
				Error:      errors.New("error praising access token from request"),
			})
	}

	context, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	resp, err := m.Client.VerifyAccessToken(context, &pb.RequestVerifyAccess{
		AccessToken: accessToken,
	})

	if err != nil {
		fmt.Println("=========Auth-Service-Down============")
		return ctx.Status(fiber.StatusUnauthorized).
			JSON(responsemodel_auth.CommonResponse{
				StatusCode: fiber.StatusUnauthorized,
				Message:    "failed to verify accesstoken",
				Error:      err.Error(),
			})
	}

	if resp.ErrorMessage != "" {
		return ctx.Status(fiber.StatusBadRequest).
			JSON(responsemodel_auth.CommonResponse{
				StatusCode: fiber.StatusBadRequest,
				Message:    "failed to verify accesstoken",
				Data:       resp,
				Error:      resp.ErrorMessage,
			})
	}

	ctx.Locals("userId", resp.UserId)
	return ctx.Next()
}

func (m *Middleware) AdminAuthMiddleware(ctx *fiber.Ctx) error {
	// Extract token from the "x-admin-token" header
	adminToken := ctx.Get("x-admin-token")

	if adminToken == "" || len(adminToken) <= 20 {
		return ctx.Status(fiber.StatusUnauthorized).
			JSON(responsemodel_auth.CommonResponse{
				StatusCode: fiber.StatusUnauthorized,
				Message:    "Error parsing admin token from request",
				Error:      errors.New("admin token missing or invalid"),
			})
	}

	// Set a timeout for the gRPC request to the Auth service
	context, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Send the request to the gRPC service to verify the admin token
	resp, err := m.Client.VerifyAdminToken(context, &pb.RequestVerifyAdmin{
		AdminToken: adminToken,
	})

	if err != nil {
		fmt.Println("=========Auth-Service-Down============")
		return ctx.Status(fiber.StatusUnauthorized).
			JSON(responsemodel_auth.CommonResponse{
				StatusCode: fiber.StatusUnauthorized,
				Message:    "Failed to verify admin token",
				Error:      err.Error(),
			})
	}

	// Check if the token verification service returned an error
	if resp.ErrorMessage != "" {
		return ctx.Status(fiber.StatusBadRequest).
			JSON(responsemodel_auth.CommonResponse{
				StatusCode: fiber.StatusBadRequest,
				Message:    "Failed to verify admin token",
				Error:      resp.ErrorMessage,
			})
	}

	// Set the admin's ID or email in the context for use in subsequent handlers
	ctx.Locals("adminEmail", resp.AdminEmail)

	return ctx.Next() // Continue to the next handler
}

func (m *Middleware) AccessRegenerator(ctx *fiber.Ctx) error {

	accessToken := ctx.Get("x-access-token")
	refreshToken := ctx.Get("x-refresh-token")

	// Check if tokens are present and of valid length
	if accessToken == "" || refreshToken == "" || len(accessToken) < 20 || len(refreshToken) < 20 {
		return ctx.Status(fiber.StatusUnauthorized).
			JSON(responsemodel_auth.CommonResponse{
				StatusCode: fiber.StatusUnauthorized,
				Message:    "error parsing access and refresh tokens from request",
				Error:      "invalid or missing tokens",
			})
	}
	fmt.Println("tokens,,", accessToken, refreshToken)
	// Create a context with timeout for the gRPC call
	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Second*100)
	defer cancel()

	// Call the AccessRegenerator gRPC method
	resp, err := m.Client.AccessRegenerator(ctxTimeout, &pb.RequestAccessGenerator{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})

	// Handle errors from the gRPC call
	if err != nil {
		fmt.Println("----------auth service down--------")

		return ctx.Status(fiber.StatusServiceUnavailable).
			JSON(responsemodel_auth.CommonResponse{
				StatusCode: fiber.StatusServiceUnavailable,
				Message:    "failed to generate access token, auth service unavailable",
				Error:      err.Error(),
			})
	}

	// Check if the response contains an error message
	if resp.ErrorMessage != "" {
		return ctx.Status(fiber.StatusBadRequest).
			JSON(responsemodel_auth.CommonResponse{
				StatusCode: fiber.StatusBadRequest,
				Message:    "failed to generate access token",
				Data:       resp,
				Error:      resp.ErrorMessage,
			})
	}

	// Return success response
	return ctx.Status(fiber.StatusOK).
		JSON(responsemodel_auth.CommonResponse{
			StatusCode: fiber.StatusOK,
			Message:    "new access token generated successfully",
			Data:       resp,
			Error:      nil,
		})
}

func Logger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Start timer
		startTime := time.Now()

		// Process the request
		err := c.Next()

		// Calculate latency
		latency := time.Since(startTime)

		// Log request details
		log.Printf("Status: %d | Latency: %v | ClientIP: %s | Method: %s | Path: %s | UserAgent: %s",
			c.Response().StatusCode(),
			latency,
			c.IP(),
			c.Method(),
			c.Path(),
			c.Get("User-Agent"),
		)

		// Log any errors
		if err != nil {
			log.Printf("Error: %v", err)
		}

		return err
	}
}
