package main

import (
	"fmt"
	"log"
	"net"

	"github.com/ShahabazSulthan/Friendzy_Auth/pkg/config"
	"github.com/ShahabazSulthan/Friendzy_Auth/pkg/di"
	"github.com/ShahabazSulthan/Friendzy_Auth/pkg/pb"
	"google.golang.org/grpc"
)

func main() {
	// Load configuration
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Error loading config:", err)
	}

	// Initialize the Auth Service with Dependency Injection (DI)
	server, err := di.InitializeAuthService(config)
	if err != nil {
		log.Fatal("Error initializing auth service:", err)
	}

	// Ensure port is prefixed with a colon (:)
	port := fmt.Sprintf(":%s", config.PortMngr.PortNo)

	// Start a TCP listener on the specified port
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal("Failed to start TCP listener:", err)
	}
	fmt.Println("Auth Service started on port", config.PortMngr.PortNo)

	// Create a new gRPC server
	grpcServer := grpc.NewServer()

	// Register the AuthService with the gRPC server
	pb.RegisterAuthServiceServer(grpcServer, server)


	// Start serving requests
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal("Failed to start gRPC server:", err)
	}
}
