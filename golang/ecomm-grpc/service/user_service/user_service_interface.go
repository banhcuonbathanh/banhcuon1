package service

import (
	// ... other imports
	"context"
	"english-ai-full/ecomm-grpc/data/user"

	"english-ai-full/ecomm-grpc/models"
)

type UserServiceInterface interface {
    Create(userRequest user.CreateUserRequest)

    Delete(userID int)
    FindByID(userID int) (user.UserResponse, error)
    FindAll() ([]user.UserResponse, error)
    FindByEmail(email string) (*models.User, error)
    FindUsersByPage(pageNumber, pageSize int) ([]user.UserResponse, error)
    Update(userRequest user.UpdateUserRequest) (*models.User, error)
    // Newly added functions:
    Login(ctx context.Context, email, password string) (*models.User, error)
    Register(userRequest user.CreateUserRequest) (bool, error)

       // Additional considerations:
    // - Methods for password management (change, reset)
    // - Role-based access control (if applicable)
    // - Email verification and validation
    // - User activation and deactivation
    // - Audit logging
}
