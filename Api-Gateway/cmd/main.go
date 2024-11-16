package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	di_auth "github.com/ShahabazSulthan/Friendzy_apiGateway/pkg/Auth_Service/di"
	di_chat "github.com/ShahabazSulthan/Friendzy_apiGateway/pkg/Chat_Call_Service/di"
	config "github.com/ShahabazSulthan/Friendzy_apiGateway/pkg/Config"
	di_notif "github.com/ShahabazSulthan/Friendzy_apiGateway/pkg/Notification_Service/di"
	di_post "github.com/ShahabazSulthan/Friendzy_apiGateway/pkg/post_relation_service/di"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/template/html/v2"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/peer"
)

const serverID = "SERVER-8000"

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("[%s] Failed to load config: %v", serverID, err)
	}

	// Open log file
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("[%s] Failed to open log file: %v", serverID, err)
	}
	defer file.Close()

	// Set up Logrus to write to both file and console without color
	logrus.SetOutput(io.MultiWriter(file, os.Stdout))

	// Configure log formatter without colors
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})

	// Log server start
	logrus.WithField("serverID", serverID).Infof("Server starting on port %s", cfg.Port)

	// Initialize Fiber app with HTML template engine
	engine := html.New("D:/BROTOTYPE/WEEK 26/Friendzy/Api-Gateway/template", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// Set up Prometheus middleware using monitor
	app.Get("/metrics", monitor.New(monitor.Config{Title: "Friendzy API Gateway Metrics"}))

	// Use logging middleware
	app.Use(FiberLogger())

	// Initialize services with dependency injection
	middleware, err := di_auth.InitAuthClient(app, cfg)
	if err != nil {
		log.Fatal(err)
	}

	if err := di_post.InitPostNrelClient(app, cfg, middleware); err != nil {
		log.Fatal(err)
	}

	if err := di_notif.InitNotificationClient(app, cfg, middleware); err != nil {
		log.Fatal(err)
	}

	if err := di_chat.InitChatNcallClient(app, cfg, middleware); err != nil {
		log.Fatal(err)
	}

	// Log running port to the console
	fmt.Printf("Server running on port %s\n", cfg.Port)

	// Start the Fiber app
	if err := app.Listen(cfg.Port); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}

func FiberLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		startTime := time.Now()

		// Process the HTTP request
		err := c.Next()

		// Calculate request details
		latency := time.Since(startTime)
		method := c.Method()
		statusCode := c.Response().StatusCode()
		clientIP := c.IP()
		userAgent := c.Get("User-Agent")
		path := c.Path()
		if path == "" {
			path = "/"
		}

		// Set up Logrus fields
		logger := logrus.WithFields(logrus.Fields{
			"status":    statusCode,
			"latency":   latency,
			"clientIP":  clientIP,
			"method":    method,
			"path":      path,
			"userAgent": userAgent,
			"serverID":  serverID,
		})

		// Log based on status code
		switch {
		case statusCode >= 500:
			logger.Error("Internal Server Error")
		case statusCode >= 400:
			logger.Warn("Client Error")
		case statusCode >= 300:
			logger.Info("Redirection")
		default:
			logger.Info("Successful Request")
		}

		// Log errors if any occurred
		if err != nil {
			logger.WithField("error", err.Error()).Error("An error occurred")
		}

		return err
	}
}

// GetClientIP retrieves the client IP address from the context
func GetClientIP(ctx context.Context) string {
	clientIP := "unknown"
	if p, ok := peer.FromContext(ctx); ok {
		clientIP = p.Addr.String()
	}
	return clientIP
}
