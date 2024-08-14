package service

import (
	"context"
	"english-ai-full/ecomm-grpc/data/user"
	"english-ai-full/ecomm-grpc/models"
)

type UserServiceInterface interface {
    CreateUser(userRequest user.CreateUserRequest) error
    Save(user models.User) error
    Update(userRequest user.UpdateUserRequest) (*models.User, error)
    Delete(userID int) error
    // FindByID(userID int) (user.UserResponse, error)
    FindAll() ([]user.UserResponse, error)
    FindByEmail(email string) (*models.User, error)
    FindUsersByPage(pageNumber, pageSize int) ([]user.UserResponse, error)
    
    // Newly added functions:
    Login(ctx context.Context, email, password string) (*models.User, error)
    Register(userRequest user.CreateUserRequest) (bool, error)


}



    // Additional considerations:
    // - Methods for password management (change, reset)
    // - Role-based access control (if applicable)
    // - Email verification and validation
    // - User activation and deactivation
    // - Audit logging