package repository

import (
	"context"
	"english-ai-full/ecomm-grpc/models"
	"fmt"

	"english-ai-full/util"

	"github.com/jackc/pgx/v4/pgxpool"
)



type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepositoryInterface {
    return &UserRepository{
        db: db,
    }
}

func (us *UserRepository) CreateUser(ctx context.Context, u *models.User) error {
	query := `INSERT INTO users (username, email) VALUES ($1, $2) RETURNING id`
	err := us.db.QueryRow(ctx, query, u.Username, u.Email).Scan(&u.ID)
	if err != nil {
		return fmt.Errorf("error inserting user: %w", err)
	}
	return nil
}

func (us *UserRepository) GetUserByID(ctx context.Context, id int64) (*models.User, error) {
	var u models.User
	query := `SELECT id, username, email FROM users WHERE id = $1`
	err := us.db.QueryRow(ctx, query, id).Scan(&u.ID, &u.Username, &u.Email)
	if err != nil {
		return nil, fmt.Errorf("error getting user: %w", err)
	}
	return &u, nil
}

// Implement the remaining methods of the UserRepositoryInterface

func (us *UserRepository) Save(user models.User) error {
	query := `
		INSERT INTO users (username, email) 
		VALUES ($1, $2) 
		ON CONFLICT (id) 
		DO UPDATE SET username = EXCLUDED.username, email = EXCLUDED.email 
		RETURNING id`
	err := us.db.QueryRow(context.Background(), query, user.Username, user.Email).Scan(&user.ID)
	if err != nil {
		return fmt.Errorf("error saving user: %w", err)
	}
	return nil
}



func (us *UserRepository) Update(user models.User) (models.User, error) {
	query := `
		UPDATE users 
		SET username = $1, email = $2 
		WHERE id = $3
		RETURNING id, username, email`
	
	var updatedUser models.User
	err := us.db.QueryRow(context.Background(), query, user.Username, user.Email, user.ID).Scan(&updatedUser.ID, &updatedUser.Username, &updatedUser.Email)
	if err != nil {
		return models.User{}, fmt.Errorf("error updating user: %w", err)
	}
	return updatedUser, nil
}


func (us *UserRepository) Delete(userID int) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := us.db.Exec(context.Background(), query, userID)
	if err != nil {
		return fmt.Errorf("error deleting user: %w", err)
	}
	return nil
}

func (us *UserRepository) FindByID(userID int) (models.User, error) {
	var user models.User
	query := `SELECT id, username, email FROM users WHERE id = $1`
	err := us.db.QueryRow(context.Background(), query, userID).Scan(&user.ID, &user.Username, &user.Email)
	if err != nil {
		return models.User{}, fmt.Errorf("error finding user by ID: %w", err)
	}
	return user, nil
}

func (us *UserRepository) FindAll() ([]models.User, error) {
	query := `SELECT id, username, email FROM users`
	rows, err := us.db.Query(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("error finding all users: %w", err)
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Username, &user.Email); err != nil {
			return nil, fmt.Errorf("error scanning user: %w", err)
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over users: %w", err)
	}
	return users, nil
}

func (us *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	query := `SELECT id, username, email FROM users WHERE email = $1`
	err := us.db.QueryRow(context.Background(), query, email).Scan(&user.ID, &user.Username, &user.Email)
	if err != nil {
		return nil, fmt.Errorf("error finding user by email: %w", err)
	}
	return &user, nil
}

func (us *UserRepository) FindUsersByPage(pageNumber, pageSize int) ([]models.User, error) {
	offset := (pageNumber - 1) * pageSize
	query := `SELECT id, username, email FROM users LIMIT $1 OFFSET $2`
	rows, err := us.db.Query(context.Background(), query, pageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("error finding users by page: %w", err)
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Username, &user.Email); err != nil {
			return nil, fmt.Errorf("error scanning user: %w", err)
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over users: %w", err)
	}
	return users, nil
}

func (us *UserRepository) Login(ctx context.Context, email, password string) (*models.User, error) {
	var user models.User
	query := `SELECT id, username, email, password FROM users WHERE email = $1`
	err := us.db.QueryRow(ctx, query, email).Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	if err != nil {
		return nil, fmt.Errorf("error finding user by email: %w", err)
	}

	// Check if the password is correct
	if err := util.CheckPassword(password, user.Password); err != nil {
		return nil, fmt.Errorf("invalid password")
	}

	return &user, nil
}


func (us *UserRepository) Register(user models.User) (models.User, error) {
	query := `
		INSERT INTO users (username, email, password) 
		VALUES ($1, $2, $3) 
		RETURNING id`
	err := us.db.QueryRow(context.Background(), query, user.Username, user.Email, user.Password).Scan(&user.ID)
	if err != nil {
		return models.User{}, fmt.Errorf("error registering user: %w", err)
	}
	return user, nil
}
