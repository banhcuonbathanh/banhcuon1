package repository

import (
	"context"
	"database/sql"
	"log"

	"fmt"

	"english-ai-full/util"

	"github.com/jackc/pgx/v4/pgxpool"

	"english-ai-full/ecomm-api/types"
)



type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
    return &UserRepository{
        db: db,
    }
}
// err := us.userRepo.CreateUser(ctx, newUser)
func (us *UserRepository) CreateUser(ctx context.Context, u *types.UserReqModel) (*types.UserReqModel, error) {
	log.Println("Inserting new user into database:", u.Name, u.Email)

	query := `INSERT INTO users (name, email, password, is_admin, phone, image, address, created_at, updated_at) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) 
			  RETURNING id`

	err := us.db.QueryRow(ctx, query, u.Name, u.Email, u.Password, u.IsAdmin, u.Phone, u.Image, u.Address, u.CreatedAt, u.UpdatedAt).Scan(&u.ID)
	if err != nil {
		log.Println("Error inserting user:", err)
		return nil, fmt.Errorf("error inserting user: %w", err)
	}

	
	log.Println("User inserted successfully. ID:", u.ID, 
		"Name:", u.Name, 
		"Email:", u.Email, 
		"Password:", u.Password, 
		"IsAdmin:", u.IsAdmin, 
		"Phone:", u.Phone, 
		"Image:", u.Image, 
		"Address:", u.Address, 
		"CreatedAt:", u.CreatedAt, 
		"UpdatedAt:", u.UpdatedAt)
	return u, nil
}

func (us *UserRepository) GetUserByID(ctx context.Context, id int64) (*types.UserReqModel, error) {
	var u types.UserReqModel
	query := `SELECT id, name, email FROM users WHERE id = $1`
	err := us.db.QueryRow(ctx, query, id).Scan(&u.ID, &u.Name, &u.Email)
	if err != nil {
		return nil, fmt.Errorf("error getting user: %w", err)
	}
	return &u, nil
}

// Implement the remaining methods of the UserRepositoryInterface

func (us *UserRepository) Save(user types.UserReqModel) error {
	query := `
		INSERT INTO users (name, email) 
		VALUES ($1, $2) 
		ON CONFLICT (id) 
		DO UPDATE SET name = EXCLUDED.name, email = EXCLUDED.email 
		RETURNING id`
	err := us.db.QueryRow(context.Background(), query, user.Name, user.Email).Scan(&user.ID)
	if err != nil {
		return fmt.Errorf("error saving user: %w", err)
	}
	return nil
}



func (us *UserRepository) Update(user types.UserReqModel) (types.UserReqModel, error) {
	query := `
		UPDATE users 
		SET name = $1, email = $2 
		WHERE id = $3
		RETURNING id, name, email`
	
	var updatedUser types.UserReqModel
	err := us.db.QueryRow(context.Background(), query, user.Name, user.Email, user.ID).Scan(&updatedUser.ID, &updatedUser.Name, &updatedUser.Email)
	if err != nil {
		return types.UserReqModel{}, fmt.Errorf("error updating user: %w", err)
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

func (us *UserRepository) FindByID(userID int) (types.UserReqModel, error) {
	var user types.UserReqModel
	query := `SELECT id, name, email FROM users WHERE id = $1`
	err := us.db.QueryRow(context.Background(), query, userID).Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		return types.UserReqModel{}, fmt.Errorf("error finding user by ID: %w", err)
	}
	return user, nil
}

func (us *UserRepository) FindAll() ([]types.UserReqModel, error) {
    query := `
        SELECT id, name, email, password, is_admin, phone, image, address, created_at, updated_at 
        FROM users
    `
    rows, err := us.db.Query(context.Background(), query)
    if err != nil {
        return nil, fmt.Errorf("error finding all users: %w", err)
    }
    defer rows.Close()

    var users []types.UserReqModel
    for rows.Next() {
        var user types.UserReqModel
        var phoneNull sql.NullInt64  // Use sql.NullInt64 to handle NULL values
        err := rows.Scan(
            &user.ID,
            &user.Name,
            &user.Email,
            &user.Password,
            &user.IsAdmin,
            &phoneNull,  // Scan into phoneNull
            &user.Image,
            &user.Address,
            &user.CreatedAt,
            &user.UpdatedAt,
        )
        if err != nil {
            return nil, fmt.Errorf("error scanning user: %w", err)
        }
        
        // Handle the NULL case for phone
        if phoneNull.Valid {
            user.Phone = &phoneNull.Int64
        } else {
            user.Phone = nil
        }
        
        users = append(users, user)
    }
    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("error iterating over users: %w", err)
    }

    // Log the details of the last user
    // if len(users) > 0 {
    //     lastUser := users[len(users)-1]
    //     log.Println("User find all repository", lastUser.ID, 
    //         "Name:", lastUser.Name, 
    //         "Email:", lastUser.Email, 
    //         "Password:", lastUser.Password, 
    //         "IsAdmin:", lastUser.IsAdmin, 
    //         "Phone:", lastUser.Phone, 
    //         "Image:", lastUser.Image, 
    //         "Address:", lastUser.Address, 
    //         "CreatedAt:", lastUser.CreatedAt, 
    //         "UpdatedAt:", lastUser.UpdatedAt)
    // }

    return users, nil
}


func (us *UserRepository) FindByEmail(email string) (*types.UserReqModel, error) {
	var user types.UserReqModel
	query := `SELECT id, name, email FROM users WHERE email = $1`
	err := us.db.QueryRow(context.Background(), query, email).Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		return nil, fmt.Errorf("error finding user by email: %w", err)
	}
	return &user, nil
}

func (us *UserRepository) FindUsersByPage(pageNumber, pageSize int) ([]types.UserReqModel, error) {
	offset := (pageNumber - 1) * pageSize
	query := `SELECT id, name, email FROM users LIMIT $1 OFFSET $2`
	rows, err := us.db.Query(context.Background(), query, pageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("error finding users by page: %w", err)
	}
	defer rows.Close()

	var users []types.UserReqModel
	for rows.Next() {
		var user types.UserReqModel
		if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
			return nil, fmt.Errorf("error scanning user: %w", err)
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over users: %w", err)
	}
	return users, nil
}

func (us *UserRepository) Login(ctx context.Context, email, password string) (*types.UserReqModel, error) {
	var storedPassword string
	var userModel types.UserReqModel

	// Fetch user data based on email
	query := `SELECT id, name, email, password, is_admin, phone, image, address, created_at, updated_at FROM users WHERE email = $1`
	err := us.db.QueryRow(ctx, query, email).Scan(
		&userModel.ID,
		&userModel.Name,
		&userModel.Email,
		&storedPassword,
		&userModel.IsAdmin,
		&userModel.Phone,
		&userModel.Image,
		&userModel.Address,
		&userModel.CreatedAt,
		&userModel.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("error finding user by email: %w", err)
	}

	// Check if the password is correct
	if err := util.CheckPassword(password, storedPassword); err != nil {
		return nil, fmt.Errorf("invalid password")
	}

	return &userModel, nil
}


