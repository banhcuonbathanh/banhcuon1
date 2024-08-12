package main

import (
    "context"
    "log"
    "time"

    pb "english-ai-full/pkg/proto"

    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", 
    grpc.WithTransportCredentials(insecure.NewCredentials()),
    grpc.WithBlock())
    if err != nil {
        log.Fatalf("did not connect: %v", err)
    }
    defer conn.Close()
    c := pb.NewUserServiceClient(conn)

    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()

    // Create User
    r, err := c.CreateUser(ctx, &pb.CreateUserRequest{Username: "testuser", Email: "testuser@example.com"})
    if err != nil {
        log.Fatalf("could not create user: %v", err)
    }
    log.Printf("Created user: %v", r.GetUser())

    // Get User
    user, err := c.GetUser(ctx, &pb.GetUserRequest{Id: r.GetUser().Id})
    if err != nil {
        log.Fatalf("could not get user: %v", err)
    }
    log.Printf("Retrieved user: %v", user.GetUser())

    // Add more test cases here (UpdateUser, DeleteUser, etc.)
}