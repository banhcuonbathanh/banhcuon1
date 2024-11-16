package repository

import (
	"context"
	"database/sql"
	"strconv"

	logg "english-ai-full/logger"
	"fmt"



	"github.com/jackc/pgx/v4/pgxpool"

	"english-ai-full/ecomm-api/types"

	 "english-ai-full/util"
)



type UserRepository struct {
	db     *pgxpool.Pool
	logger *logg.Logger
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		db:     db,
		logger: logg.NewLogger(),
	}
}

func (us *UserRepository) CreateUser(ctx context.Context, u *types.UserReqModel) (*types.UserReqModel, error) {
	formattedName := util.GenerateFormattedName(u.Name)
	us.logger.Info("Inserting new user into database: " + formattedName + ", " + u.Email)

	query := `INSERT INTO users (name, email, password, role, phone, image, address, created_at, updated_at) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) 
			  RETURNING id`

	err := us.db.QueryRow(ctx, query, formattedName, u.Email, u.Password, u.Role, u.Phone, u.Image, u.Address, u.CreatedAt, u.UpdatedAt).Scan(&u.ID)
	if err != nil {
		us.logger.Error("Error inserting user: " + err.Error())
		return nil, fmt.Errorf("error inserting user: %w", err)
	}

	u.Name = formattedName // Update the name in the model
	
	us.logger.Info("User inserted successfully. ID: " + fmt.Sprint(u.ID) +
		", Name: " + u.Name +
		", Email: " + u.Email +
		", Role: " + u.Role +
		", Phone: " + u.Phone +
		", Image: " + u.Image +
		", Address: " + u.Address +
		", CreatedAt: " + u.CreatedAt.String() +
		", UpdatedAt: " + u.UpdatedAt.String())
	return u, nil
}

func (us *UserRepository) FindAll() ([]types.UserReqModel, error) {
	query := `
		SELECT id, name, email, password, role, phone, image, address, created_at, updated_at 
		FROM users
	`
	rows, err := us.db.Query(context.Background(), query)
	if err != nil {
		us.logger.Error("Error finding all users: " + err.Error())
		return nil, fmt.Errorf("error finding all users: %w", err)
	}
	defer rows.Close()

	var users []types.UserReqModel
	for rows.Next() {
		var user types.UserReqModel
		var phoneNull sql.NullInt64
		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Password,
			&user.Role,
			&phoneNull,
			&user.Image,
			&user.Address,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			us.logger.Error("Error scanning user: " + err.Error())
			return nil, fmt.Errorf("error scanning user: %w", err)
		}
		
		// Check if phone is valid and assign it to user.Phone
		if phoneNull.Valid {
			user.Phone = strconv.FormatInt(phoneNull.Int64, 10)
		}
		
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		us.logger.Error("Error iterating over users: " + err.Error())
		return nil, fmt.Errorf("error iterating over users: %w", err)
	}

	us.logger.Info("Found " + fmt.Sprintf("%d", len(users)) + " users")
	return users, nil
}

func (us *UserRepository) FindByEmail(email string) (*types.UserReqModel, error) {
	var user types.UserReqModel
	query := `SELECT id, name, email, password, role, phone, image, address, created_at, updated_at FROM users WHERE email = $1`
	err := us.db.QueryRow(context.Background(), query, email).Scan(
		&user.ID, &user.Name, &user.Email, &user.Password, &user.Role, &user.Phone, &user.Image, &user.Address, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		us.logger.Error("Error finding user by email: " + err.Error())
		return nil, fmt.Errorf("error finding user by email: %w", err)
	}

	us.logger.Info("User found by email: " + email)
	return &user, nil
}

func (us *UserRepository) Login(ctx context.Context, email, password string) (*types.UserReqModel, error) {
	var storedPassword string
	var userModel types.UserReqModel

	query := `SELECT id, name, email, password, role, phone, image, address, created_at, updated_at FROM users WHERE email = $1`
	err := us.db.QueryRow(ctx, query, email).Scan(
		&userModel.ID,
		&userModel.Name,
		&userModel.Email,
		&storedPassword,
		&userModel.Role,
		&userModel.Phone,
		&userModel.Image,
		&userModel.Address,
		&userModel.CreatedAt,
		&userModel.UpdatedAt,
	)
	if err != nil {
		us.logger.Error("Error finding user by email during login: " + err.Error())
		return nil, fmt.Errorf("error finding user by email: %w", err)
	}

	if err := util.CheckPassword(password, storedPassword); err != nil {
		us.logger.Warning("Invalid password attempt for user: " + email)
		return nil, fmt.Errorf("invalid password")
	}

	us.logger.Info("User logged in successfully: " + email)
	return &userModel, nil
}

func (us *UserRepository) Save(user types.UserReqModel) error {
	query := `
		INSERT INTO users (name, email) 
		VALUES ($1, $2) 
		ON CONFLICT (id) 
		DO UPDATE SET name = EXCLUDED.name, email = EXCLUDED.email 
		RETURNING id`
	err := us.db.QueryRow(context.Background(), query, user.Name, user.Email).Scan(&user.ID)
	if err != nil {
		us.logger.Error("Error saving user: " + err.Error())
		return fmt.Errorf("error saving user: %w", err)
	}
	us.logger.Info("User saved successfully. ID: " )
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
		us.logger.Error("Error updating user: " + err.Error())
		return types.UserReqModel{}, fmt.Errorf("error updating user: %w", err)
	}
	us.logger.Info("User updated successfully. ID: " )
	return updatedUser, nil
}

func (us *UserRepository) Delete(userID int) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := us.db.Exec(context.Background(), query, userID)
	if err != nil {
		us.logger.Error("Error deleting user: " + err.Error())
		return fmt.Errorf("error deleting user: %w", err)
	}
	us.logger.Info("User deleted successfully. ID: " + fmt.Sprintf("%d", userID))
	return nil
}

func (us *UserRepository) FindUsersByPage(pageNumber, pageSize int) ([]types.UserReqModel, error) {
	offset := (pageNumber - 1) * pageSize
	query := `SELECT id, name, email FROM users LIMIT $1 OFFSET $2`
	rows, err := us.db.Query(context.Background(), query, pageSize, offset)
	if err != nil {
		us.logger.Error("Error finding users by page: " + err.Error())
		return nil, fmt.Errorf("error finding users by page: %w", err)
	}
	defer rows.Close()

	var users []types.UserReqModel
	for rows.Next() {
		var user types.UserReqModel
		if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
			us.logger.Error("Error scanning user: " + err.Error())
			return nil, fmt.Errorf("error scanning user: %w", err)
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		us.logger.Error("Error iterating over users: " + err.Error())
		return nil, fmt.Errorf("error iterating over users: %w", err)
	}
	us.logger.Info("Found " + fmt.Sprintf("%d", len(users)) + " users on page " + fmt.Sprintf("%d", pageNumber))
	return users, nil
}

func (us *UserRepository) CreateSession(ctx context.Context, s *types.Session) (*types.Session, error) {
	query := `INSERT INTO sessions (id, user_email, refresh_token, is_revoked, expires_at) 
			  VALUES ($1, $2, $3, $4, $5) 
			  RETURNING id`

	err := us.db.QueryRow(ctx, query, s.ID, s.UserEmail, s.RefreshToken, s.IsRevoked, s.ExpiresAt).Scan(&s.ID)
	if err != nil {
		us.logger.Error("Error inserting session: " + err.Error())
		return nil, fmt.Errorf("error inserting session: %w", err)
	}
	us.logger.Info("Session created successfully for user: " + s.UserEmail)
	return s, nil
}

func (us *UserRepository) GetSession(ctx context.Context, id string) (*types.Session, error) {
	var s types.Session
	query := "SELECT * FROM sessions WHERE id=$1"
	err := us.db.QueryRow(ctx, query, id).Scan(&s.ID, &s.UserEmail, &s.RefreshToken, &s.IsRevoked, &s.ExpiresAt)
	if err != nil {
		us.logger.Error("Error getting session: " + err.Error())
		return nil, fmt.Errorf("error getting session: %w", err)
	}
	us.logger.Info("Session retrieved successfully. ID: " + id)
	return &s, nil
}

func (us *UserRepository) RevokeSession(ctx context.Context, id string) error {
	query := "UPDATE sessions SET is_revoked=true WHERE id=$1"
	_, err := us.db.Exec(ctx, query, id)
	if err != nil {
		us.logger.Error("Error revoking session: " + err.Error())
		return fmt.Errorf("error revoking session: %w", err)
	}
	us.logger.Info("Session revoked successfully. ID: " + id)
	return nil
}

func (us *UserRepository) DeleteSession(ctx context.Context, id string) error {
	query := "DELETE FROM sessions WHERE id=$1"
	_, err := us.db.Exec(ctx, query, id)
	if err != nil {
		us.logger.Error("Error deleting session: " + err.Error())
		return fmt.Errorf("error deleting session: %w", err)
	}
	us.logger.Info("Session deleted successfully. ID: " + id)
	return nil
}