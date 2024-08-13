package main

import (
	"context"
	"log"
	"testing"
	"time"

	pb "english-ai-full/ecomm-grpc/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var client pb.UserServiceClient

func init() {
    conn, err := grpc.Dial("localhost:50051", 
        grpc.WithTransportCredentials(insecure.NewCredentials()),
        grpc.WithBlock())
    if err != nil {
        log.Fatalf("did not connect: %v", err)
    }
    client = pb.NewUserServiceClient(conn)
}

func TestCreateUser(t *testing.T) {
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()

    r, err := client.CreateUser(ctx, &pb.CreateUserRequest{Username: "testuser", Email: "testuser@example.com"})
    if err != nil {
        t.Fatalf("could not create user: %v", err)
    }
    if r.GetUser().Username != "testuser" || r.GetUser().Email != "testuser@example.com" {
        t.Errorf("unexpected user data: %v", r.GetUser())
    }
    t.Logf("Created user: %v", r.GetUser())
}

func TestGetUser(t *testing.T) {
    // First, create a user
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()

    r, err := client.CreateUser(ctx, &pb.CreateUserRequest{Username: "getuser", Email: "getuser@example.com"})
    if err != nil {
        t.Fatalf("could not create user for get test: %v", err)
    }

    // Now try to get the user
    user, err := client.GetUser(ctx, &pb.GetUserRequest{Id: r.GetUser().Id})
    if err != nil {
        t.Fatalf("could not get user: %v", err)
    }
    if user.GetUser().Username != "getuser" || user.GetUser().Email != "getuser@example.com" {
        t.Errorf("unexpected user data: %v", user.GetUser())
    }
    t.Logf("Retrieved user: %v", user.GetUser())
}

// Add more test functions for other methods (UpdateUser, DeleteUser, etc.)