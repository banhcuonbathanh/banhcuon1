package order_grpc

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/protobuf/types/known/timestamppb"

	"english-ai-full/logger"
	"english-ai-full/quanqr/proto_qr/order"
)

type OrderRepository struct {
	db     *pgxpool.Pool
	logger *logger.Logger
}

func NewOrderRepository(db *pgxpool.Pool) *OrderRepository {
	return &OrderRepository{
		db:     db,
		logger: logger.NewLogger(),
	}
}
func (or *OrderRepository) CreateOrder(ctx context.Context, req *order.CreateOrderRequest) (*order.Order, error) {
    or.logger.Info(fmt.Sprintf("Creating new order: %+v", req))
    tx, err := or.db.Begin(ctx)
    if err != nil {
        or.logger.Error("Error starting transaction: " + err.Error())
        return nil, fmt.Errorf("error starting transaction: %w", err)
    }
    defer tx.Rollback(ctx)

    // Check if the table exists
    var tableExists bool
    err = tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM tables WHERE number = $1)", req.TableNumber).Scan(&tableExists)
    if err != nil {
        or.logger.Error("Error checking table existence: " + err.Error())
        return nil, fmt.Errorf("error checking table existence: %w", err)
    }
    if !tableExists {
        or.logger.Error(fmt.Sprintf("Table number %d does not exist", req.TableNumber))
        return nil, fmt.Errorf("table number %d does not exist", req.TableNumber)
    }

	// check set

    // Create the order
    query := `
        INSERT INTO orders 
            (guest_id, user_id, is_guest, table_number, order_handler_id, status, created_at, updated_at, total_price)
        VALUES 
            ($1, $2, $3, $4, $5, $6, $7, $7, $8)
        RETURNING id, created_at, updated_at
    `
    
    var o order.Order
    var createdAt, updatedAt time.Time
    var guestId, userId sql.NullInt64

    if req.IsGuest {
        guestId = sql.NullInt64{Int64: req.GuestId, Valid: true}
        userId = sql.NullInt64{Valid: false}
    } else {
        userId = sql.NullInt64{Int64: req.UserId, Valid: true}
        guestId = sql.NullInt64{Valid: false}
    }
    err = tx.QueryRow(ctx, query,
        guestId,
        userId,
        req.IsGuest,
        req.TableNumber,
        req.OrderHandlerId,
        req.Status,
        time.Now(),
        req.TotalPrice,
    ).Scan(&o.Id, &createdAt, &updatedAt)
    
    if err != nil {
        or.logger.Error("Error creating order: " + err.Error())
        return nil, fmt.Errorf("error creating order: %w", err)
    }

    // Insert dish order items
    for _, item := range req.DishItems {
        // Create a dish snapshot
        var dishSnapshotId int64
        err := tx.QueryRow(ctx, `
            INSERT INTO dish_snapshots (name, price, description, image, status, dish_id)
            SELECT name, price, description, image, status, id
            FROM dishes
            WHERE id = $1
            RETURNING id
        `, item.Dish.Id).Scan(&dishSnapshotId)
        if err != nil {
            or.logger.Error("Error creating dish snapshot: " + err.Error())
            return nil, fmt.Errorf("error creating dish snapshot: %w", err)
        }

        // Insert the dish order item using the snapshot
        _, err = tx.Exec(ctx, `
            INSERT INTO dish_order_items (order_id, dish_snapshot_id, quantity)
            VALUES ($1, $2, $3)
        `, o.Id, dishSnapshotId, item.Quantity)
        if err != nil {
            or.logger.Error("Error inserting dish order item: " + err.Error())
            return nil, fmt.Errorf("error inserting dish order item: %w", err)
        }
    }

    // Insert set order items
    for _, item := range req.SetItems {
        // Create a set snapshot
        var setSnapshotId int64

		var setExists bool
		err = tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM sets WHERE id = $1)", item.Set.Id).Scan(&setExists)
		if err != nil {
			or.logger.Error("Error checking set existence: " + err.Error())
			return nil, fmt.Errorf("error checking set existence: %w", err)
		}
		if !setExists {
			or.logger.Error(fmt.Sprintf("Set with ID %d does not exist", item.Set.Id))
			return nil, fmt.Errorf("set with ID %d does not exist", item.Set.Id)
		}

        err := tx.QueryRow(ctx, `
            INSERT INTO set_snapshots (original_set_id, name, description, user_id, is_public, image)
            SELECT id, name, description, user_id, is_public, image
            FROM sets
            WHERE id = $1
            RETURNING id
        `, item.Set.Id).Scan(&setSnapshotId)
        if err != nil {
            or.logger.Error("Error creating set snapshot: " + err.Error())
            return nil, fmt.Errorf("error creating set snapshot: %w", err)
        }

        // Insert the set order item using the snapshot
        _, err = tx.Exec(ctx, `
            INSERT INTO set_order_items (order_id, set_snapshot_id, quantity)
            VALUES ($1, $2, $3)
        `, o.Id, setSnapshotId, item.Quantity)
        if err != nil {
            or.logger.Error("Error inserting set order item: " + err.Error())
            return nil, fmt.Errorf("error inserting set order item: %w", err)
        }

        // Create dish snapshots for each dish in the set and link them to the set snapshot
        for _, dish := range item.Set.Dishes {
            var dishSnapshotId int64
            err := tx.QueryRow(ctx, `
                INSERT INTO dish_snapshots (name, price, description, image, status, dish_id)
                SELECT name, price, description, image, status, id
                FROM dishes
                WHERE id = $1
                RETURNING id
            `, dish.Id).Scan(&dishSnapshotId)
            if err != nil {
                or.logger.Error("Error creating dish snapshot for set: " + err.Error())
                return nil, fmt.Errorf("error creating dish snapshot for set: %w", err)
            }

            _, err = tx.Exec(ctx, `
                INSERT INTO set_snapshot_dishes (set_snapshot_id, dish_snapshot_id, quantity)
                VALUES ($1, $2, 1)
            `, setSnapshotId, dishSnapshotId)
            if err != nil {
                or.logger.Error("Error linking dish snapshot to set snapshot: " + err.Error())
                return nil, fmt.Errorf("error linking dish snapshot to set snapshot: %w", err)
            }
        }
    }

    if err := tx.Commit(ctx); err != nil {
        or.logger.Error("Error committing transaction: " + err.Error())
        return nil, fmt.Errorf("error committing transaction: %w", err)
    }

    // Populate the response object
    o.GuestId = req.GuestId
    o.UserId = req.UserId
    o.IsGuest = req.IsGuest
    o.TableNumber = req.TableNumber
    o.OrderHandlerId = req.OrderHandlerId
    o.Status = req.Status
    o.CreatedAt = timestamppb.New(createdAt)
    o.UpdatedAt = timestamppb.New(updatedAt)
    o.TotalPrice = req.TotalPrice
    o.DishItems = req.DishItems
    o.SetItems = req.SetItems

    or.logger.Info(fmt.Sprintf("Successfully created new order with ID: %d", o.Id))
    return &o, nil
}
// func (or *OrderRepository) CreateOrder(ctx context.Context, req *order.CreateOrderRequest) (*order.Order, error) {
// 	or.logger.Info(fmt.Sprintf("Creating new order: %+v", req))
// 	tx, err := or.db.Begin(ctx)
// 	if err != nil {
// 		or.logger.Error("Error starting transaction: " + err.Error())
// 		return nil, fmt.Errorf("error starting transaction: %w", err)
// 	}
// 	defer tx.Rollback(ctx)

// 	query := `
// 		INSERT INTO orders (guest_id, user_id, is_guest, table_number, order_handler_id, status, created_at, updated_at, total_price)
// 		VALUES ($1, $2, $3, $4, $5, $6, $7, $7, $8)
// 		RETURNING id, created_at, updated_at
// 	`
// 	var o order.Order
// 	var createdAt, updatedAt time.Time
// 	err = tx.QueryRow(ctx, query,
// 		req.GuestId,
// 		req.UserId,
// 		req.IsGuest,
// 		req.TableNumber,
// 		req.OrderHandlerId,
// 		req.Status,
// 		time.Now(),
// 		req.TotalPrice,
// 	).Scan(&o.Id, &createdAt, &updatedAt)
// 	if err != nil {
// 		or.logger.Error("Error creating order: " + err.Error())
// 		return nil, fmt.Errorf("error creating order: %w", err)
// 	}

// 	o.GuestId = req.GuestId
// 	o.UserId = req.UserId
// 	o.IsGuest = req.IsGuest
// 	o.TableNumber = req.TableNumber
// 	o.OrderHandlerId = req.OrderHandlerId
// 	o.Status = req.Status
// 	o.CreatedAt = timestamppb.New(createdAt)
// 	o.UpdatedAt = timestamppb.New(updatedAt)
// 	o.TotalPrice = req.TotalPrice

// 	// Insert dish items
// 	for _, item := range req.DishItems {
// 		err := or.insertDishOrderItem(ctx, tx, o.Id, item)
// 		if err != nil {
// 			return nil, err
// 		}
// 	}

// 	// Insert set items
// 	for _, item := range req.SetItems {
// 		err := or.insertSetOrderItem(ctx, tx, o.Id, item)
// 		if err != nil {
// 			return nil, err
// 		}
// 	}

// 	if err := tx.Commit(ctx); err != nil {
// 		or.logger.Error("Error committing transaction: " + err.Error())
// 		return nil, fmt.Errorf("error committing transaction: %w", err)
// 	}

// 	o.DishItems = req.DishItems
// 	o.SetItems = req.SetItems
// 	or.logger.Info(fmt.Sprintf("Successfully created new order with ID: %d", o.Id))
// 	return &o, nil
// }

func (or *OrderRepository) GetOrders(ctx context.Context, req *order.GetOrdersRequest) ([]*order.Order, error) {
	or.logger.Info("Fetching orders")
	query := `
		SELECT id, guest_id, user_id, is_guest, table_number, order_handler_id, status, created_at, updated_at, total_price
		FROM orders
		WHERE created_at BETWEEN $1 AND $2
	`
	args := []interface{}{req.FromDate.AsTime(), req.ToDate.AsTime()}
	if req.UserId != nil {
		query += " AND user_id = $3"
		args = append(args, *req.UserId)
	} else if req.GuestId != nil {
		query += " AND guest_id = $3"
		args = append(args, *req.GuestId)
	}

	rows, err := or.db.Query(ctx, query, args...)
	if err != nil {
		or.logger.Error("Error fetching orders: " + err.Error())
		return nil, fmt.Errorf("error fetching orders: %w", err)
	}
	defer rows.Close()

	var orders []*order.Order
	for rows.Next() {
		var o order.Order
		var createdAt, updatedAt time.Time
		err := rows.Scan(
			&o.Id, &o.GuestId, &o.UserId, &o.IsGuest, &o.TableNumber, &o.OrderHandlerId,
			&o.Status, &createdAt, &updatedAt, &o.TotalPrice,
		)
		if err != nil {
			or.logger.Error("Error scanning order: " + err.Error())
			return nil, fmt.Errorf("error scanning order: %w", err)
		}
		o.CreatedAt = timestamppb.New(createdAt)
		o.UpdatedAt = timestamppb.New(updatedAt)

		dishItems, err := or.getDishItemsForOrder(ctx, o.Id)
		if err != nil {
			return nil, err
		}
		o.DishItems = dishItems

		setItems, err := or.getSetItemsForOrder(ctx, o.Id)
		if err != nil {
			return nil, err
		}
		o.SetItems = setItems

		orders = append(orders, &o)
	}

	if err := rows.Err(); err != nil {
		or.logger.Error("Error iterating over orders: " + err.Error())
		return nil, fmt.Errorf("error iterating over orders: %w", err)
	}

	or.logger.Info(fmt.Sprintf("Successfully fetched %d orders", len(orders)))
	return orders, nil
}

func (or *OrderRepository) GetOrderDetail(ctx context.Context, id int64) (*order.Order, error) {
	or.logger.Info(fmt.Sprintf("Fetching order detail for ID: %d", id))
	query := `
		SELECT id, guest_id, user_id, is_guest, table_number, order_handler_id, status, created_at, updated_at, total_price
		FROM orders
		WHERE id = $1
	`
	var o order.Order
	var createdAt, updatedAt time.Time
	err := or.db.QueryRow(ctx, query, id).Scan(
		&o.Id, &o.GuestId, &o.UserId, &o.IsGuest, &o.TableNumber, &o.OrderHandlerId,
		&o.Status, &createdAt, &updatedAt, &o.TotalPrice,
	)
	if err != nil {
		or.logger.Error(fmt.Sprintf("Error fetching order detail for ID %d: %s", id, err.Error()))
		return nil, fmt.Errorf("error fetching order detail: %w", err)
	}
	o.CreatedAt = timestamppb.New(createdAt)
	o.UpdatedAt = timestamppb.New(updatedAt)

	dishItems, err := or.getDishItemsForOrder(ctx, o.Id)
	if err != nil {
		return nil, err
	}
	o.DishItems = dishItems

	setItems, err := or.getSetItemsForOrder(ctx, o.Id)
	if err != nil {
		return nil, err
	}
	o.SetItems = setItems

	or.logger.Info(fmt.Sprintf("Successfully fetched order detail for ID: %d", id))
	return &o, nil
}

func (or *OrderRepository) UpdateOrder(ctx context.Context, req *order.UpdateOrderRequest) (*order.Order, error) {
	or.logger.Info(fmt.Sprintf("Updating order with ID: %d", req.Id))
	tx, err := or.db.Begin(ctx)
	if err != nil {
		or.logger.Error("Error starting transaction: " + err.Error())
		return nil, fmt.Errorf("error starting transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	query := `
		UPDATE orders
		SET guest_id = $2, user_id = $3, table_number = $4, order_handler_id = $5, status = $6, updated_at = $7, total_price = $8, is_guest = $9
		WHERE id = $1
		RETURNING created_at, updated_at
	`
	var o order.Order
	var createdAt, updatedAt time.Time
	err = tx.QueryRow(ctx, query,
		req.Id, req.GuestId, req.UserId, req.TableNumber, req.OrderHandlerId, req.Status, time.Now(), req.TotalPrice, req.IsGuest,
	).Scan(&createdAt, &updatedAt)
	if err != nil {
		or.logger.Error(fmt.Sprintf("Error updating order with ID %d: %s", req.Id, err.Error()))
		return nil, fmt.Errorf("error updating order: %w", err)
	}

	o.Id = req.Id
	o.GuestId = req.GuestId
	o.UserId = req.UserId
	o.TableNumber = req.TableNumber
	o.OrderHandlerId = req.OrderHandlerId
	o.Status = req.Status
	o.CreatedAt = timestamppb.New(createdAt)
	o.UpdatedAt = timestamppb.New(updatedAt)
	o.TotalPrice = req.TotalPrice
	o.IsGuest = req.IsGuest

	// Update dish items
	err = or.updateOrderItems(ctx, tx, o.Id, req.DishItems, req.SetItems)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		or.logger.Error("Error committing transaction: " + err.Error())
		return nil, fmt.Errorf("error committing transaction: %w", err)
	}

	o.DishItems = req.DishItems
	o.SetItems = req.SetItems
	or.logger.Info(fmt.Sprintf("Successfully updated order with ID: %d", o.Id))
	return &o, nil
}

func (or *OrderRepository) PayOrders(ctx context.Context, req *order.PayOrdersRequest) ([]*order.Order, error) {
	or.logger.Info("Processing payment for orders")
	tx, err := or.db.Begin(ctx)
	if err != nil {
		or.logger.Error("Error starting transaction: " + err.Error())
		return nil, fmt.Errorf("error starting transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	var query string
	var arg interface{}
	if req.GetGuestId() != 0 {
		query = "UPDATE orders SET status = 'paid', updated_at = $1 WHERE guest_id = $2 AND status != 'paid' RETURNING id"
		arg = req.GetGuestId()
	} else if req.GetUserId() != 0 {
		query = "UPDATE orders SET status = 'paid', updated_at = $1 WHERE user_id = $2 AND status != 'paid' RETURNING id"
		arg = req.GetUserId()
	} else {
		return nil, fmt.Errorf("either guest_id or user_id must be provided")
	}

	rows, err := tx.Query(ctx, query, time.Now(), arg)
	if err != nil {
		or.logger.Error("Error updating orders: " + err.Error())
		return nil, fmt.Errorf("error updating orders: %w", err)
	}
	defer rows.Close()

	var orderIDs []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			or.logger.Error("Error scanning order ID: " + err.Error())
			return nil, fmt.Errorf("error scanning order ID: %w", err)
		}
		orderIDs = append(orderIDs, id)
	}

	if err := tx.Commit(ctx); err != nil {
		or.logger.Error("Error committing transaction: " + err.Error())
		return nil, fmt.Errorf("error committing transaction: %w", err)
	}

	var paidOrders []*order.Order
	for _, id := range orderIDs {
		o, err := or.GetOrderDetail(ctx, id)
		if err != nil {
			return nil, err
		}
		paidOrders = append(paidOrders, o)
	}

	or.logger.Info(fmt.Sprintf("Successfully processed payment for %d orders", len(paidOrders)))
	return paidOrders, nil
}

func (or *OrderRepository) insertDishOrderItem(ctx context.Context, tx pgx.Tx, orderID int64, item *order.DishOrderItem) error {
	query := `
		INSERT INTO order_dish_items (order_id, dish_id, quantity)
		VALUES ($1, $2, $3)
		RETURNING id
	`
	var itemID int64
	err := tx.QueryRow(ctx, query, orderID, item.Dish.Id, item.Quantity).Scan(&itemID)
	if err != nil {
		or.logger.Error(fmt.Sprintf("Error inserting dish order item: %s", err.Error()))
		return fmt.Errorf("error inserting dish order item: %w", err)
	}
	item.Id = itemID
	return nil
}

func (or *OrderRepository) insertSetOrderItem(ctx context.Context, tx pgx.Tx, orderID int64, item *order.SetOrderItem) error {
	query := `
		INSERT INTO order_set_items (order_id, set_id, quantity)
		VALUES ($1, $2, $3)
		RETURNING id
	`
	var itemID int64
	err := tx.QueryRow(ctx, query, orderID, item.Set.Id, item.Quantity).Scan(&itemID)
	if err != nil {
		or.logger.Error(fmt.Sprintf("Error inserting set order item: %s", err.Error()))
		return fmt.Errorf("error inserting set order item: %w", err)
	}
	item.Id = itemID
	return nil
}

func (or *OrderRepository) getDishItemsForOrder(ctx context.Context, orderID int64) ([]*order.DishOrderItem, error) {
	query := `
		SELECT odi.id, odi.quantity, d.id, d.name, d.price, d.description, d.image, d.status, d.created_at, d.updated_at
		FROM order_dish_items odi
		JOIN dishes d ON odi.dish_id = d.id
		WHERE odi.order_id = $1
	`
	rows, err := or.db.Query(ctx, query, orderID)
	if err != nil {
		return nil, fmt.Errorf("error fetching dish items for order: %w", err)


	}
	defer rows.Close()

	var items []*order.DishOrderItem
	for rows.Next() {
		var item order.DishOrderItem
		var dish order.DishOrder
		var createdAt, updatedAt time.Time
		err := rows.Scan(
			&item.Id, &item.Quantity,
			&dish.Id, &dish.Name, &dish.Price, &dish.Description, &dish.Image, &dish.Status,
			&createdAt, &updatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning dish order item: %w", err)
		}
		dish.CreatedAt = timestamppb.New(createdAt)
		dish.UpdatedAt = timestamppb.New(updatedAt)
		item.Dish = &dish
		items = append(items, &item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over dish items: %w", err)
	}

	return items, nil
}

func (or *OrderRepository) getSetItemsForOrder(ctx context.Context, orderID int64) ([]*order.SetOrderItem, error) {
	query := `
		SELECT osi.id, osi.quantity, s.id, s.name, s.description, s.user_id, s.created_at, s.updated_at, s.is_favourite, s.is_public, s.image
		FROM order_set_items osi
		JOIN sets s ON osi.set_id = s.id
		WHERE osi.order_id = $1
	`
	rows, err := or.db.Query(ctx, query, orderID)
	if err != nil {
		return nil, fmt.Errorf("error fetching set items for order: %w", err)
	}
	defer rows.Close()

	var items []*order.SetOrderItem
	for rows.Next() {
		var item order.SetOrderItem
		var set order.SetProto
		var createdAt, updatedAt time.Time
		var userID int64
		err := rows.Scan(
			&item.Id, &item.Quantity,
			&set.Id, &set.Name, &set.Description, &userID, &createdAt, &updatedAt,
			&set.IsFavourite, &set.IsPublic, &set.Image,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning set order item: %w", err)
		}
		set.CreatedAt = timestamppb.New(createdAt)
		set.UpdatedAt = timestamppb.New(updatedAt)
		set.UserId = userID
		
		// Fetch dishes for the set
		set.Dishes, err = or.getDishesForSet(ctx, set.Id)
		if err != nil {
			return nil, err
		}

		item.Set = &set
		items = append(items, &item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over set items: %w", err)
	}

	return items, nil
}

func (or *OrderRepository) getDishesForSet(ctx context.Context, setID int32) ([]*order.SetProtoDish, error) {
	query := `
		SELECT d.id, d.name, d.price, d.description, d.image, d.status, d.created_at, d.updated_at
		FROM set_dishes sd
		JOIN dishes d ON sd.dish_id = d.id
		WHERE sd.set_id = $1
	`
	rows, err := or.db.Query(ctx, query, setID)
	if err != nil {
		return nil, fmt.Errorf("error fetching dishes for set: %w", err)
	}
	defer rows.Close()

	var dishes []*order.SetProtoDish
	for rows.Next() {
		var dish order.SetProtoDish
		var createdAt, updatedAt time.Time
		err := rows.Scan(
			&dish.Id, &dish.Name, &dish.Price, &dish.Description, &dish.Image, &dish.Status,
			&createdAt, &updatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning dish for set: %w", err)
		}
		dish.CreatedAt = timestamppb.New(createdAt)
		dish.UpdatedAt = timestamppb.New(updatedAt)
		dishes = append(dishes, &dish)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over dishes for set: %w", err)
	}

	return dishes, nil
}

func (or *OrderRepository) updateOrderItems(ctx context.Context, tx pgx.Tx, orderID int64, dishItems []*order.DishOrderItem, setItems []*order.SetOrderItem) error {
	// Delete existing items
	_, err := tx.Exec(ctx, "DELETE FROM order_dish_items WHERE order_id = $1", orderID)
	if err != nil {
		return fmt.Errorf("error deleting existing dish items: %w", err)
	}
	_, err = tx.Exec(ctx, "DELETE FROM order_set_items WHERE order_id = $1", orderID)
	if err != nil {
		return fmt.Errorf("error deleting existing set items: %w", err)
	}

	// Insert new dish items
	for _, item := range dishItems {
		err := or.insertDishOrderItem(ctx, tx, orderID, item)
		if err != nil {
			return err
		}
	}

	// Insert new set items
	for _, item := range setItems {
		err := or.insertSetOrderItem(ctx, tx, orderID, item)
		if err != nil {
			return err
		}
	}

	return nil
}