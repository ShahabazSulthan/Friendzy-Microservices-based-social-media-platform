package main

import (
	"fmt"
	"log"
	"net"

	"github.com/ShahabazSulthan/friendzy_post/pkg/config"
	"github.com/ShahabazSulthan/friendzy_post/pkg/di"
	"github.com/ShahabazSulthan/friendzy_post/pkg/pb"
	"google.golang.org/grpc"
)

func main() {
	// Load the configuration
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize the ChatNCall service
	server, err := di.InitializeChatNCallSvc(config)
	if err != nil {
		log.Fatalf("Failed to initialize ChatNCall service: %v", err)
	}

	// Start listening on the specified port
	listener, err := net.Listen("tcp", config.PortMngr.RunnerPort)
	if err != nil {
		log.Fatalf("Failed to start listener on port %s: %v", config.PortMngr.RunnerPort, err)
	}
	fmt.Printf("ChatNCall Service started on: %s\n", config.PortMngr.RunnerPort)

	// Create a new gRPC server
	grpcServer := grpc.NewServer()
	pb.RegisterChatNCallServiceServer(grpcServer, server)

	// Log incoming connections
	go logIncomingConnections(listener)

	// Serve the gRPC server
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to start ChatNCall gRPC server: %v", err)
	}
}

// logIncomingConnections logs each new connection to the server.
func logIncomingConnections(listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}
		log.Printf("New connection from: %s", conn.RemoteAddr())
		conn.Close() // Close the connection after logging
	}
}
