package qr_guests

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/protobuf/types/known/timestamppb"

	proto "english-ai-full/quanqr/proto_qr/guest"
)

type GuestRepository struct {
	db *pgxpool.Pool
}

func NewGuestRepository(db *pgxpool.Pool) *GuestRepository {
	return &GuestRepository{
		db: db,
	}
}

func (gr *GuestRepository) GuestLogin(ctx context.Context, req *proto.GuestLoginRequest) (*proto.GuestLoginResponse, error) {
	query := `
		INSERT INTO guests (name, table_number, refresh_token, refresh_token_expires_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $5)
		RETURNING id, created_at, updated_at
	`
	var guest proto.GuestInfo
	var createdAt, updatedAt time.Time
	refreshTokenExpiresAt := time.Now().Add(24 * time.Hour) // Set refresh token to expire in 24 hours

	err := gr.db.QueryRow(ctx, query,
		req.Name,
		req.TableNumber,
		req.Token,
		refreshTokenExpiresAt,
		time.Now(),
	).Scan(&guest.Id, &createdAt, &updatedAt)
	if err != nil {
		return nil, fmt.Errorf("error creating guest: %w", err)
	}

	guest.Name = req.Name
	guest.Role = "guest" // Assuming all are guests for now
	guest.TableNumber = req.TableNumber
	guest.CreatedAt = timestamppb.New(createdAt)
	guest.UpdatedAt = timestamppb.New(updatedAt)

	// In a real-world scenario, you'd generate these tokens securely
	accessToken := "sample_access_token"
	refreshToken := req.Token

	return &proto.GuestLoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Guest:        &guest,
		Message:      "Guest logged in successfully",
	}, nil
}

func (gr *GuestRepository) GuestLogout(ctx context.Context, req *proto.LogoutRequest) error {
	query := `
		UPDATE guests
		SET refresh_token = NULL, refresh_token_expires_at = NULL
		WHERE refresh_token = $1
	`
	_, err := gr.db.Exec(ctx, query, req.RefreshToken)
	if err != nil {
		return fmt.Errorf("error logging out guest: %w", err)
	}
	return nil
}

func (gr *GuestRepository) GuestRefreshToken(ctx context.Context, req *proto.RefreshTokenRequest) (*proto.RefreshTokenResponse, error) {
	query := `
		SELECT id FROM guests
		WHERE refresh_token = $1 AND refresh_token_expires_at > $2
	`
	var guestID int64
	err := gr.db.QueryRow(ctx, query, req.RefreshToken, time.Now()).Scan(&guestID)
	if err != nil {
		return nil, fmt.Errorf("error refreshing token: %w", err)
	}

	// In a real-world scenario, you'd generate these tokens securely
	newAccessToken := "new_sample_access_token"
	newRefreshToken := "new_sample_refresh_token"

	updateQuery := `
		UPDATE guests
		SET refresh_token = $1, refresh_token_expires_at = $2
		WHERE id = $3
	`
	_, err = gr.db.Exec(ctx, updateQuery, newRefreshToken, time.Now().Add(24*time.Hour), guestID)
	if err != nil {
		return nil, fmt.Errorf("error updating refresh token: %w", err)
	}

	return &proto.RefreshTokenResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		Message:      "Tokens refreshed successfully",
	}, nil
}

func (gr *GuestRepository) GuestCreateOrders(ctx context.Context, req *proto.GuestCreateOrderRequest) (*proto.OrdersResponse, error) {

	log.Print("golang/quanqr/qr_guests/qr_guests_repository.go")
	tx, err := gr.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("error starting transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	var orders []*proto.Order
	for _, item := range req.Items {
		query := `
			INSERT INTO orders (guest_id, table_number, dish_id, quantity, status, created_at, updated_at)
			SELECT $1, table_number, $2, $3, 'pending', $4, $4
			FROM guests
			WHERE id = $1
			RETURNING id, guest_id, table_number, dish_id, quantity, status, created_at, updated_at
		`
		var order proto.Order
		var createdAt, updatedAt time.Time
		err := tx.QueryRow(ctx, query,
			item.GuestId,
			item.DishId,
			item.Quantity,
			time.Now(),
		).Scan(
			&order.Id,
			&order.GuestId,
			&order.TableNumber,
			&order.DishId,
			&order.Quantity,
			&order.Status,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error creating order: error %w", err)
		}
		order.CreatedAt = timestamppb.New(createdAt)
		order.UpdatedAt = timestamppb.New(updatedAt)
		orders = append(orders, &order)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, fmt.Errorf("error committing transaction: %w", err)
	}

	return &proto.OrdersResponse{
		Data:    orders,
		Message: "Orders created successfully",
	}, nil
}

func (gr *GuestRepository) GuestGetOrders(ctx context.Context, req *proto.GuestGetOrdersGRPCRequest) (*proto.ListOrdersResponse, error) {
	query := `
		SELECT id, guest_id, table_number, dish_id, quantity, status, created_at, updated_at
		FROM orders
		WHERE guest_id = $1
		ORDER BY created_at DESC
	`
	rows, err := gr.db.Query(ctx, query, req.GuestId)
	if err != nil {
		return nil, fmt.Errorf("error fetching orders: %w", err)
	}
	defer rows.Close()

	var orders []*proto.Order
	for rows.Next() {
		var order proto.Order
		var createdAt, updatedAt time.Time
		err := rows.Scan(
			&order.Id,
			&order.GuestId,
			&order.TableNumber,
			&order.DishId,
			&order.Quantity,
			&order.Status,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning order: %w", err)
		}
		order.CreatedAt = timestamppb.New(createdAt)
		order.UpdatedAt = timestamppb.New(updatedAt)
		orders = append(orders, &order)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over orders: %w", err)
	}

	return &proto.ListOrdersResponse{
		Orders:  orders,
		Message: "Orders fetched successfully",
	}, nil
}