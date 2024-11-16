package main

import (
	"context"
	"fmt"
	"net"
	"os"

	"github.com/ShahabazSulthan/friendzy_post/pkg/config"
	"github.com/ShahabazSulthan/friendzy_post/pkg/di"
	"github.com/ShahabazSulthan/friendzy_post/pkg/pb"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

func main() {

	logFile, err := os.OpenFile("post.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize the post and relation server using DI
	postAndRelationServer, err := di.InitializePostAndRelationServer(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize Post and Relation server: %v", err)
	}

	// Listen on the specified port from the configuration
	lis, err := net.Listen("tcp", cfg.PortMngr.RunnerPort)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", cfg.PortMngr.RunnerPort, err)
	}
	defer lis.Close()

	fmt.Println("Post and Relation Service started on port:", cfg.PortMngr.RunnerPort)
	log.WithField("port", cfg.PortMngr.RunnerPort).Info("Post and Relation Service started")

	// Create a new gRPC server
	grpcServer := grpc.NewServer(loggingInterceptor(log))

	// Register the PostNrelServiceServer
	pb.RegisterPostNrelServiceServer(grpcServer, postAndRelationServer)

	// Start serving the gRPC server
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to start gRPC server: %v", err)
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
