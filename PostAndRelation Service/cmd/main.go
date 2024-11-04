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

	// Create a new gRPC server
	grpcServer := grpc.NewServer()

	// Register the PostNrelServiceServer
	pb.RegisterPostNrelServiceServer(grpcServer, postAndRelationServer)

	// Start serving the gRPC server
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to start gRPC server: %v", err)
	}
}
