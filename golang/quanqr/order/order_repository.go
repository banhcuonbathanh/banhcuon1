package order_grpc

import (
	"context"
	"database/sql"
	"fmt"
	"time"

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

    query := `
        INSERT INTO orders (
            guest_id, user_id, is_guest, table_number, order_handler_id,
            status, created_at, updated_at, total_price, bow_chili, bow_no_chili
        )
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
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

    now := time.Now()
    err = tx.QueryRow(ctx, query,
        guestId,
        userId,
        req.IsGuest,
        req.TableNumber,
        req.OrderHandlerId,
        req.Status,
        now,          // created_at
        now,          // updated_at
        req.TotalPrice,
        req.BowChili,
        req.BowNoChili,
    ).Scan(&o.Id, &createdAt, &updatedAt)

    if err != nil {
        or.logger.Error("Error creating order: " + err.Error())
        return nil, fmt.Errorf("error creating order: %w", err)
    }

    // Verify dishes exist before inserting
    for _, dish := range req.DishItems {
        // First verify the dish exists
        var exists bool
        err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM dishes WHERE id = $1)", dish.DishId).Scan(&exists)
        if err != nil {
            or.logger.Error(fmt.Sprintf("Error verifying dish existence: %s", err.Error()))
            return nil, fmt.Errorf("error verifying dish existence: %w", err)
        }
        if !exists {
            or.logger.Error(fmt.Sprintf("Dish with id %d does not exist", dish.DishId))
            return nil, fmt.Errorf("dish with id %d does not exist", dish.DishId)
        }

        // Then insert the order item
        _, err = tx.Exec(ctx, 
            "INSERT INTO dish_order_items (order_id, dish_id, quantity) VALUES ($1, $2, $3)",
            o.Id, dish.DishId, dish.Quantity)
        if err != nil {
            or.logger.Error(fmt.Sprintf("Error inserting order dish: %s", err.Error()))
            return nil, fmt.Errorf("error inserting order dish: %w", err)
        }
    }

    // Verify sets exist before inserting
    for _, set := range req.SetItems {
        // First verify the set exists
        var exists bool
        err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM sets WHERE id = $1)", set.SetId).Scan(&exists)
        if err != nil {
            or.logger.Error(fmt.Sprintf("Error verifying set existence: %s", err.Error()))
            return nil, fmt.Errorf("error verifying set existence: %w", err)
        }
        if !exists {
            or.logger.Error(fmt.Sprintf("Set with id %d does not exist", set.SetId))
            return nil, fmt.Errorf("set with id %d does not exist", set.SetId)
        }

        // Then insert the set item
        _, err = tx.Exec(ctx, 
            "INSERT INTO set_order_items (order_id, set_id, quantity) VALUES ($1, $2, $3)",
            o.Id, set.SetId, set.Quantity)
        if err != nil {
            or.logger.Error(fmt.Sprintf("Error inserting order set: %s", err.Error()))
            return nil, fmt.Errorf("error inserting order set: %w", err)
        }
    }

    if err := tx.Commit(ctx); err != nil {
        or.logger.Error("Error committing transaction: " + err.Error())
        return nil, fmt.Errorf("error committing transaction: %w", err)
    }

    // Populate response
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
    o.BowChili = req.BowChili
    o.BowNoChili = req.BowNoChili

    return &o, nil
}
func (or *OrderRepository) GetOrders(ctx context.Context, req *order.GetOrdersRequest) ([]*order.Order, error) {
    or.logger.Info("Fetching orders with filters")
    
    query := `
        SELECT 
            id, guest_id, user_id, is_guest, table_number, order_handler_id,
            status, created_at, updated_at, total_price, bow_chili, bow_no_chili
        FROM orders
        WHERE ($1::timestamp IS NULL OR created_at >= $1)
        AND ($2::timestamp IS NULL OR created_at <= $2)
        AND ($3::bigint IS NULL OR user_id = $3)
        AND ($4::bigint IS NULL OR guest_id = $4)
    `

    rows, err := or.db.Query(ctx, query, 
        req.FromDate.AsTime(),
        req.ToDate.AsTime(),
        req.UserId,
        req.GuestId)
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
            &o.Id,
            &o.GuestId,
            &o.UserId,
            &o.IsGuest,
            &o.TableNumber,
            &o.OrderHandlerId,
            &o.Status,
            &createdAt,
            &updatedAt,
            &o.TotalPrice,
            &o.BowChili,
            &o.BowNoChili,
        )
        if err != nil {
            or.logger.Error("Error scanning order: " + err.Error())
            return nil, fmt.Errorf("error scanning order: %w", err)
        }

        o.CreatedAt = timestamppb.New(createdAt)
        o.UpdatedAt = timestamppb.New(updatedAt)

        // Get dish items
        dishItems, err := or.GetOrderDishItems(ctx, o.Id)
        if err != nil {
            return nil, err
        }
        o.DishItems = dishItems

        // Get set items
        setItems, err := or.GetOrderSetItems(ctx, o.Id)
        if err != nil {
            return nil, err
        }
        o.SetItems = setItems

        orders = append(orders, &o)
    }

    return orders, nil
}

func (or *OrderRepository) GetOrderDetail(ctx context.Context, id int64) (*order.Order, error) {
    or.logger.Info(fmt.Sprintf("Fetching order detail for ID: %d", id))
    
    query := `
        SELECT 
            id, guest_id, user_id, is_guest, table_number, order_handler_id,
            status, created_at, updated_at, total_price, bow_chili, bow_no_chili
        FROM orders
        WHERE id = $1
    `

    var o order.Order
    var createdAt, updatedAt time.Time

    err := or.db.QueryRow(ctx, query, id).Scan(
        &o.Id,
        &o.GuestId,
        &o.UserId,
        &o.IsGuest,
        &o.TableNumber,
        &o.OrderHandlerId,
        &o.Status,
        &createdAt,
        &updatedAt,
        &o.TotalPrice,
        &o.BowChili,
        &o.BowNoChili,
    )
    if err != nil {
        or.logger.Error(fmt.Sprintf("Error fetching order detail: %s", err.Error()))
        return nil, fmt.Errorf("error fetching order detail: %w", err)
    }

    o.CreatedAt = timestamppb.New(createdAt)
    o.UpdatedAt = timestamppb.New(updatedAt)

    // Get dish items
    dishItems, err := or.GetOrderDishItems(ctx, o.Id)
    if err != nil {
        return nil, err
    }
    o.DishItems = dishItems

    // Get set items
    setItems, err := or.GetOrderSetItems(ctx, o.Id)
    if err != nil {
        return nil, err
    }
    o.SetItems = setItems

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
        SET guest_id = $2, user_id = $3, table_number = $4, order_handler_id = $5,
            status = $6, updated_at = $7, total_price = $8, is_guest = $9,
            bow_chili = $10, bow_no_chili = $11
        WHERE id = $1
        RETURNING created_at, updated_at
    `

    var o order.Order
    var createdAt, updatedAt time.Time

    err = tx.QueryRow(ctx, query,
        req.Id,
        req.GuestId,
        req.UserId,
        req.TableNumber,
        req.OrderHandlerId,
        req.Status,
        time.Now(),
        req.TotalPrice,
        req.IsGuest,
        req.BowChili,
        req.BowNoChili,
    ).Scan(&createdAt, &updatedAt)

    if err != nil {
        or.logger.Error(fmt.Sprintf("Error updating order: %s", err.Error()))
        return nil, fmt.Errorf("error updating order: %w", err)
    }

    // Update dish items
    _, err = tx.Exec(ctx, "DELETE FROM order_dishes WHERE order_id = $1", req.Id)
    if err != nil {
        return nil, fmt.Errorf("error deleting order dishes: %w", err)
    }

    for _, dish := range req.DishItems {
        _, err := tx.Exec(ctx, 
            "INSERT INTO order_dishes (order_id, dish_id, quantity) VALUES ($1, $2, $3)",
            req.Id, dish.DishId, dish.Quantity)
        if err != nil {
            return nil, fmt.Errorf("error inserting order dish: %w", err)
        }
    }

    // Update set items
    _, err = tx.Exec(ctx, "DELETE FROM order_sets WHERE order_id = $1", req.Id)
    if err != nil {
        return nil, fmt.Errorf("error deleting order sets: %w", err)
    }

    for _, set := range req.SetItems {
        _, err := tx.Exec(ctx, 
            "INSERT INTO order_sets (order_id, set_id, quantity) VALUES ($1, $2, $3)",
            req.Id, set.SetId, set.Quantity)
        if err != nil {
            return nil, fmt.Errorf("error inserting order set: %w", err)
        }
    }

    if err := tx.Commit(ctx); err != nil {
        or.logger.Error("Error committing transaction: " + err.Error())
        return nil, fmt.Errorf("error committing transaction: %w", err)
    }

    // Populate response
    o.Id = req.Id
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
    o.BowChili = req.BowChili
    o.BowNoChili = req.BowNoChili

    return &o, nil
}

func (or *OrderRepository) PayOrders(ctx context.Context, req *order.PayOrdersRequest) ([]*order.Order, error) {
    or.logger.Info("Processing payment for orders")
    
    var userIDFilter, guestIDFilter interface{}
    if req.Identifier != nil {
        switch v := req.Identifier.(type) {
        case *order.PayOrdersRequest_UserId:
            userIDFilter = v.UserId
        case *order.PayOrdersRequest_GuestId:
            guestIDFilter = v.GuestId
        }
    }

    tx, err := or.db.Begin(ctx)
    if err != nil {
        return nil, fmt.Errorf("error starting transaction: %w", err)
    }
    defer tx.Rollback(ctx)

    query := `
        UPDATE orders
        SET status = 'paid', updated_at = $1
        WHERE (user_id = $2 OR $2 IS NULL)
        AND (guest_id = $3 OR $3 IS NULL)
        AND status = 'pending'
        RETURNING id
    `

    rows, err := tx.Query(ctx, query, time.Now(), userIDFilter, guestIDFilter)
    if err != nil {
        return nil, fmt.Errorf("error updating orders: %w", err)
    }
    defer rows.Close()

    var orderIDs []int64
    for rows.Next() {
        var orderID int64
        if err := rows.Scan(&orderID); err != nil {
            return nil, fmt.Errorf("error scanning order ID: %w", err)
        }
        orderIDs = append(orderIDs, orderID)
    }

    if err := tx.Commit(ctx); err != nil {
        return nil, fmt.Errorf("error committing transaction: %w", err)
    }

    // Fetch updated orders
    var orders []*order.Order
    for _, orderID := range orderIDs {
        order, err := or.GetOrderDetail(ctx, orderID)
        if err != nil {
            return nil, err
        }
        orders = append(orders, order)
    }

    return orders, nil
}

func (or *OrderRepository) GetOrderDishItems(ctx context.Context, orderID int64) ([]*order.DishOrderItem, error) {
    query := `
        SELECT dish_id, quantity
        FROM order_dishes
        WHERE order_id = $1
    `
    rows, err := or.db.Query(ctx, query, orderID)
    if err != nil {
        return nil, fmt.Errorf("error fetching order dish items: %w", err)
    }
    defer rows.Close()

    var items []*order.DishOrderItem
    for rows.Next() {
        item := &order.DishOrderItem{}
        if err := rows.Scan(&item.DishId, &item.Quantity); err != nil {
            or.logger.Error(fmt.Sprintf("Error scanning order dish item: %s", err.Error()))
            return nil, fmt.Errorf("error scanning order dish item: %w", err)
        }
        items = append(items, item)
    }

    if err = rows.Err(); err != nil {
        or.logger.Error(fmt.Sprintf("Error iterating order dish items: %s", err.Error()))
        return nil, fmt.Errorf("error iterating order dish items: %w", err)
    }

    return items, nil
}

func (or *OrderRepository) GetOrderSetItems(ctx context.Context, orderID int64) ([]*order.SetOrderItem, error) {
    query := `
        SELECT set_id, quantity
        FROM order_sets
        WHERE order_id = $1
    `
    rows, err := or.db.Query(ctx, query, orderID)
    if err != nil {
        or.logger.Error(fmt.Sprintf("Error fetching order set items: %s", err.Error()))
        return nil, fmt.Errorf("error fetching order set items: %w", err)
    }
    defer rows.Close()

    var items []*order.SetOrderItem
    for rows.Next() {
        item := &order.SetOrderItem{}
        if err := rows.Scan(&item.SetId, &item.Quantity); err != nil {
            or.logger.Error(fmt.Sprintf("Error scanning order set item: %s", err.Error()))
            return nil, fmt.Errorf("error scanning order set item: %w", err)
        }
        items = append(items, item)
    }

    if err = rows.Err(); err != nil {
        or.logger.Error(fmt.Sprintf("Error iterating order set items: %s", err.Error()))
        return nil, fmt.Errorf("error iterating order set items: %w", err)
    }

    return items, nil
}