package main

import (

	"english-ai-full/internal/config"
	"english-ai-full/internal/db"
	"english-ai-full/internal/repository"
	"english-ai-full/internal/service"
	pb "english-ai-full/pkg/proto"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	dbConn, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbConn.Close()

	if err := db.Migrate(cfg.DatabaseURL); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	userRepo := repository.NewUserRepository(dbConn)
	userService := service.NewUserService(userRepo)

	lis, err := net.Listen("tcp", cfg.GRPCAddress)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, userService)

	go func() {
		log.Printf("Starting gRPC server on %s", cfg.GRPCAddress)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down gRPC server...")
	grpcServer.GracefulStop()
	log.Println("Server stopped")
}