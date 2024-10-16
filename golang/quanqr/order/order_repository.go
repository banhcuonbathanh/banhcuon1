package order_grpc

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/protobuf/types/known/timestamppb"

	"english-ai-full/quanqr/proto_qr/order"
)
type OrderRepository struct {
	db *pgxpool.Pool
}

func NewOrderRepository(db *pgxpool.Pool) *OrderRepository {
	return &OrderRepository{
		db: db,
	}
}

func (or *OrderRepository) CreateOrders(ctx context.Context, req *order.CreateOrderRequest) (*order.OrderListResponse, error) {
	tx, err := or.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("error starting transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	var orders []*order.Order

	// Insert main order
	query := `
		INSERT INTO orders (guest_id, table_number, order_handler_id, status, total_price, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $6)
		RETURNING id, guest_id, table_number, order_handler_id, status, total_price, created_at, updated_at
	`
	var o order.Order
	var createdAt, updatedAt time.Time
	err = tx.QueryRow(ctx, query,
		req.GuestId,
		req.TableNumber,
		req.OrderHandlerId,
		req.Status,
		req.TotalPrice,
		time.Now(),
	).Scan(
		&o.Id,
		&o.GuestId,
		&o.TableNumber,
		&o.OrderHandlerId,
		&o.Status,
		&o.TotalPrice,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("error creating order: %w", err)
	}

	o.CreatedAt = timestamppb.New(createdAt)
	o.UpdatedAt = timestamppb.New(updatedAt)

	// Insert dish order items
	for _, item := range req.DishItems {
		_, err := tx.Exec(ctx, `
			INSERT INTO dish_order_items (order_id, dish_snapshot_id, quantity)
			VALUES ($1, $2, $3)
		`, o.Id, item.DishSnapshotId, item.Quantity)
		if err != nil {
			return nil, fmt.Errorf("error inserting dish order item: %w", err)
		}
		o.DishItems = append(o.DishItems, item)
	}

	// Insert set order items
	for _, item := range req.SetItems {
		_, err := tx.Exec(ctx, `
			INSERT INTO set_order_items (order_id, set_snapshot_id, quantity)
			VALUES ($1, $2, $3)
		`, o.Id, item.SetSnapshotId, item.Quantity)
		if err != nil {
			return nil, fmt.Errorf("error inserting set order item: %w", err)
		}
		o.SetItems = append(o.SetItems, item)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("error committing transaction: %w", err)
	}

	orders = append(orders, &o)

	return &order.OrderListResponse{Data: orders}, nil
}

func (or *OrderRepository) GetOrders(ctx context.Context, req *order.GetOrdersRequest) (*order.OrderListResponse, error) {
	query := `
		SELECT id, guest_id, table_number, order_handler_id, status, total_price, created_at, updated_at
		FROM orders
		WHERE created_at BETWEEN $1 AND $2
	`
	rows, err := or.db.Query(ctx, query, req.FromDate.AsTime(), req.ToDate.AsTime())
	if err != nil {
		return nil, fmt.Errorf("error fetching orders: %w", err)
	}
	defer rows.Close()

	var orders []*order.Order
	for rows.Next() {
		var o order.Order
		var createdAt, updatedAt time.Time
		err := rows.Scan(
			&o.Id,
			&o.GuestId,
			&o.TableNumber,
			&o.OrderHandlerId,
			&o.Status,
			&o.TotalPrice,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning order: %w", err)
		}
		o.CreatedAt = timestamppb.New(createdAt)
		o.UpdatedAt = timestamppb.New(updatedAt)

		if err := or.fetchOrderItems(ctx, &o); err != nil {
			return nil, err
		}

		orders = append(orders, &o)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over orders: %w", err)
	}

	return &order.OrderListResponse{Data: orders}, nil
}

func (or *OrderRepository) GetOrderDetail(ctx context.Context, req *order.OrderIdParam) (*order.OrderResponse, error) {
	query := `
		SELECT id, guest_id, table_number, order_handler_id, status, total_price, created_at, updated_at
		FROM orders
		WHERE id = $1
	`
	var o order.Order
	var createdAt, updatedAt time.Time
	err := or.db.QueryRow(ctx, query, req.Id).Scan(
		&o.Id,
		&o.GuestId,
		&o.TableNumber,
		&o.OrderHandlerId,
		&o.Status,
		&o.TotalPrice,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("error fetching order detail: %w", err)
	}
	o.CreatedAt = timestamppb.New(createdAt)
	o.UpdatedAt = timestamppb.New(updatedAt)

	if err := or.fetchOrderItems(ctx, &o); err != nil {
		return nil, err
	}

	return &order.OrderResponse{Data: &o}, nil
}

func (or *OrderRepository) UpdateOrder(ctx context.Context, req *order.UpdateOrderRequest) (*order.OrderResponse, error) {
	tx, err := or.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("error starting transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	query := `
		UPDATE orders
		SET guest_id = $2, table_number = $3, order_handler_id = $4, status = $5, total_price = $6, updated_at = $7
		WHERE id = $1
		RETURNING id, guest_id, table_number, order_handler_id, status, total_price, created_at, updated_at
	`
	var o order.Order
	var createdAt, updatedAt time.Time
	err = tx.QueryRow(ctx, query,
		req.Id,
		req.GuestId,
		req.TableNumber,
		req.OrderHandlerId,
		req.Status,
		req.TotalPrice,
		time.Now(),
	).Scan(
		&o.Id,
		&o.GuestId,
		&o.TableNumber,
		&o.OrderHandlerId,
		&o.Status,
		&o.TotalPrice,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("error updating order: %w", err)
	}
	o.CreatedAt = timestamppb.New(createdAt)
	o.UpdatedAt = timestamppb.New(updatedAt)

	// Update dish order items
	if err := or.updateOrderItems(ctx, tx, &o, req.DishItems, req.SetItems); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("error committing transaction: %w", err)
	}

	return &order.OrderResponse{Data: &o}, nil
}

func (or *OrderRepository) PayGuestOrders(ctx context.Context, req *order.PayGuestOrdersRequest) (*order.OrderListResponse, error) {
	tx, err := or.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("error starting transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	query := `
		UPDATE orders
		SET status = 'paid', updated_at = $2
		WHERE guest_id = $1 AND status = 'pending'
		RETURNING id, guest_id, table_number, order_handler_id, status, total_price, created_at, updated_at
	`
	rows, err := tx.Query(ctx, query, req.GuestId, time.Now())
	if err != nil {
		return nil, fmt.Errorf("error paying guest orders: %w", err)
	}
	defer rows.Close()

	var orders []*order.Order
	for rows.Next() {
		var o order.Order
		var createdAt, updatedAt time.Time
		err := rows.Scan(
			&o.Id,
			&o.GuestId,
			&o.TableNumber,
			&o.OrderHandlerId,
			&o.Status,
			&o.TotalPrice,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning paid order: %w", err)
		}
		o.CreatedAt = timestamppb.New(createdAt)
		o.UpdatedAt = timestamppb.New(updatedAt)

		if err := or.fetchOrderItems(ctx, &o); err != nil {
			return nil, err
		}

		orders = append(orders, &o)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over paid orders: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("error committing transaction: %w", err)
	}

	return &order.OrderListResponse{Data: orders}, nil
}

func (or *OrderRepository) fetchOrderItems(ctx context.Context, o *order.Order) error {
	// Fetch dish order items
	dishRows, err := or.db.Query(ctx, `
		SELECT dish_snapshot_id, quantity
		FROM dish_order_items
		WHERE order_id = $1
	`, o.Id)
	if err != nil {
		return fmt.Errorf("error fetching dish order items: %w", err)
	}
	defer dishRows.Close()

	for dishRows.Next() {
		var item order.DishOrderItem
		if err := dishRows.Scan(&item.DishSnapshotId, &item.Quantity); err != nil {
			return fmt.Errorf("error scanning dish order item: %w", err)
		}
		o.DishItems = append(o.DishItems, &item)
	}

	// Fetch set order items
	setRows, err := or.db.Query(ctx, `
		SELECT set_snapshot_id, quantity
		FROM set_order_items
		WHERE order_id = $1
	`, o.Id)
	if err != nil {
		return fmt.Errorf("error fetching set order items: %w", err)
	}
	defer setRows.Close()

	for setRows.Next() {
		var item order.SetOrderItem
		if err := setRows.Scan(&item.SetSnapshotId, &item.Quantity); err != nil {
			return fmt.Errorf("error scanning set order item: %w", err)
		}
		o.SetItems = append(o.SetItems, &item)
	}

	return nil
}

func (or *OrderRepository) updateOrderItems(ctx context.Context, tx pgx.Tx, o *order.Order, dishItems []*order.DishOrderItem, setItems []*order.SetOrderItem) error {
	// Delete existing items
	if _, err := tx.Exec(ctx, "DELETE FROM dish_order_items WHERE order_id = $1", o.Id); err != nil {
		return fmt.Errorf("error deleting existing dish order items: %w", err)
	}
	if _, err := tx.Exec(ctx, "DELETE FROM set_order_items WHERE order_id = $1", o.Id); err != nil {
		return fmt.Errorf("error deleting existing set order items: %w", err)
	}

	// Insert new dish order items
	for _, item := range dishItems {
		_, err := tx.Exec(ctx, `
			INSERT INTO dish_order_items (order_id, dish_snapshot_id, quantity)
			VALUES ($1, $2, $3)
		`, o.Id, item.DishSnapshotId, item.Quantity)
		if err != nil {
			return fmt.Errorf("error inserting dish order item: %w", err)
		}
	}

	// Insert new set order items
	for _, item := range setItems {
		_, err := tx.Exec(ctx, `
			INSERT INTO set_order_items (order_id, set_snapshot_id, quantity)
			VALUES ($1, $2, $3)
		`, o.Id, item.SetSnapshotId, item.Quantity)
		if err != nil {
			return fmt.Errorf("error inserting set order item: %w", err)
		}
	}

	o.DishItems = dishItems
	o.SetItems = setItems

	return nil
}