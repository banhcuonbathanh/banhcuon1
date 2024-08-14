package repository

import (
	"context"
	"english-ai-full/ecomm-grpc/models"
)

type UserRepositoryInterface interface {
    CreateUser(ctx context.Context, u *models.User) error
    Save(user models.User) error
    Update(user models.User) (models.User, error)
    Delete(userID int) error
    FindByID(userID int) (models.User, error)
    FindAll() ([]models.User, error)
    FindByEmail(email string) (*models.User, error)
    FindUsersByPage(pageNumber, pageSize int) ([]models.User, error)

    // Newly added functions:
    Login(ctx context.Context, email, password string) (*models.User, error) 
    Register(user models.User) (models.User, error)
}
