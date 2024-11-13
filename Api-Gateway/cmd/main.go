package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	di_auth "github.com/ShahabazSulthan/Friendzy_apiGateway/pkg/Auth_Service/di"
	di_chat "github.com/ShahabazSulthan/Friendzy_apiGateway/pkg/Chat_Call_Service/di"
	config "github.com/ShahabazSulthan/Friendzy_apiGateway/pkg/Config"
	di_notif "github.com/ShahabazSulthan/Friendzy_apiGateway/pkg/Notification_Service/di"
	di_post "github.com/ShahabazSulthan/Friendzy_apiGateway/pkg/post_relation_service/di"
	"github.com/fatih/color"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
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

	// Set log output to file and configure flags
	log.SetOutput(file)
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	// Log server start
	log.Printf("[%s] Server starting on port %s", serverID, cfg.Port)

	// Initialize Fiber app with HTML template engine
	engine := html.New("D:/BROTOTYPE/WEEK 26/Friendzy/Api-Gateway/template", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})
	app.Use(FiberLogger()) // Add custom logging middleware

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

// FiberLogger logs HTTP requests handled by Fiber with bold colors
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

		// Define colors for status codes
		statusColor := color.New(color.Bold)
		switch {
		case statusCode >= 500:
			statusColor.Add(color.FgRed)
		case statusCode >= 400:
			statusColor.Add(color.FgYellow)
		case statusCode >= 300:
			statusColor.Add(color.FgCyan)
		default:
			statusColor.Add(color.FgHiMagenta)
		}

		// Log HTTP request details to both file and terminal with server ID
		logLine := fmt.Sprintf(
			"[%s] [HTTP] status=%d latency=%v clientIP=%s method=%s path=%s userAgent=%s",
			serverID, statusCode, latency, clientIP, method, path, userAgent,
		)
		log.Println(logLine)

		// Output to terminal with colors
		statusColor.Printf("%s\n", logLine)

		// Log errors if any occurred
		if err != nil {
			errorLog := fmt.Sprintf("[%s] [HTTP] error=%s", serverID, err.Error())
			log.Println(errorLog)
			color.New(color.FgRed, color.Bold).Printf("%s\n", errorLog)
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
