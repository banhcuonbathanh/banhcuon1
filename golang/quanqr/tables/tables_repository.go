package tables_test
import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/protobuf/types/known/timestamppb"

	"english-ai-full/quanqr/proto_qr/table"
)

type TableRepository struct {
	db *pgxpool.Pool
}

func NewTableRepository(db *pgxpool.Pool) *TableRepository {
	return &TableRepository{
		db: db,
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
	err := tr.db.QueryRow(ctx, query, number).Scan(
		&t.Number,
		&t.Capacity,
		&t.Status,
		&t.Token,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("error fetching table detail: %w", err)
	}
	t.CreatedAt = timestamppb.New(createdAt)
	t.UpdatedAt = timestamppb.New(updatedAt)

	return &t, nil
}

func (tr *TableRepository) CreateTable(ctx context.Context, req *table.CreateTableRequest) (*table.Table, error) {
	query := `
		INSERT INTO tables (number, capacity, status, token, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $5)
		RETURNING number, capacity, status, token, created_at, updated_at
	`
	var t table.Table
	var createdAt, updatedAt time.Time
	err := tr.db.QueryRow(ctx, query,
		req.Number,
		req.Capacity,
		req.Status,
		generateToken(), // You need to implement this function
		time.Now(),
	).Scan(
		&t.Number,
		&t.Capacity,
		&t.Status,
		&t.Token,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("error creating table: %w", err)
	}
	t.CreatedAt = timestamppb.New(createdAt)
	t.UpdatedAt = timestamppb.New(updatedAt)

	return &t, nil
}

func (tr *TableRepository) UpdateTable(ctx context.Context, req *table.UpdateTableRequest) (*table.Table, error) {
	query := `
		UPDATE tables
		SET capacity = $2, status = $3, token = CASE WHEN $4 THEN $5 ELSE token END, updated_at = $6
		WHERE number = $1
		RETURNING number, capacity, status, token, created_at, updated_at
	`
	var t table.Table
	var createdAt, updatedAt time.Time
	err := tr.db.QueryRow(ctx, query,
		req.Number,
		req.Capacity,
		req.Status,
		req.ChangeToken,
		generateToken(), // You need to implement this function
		time.Now(),
	).Scan(
		&t.Number,
		&t.Capacity,
		&t.Status,
		&t.Token,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("error updating table: %w", err)
	}
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
func generateToken() string {
	// Implement token generation logic here
	return "generated-token"
}