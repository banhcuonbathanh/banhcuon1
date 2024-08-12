package repository

import (
	"context"

	"english-ai-full/internal/models"

	"github.com/jackc/pgx/v4/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *models.User) error {
	query := `INSERT INTO users (username, email) VALUES ($1, $2) RETURNING id`
	return r.db.QueryRow(ctx, query, user.Username, user.Email).Scan(&user.ID)
}
func (r *UserRepository) GetUserByID(ctx context.Context, id int64) (*models.User, error) {
    query := `SELECT id, username, email FROM users WHERE id = $1`
    user := &models.User{}
    err := r.db.QueryRow(ctx, query, id).Scan(&user.ID, &user.Username, &user.Email)
    if err != nil {
        return nil, err
    }
    return user, nil
}
// Implement other repository methods (GetUser, UpdateUser, DeleteUser, etc.)