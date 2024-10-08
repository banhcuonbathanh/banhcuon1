package set_qr



import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/protobuf/types/known/timestamppb"

	"english-ai-full/quanqr/proto_qr/set"
)

type SetRepository struct {
	db *pgxpool.Pool
}

func NewSetRepository(db *pgxpool.Pool) *SetRepository {
	return &SetRepository{
		db: db,
	}
}

func (sr *SetRepository) GetSetProtoList(ctx context.Context) ([]*set.SetProto, error) {
	query := `
		SELECT id, name, description, user_id, created_at, updated_at 
		FROM sets
	`
	rows, err := sr.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error fetching sets: %w", err)
	}
	defer rows.Close()

	var sets []*set.SetProto
	for rows.Next() {
		var s set.SetProto
		var createdAt, updatedAt time.Time
		err := rows.Scan(
			&s.Id,
			&s.Name,
			&s.Description,
			&s.UserId,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning set: %w", err)
		}
		s.CreatedAt = timestamppb.New(createdAt)
		s.UpdatedAt = timestamppb.New(updatedAt)

		// Fetch dishes for this set
		s.Dishes, err = sr.getSetDishes(ctx, s.Id)
		if err != nil {
			return nil, fmt.Errorf("error fetching dishes for set %d: %w", s.Id, err)
		}

		sets = append(sets, &s)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over sets: %w", err)
	}

	return sets, nil
}

func (sr *SetRepository) GetSetProtoDetail(ctx context.Context, id int32) (*set.SetProto, error) {
	query := `
		SELECT id, name, description, user_id, created_at, updated_at 
		FROM sets 
		WHERE id = $1
	`
	var s set.SetProto
	var createdAt, updatedAt time.Time
	err := sr.db.QueryRow(ctx, query, id).Scan(
		&s.Id,
		&s.Name,
		&s.Description,
		&s.UserId,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("error fetching set detail: %w", err)
	}
	s.CreatedAt = timestamppb.New(createdAt)
	s.UpdatedAt = timestamppb.New(updatedAt)

	// Fetch dishes for this set
	s.Dishes, err = sr.getSetDishes(ctx, s.Id)
	if err != nil {
		return nil, fmt.Errorf("error fetching dishes for set %d: %w", s.Id, err)
	}

	return &s, nil
}

func (sr *SetRepository) CreateSetProto(ctx context.Context, req *set.CreateSetProtoRequest) (*set.SetProto, error) {
	tx, err := sr.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("error starting transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	query := `
		INSERT INTO sets (name, description, user_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $4)
		RETURNING id, created_at, updated_at
	`
	var s set.SetProto
	var createdAt, updatedAt time.Time
	err = tx.QueryRow(ctx, query,
		req.Name,
		req.Description,
		req.UserId,
		time.Now(),
	).Scan(&s.Id, &createdAt, &updatedAt)
	if err != nil {
		return nil, fmt.Errorf("error creating set: %w", err)
	}

	s.Name = req.Name
	s.Description = req.Description
	s.UserId = req.UserId
	s.CreatedAt = timestamppb.New(createdAt)
	s.UpdatedAt = timestamppb.New(updatedAt)

	// Insert set dishes
	for _, dish := range req.Dishes {
		_, err := tx.Exec(ctx, "INSERT INTO set_dishes (set_id, dish_id, quantity) VALUES ($1, $2, $3)",
			s.Id, dish.Dish.Id, dish.Quantity)
		if err != nil {
			return nil, fmt.Errorf("error inserting set dish: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("error committing transaction: %w", err)
	}

	s.Dishes = req.Dishes
	return &s, nil
}

func (sr *SetRepository) UpdateSetProto(ctx context.Context, req *set.UpdateSetProtoRequest) (*set.SetProto, error) {
	tx, err := sr.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("error starting transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	query := `
		UPDATE sets
		SET name = $2, description = $3, updated_at = $4
		WHERE id = $1
		RETURNING user_id, created_at, updated_at
	`
	var s set.SetProto
	var createdAt, updatedAt time.Time
	err = tx.QueryRow(ctx, query,
		req.Id,
		req.Name,
		req.Description,
		time.Now(),
	).Scan(&s.UserId, &createdAt, &updatedAt)
	if err != nil {
		return nil, fmt.Errorf("error updating set: %w", err)
	}

	s.Id = req.Id
	s.Name = req.Name
	s.Description = req.Description
	s.CreatedAt = timestamppb.New(createdAt)
	s.UpdatedAt = timestamppb.New(updatedAt)

	// Delete existing set dishes
	_, err = tx.Exec(ctx, "DELETE FROM set_dishes WHERE set_id = $1", req.Id)
	if err != nil {
		return nil, fmt.Errorf("error deleting existing set dishes: %w", err)
	}

	// Insert updated set dishes
	for _, dish := range req.Dishes {
		_, err := tx.Exec(ctx, "INSERT INTO set_dishes (set_id, dish_id, quantity) VALUES ($1, $2, $3)",
			s.Id, dish.Dish.Id, dish.Quantity)
		if err != nil {
			return nil, fmt.Errorf("error inserting updated set dish: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("error committing transaction: %w", err)
	}

	s.Dishes = req.Dishes
	return &s, nil
}

func (sr *SetRepository) DeleteSetProto(ctx context.Context, id int32) (*set.SetProto, error) {
	tx, err := sr.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("error starting transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Delete set dishes first
	_, err = tx.Exec(ctx, "DELETE FROM set_dishes WHERE set_id = $1", id)
	if err != nil {
		return nil, fmt.Errorf("error deleting set dishes: %w", err)
	}

	query := `
		DELETE FROM sets
		WHERE id = $1
		RETURNING id, name, description, user_id, created_at, updated_at
	`
	var s set.SetProto
	var createdAt, updatedAt time.Time
	err = tx.QueryRow(ctx, query, id).Scan(
		&s.Id,
		&s.Name,
		&s.Description,
		&s.UserId,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("error deleting set: %w", err)
	}
	s.CreatedAt = timestamppb.New(createdAt)
	s.UpdatedAt = timestamppb.New(updatedAt)

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("error committing transaction: %w", err)
	}

	return &s, nil
}

func (sr *SetRepository) getSetDishes(ctx context.Context, setID int32) ([]*set.SetProtoDish, error) {
	query := `
		SELECT d.id, d.name, d.price, d.description, d.image, d.status, d.created_at, d.updated_at, sd.quantity
		FROM set_dishes sd
		JOIN dishes d ON sd.dish_id = d.id
		WHERE sd.set_id = $1
	`
	rows, err := sr.db.Query(ctx, query, setID)
	if err != nil {
		return nil, fmt.Errorf("error fetching set dishes: %w", err)
	}
	defer rows.Close()

	var dishes []*set.SetProtoDish
	for rows.Next() {
		var sd set.SetProtoDish
		var d set.Dish
		var createdAt, updatedAt time.Time
		err := rows.Scan(
			&d.Id,
			&d.Name,
			&d.Price,
			&d.Description,
			&d.Image,
			&d.Status,
			&createdAt,
			&updatedAt,
			&sd.Quantity,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning set dish: %w", err)
		}
		d.CreatedAt = timestamppb.New(createdAt)
		d.UpdatedAt = timestamppb.New(updatedAt)
		sd.Dish = &d
		dishes = append(dishes, &sd)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over set dishes: %w", err)
	}

	return dishes, nil
}