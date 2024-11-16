package main

import (
	"context"
	"fmt"
	"net"
	"os"

	"github.com/ShahabazSulthan/Friendzy_Auth/pkg/config"
	"github.com/ShahabazSulthan/Friendzy_Auth/pkg/di"
	"github.com/ShahabazSulthan/Friendzy_Auth/pkg/pb"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

func main() {
	// Open or create auth.log file for logging
	logFile, err := os.OpenFile("auth.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Error opening log file: %v\n", err)
		return
	}
	defer logFile.Close()

	// Set up Logrus to output to auth.log
	log := logrus.New()
	log.SetOutput(logFile)
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	// Load configuration
	config, err := config.LoadConfig()
	if err != nil {
		log.WithError(err).Fatal("Error loading config")
	}

	// Initialize the Auth Service with Dependency Injection (DI)
	server, err := di.InitializeAuthService(config)
	if err != nil {
		log.WithError(err).Fatal("Error initializing auth service")
	}

	// Ensure port is prefixed with a colon (:)
	port := fmt.Sprintf(":%s", config.PortMngr.PortNo)

	// Start a TCP listener on the specified port
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.WithError(err).Fatal("Failed to start TCP listener")
	}
	log.WithField("port", config.PortMngr.PortNo).Info("Auth Service started")
	fmt.Println("Auth Service started on port:", config.PortMngr.PortNo)
	// Create a new gRPC server with logging interceptor
	grpcServer := grpc.NewServer(loggingInterceptor(log))

	// Register the AuthService with the gRPC server
	pb.RegisterAuthServiceServer(grpcServer, server)

	// Start serving requests
	if err := grpcServer.Serve(lis); err != nil {
		log.WithError(err).Fatal("Failed to start gRPC server")
	}
}

// Interceptor to log gRPC requests and responses
func loggingInterceptor(log *logrus.Logger) grpc.ServerOption {
	return grpc.UnaryInterceptor(func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// Log incoming request
		log.WithFields(logrus.Fields{
			"method":   info.FullMethod,
			"request":  req,
			"clientIP": getClientIP(ctx),
		}).Info("Received gRPC request")

		// Call the handler
		resp, err := handler(ctx, req)

		// Log response and error if any
		if err != nil {
			log.WithFields(logrus.Fields{
				"method":   info.FullMethod,
				"response": resp,
				"error":    err.Error(),
				"clientIP": getClientIP(ctx),
			}).Error("Error processing request")
		} else {
			log.WithFields(logrus.Fields{
				"method":   info.FullMethod,
				"response": resp,
				"clientIP": getClientIP(ctx),
			}).Info("Successfully processed request")
		}

		return resp, err
	})
}

// Helper function to extract the client IP address
func getClientIP(ctx context.Context) string {
	clientIP := "unknown"
	if p, ok := peer.FromContext(ctx); ok {
		clientIP = p.Addr.String()
	}
	return clientIP
}