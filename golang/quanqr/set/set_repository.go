package set_qr

import (
	"context"
	"english-ai-full/logger"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/protobuf/types/known/timestamppb"

	"english-ai-full/quanqr/proto_qr/set"
)





type SetRepository struct {
    db     *pgxpool.Pool
    logger *logger.Logger
}

func NewSetRepository(db *pgxpool.Pool) *SetRepository {
    return &SetRepository{
        db:     db,
        logger: logger.NewLogger(),
    }
}

func (sr *SetRepository) GetSetProtoList(ctx context.Context) ([]*set.SetProto, error) {
    sr.logger.Info("Fetching set list GetSetProtoList golang/quanqr/set/set_repository.go")
    query := `
        SELECT id, name, description, user_id, created_at, updated_at, is_favourite, like_by, is_public
        FROM sets
    `
    rows, err := sr.db.Query(ctx, query)
    if err != nil {
        sr.logger.Error("Error fetching sets: " + err.Error())
        return nil, fmt.Errorf("error fetching sets: %w", err)
    }
    defer rows.Close()
    
    var sets []*set.SetProto
    for rows.Next() {
        var s set.SetProto
        var createdAt, updatedAt time.Time
        var userID *int32
        var likeBy []int64
        
        err := rows.Scan(
            &s.Id,
            &s.Name,
            &s.Description,
            &userID,
            &createdAt,
            &updatedAt,
            &s.IsFavourite,
            &likeBy,
            &s.IsPublic,
        )
        if err != nil {
            sr.logger.Error("Error scanning set: " + err.Error())
            return nil, fmt.Errorf("error scanning set: %w", err)
        }
        
        s.CreatedAt = timestamppb.New(createdAt)
        s.UpdatedAt = timestamppb.New(updatedAt)
        s.UserId = userID
        s.LikeBy = likeBy
        
        // Fetch dishes for the current set
        dishes, err := sr.GetDishesForSet(ctx, s.Id)
        sr.logger.Info(fmt.Sprintf("golang/quanqr/set/set_repository.go dishes: %+v", dishes))
        if err != nil {
            sr.logger.Error(fmt.Sprintf("Error fetching dishes for set %d: %s", s.Id, err.Error()))
            return nil, fmt.Errorf("error fetching dishes for set %d: %w", s.Id, err)
        }
        
        s.Dishes = dishes
        sets = append(sets, &s)
    }
    
    if err := rows.Err(); err != nil {
        sr.logger.Error("Error iterating over sets: " + err.Error())
        return nil, fmt.Errorf("error iterating over sets: %w", err)
    }
    
    sr.logger.Info(fmt.Sprintf("Successfully fetched %d sets", len(sets)))
    return sets, nil
}




func (sr *SetRepository) GetSetProtoDetail(ctx context.Context, id int32) (*set.SetProto, error) {
    sr.logger.Info(fmt.Sprintf("Fetching set detail for ID: %d", id))
    query := `
        SELECT id, name, description, user_id, created_at, updated_at, is_favourite, like_by
        FROM sets 
        WHERE id = $1
    `
    var s set.SetProto
    var createdAt, updatedAt time.Time
    var userID *int32
    var likeBy []int64
    err := sr.db.QueryRow(ctx, query, id).Scan(
        &s.Id,
        &s.Name,
        &s.Description,
        &userID,
        &createdAt,
        &updatedAt,
        &s.IsFavourite,
        &likeBy,
    )
    if err != nil {
        sr.logger.Error(fmt.Sprintf("Error fetching set detail for ID %d: %s", id, err.Error()))
        return nil, fmt.Errorf("error fetching set detail: %w", err)
    }
    s.CreatedAt = timestamppb.New(createdAt)
    s.UpdatedAt = timestamppb.New(updatedAt)
    s.UserId = userID
    s.LikeBy = likeBy

    s.Dishes, err = sr.getSetDishes(ctx, s.Id)
    if err != nil {
        sr.logger.Error(fmt.Sprintf("Error fetching dishes for set %d: %s", s.Id, err.Error()))
        return nil, fmt.Errorf("error fetching dishes for set %d: %w", s.Id, err)
    }

    sr.logger.Info(fmt.Sprintf("Successfully fetched set detail for ID: %d", id))
    return &s, nil
}

func (sr *SetRepository) CreateSetProto(ctx context.Context, req *set.CreateSetProtoRequest) (*set.SetProto, error) {
    sr.logger.Info(fmt.Sprintf("Creating new set:CreateSetProto repository  %+v", req)) 
    tx, err := sr.db.Begin(ctx)
    if err != nil {
        sr.logger.Error("Error starting transaction: " + err.Error())
        return nil, fmt.Errorf("error starting transaction: %w", err)
    }
    defer tx.Rollback(ctx)

    query := `
        INSERT INTO sets (name, description, user_id, created_at, updated_at, is_favourite, like_by, is_public)
        VALUES ($1, $2, $3, $4, $4, $5, $6, $7)
        RETURNING id, created_at, updated_at
    `
    var s set.SetProto
    var createdAt, updatedAt time.Time
    err = tx.QueryRow(ctx, query,
        req.Name,
        req.Description,
        req.UserId,
        time.Now(),
        false,
        []int64{}, // Empty array for like_by
        req.IsPublic,
    ).Scan(&s.Id, &createdAt, &updatedAt)
    if err != nil {
        sr.logger.Error("Error creating set: " + err.Error())
        return nil, fmt.Errorf("error creating set: %w", err)
    }

    s.Name = req.Name
    s.Description = req.Description
    s.UserId = &req.UserId
    s.CreatedAt = timestamppb.New(createdAt)
    s.UpdatedAt = timestamppb.New(updatedAt)
    s.IsFavourite = false
    s.LikeBy = []int64{}
    s.IsPublic = req.IsPublic

    for _, dish := range req.Dishes {
        _, err := tx.Exec(ctx, "INSERT INTO set_dishes (set_id, dish_id, quantity) VALUES ($1, $2, $3)",
            s.Id, dish.DishId, dish.Quantity)
        if err != nil {
            sr.logger.Error(fmt.Sprintf("Error inserting set dish: %s", err.Error()))
            return nil, fmt.Errorf("error inserting set dish: %w", err)
        }
    }

    if err := tx.Commit(ctx); err != nil {
        sr.logger.Error("Error committing transaction: " + err.Error())
        return nil, fmt.Errorf("error committing transaction: %w", err)
    }

    s.Dishes = req.Dishes
    sr.logger.Info(fmt.Sprintf("Successfully created new set with ID: %d", s.Id))
    return &s, nil
}

func (sr *SetRepository) UpdateSetProto(ctx context.Context, req *set.UpdateSetProtoRequest) (*set.SetProto, error) {
    sr.logger.Info(fmt.Sprintf("Updating set with ID: %d", req.Id))
    tx, err := sr.db.Begin(ctx)
    if err != nil {
        sr.logger.Error("Error starting transaction: " + err.Error())
        return nil, fmt.Errorf("error starting transaction: %w", err)
    }
    defer tx.Rollback(ctx)

    query := `
        UPDATE sets
        SET name = $2, description = $3, updated_at = $4
        WHERE id = $1
        RETURNING user_id, created_at, updated_at, is_favourite, like_by
    `
    var s set.SetProto
    var createdAt, updatedAt time.Time
    var userID *int32
    var likeBy []int64
    err = tx.QueryRow(ctx, query,
        req.Id,
        req.Name,
        req.Description,
        time.Now(),
    ).Scan(&userID, &createdAt, &updatedAt, &s.IsFavourite, &likeBy)
    if err != nil {
        sr.logger.Error(fmt.Sprintf("Error updating set with ID %d: %s", req.Id, err.Error()))
        return nil, fmt.Errorf("error updating set: %w", err)
    }

    s.Id = req.Id
    s.Name = req.Name
    s.Description = req.Description
    s.UserId = userID
    s.CreatedAt = timestamppb.New(createdAt)
    s.UpdatedAt = timestamppb.New(updatedAt)
    s.LikeBy = likeBy

    _, err = tx.Exec(ctx, "DELETE FROM set_dishes WHERE set_id = $1", req.Id)
    if err != nil {
        sr.logger.Error(fmt.Sprintf("Error deleting existing set dishes for set ID %d: %s", req.Id, err.Error()))
        return nil, fmt.Errorf("error deleting existing set dishes: %w", err)
    }

    for _, dish := range req.Dishes {
        _, err := tx.Exec(ctx, "INSERT INTO set_dishes (set_id, dish_id, quantity) VALUES ($1, $2, $3)",
            s.Id, dish.DishId, dish.Quantity)
        if err != nil {
            sr.logger.Error(fmt.Sprintf("Error inserting updated set dish for set ID %d: %s", s.Id, err.Error()))
            return nil, fmt.Errorf("error inserting updated set dish: %w", err)
        }
    }

    if err := tx.Commit(ctx); err != nil {
        sr.logger.Error("Error committing transaction: " + err.Error())
        return nil, fmt.Errorf("error committing transaction: %w", err)
    }

    s.Dishes = req.Dishes
    sr.logger.Info(fmt.Sprintf("Successfully updated set with ID: %d", s.Id))
    return &s, nil
}

func (sr *SetRepository) DeleteSetProto(ctx context.Context, id int32) (*set.SetProto, error) {
    sr.logger.Info(fmt.Sprintf("Deleting set with ID: %d", id))
    tx, err := sr.db.Begin(ctx)
    if err != nil {
        sr.logger.Error("Error starting transaction: " + err.Error())
        return nil, fmt.Errorf("error starting transaction: %w", err)
    }
    defer tx.Rollback(ctx)

    _, err = tx.Exec(ctx, "DELETE FROM set_dishes WHERE set_id = $1", id)
    if err != nil {
        sr.logger.Error(fmt.Sprintf("Error deleting set dishes for set ID %d: %s", id, err.Error()))
        return nil, fmt.Errorf("error deleting set dishes: %w", err)
    }

    query := `
        DELETE FROM sets
        WHERE id = $1
        RETURNING id, name, description, user_id, created_at, updated_at, is_favourite, like_by
    `
    var s set.SetProto
    var createdAt, updatedAt time.Time
    var userID *int32
    var likeBy []int64
    err = tx.QueryRow(ctx, query, id).Scan(
        &s.Id,
        &s.Name,
        &s.Description,
        &userID,
        &createdAt,
        &updatedAt,
        &s.IsFavourite,
        &likeBy,
    )
    if err != nil {
        sr.logger.Error(fmt.Sprintf("Error deleting set with ID %d: %s", id, err.Error()))
        return nil, fmt.Errorf("error deleting set: %w", err)
    }
    s.CreatedAt = timestamppb.New(createdAt)
    s.UpdatedAt = timestamppb.New(updatedAt)
    s.UserId = userID
    s.LikeBy = likeBy

    if err := tx.Commit(ctx); err != nil {
        sr.logger.Error("Error committing transaction: " + err.Error())
        return nil, fmt.Errorf("error committing transaction: %w", err)
    }

    sr.logger.Info(fmt.Sprintf("Successfully deleted set with ID: %d", id))
    return &s, nil
}

func (sr *SetRepository) getSetDishes(ctx context.Context, setID int32) ([]*set.SetProtoDish, error) {
    sr.logger.Info(fmt.Sprintf("Fetching dishes for set ID: %d", setID))
    query := `
        SELECT dish_id, quantity
        FROM set_dishes
        WHERE set_id = $1
    `
    rows, err := sr.db.Query(ctx, query, setID)
    if err != nil {
        sr.logger.Error(fmt.Sprintf("Error fetching set dishes for set ID %d: %s", setID, err.Error()))
        return nil, fmt.Errorf("error fetching set dishes: %w", err)
    }
    defer rows.Close()

    var dishes []*set.SetProtoDish
    for rows.Next() {
        var sd set.SetProtoDish
        err := rows.Scan(
            &sd.DishId,
            &sd.Quantity,
        )
        if err != nil {
            sr.logger.Error(fmt.Sprintf("Error scanning set dish for set ID %d: %s", setID, err.Error()))
            return nil, fmt.Errorf("error scanning set dish: %w", err)
        }
        dishes = append(dishes, &sd)
    }
    if err := rows.Err(); err != nil {
        sr.logger.Error(fmt.Sprintf("Error iterating over set dishes for set ID %d: %s", setID, err.Error()))
        return nil, fmt.Errorf("error iterating over set dishes: %w", err)
    }

    sr.logger.Info(fmt.Sprintf("Successfully fetched %d dishes for set ID: %d", len(dishes), setID))
    return dishes, nil
}

func (sr *SetRepository) GetDishesForSet(ctx context.Context, setID int32) ([]*set.SetProtoDish, error) {
    query := `
        SELECT sd.dish_id, sd.quantity, d.id, d.name, d.price, d.description, d.image, d.status, d.created_at, d.updated_at
        FROM set_dishes sd
        JOIN dishes d ON sd.dish_id = d.id
        WHERE sd.set_id = $1
    `
    rows, err := sr.db.Query(ctx, query, setID)
    if err != nil {
        return nil, fmt.Errorf("error fetching dishes for set: %w", err)
    }
    defer rows.Close()

    var dishes []*set.SetProtoDish
    for rows.Next() {
        var spd set.SetProtoDish
        var d set.Dish
        var createdAt, updatedAt time.Time
        err := rows.Scan(
            &spd.DishId,
            &spd.Quantity,
            &d.Id,
            &d.Name,
            &d.Price,
            &d.Description,
            &d.Image,
            &d.Status,
            &createdAt,
            &updatedAt,
        )
        if err != nil {
            return nil, fmt.Errorf("error scanning dish: %w", err)
        }
        d.CreatedAt = timestamppb.New(createdAt)
        d.UpdatedAt = timestamppb.New(updatedAt)
        spd.Dish = &d
        dishes = append(dishes, &spd)
    }
    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("error iterating over dishes: %w", err)
    }
    return dishes, nil
}