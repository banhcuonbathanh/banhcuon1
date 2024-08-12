package service

import (
	"context"

	"english-ai-full/internal/models"
	"english-ai-full/internal/repository"
	pb "english-ai-full/pkg/proto"
)

type UserService struct {
	pb.UnimplementedUserServiceServer
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	user := &models.User{
		Username: req.Username,
		Email:    req.Email,
	}

	if err := s.repo.CreateUser(ctx, user); err != nil {
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


func (s *UserService) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
    user, err := s.repo.GetUserByID(ctx, req.Id)
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
// Implement other service methods (GetUser, UpdateUser, DeleteUser, etc.)