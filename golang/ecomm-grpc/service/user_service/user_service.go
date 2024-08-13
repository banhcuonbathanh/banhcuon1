package service

import (
	"context"
	"english-ai-full/ecomm-grpc/models"
	"english-ai-full/ecomm-grpc/repository/user_repository"
	pb "english-ai-full/ecomm-grpc/proto"
)

type Server struct {
	userRepo *repository.UserRepository
	pb.UnimplementedUserServiceServer
}

func NewServer(userRepo *repository.UserRepository) *Server {
	return &Server{
		userRepo: userRepo,
	}
}

func (s *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	user := &models.User{
		Username: req.GetUsername(),
		Email:    req.GetEmail(),
	}

	err := s.userRepo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return &pb.CreateUserResponse{
		User: &pb.User{
			Id:       user.ID,
			Username: user.Username,
			Email:    user.Email,
		},
	}, nil
}

func (s *Server) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	user, err := s.userRepo.GetUserByID(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return &pb.GetUserResponse{
		User: &pb.User{
			Id:       user.ID,
			Username: user.Username,
			Email:    user.Email,
		},
	}, nil
}

// Implement other methods (UpdateUser, DeleteUser, etc.) as needed