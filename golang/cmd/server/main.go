package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"

	// "strconv"
	"syscall"
	"time"

	// pb "english-ai-full/ecomm-grpc/proto"
	"english-ai-full/internal/config"

	"github.com/go-chi/chi/v5"
	// "github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/ianschenck/envflag"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	// "english-ai-full/internal/handler"
)

const minSecretKeySize = 32

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	var (
		secretKey = envflag.String("SECRET_KEY", "01234567890123456789012345678901", "secret key for JWT signing")
		svcAddr   = envflag.String("GRPC_SVC_ADDR", cfg.GRPCAddress, "address where the ecomm-grpc service is listening on")
	)
	envflag.Parse()

	if len(*secretKey) < minSecretKeySize {
		log.Fatalf("SECRET_KEY must be at least %d characters", minSecretKeySize)
	}

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	conn, err := grpc.Dial(*svcAddr, opts...)
	if err != nil {
		log.Fatalf("failed to connect to gRPC server: %v", err)
	}
	defer conn.Close()

	// client := pb.NewUserServiceClient(conn)

	// hdl := handler.NewHandler(client, *secretKey)
	// r := handler.RegisterRoutes(hdl)
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Define your HTTP endpoints
	r.Get("/api/hello", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{"message": "Hello from Chi HTTP server!"})
	})

	r.Get("/api/user/{id}", func(w http.ResponseWriter, r *http.Request) {
		// idStr := chi.URLParam(r, "id")
		// id, err := strconv.ParseInt(idStr, 10, 64)
		// if err != nil {
		// 	http.Error(w, "Invalid user ID", http.StatusBadRequest)
		// 	return
		// }

		// // user, err := hdl.GetUser(r.Context(), id)
		// // if err != nil {
		// // 	http.Error(w, "User not found", http.StatusNotFound)
		// // 	return
		// // }

		// json.NewEncoder(w).Encode(user)
	})

	// Add more routes as needed
	// handler.RegisterRoutes(r, hdl)

	server := &http.Server{
		Addr:    cfg.HTTPAddress,
		Handler: r,
	}

	go func() {
		log.Printf("Starting HTTP server on %s", cfg.HTTPAddress)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start Chi server: %v", err)
		}
	}()

	// Graceful shutdown
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