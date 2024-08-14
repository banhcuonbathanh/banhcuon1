package main

import (
	"log"
	"net"
"english-ai-full/ecomm-grpc/config"
"english-ai-full/ecomm-grpc/db"
	"english-ai-full/ecomm-grpc/repository/user_repository"
	"english-ai-full/ecomm-grpc/service/user_service"
	pb "english-ai-full/ecomm-grpc/proto"

	"google.golang.org/grpc"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
// connect 
	dbConn, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbConn.Close()

	if err := db.Migrate(cfg.DatabaseURL); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	userRepo := repository.NewUserRepository(dbConn)
	userService := service.NewUserServer(userRepo)

	lis, err := net.Listen("tcp", cfg.GRPCAddress)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterEcommUserServer(grpcServer, userService)

	log.Printf("Starting gRPC server on %s", cfg.GRPCAddress)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}