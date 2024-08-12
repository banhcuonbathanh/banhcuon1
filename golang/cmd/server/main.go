package main

import (
	"context"
	"encoding/json"
	"english-ai-full/internal/config"
	"english-ai-full/internal/db"
	"english-ai-full/internal/repository"
	"english-ai-full/internal/service"
	pb "english-ai-full/pkg/proto"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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

	// Start gRPC server
	go startGRPCServer(cfg.GRPCAddress, userService)

	// Start HTTP server with Chi
	go startChiServer(cfg.HTTPAddress, userService)

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down servers...")
}

func startGRPCServer(address string, userService pb.UserServiceServer) {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, userService)

	log.Printf("Starting gRPC server on %s", address)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}

func startChiServer(address string, userService *service.UserService) {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Define your HTTP endpoints
	r.Get("/api/hello", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{"message": "Hello from Chi HTTP server!"})
	})

	// Example of how to use the userService in a Chi handler
	r.Get("/api/user/{id}", func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}
	
		req := &pb.GetUserRequest{Id: id}
		user, err := userService.GetUser(r.Context(), req)
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
	
		json.NewEncoder(w).Encode(user)
	})
	

	// Add more routes as needed

	server := &http.Server{
		Addr:    address,
		Handler: r,
	}

	go func() {
		log.Printf("Starting HTTP server on %s", address)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start Chi server: %v", err)
		}
	}()

	// Graceful shutdown for Chi server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("HTTP server stopped")
}