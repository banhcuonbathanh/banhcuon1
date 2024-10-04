package order_grpc

import (
	"context"
	"fmt"
	"time"


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

func (or *OrderRepository) CreateOrders(ctx context.Context, req *order.CreateOrdersRequest) ([]*order.OrderDetail, error) {
	tx, err := or.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("error starting transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	var orderDetails []*order.OrderDetail

	for _, item := range req.Orders {
		query := `
			INSERT INTO orders (guest_id, dish_snapshot_id, quantity, status, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $5)
			RETURNING id, guest_id, table_number, dish_snapshot_id, quantity, order_handler_id, status, created_at, updated_at
		`
		var o order.Order
		var createdAt, updatedAt time.Time
		err := tx.QueryRow(ctx, query,
			req.GuestId,
			item.DishId,
			item.Quantity,
			"pending",
			time.Now(),
		).Scan(
			&o.Id,
			&o.GuestId,
			&o.TableNumber,
			&o.DishSnapshotId,
			&o.Quantity,
			&o.OrderHandlerId,
			&o.Status,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error creating order: %w", err)
		}

		o.CreatedAt = timestamppb.New(createdAt)
		o.UpdatedAt = timestamppb.New(updatedAt)
		orderDetail, err := or.getOrderDetail(ctx, or.db, &o)

		// Fetch additional details (guest, dish_snapshot, order_handler, table)
	
		if err != nil {
			return nil, err
		}

		orderDetails = append(orderDetails, orderDetail)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("error committing transaction: %w", err)
	}

	return orderDetails, nil
}

func (or *OrderRepository) GetOrders(ctx context.Context, req *order.GetOrdersRequest) ([]*order.OrderDetail, error) {
	query := `
		SELECT id, guest_id, table_number, dish_snapshot_id, quantity, order_handler_id, status, created_at, updated_at
		FROM orders
		WHERE created_at BETWEEN $1 AND $2
	`
	rows, err := or.db.Query(ctx, query, req.FromDate.AsTime(), req.ToDate.AsTime())
	if err != nil {
		return nil, fmt.Errorf("error fetching orders: %w", err)
	}
	defer rows.Close()

	var orderDetails []*order.OrderDetail
	for rows.Next() {
		var o order.Order
		var createdAt, updatedAt time.Time
		err := rows.Scan(
			&o.Id,
			&o.GuestId,
			&o.TableNumber,
			&o.DishSnapshotId,
			&o.Quantity,
			&o.OrderHandlerId,
			&o.Status,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning order: %w", err)
		}
		o.CreatedAt = timestamppb.New(createdAt)
		o.UpdatedAt = timestamppb.New(updatedAt)

		orderDetail, err := or.getOrderDetail(ctx, or.db, &o)
		if err != nil {
			return nil, err
		}

		orderDetails = append(orderDetails, orderDetail)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over orders: %w", err)
	}

	return orderDetails, nil
}

func (or *OrderRepository) GetOrderDetail(ctx context.Context, id int64) (*order.OrderDetail, error) {
	query := `
		SELECT id, guest_id, table_number, dish_snapshot_id, quantity, order_handler_id, status, created_at, updated_at
		FROM orders
		WHERE id = $1
	`
	var o order.Order
	var createdAt, updatedAt time.Time
	err := or.db.QueryRow(ctx, query, id).Scan(
		&o.Id,
		&o.GuestId,
		&o.TableNumber,
		&o.DishSnapshotId,
		&o.Quantity,
		&o.OrderHandlerId,
		&o.Status,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("error fetching order detail: %w", err)
	}
	o.CreatedAt = timestamppb.New(createdAt)
	o.UpdatedAt = timestamppb.New(updatedAt)

	return or.getOrderDetail(ctx, or.db, &o)
}

func (or *OrderRepository) UpdateOrder(ctx context.Context, req *order.UpdateOrderRequest) (*order.OrderDetail, error) {
	query := `
		UPDATE orders
		SET status = $2, dish_snapshot_id = $3, quantity = $4, updated_at = $5
		WHERE id = $1
		RETURNING id, guest_id, table_number, dish_snapshot_id, quantity, order_handler_id, status, created_at, updated_at
	`
	var o order.Order
	var createdAt, updatedAt time.Time
	err := or.db.QueryRow(ctx, query,
		req.OrderId,
		req.Status,
		req.DishId,
		req.Quantity,
		time.Now(),
	).Scan(
		&o.Id,
		&o.GuestId,
		&o.TableNumber,
		&o.DishSnapshotId,
		&o.Quantity,
		&o.OrderHandlerId,
		&o.Status,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("error updating order: %w", err)
	}
	o.CreatedAt = timestamppb.New(createdAt)
	o.UpdatedAt = timestamppb.New(updatedAt)

	return or.getOrderDetail(ctx, or.db, &o)
}

func (or *OrderRepository) PayGuestOrders(ctx context.Context, guestId int64) ([]*order.OrderDetail, error) {
	tx, err := or.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("error starting transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	query := `
		UPDATE orders
		SET status = 'paid', updated_at = $2
		WHERE guest_id = $1 AND status = 'pending'
		RETURNING id, guest_id, table_number, dish_snapshot_id, quantity, order_handler_id, status, created_at, updated_at
	`
	rows, err := tx.Query(ctx, query, guestId, time.Now())
	if err != nil {
		return nil, fmt.Errorf("error paying guest orders: %w", err)
	}
	defer rows.Close()

	var orderDetails []*order.OrderDetail
	for rows.Next() {
		var o order.Order
		var createdAt, updatedAt time.Time
		err := rows.Scan(
			&o.Id,
			&o.GuestId,
			&o.TableNumber,
			&o.DishSnapshotId,
			&o.Quantity,
			&o.OrderHandlerId,
			&o.Status,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning paid order: %w", err)
		}
		o.CreatedAt = timestamppb.New(createdAt)
		o.UpdatedAt = timestamppb.New(updatedAt)

		orderDetail, err := or.getOrderDetail(ctx, or.db, &o)

		if err != nil {
			return nil, err
		}

		orderDetails = append(orderDetails, orderDetail)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over paid orders: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("error committing transaction: %w", err)
	}

	return orderDetails, nil
}





func (or *OrderRepository) getOrderDetail(ctx context.Context, db *pgxpool.Pool, o *order.Order) (*order.OrderDetail, error) {
	var orderDetail order.OrderDetail
	orderDetail.Order = o

	// Fetch Guest
	guestQuery := `
		SELECT id, name, table_number, created_at, updated_at
		FROM guests
		WHERE id = $1
	`
	var guest order.Guest
	var guestCreatedAt, guestUpdatedAt time.Time
	err := db.QueryRow(ctx, guestQuery, o.GuestId).Scan(
		&guest.Id,
		&guest.Name,
		&guest.TableNumber,
		&guestCreatedAt,
		&guestUpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("error fetching guest: %w", err)
	}
	guest.CreatedAt = timestamppb.New(guestCreatedAt)
	guest.UpdatedAt = timestamppb.New(guestUpdatedAt)
	orderDetail.Guest = &guest

	// Fetch DishSnapshot
	dishSnapshotQuery := `
		SELECT id, name, price, image, description, status, dish_id, created_at, updated_at
		FROM dish_snapshots
		WHERE id = $1
	`
	var dishSnapshot order.DishSnapshot
	var dishSnapshotCreatedAt, dishSnapshotUpdatedAt time.Time
	err = db.QueryRow(ctx, dishSnapshotQuery, o.DishSnapshotId).Scan(
		&dishSnapshot.Id,
		&dishSnapshot.Name,
		&dishSnapshot.Price,
		&dishSnapshot.Image,
		&dishSnapshot.Description,
		&dishSnapshot.Status,
		&dishSnapshot.DishId,
		&dishSnapshotCreatedAt,
		&dishSnapshotUpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("error fetching dish snapshot: %w", err)
	}
	dishSnapshot.CreatedAt = timestamppb.New(dishSnapshotCreatedAt)
	dishSnapshot.UpdatedAt = timestamppb.New(dishSnapshotUpdatedAt)
	orderDetail.DishSnapshot = &dishSnapshot

	// Fetch Account (OrderHandler)
	accountQuery := `
		SELECT id, name, email, role, avatar
		FROM accounts
		WHERE id = $1
	`
	var account order.Account
	err = db.QueryRow(ctx, accountQuery, o.OrderHandlerId).Scan(
		&account.Id,
		&account.Name,
		&account.Email,
		&account.Role,
		&account.Avatar,
	)
	if err != nil {
		return nil, fmt.Errorf("error fetching order handler account: %w", err)
	}
	orderDetail.OrderHandler = &account

	// Fetch Table
	tableQuery := `
		SELECT number, capacity, status, token, created_at, updated_at
		FROM tables
		WHERE number = $1
	`
	var table order.Table
	var tableCreatedAt, tableUpdatedAt time.Time
	err = db.QueryRow(ctx, tableQuery, o.TableNumber).Scan(
		&table.Number,
		&table.Capacity,
		&table.Status,
		&table.Token,
		&tableCreatedAt,
		&tableUpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("error fetching table: %w", err)
	}
	table.CreatedAt = timestamppb.New(tableCreatedAt)
	table.UpdatedAt = timestamppb.New(tableUpdatedAt)
	orderDetail.Table = &table

	return &orderDetail, nil
}
