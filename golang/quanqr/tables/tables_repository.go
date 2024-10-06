package tables_test

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/protobuf/types/known/timestamppb"

	"english-ai-full/quanqr/proto_qr/table"
	"english-ai-full/token"
)

type TableRepository struct {
	db *pgxpool.Pool
	jwtMaker *token.JWTMaker
}

func NewTableRepository(db *pgxpool.Pool, secretKey string) *TableRepository {
	return &TableRepository{
		db: db,
		jwtMaker: token.NewJWTMaker(secretKey),
	}
}

func (tr *TableRepository) GetTableList(ctx context.Context) ([]*table.Table, error) {
	query := `
		SELECT number, capacity, status, token, created_at, updated_at
		FROM tables
	`
	rows, err := tr.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error fetching tables: %w", err)
	}
	defer rows.Close()

	var tables []*table.Table
	for rows.Next() {
		var t table.Table
		var createdAt, updatedAt time.Time
		err := rows.Scan(
			&t.Number,
			&t.Capacity,
			&t.Status,
			&t.Token,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning table: %w", err)
		}
		t.CreatedAt = timestamppb.New(createdAt)
		t.UpdatedAt = timestamppb.New(updatedAt)
		tables = append(tables, &t)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over tables: %w", err)
	}

	return tables, nil
}

func (tr *TableRepository) GetTableDetail(ctx context.Context, number int32) (*table.Table, error) {
	query := `
		SELECT number, capacity, status, token, created_at, updated_at
		FROM tables
		WHERE number = $1
	`
	var t table.Table
	var createdAt, updatedAt time.Time
	var statusStr string
	err := tr.db.QueryRow(ctx, query, number).Scan(
		&t.Number,
		&t.Capacity,
		&statusStr,
		&t.Token,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("error fetching table detail: %w", err)
	}
	
	t.Status = table.TableStatus(table.TableStatus_value[statusStr])
	t.CreatedAt = timestamppb.New(createdAt)
	t.UpdatedAt = timestamppb.New(updatedAt)

	return &t, nil
}

func (tr *TableRepository) CreateTable(ctx context.Context, req *table.CreateTableRequest) (*table.Table, error) {
	log.Print("golang/quanqr/tables/tables_repository.go 1 ")
	token, err := tr.generateToken(req.Number)
	if err != nil {
		return nil, fmt.Errorf("error generating token: %w", err)
	}

	// Truncate token if it's longer than 255 characters
	if len(token) > 255 {
		token = token[:255]
	}

	query := `
		INSERT INTO tables (number, capacity, status, token, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $5)
		RETURNING number, capacity, status, token, created_at, updated_at
	`
	var t table.Table
	var createdAt, updatedAt time.Time
	var statusStr string
	err = tr.db.QueryRow(ctx, query,
		req.Number,
		req.Capacity,
		req.Status.String(), // Convert enum to string
		token,
		time.Now(),
	).Scan(
		&t.Number,
		&t.Capacity,
		&statusStr,
		&t.Token,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		log.Printf("Token: %s", t.Token)
		log.Printf("Token pointer: %p", &t.Token)
		return nil, fmt.Errorf("error creating table: %w", err)
	}
	
	// Convert string status back to TableStatus enum
	t.Status = table.TableStatus(table.TableStatus_value[statusStr])
	t.CreatedAt = timestamppb.New(createdAt)
	t.UpdatedAt = timestamppb.New(updatedAt)

	return &t, nil
}

func (tr *TableRepository) UpdateTable(ctx context.Context, req *table.UpdateTableRequest) (*table.Table, error) {

	log.Print("golang/quanqr/tables/tables_repository.go ")
	var newToken string
	var err error
	if req.ChangeToken {
		newToken, err = tr.generateToken(req.Number)
		if err != nil {
			return nil, fmt.Errorf("error generating new token: %w", err)
		}
	}
	log.Print("golang/quanqr/tables/tables_repository.go 111 ")
	query := `
		UPDATE tables
		SET capacity = $2, status = $3, token = CASE WHEN $4 THEN $5 ELSE token END, updated_at = $6
		WHERE number = $1
		RETURNING number, capacity, status, token, created_at, updated_at
	`

	log.Print("golang/quanqr/tables/tables_repository.go 222 ")
	var t table.Table
	var createdAt, updatedAt time.Time
	var statusStr string
	err = tr.db.QueryRow(ctx, query,
		req.Number,
		req.Capacity,
		req.Status.String(),
		req.ChangeToken,
		newToken,
		time.Now(),
	).Scan(
		&t.Number,
		&t.Capacity,
		&statusStr,
		&t.Token,
		&createdAt,
		&updatedAt,
	)


	if err != nil {
		return nil, fmt.Errorf("error updating table: %w", err)
	}

	t.Status = table.TableStatus(table.TableStatus_value[statusStr])
	t.CreatedAt = timestamppb.New(createdAt)
	t.UpdatedAt = timestamppb.New(updatedAt)

	return &t, nil
}
func (tr *TableRepository) DeleteTable(ctx context.Context, number int32) (*table.Table, error) {
	query := `
		DELETE FROM tables
		WHERE number = $1
		RETURNING number, capacity, status, token, created_at, updated_at
	`
	var t table.Table
	var createdAt, updatedAt time.Time
	err := tr.db.QueryRow(ctx, query, number).Scan(
		&t.Number,
		&t.Capacity,
		&t.Status,
		&t.Token,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("error deleting table: %w", err)
	}
	t.CreatedAt = timestamppb.New(createdAt)
	t.UpdatedAt = timestamppb.New(updatedAt)

	return &t, nil
}

// Helper function to generate a token (you need to implement this)
func (tr *TableRepository) generateToken(tableNumber int32) (string, error) {
	// Create a token with the table number as the subject
	tokenString, _, err := tr.jwtMaker.CreateToken(
		int64(tableNumber),
		fmt.Sprintf("table_%d@example.com", tableNumber),
		"table",
		100*365*24*time.Hour, // Token valid for 100 years, adjust as needed
	)
	if err != nil {
		return "", fmt.Errorf("error creating token: %w", err)
	}

	return tokenString, nil
}