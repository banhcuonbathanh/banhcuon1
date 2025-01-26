package order_grpc

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"sort"

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

// ----------------------------------

func (or *OrderRepository) GetOrderProtoListDetail(ctx context.Context, page, pageSize int32) (*order.OrderDetailedListResponse, error) {
    // or.logger.Info("Fetching detailed order list with pagination")
    
    // Get total count for pagination
    countQuery := `SELECT COUNT(*) FROM orders`
    
    var totalItems int64
    err := or.db.QueryRow(ctx, countQuery).Scan(&totalItems)
    if err != nil {
        or.logger.Error("Error counting orders: " + err.Error())
        return nil, fmt.Errorf("error counting orders: %w", err)
    }

    // Calculate pagination info
    totalPages := int32(math.Ceil(float64(totalItems) / float64(pageSize)))
    offset := (page - 1) * pageSize

    // Main order query
    query := `
        SELECT 
            o.id, 
            o.guest_id, 
            o.user_id, 
            o.is_guest,
            o.table_number, 
            o.order_handler_id,
            COALESCE(o.status, 'Pending') as status, 
            o.total_price,
            COALESCE(o.topping, '') as topping,
            COALESCE(o.tracking_order, '') as tracking_order,
            COALESCE(o.take_away, false) as take_away,
            COALESCE(o.chili_number, 0) as chili_number,
              o.table_token,
            COALESCE(o.order_name, '') as order_name
        FROM orders o
        ORDER BY o.created_at DESC
        LIMIT $1 OFFSET $2
    `

    rows, err := or.db.Query(ctx, query, pageSize, offset)
    if err != nil {
        or.logger.Error("Error fetching orders: " + err.Error())
        return nil, fmt.Errorf("error fetching orders: %w", err)
    }
    defer rows.Close()

    var detailedOrders []*order.OrderDetailedResponse
    for rows.Next() {
        var o order.OrderDetailedResponse
        
        // Create nullable variables for fields that can be NULL
        var (
            guestId        sql.NullInt64
            userId         sql.NullInt64
            tableNumber    sql.NullInt64
            orderHandlerId sql.NullInt64
            totalPrice     sql.NullInt32
            status         sql.NullString
            topping       sql.NullString
            trackingOrder     sql.NullString
            chiliNumber    sql.NullInt64
            orderName      sql.NullString
        )

        err := rows.Scan(
            &o.Id,
            &guestId,
            &userId,
            &o.IsGuest,
            &tableNumber,
            &orderHandlerId,
            &status,
            &totalPrice,
            &topping,
            &trackingOrder,
            &o.TakeAway,
            &chiliNumber,
            &o.TableToken,
            &orderName,
        )
        if err != nil {
            or.logger.Error("Error scanning order: " + err.Error())
            return nil, fmt.Errorf("error scanning order: %w", err)
        }

        // Handle NULL values
        if guestId.Valid {
            o.GuestId = guestId.Int64
        }
        if userId.Valid {
            o.UserId = userId.Int64
        }
        if tableNumber.Valid {
            o.TableNumber = tableNumber.Int64
        }
        if orderHandlerId.Valid {
            o.OrderHandlerId = orderHandlerId.Int64
        }
        if totalPrice.Valid {
            o.TotalPrice = totalPrice.Int32
        }
        if status.Valid {
            o.Status = status.String
        }
        if topping.Valid {
            o.Topping = topping.String
        }
        if trackingOrder.Valid {
            o.TrackingOrder = trackingOrder.String
        }
        if chiliNumber.Valid {
            o.ChiliNumber = chiliNumber.Int64
        }
        if orderName.Valid {
            o.OrderName = orderName.String
        }
        // Fetch detailed dish items
        dishQuery := `
            SELECT 
                d.id,
                doi.quantity,
                d.name,
                d.price,
                d.description,
                d.image,
                d.status
            FROM dish_order_items doi
            JOIN dishes d ON doi.dish_id = d.id
            WHERE doi.order_id = $1
        `
        dishRows, err := or.db.Query(ctx, dishQuery, o.Id)
        if err != nil {
            or.logger.Error("Error fetching dish details: " + err.Error())
            return nil, fmt.Errorf("error fetching dish details: %w", err)
        }
        defer dishRows.Close()

        var dishItems []*order.OrderDetailedDish
        for dishRows.Next() {
            var dish order.OrderDetailedDish
            err := dishRows.Scan(
                &dish.DishId,
                &dish.Quantity,
                &dish.Name,
                &dish.Price,
                &dish.Description,
                &dish.Image,
                &dish.Status,
            )
            if err != nil {
                or.logger.Error("Error scanning dish detail: " + err.Error())
                return nil, fmt.Errorf("error scanning dish detail: %w", err)
            }
            dishItems = append(dishItems, &dish)
        }
        o.DataDish = dishItems

        // Fetch detailed set items
        setQuery := `
            SELECT 
                s.id,
                s.name,
                s.description,
                s.user_id,
                s.is_favourite,
                s.is_public,
                s.image,
                s.price,
                soi.quantity,
                s.created_at,
                s.updated_at
            FROM set_order_items soi
            JOIN sets s ON soi.set_id = s.id
            WHERE soi.order_id = $1
        `
        setRows, err := or.db.Query(ctx, setQuery, o.Id)
        if err != nil {
            or.logger.Error("Error fetching set details: " + err.Error())
            return nil, fmt.Errorf("error fetching set details: %w", err)
        }
        defer setRows.Close()

        var setItems []*order.OrderSetDetailed
        for setRows.Next() {
            var set order.OrderSetDetailed
            var createdAt, updatedAt time.Time
            var userID sql.NullInt32
            
            err := setRows.Scan(
                &set.Id,
                &set.Name,
                &set.Description,
                &userID,
                &set.IsFavourite,
                &set.IsPublic,
                &set.Image,
                &set.Price,
                &set.Quantity,
                &createdAt,
                &updatedAt,
            )
            if err != nil {
                or.logger.Error("Error scanning set detail: " + err.Error())
                return nil, fmt.Errorf("error scanning set detail: %w", err)
            }

            if userID.Valid {
                set.UserId = userID.Int32
            }
            
            set.CreatedAt = timestamppb.New(createdAt)
            set.UpdatedAt = timestamppb.New(updatedAt)

            // Fetch dishes for this set
            setDishQuery := `
                SELECT 
                    d.id,
                    sd.quantity,
                    d.name,
                    d.price,
                    d.description,
                    d.image,
                    d.status
                FROM set_dishes sd
                JOIN dishes d ON sd.dish_id = d.id
                WHERE sd.set_id = $1
            `
            setDishRows, err := or.db.Query(ctx, setDishQuery, set.Id)
            if err != nil {
                or.logger.Error("Error fetching set dish details: " + err.Error())
                return nil, fmt.Errorf("error fetching set dish details: %w", err)
            }
            defer setDishRows.Close()

            var setDishes []*order.OrderDetailedDish
            for setDishRows.Next() {
                var dish order.OrderDetailedDish
                err := setDishRows.Scan(
                    &dish.DishId,
                    &dish.Quantity,
                    &dish.Name,
                    &dish.Price,
                    &dish.Description,
                    &dish.Image,
                    &dish.Status,
                )
                if err != nil {
                    or.logger.Error("Error scanning set dish detail: " + err.Error())
                    return nil, fmt.Errorf("error scanning set dish detail: %w", err)
                }
                setDishes = append(setDishes, &dish)
            }
            set.Dishes = setDishes
            setItems = append(setItems, &set)
        }
        o.DataSet = setItems

        detailedOrders = append(detailedOrders, &o)
    }

    response := &order.OrderDetailedListResponse{
        Data: detailedOrders,
        Pagination: &order.PaginationInfo{
            CurrentPage: page,
            TotalPages: totalPages,
            TotalItems: totalItems,
            PageSize:   pageSize,
        },
    }

    return response, nil
}





// ---------------------------





func (or *OrderRepository) GetOrders(ctx context.Context, page, pageSize int32) ([]*order.Order, int64, error) {
    or.logger.Info("Fetching orders with pagination")
    
    // Get total count for pagination
    countQuery := `SELECT COUNT(*) FROM orders`
    
    var totalItems int64
    err := or.db.QueryRow(ctx, countQuery).Scan(&totalItems)
    if err != nil {
        or.logger.Error("Error counting orders: " + err.Error())
        return nil, 0, fmt.Errorf("error counting orders: %w", err)
    }

    // Calculate offset
    offset := (page - 1) * pageSize
    
    // Main order query
    query := `
        SELECT 
            o.id, 
            o.guest_id, 
            o.user_id, 
            o.is_guest, 
            o.table_number, 
            o.order_handler_id,
            COALESCE(o.status, 'Pending') as status, 
            o.created_at, 
            o.updated_at, 
            o.total_price, 
            COALESCE(o.topping, '') as topping, 
            COALESCE(o.tracking_order, '') as tracking_order,
            COALESCE(o.take_away, false) as take_away, 
            COALESCE(o.chili_number, 0) as chili_number,
            o.table_token,
            COALESCE(o.order_name, '') as order_name
        FROM orders o
        ORDER BY o.created_at DESC
        LIMIT $1 OFFSET $2
    `

    rows, err := or.db.Query(ctx, query, pageSize, offset)
    if err != nil {
        or.logger.Error("Error fetching orders: " + err.Error())
        return nil, 0, fmt.Errorf("error fetching orders: %w", err)
    }
    defer rows.Close()

    var orders []*order.Order
    for rows.Next() {
        var o order.Order
        var createdAt, updatedAt time.Time

        // Create nullable variables for fields that can be NULL in the database
        var (
            guestId        sql.NullInt64
            userId         sql.NullInt64
            tableNumber    sql.NullInt64
            orderHandlerId sql.NullInt64
            totalPrice     sql.NullInt32
            status         sql.NullString
            topping       sql.NullString
            trackingOrder     sql.NullString
            chiliNumber    sql.NullInt64
            orderName      sql.NullString
        )
        err = rows.Scan(
            &o.Id,
            &guestId,
            &userId,
            &o.IsGuest,
            &tableNumber,
            &orderHandlerId,
            &status,
            &createdAt,
            &updatedAt,
            &totalPrice,
            &topping,
            &trackingOrder,
            &o.TakeAway,
            &chiliNumber,
            &o.TableToken,
            &orderName,
        )
        if err != nil {
            or.logger.Error("Error scanning order: " + err.Error())
            return nil, 0, fmt.Errorf("error scanning order: %w", err)
        }

        // Convert nullable fields
        o.GuestId = guestId.Int64
        o.UserId = userId.Int64
        if tableNumber.Valid {
            o.TableNumber = tableNumber.Int64
        }
        o.OrderHandlerId = orderHandlerId.Int64
        o.Status = status.String
        o.TotalPrice = totalPrice.Int32
        o.Topping = topping.String
        o.TrackingOrder = trackingOrder.String
        o.ChiliNumber = chiliNumber.Int64
        if orderName.Valid {
            o.OrderName = orderName.String
        }

        // Handle timestamps
        o.CreatedAt = timestamppb.New(createdAt)
        o.UpdatedAt = timestamppb.New(updatedAt)

        orders = append(orders, &o)
    }

    return orders, totalItems, nil
}


//--------------


func (or *OrderRepository) GetOrderDetail(ctx context.Context, id int64) (*order.Order, error) {
    // or.logger.Info(fmt.Sprintf("Fetching order detail for ID: %d", id))
    
    query := `
        SELECT 
            id, guest_id, user_id, is_guest, table_number, order_handler_id,
            status, created_at, updated_at, total_price, topping, tracking_order,
            take_away, chili_number, table_token, order_name
        FROM orders
        WHERE id = $1
    `

    var o order.Order
    var createdAt, updatedAt time.Time
    var orderName sql.NullString

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
        &o.Topping,
        &o.TrackingOrder,
        &o.TakeAway,
        &o.ChiliNumber,
        &o.TableToken,
        &orderName,
    )
    if err != nil {
        or.logger.Error(fmt.Sprintf("Error fetching order detail: %s", err.Error()))
        return nil, fmt.Errorf("error fetching order detail: %w", err)
    }

    if orderName.Valid {
        o.OrderName = orderName.String
    }

    o.CreatedAt = timestamppb.New(createdAt)
    o.UpdatedAt = timestamppb.New(updatedAt)

    return &o, nil
}






func getOrdinalSuffix(n int64) string {
    if n%100 >= 11 && n%100 <= 13 {
        return "th"
    }
    switch n % 10 {
    case 1:
        return "st"
    case 2:
        return "nd"
    case 3:
        return "rd"
    default:
        return "th"
    }
}





// 






// Add this function to your OrderRepository struct

func (or *OrderRepository) FetchOrdersByCriteria(ctx context.Context, req *order.FetchOrdersByCriteriaRequest) (*order.OrderDetailedListResponse, error) {
    or.logger.Info("Fetching orders by criteria repository golang/quanqr/order/order_repository.go")

    // Build dynamic conditions and params
    var conditions []string
    var params []interface{}
    paramCount := 1

    // Add order IDs condition if provided
    if len(req.OrderIds) > 0 {
        conditions = append(conditions, fmt.Sprintf("o.id = ANY($%d)", paramCount))
        params = append(params, req.OrderIds)
        paramCount++
    }

    // Add order name condition if provided
    if req.OrderName != "" {
        conditions = append(conditions, fmt.Sprintf("o.order_name ILIKE $%d", paramCount))
        params = append(params, "%"+req.OrderName+"%")
        paramCount++
    }

    // Add date range conditions if provided
    if req.StartDate != nil {
        conditions = append(conditions, fmt.Sprintf("o.created_at >= $%d", paramCount))
        params = append(params, req.StartDate.AsTime())
        paramCount++
    }
    if req.EndDate != nil {
        conditions = append(conditions, fmt.Sprintf("o.created_at <= $%d", paramCount))
        params = append(params, req.EndDate.AsTime())
        paramCount++
    }

    // Build WHERE clause
    whereClause := "WHERE 1=1"
    for _, condition := range conditions {
        whereClause += " AND " + condition
    }

    // Calculate pagination
    offset := (req.Page - 1) * req.PageSize

    // Combined query for both count and data
    query := fmt.Sprintf(`
        WITH filtered_orders AS (
            SELECT 
                o.id, 
                o.guest_id, 
                o.user_id, 
                o.is_guest,
                o.table_number, 
                o.order_handler_id,
                COALESCE(o.status, 'Pending') as status, 
                o.total_price,
                COALESCE(o.topping, '') as topping,
                COALESCE(o.tracking_order, '') as tracking_order,
                COALESCE(o.take_away, false) as take_away,
                COALESCE(o.chili_number, 0) as chili_number,
                o.table_token,
                COALESCE(o.order_name, '') as order_name,
                COUNT(*) OVER() as total_count
            FROM orders o
            %s
            ORDER BY o.created_at DESC
            LIMIT $%d OFFSET $%d
        )
        SELECT * FROM filtered_orders`,
        whereClause,
        paramCount,
        paramCount+1,
    )

    // Add pagination parameters
    params = append(params, req.PageSize, offset)

    // Execute the query
    rows, err := or.db.Query(ctx, query, params...)
    if err != nil {
        or.logger.Error("Error executing fetch orders query: " + err.Error())
        return nil, fmt.Errorf("error executing fetch orders query: %w", err)
    }
    defer rows.Close()

    var detailedOrders []*order.OrderDetailedResponse
    var totalItems int64
    for rows.Next() {
        var o order.OrderDetailedResponse
        var (
            guestId        sql.NullInt64
            userId         sql.NullInt64
            tableNumber    sql.NullInt64
            orderHandlerId sql.NullInt64
            totalPrice     sql.NullInt32
            status         sql.NullString
            topping        sql.NullString
            trackingOrder  sql.NullString
            chiliNumber    sql.NullInt64
            orderName      sql.NullString
        )

        err := rows.Scan(
            &o.Id,
            &guestId,
            &userId,
            &o.IsGuest,
            &tableNumber,
            &orderHandlerId,
            &status,
            &totalPrice,
            &topping,
            &trackingOrder,
            &o.TakeAway,
            &chiliNumber,
            &o.TableToken,
            &orderName,
            &totalItems,
        )
        if err != nil {
            or.logger.Error("Error scanning order: " + err.Error())
            return nil, fmt.Errorf("error scanning order: %w", err)
        }

        // Handle NULL values
        if guestId.Valid {
            o.GuestId = guestId.Int64
        }
        if userId.Valid {
            o.UserId = userId.Int64
        }
        if tableNumber.Valid {
            o.TableNumber = tableNumber.Int64
        }
        if orderHandlerId.Valid {
            o.OrderHandlerId = orderHandlerId.Int64
        }
        if totalPrice.Valid {
            o.TotalPrice = totalPrice.Int32
        }
        if status.Valid {
            o.Status = status.String
        }
        if topping.Valid {
            o.Topping = topping.String
        }
        if trackingOrder.Valid {
            o.TrackingOrder = trackingOrder.String
        }
        if chiliNumber.Valid {
            o.ChiliNumber = chiliNumber.Int64
        }
        if orderName.Valid {
            o.OrderName = orderName.String
        }

        // Fetch dish items for this order
        dishItems, err := or.getOrderDishDetails(ctx, o.Id)
        if err != nil {
            return nil, err
        }
        o.DataDish = dishItems

        // Fetch set items for this order
        setItems, err := or.getOrderSetDetails(ctx, o.Id)
        if err != nil {
            return nil, err
        }
        o.DataSet = setItems

        detailedOrders = append(detailedOrders, &o)
    }

    totalPages := int32(math.Ceil(float64(totalItems) / float64(req.PageSize)))

    return &order.OrderDetailedListResponse{
        Data: detailedOrders,
        Pagination: &order.PaginationInfo{
            CurrentPage: req.Page,
            TotalPages: totalPages,
            TotalItems: totalItems,
            PageSize:   req.PageSize,
        },
    }, nil
}

// Helper function to get order dish details
func (or *OrderRepository) getOrderDishDetails(ctx context.Context, orderID int64) ([]*order.OrderDetailedDish, error) {
    query := `
        SELECT 
            d.id,
            doi.quantity,
            d.name,
            d.price,
            d.description,
            d.image,
            d.status
        FROM dish_order_items doi
        JOIN dishes d ON doi.dish_id = d.id
        WHERE doi.order_id = $1
    `
    
    rows, err := or.db.Query(ctx, query, orderID)
    if err != nil {
        return nil, fmt.Errorf("error fetching dish details: %w", err)
    }
    defer rows.Close()

    var dishes []*order.OrderDetailedDish
    for rows.Next() {
        var dish order.OrderDetailedDish
        err := rows.Scan(
            &dish.DishId,
            &dish.Quantity,
            &dish.Name,
            &dish.Price,
            &dish.Description,
            &dish.Image,
            &dish.Status,
        )
        if err != nil {
            return nil, fmt.Errorf("error scanning dish detail: %w", err)
        }
        dishes = append(dishes, &dish)
    }

    return dishes, nil
}

// Helper function to get order set details
func (or *OrderRepository) getOrderSetDetails(ctx context.Context, orderID int64) ([]*order.OrderSetDetailed, error) {
    query := `
        SELECT 
            s.id,
            s.name,
            s.description,
            s.user_id,
            s.is_favourite,
            s.is_public,
            s.image,
            s.price,
            soi.quantity,
            s.created_at,
            s.updated_at
        FROM set_order_items soi
        JOIN sets s ON soi.set_id = s.id
        WHERE soi.order_id = $1
    `
    
    rows, err := or.db.Query(ctx, query, orderID)
    if err != nil {
        return nil, fmt.Errorf("error fetching set details: %w", err)
    }
    defer rows.Close()

    var sets []*order.OrderSetDetailed
    for rows.Next() {
        var set order.OrderSetDetailed
        var createdAt, updatedAt time.Time
        var userID sql.NullInt32
        
        err := rows.Scan(
            &set.Id,
            &set.Name,
            &set.Description,
            &userID,
            &set.IsFavourite,
            &set.IsPublic,
            &set.Image,
            &set.Price,
            &set.Quantity,
            &createdAt,
            &updatedAt,
        )
        if err != nil {
            return nil, fmt.Errorf("error scanning set detail: %w", err)
        }

        if userID.Valid {
            set.UserId = userID.Int32
        }
        
        set.CreatedAt = timestamppb.New(createdAt)
        set.UpdatedAt = timestamppb.New(updatedAt)

        // Get dishes for this set
        dishes, err := or.getSetDishes(ctx, set.Id)
        if err != nil {
            return nil, err
        }
        set.Dishes = dishes

        sets = append(sets, &set)
    }

    return sets, nil
}

// Helper function to get set dishes
func (or *OrderRepository) getSetDishes(ctx context.Context, setID int64) ([]*order.OrderDetailedDish, error) {
    query := `
        SELECT 
            d.id,
            sd.quantity,
            d.name,
            d.price,
            d.description,
            d.image,
            d.status
        FROM set_dishes sd
        JOIN dishes d ON sd.dish_id = d.id
        WHERE sd.set_id = $1
    `
    
    rows, err := or.db.Query(ctx, query, setID)
    if err != nil {
        return nil, fmt.Errorf("error fetching set dish details: %w", err)
    }
    defer rows.Close()

    var dishes []*order.OrderDetailedDish
    for rows.Next() {
        var dish order.OrderDetailedDish
        err := rows.Scan(
            &dish.DishId,
            &dish.Quantity,
            &dish.Name,
            &dish.Price,
            &dish.Description,
            &dish.Image,
            &dish.Status,
        )
        if err != nil {
            return nil, fmt.Errorf("error scanning set dish detail: %w", err)
        }
        dishes = append(dishes, &dish)
    }

    return dishes, nil
}

// -------------------------------------------------- create order start -----------------------

func (or *OrderRepository) CreateOrder(ctx context.Context, req *order.CreateOrderRequest) (*order.Order, error) {
    or.logger.Info(fmt.Sprintf("Creating new order: %+v", req))
    
    tx, err := or.db.Begin(ctx)
    if err != nil {
        or.logger.Error("Error starting transaction: " + err.Error())
        return nil, fmt.Errorf("error starting transaction: %w", err)
    }
    defer tx.Rollback(ctx)

    now := time.Now()
    
    // Get client ID based on whether it's a guest or user
    clientId := req.GuestId
    if !req.IsGuest {
        clientId = req.UserId
    }
    or.logger.Info(fmt.Sprintf("Creating order for clientId: %d, isGuest: %v", clientId, req.IsGuest))
    
    // Get tracking order information
    trackingOrder, err := or.getTrackingOrderInfo(ctx, or.db, now, req.IsGuest, clientId)
    if err != nil {
        or.logger.Error("Error getting tracking order info: " + err.Error())
        return nil, fmt.Errorf("error getting tracking order info: %w", err)
    }

    // Get client number for the order name
    clientNumber, err := or.getClientNumberForDay(ctx, or.db, now, req.IsGuest, clientId)
    if err != nil {
        or.logger.Error(fmt.Sprintf("Error getting client number: %v", err))
        return nil, fmt.Errorf("error getting client number: %w", err)
    }
    or.logger.Info(fmt.Sprintf("Got client number: %d", clientNumber))

    var guestId, userId sql.NullInt64
    if req.IsGuest {
        guestId = sql.NullInt64{Int64: req.GuestId, Valid: true}
        userId = sql.NullInt64{Valid: false}
    } else {
        userId = sql.NullInt64{Int64: req.UserId, Valid: true}
        guestId = sql.NullInt64{Valid: false}
    }

    // Updated insert query with new fields
    query := `
        INSERT INTO orders (
            guest_id, user_id, is_guest, table_number, order_handler_id,
            status, created_at, updated_at, total_price, topping, tracking_order,
            take_away, chili_number, table_token, order_name, version, parent_order_id
        )
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)
        RETURNING id, created_at, updated_at
    `
    or.logger.Info("Executing order insertion")
    var o order.Order
    var createdAt, updatedAt time.Time

    // Set initial version to 1 and parent_order_id to null for new orders
    err = tx.QueryRow(ctx, query,
        guestId,
        userId,
        req.IsGuest,
        req.TableNumber,
        req.OrderHandlerId,
        req.Status,
        now,
        now,
        req.TotalPrice,
        req.Topping,
        trackingOrder,
        req.TakeAway,
        req.ChiliNumber,
        req.TableToken,
        req.OrderName,
        1,  // Initial version
        nil, // No parent order for new orders
    ).Scan(&o.Id, &createdAt, &updatedAt)

    if err != nil {
        or.logger.Error("Error creating order: " + err.Error())
        return nil, fmt.Errorf("error creating order: %w", err)
    }

    // Create initial order modification record
    modificationQuery := `
        INSERT INTO order_modifications (
            order_id, modification_number, modification_type, 
            modified_by_user_id, order_name
        )
        VALUES ($1, $2, $3, $4, $5)
    `
    _, err = tx.Exec(ctx, modificationQuery,
        o.Id,
        1, // First modification
        "INITIAL",
        req.OrderHandlerId,
        req.OrderName,
    )
    if err != nil {
        or.logger.Error("Error creating order modification record: " + err.Error())
        return nil, fmt.Errorf("error creating order modification record: %w", err)
    }

    // Insert dish items with new fields
    for _, dish := range req.DishItems {
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

        _, err = tx.Exec(ctx, `
            INSERT INTO dish_order_items (
                order_id, dish_id, quantity, created_at, updated_at, 
                order_name, modification_type, modification_number
            )
            VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
            o.Id, dish.DishId, dish.Quantity, now, now, 
            req.OrderName, "INITIAL", 1)
        if err != nil {
            or.logger.Error(fmt.Sprintf("Error inserting order dish: %s", err.Error()))
            return nil, fmt.Errorf("error inserting order dish: %w", err)
        }

        // Create initial delivery record for each dish
        _, err = tx.Exec(ctx, `
            INSERT INTO dish_deliveries (
                order_id, order_name, quantity_delivered,
                delivery_status, modification_number
            )
            VALUES ($1, $2, $3, $4, $5)`,
            o.Id, req.OrderName, dish.Quantity, "PENDING", 1)
        if err != nil {
            or.logger.Error(fmt.Sprintf("Error creating dish delivery record: %s", err.Error()))
            return nil, fmt.Errorf("error creating dish delivery record: %w", err)
        }
    }

    // Insert set items with new fields
    for _, set := range req.SetItems {
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

        _, err = tx.Exec(ctx, `
            INSERT INTO set_order_items (
                order_id, set_id, quantity, created_at, updated_at, 
                order_name, modification_type, modification_number
            )
            VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
            o.Id, set.SetId, set.Quantity, now, now, 
            req.OrderName, "INITIAL", 1)
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
    o.Topping = req.Topping
    o.TrackingOrder = trackingOrder
    o.TakeAway = req.TakeAway
    o.ChiliNumber = req.ChiliNumber
    o.TableToken = req.TableToken
    o.OrderName = req.OrderName
    o.Version = 1 // Set initial version

    return &o, nil
}


func (or *OrderRepository) getTrackingOrderInfo(ctx context.Context, tx *pgxpool.Pool, currentTime time.Time, isGuest bool, clientId int64) (string, error) {
    or.logger.Info("golang/quanqr/order/order_repository.go getTrackingOrderInfo - Starting")
    
    startOfDay := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, time.UTC)
    endOfDay := startOfDay.Add(24 * time.Hour)
    
    or.logger.Info(fmt.Sprintf("golang/quanqr/order/order_repository.go getTrackingOrderInfo - Processing for clientId: %d, isGuest: %v, date range: %v to %v", 
        clientId, isGuest, startOfDay, endOfDay))

    clientPosition, err := or.getClientPosition(ctx, tx, startOfDay, endOfDay, clientId, isGuest)
    if err != nil {
        or.logger.Error(fmt.Sprintf("golang/quanqr/order/order_repository.go getTrackingOrderInfo - Error getting client position: %v", err))
        return "", fmt.Errorf("error getting client position: %w", err)
    }
    or.logger.Info(fmt.Sprintf("golang/quanqr/order/order_repository.go getTrackingOrderInfo - Got client position: %d", clientPosition))

    orderCount, err := or.getClientOrderCount(ctx, tx, startOfDay, endOfDay, clientId, isGuest)
    if err != nil {
        or.logger.Error(fmt.Sprintf("golang/quanqr/order/order_repository.go getTrackingOrderInfo - Error getting client order count: %v", err))
        return "", fmt.Errorf("error getting client order count: %w", err)
    }
    or.logger.Info(fmt.Sprintf("golang/quanqr/order/order_repository.go getTrackingOrderInfo - Got order count: %d", orderCount))

    var clientType string
    if isGuest {
        clientType = "Guest"
    } else {
        clientType = "Client"
    }

    trackingOrder := fmt.Sprintf("%d%s %s - Order #%d", 
        clientPosition,
        getOrdinalSuffix(clientPosition),
        clientType,
        orderCount)

    or.logger.Info(fmt.Sprintf("golang/quanqr/order/order_repository.go getTrackingOrderInfo - Generated tracking order: %s", trackingOrder))
    return trackingOrder, nil
}


func (or *OrderRepository) getClientNumberForDay(ctx context.Context, tx *pgxpool.Pool, currentTime time.Time, isGuest bool, clientId int64) (int64, error) {
    or.logger.Info(fmt.Sprintf("golang/quanqr/order/order_repository.go getClientNumberForDay - Starting for clientId: %d, isGuest: %v", clientId, isGuest))
    
    startOfDay := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, time.UTC)
    endOfDay := startOfDay.Add(24 * time.Hour)
    or.logger.Info(fmt.Sprintf("golang/quanqr/order/order_repository.go getClientNumberForDay - Date range: %v to %v", startOfDay, endOfDay))

    var query string
    var args []interface{}

    if isGuest {
        query = `
            SELECT COALESCE(
                (SELECT COUNT(*) 
                FROM orders 
                WHERE guest_id = $1 
                AND created_at >= $2 
                AND created_at < $3),
                0
            ) + 1
        `
        args = []interface{}{clientId, startOfDay, endOfDay}
    } else {
        query = `
            SELECT COALESCE(
                (SELECT COUNT(*) 
                FROM orders 
                WHERE user_id = $1 
                AND created_at >= $2 
                AND created_at < $3),
                0
            ) + 1
        `
        args = []interface{}{clientId, startOfDay, endOfDay}
    }

    var clientNumber int64
    err := tx.QueryRow(ctx, query, args...).Scan(&clientNumber)
    if err != nil {
        if err.Error() == "no rows in result set" {
            or.logger.Info("golang/quanqr/order/order_repository.go getClientNumberForDay - No existing client number found, returning 1")
            return 1, nil
        }
        or.logger.Error(fmt.Sprintf("golang/quanqr/order/order_repository.go getClientNumberForDay - Error getting client number: %v", err))
        return 0, fmt.Errorf("error getting client number: %w", err)
    }
    or.logger.Info(fmt.Sprintf("golang/quanqr/order/order_repository.go getClientNumberForDay - clientNumber %v ", clientNumber))
    return clientNumber, nil
}
func (or *OrderRepository) getClientPosition(ctx context.Context, tx *pgxpool.Pool, startOfDay, endOfDay time.Time, clientId int64, isGuest bool) (int64, error) {
    or.logger.Info(fmt.Sprintf("golang/quanqr/order/order_repository.go getClientPosition - Starting for clientId: %d, isGuest: %v", clientId, isGuest))
    or.logger.Info(fmt.Sprintf("golang/quanqr/order/order_repository.go getClientPosition - Date range: %v to %v", startOfDay, endOfDay))
    
    // First, let's debug the virtual table
    debugQuery := `
        WITH ordered_clients AS (
            SELECT 
                CASE 
                    WHEN guest_id IS NOT NULL THEN guest_id 
                    ELSE user_id 
                END as client_id,
                MIN(created_at) as first_order_time,
                DENSE_RANK() OVER (ORDER BY MIN(created_at)) as position
            FROM orders
            WHERE created_at >= $1 AND created_at < $2
            GROUP BY 
                CASE 
                    WHEN guest_id IS NOT NULL THEN guest_id 
                    ELSE user_id 
                END
        )
        SELECT client_id, first_order_time, position 
        FROM ordered_clients
        ORDER BY position;
    `
    
    // Execute debug query to see the virtual table
    rows, err := tx.Query(ctx, debugQuery, startOfDay, endOfDay)
    if err != nil {
        or.logger.Error(fmt.Sprintf("Debug query error: %v", err))
    } else {
        defer rows.Close()
        or.logger.Info("Virtual Table Content:")
        or.logger.Info("| Client ID | First Order Time | Position |")
        or.logger.Info("|-----------|-----------------|-----------|")
        
        for rows.Next() {
            var cid int64
            var orderTime time.Time
            var pos int64
            if err := rows.Scan(&cid, &orderTime, &pos); err != nil {
                or.logger.Error(fmt.Sprintf("Error scanning debug row: %v", err))
                continue
            }
            or.logger.Info(fmt.Sprintf("| %9d | %s | %9d |", cid, orderTime.Format("15:04:05"), pos))
        }
    }

    // Now execute the actual position query
    query := `
        WITH ordered_clients AS (
            SELECT 
                CASE 
                    WHEN guest_id IS NOT NULL THEN guest_id 
                    ELSE user_id 
                END as client_id,
                MIN(created_at) as first_order_time,
                DENSE_RANK() OVER (ORDER BY MIN(created_at)) as position
            FROM orders
            WHERE created_at >= $1 AND created_at < $2
            GROUP BY 
                CASE 
                    WHEN guest_id IS NOT NULL THEN guest_id 
                    ELSE user_id 
                END
        )
        SELECT position
        FROM ordered_clients
        WHERE client_id = $3
        UNION ALL
        SELECT COALESCE(
            (SELECT MAX(position) + 1 
            FROM ordered_clients),
            1
        )
        WHERE NOT EXISTS (
            SELECT 1 
            FROM ordered_clients 
            WHERE client_id = $3
        )
        LIMIT 1
    `

    var position int64
    err = tx.QueryRow(ctx, query, startOfDay, endOfDay, clientId).Scan(&position)
    if err != nil {
        return 0, fmt.Errorf("error getting client position: %w", err)
    }
    
    or.logger.Info(fmt.Sprintf("golang/quanqr/order/order_repository.go getClientPosition - Retrieved position for client %d: %d", clientId, position))
    return position, nil
}

func (or *OrderRepository) getClientOrderCount(ctx context.Context, tx *pgxpool.Pool, startOfDay, endOfDay time.Time, clientId int64, isGuest bool) (int64, error) {
    or.logger.Info(fmt.Sprintf("golang/quanqr/order/order_repository.go getClientOrderCount - Starting for clientId: %d, isGuest: %v", clientId, isGuest))
    or.logger.Info(fmt.Sprintf("golang/quanqr/order/order_repository.go getClientOrderCount - Date range: %v to %v", startOfDay, endOfDay))
    var query string
    if isGuest {
        query = `
            SELECT COUNT(*) + 1
            FROM orders
            WHERE guest_id = $1
            AND created_at >= $2 AND created_at < $3
        `
    } else {
        query = `
            SELECT COUNT(*) + 1
            FROM orders
            WHERE user_id = $1
            AND created_at >= $2 AND created_at < $3
        `
    }

    var count int64
    err := tx.QueryRow(ctx, query, clientId, startOfDay, endOfDay).Scan(&count)
    if err != nil {
        if err.Error() == "no rows in result set" {
            return 1, nil
        }
        return 0, fmt.Errorf("error getting order count: %w", err)
    }
    or.logger.Info(fmt.Sprintf("golang/quanqr/order/order_repository.go getClientOrderCount - Retrieved count: %d", count))
    return count, nil
}

// -------------------------------------------------- create order end  -----------------------

// -------------------------------------------------- update order start  -----------------------

func (or *OrderRepository) UpdateOrder(ctx context.Context, req *order.UpdateOrderRequest) (*order.OrderDetailedListResponse, error) {
    or.logger.Info(fmt.Sprintf("Updating order with ID: %d", req.Id))
    
    // Begin transaction using pgxpool
    tx, err := or.db.Begin(ctx)
    if err != nil {
        or.logger.Error("Error starting transaction: " + err.Error())
        return nil, fmt.Errorf("error starting transaction: %w", err)
    }
    defer tx.Rollback(ctx)

    // Version control check
    var currentVersion int32
    err = tx.QueryRow(ctx, "SELECT version FROM orders WHERE id = $1", req.Id).Scan(&currentVersion)
    if err != nil {
        or.logger.Error(fmt.Sprintf("Error fetching order version: %s", err.Error()))
        return nil, fmt.Errorf("error fetching order version: %w", err)
    }

    if currentVersion != req.Version {
        return nil, fmt.Errorf("order version mismatch: expected %d, got %d", currentVersion, req.Version)
    }

    // Update order with incremented version
    newVersion := currentVersion + 1
    query := `
        UPDATE orders
        SET guest_id = $2, user_id = $3, table_number = $4, order_handler_id = $5,
            status = $6, updated_at = $7, total_price = $8, is_guest = $9,
            topping = $10, tracking_order = $11, take_away = $12, 
            chili_number = $13, table_token = $14, order_name = $15,
            version = $16
        WHERE id = $1
        RETURNING created_at, updated_at
    `

    var createdAt, updatedAt time.Time
    err = tx.QueryRow(ctx, query,
        req.Id, req.GuestId, req.UserId, req.TableNumber, req.OrderHandlerId,
        req.Status, time.Now(), req.TotalPrice, req.IsGuest,
        req.Topping, req.TrackingOrder, req.TakeAway,
        req.ChiliNumber, req.TableToken, req.OrderName, newVersion,
    ).Scan(&createdAt, &updatedAt)

    if err != nil {
        or.logger.Error(fmt.Sprintf("Error updating order: %s", err.Error()))
        return nil, fmt.Errorf("error updating order: %w", err)
    }

    // Create modification record
    _, err = tx.Exec(ctx, `
        INSERT INTO order_modifications (
            order_id, modification_number, modification_type, 
            modified_by_user_id, order_name
        ) VALUES ($1, $2, $3, $4, $5)
    `,
        req.Id, newVersion, "UPDATE", req.OrderHandlerId, req.OrderName,
    )
    if err != nil {
        or.logger.Error(fmt.Sprintf("Error creating modification record: %s", err.Error()))
        return nil, fmt.Errorf("error creating modification record: %w", err)
    }

    // Handle dish items update
    if len(req.DishItems) > 0 {
        _, err = tx.Exec(ctx, "DELETE FROM dish_order_items WHERE order_id = $1", req.Id)
        if err != nil {
            return nil, fmt.Errorf("error deleting existing dish items: %w", err)
        }

        for _, item := range req.DishItems {
            _, err = tx.Exec(ctx, `
                INSERT INTO dish_order_items (
                    order_id, dish_id, quantity, order_name,
                    modification_type, modification_number
                ) VALUES ($1, $2, $3, $4, $5, $6)
            `,
                req.Id, item.DishId, item.Quantity, item.OrderName,
                "UPDATE", newVersion,
            )
            if err != nil {
                return nil, fmt.Errorf("error inserting updated dish item: %w", err)
            }
        }
    }

    // Handle set items update
    if len(req.SetItems) > 0 {
        _, err = tx.Exec(ctx, "DELETE FROM set_order_items WHERE order_id = $1", req.Id)
        if err != nil {
            return nil, fmt.Errorf("error deleting existing set items: %w", err)
        }

        for _, item := range req.SetItems {
            _, err = tx.Exec(ctx, `
                INSERT INTO set_order_items (
                    order_id, set_id, quantity, order_name,
                    modification_type, modification_number
                ) VALUES ($1, $2, $3, $4, $5, $6)
            `,
                req.Id, item.SetId, item.Quantity, item.OrderName,
                "UPDATE", newVersion,
            )
            if err != nil {
                return nil, fmt.Errorf("error inserting updated set item: %w", err)
            }
        }
    }

    // Fetch detailed information
    detailedDishes, err := or.fetchDetailedDishes(ctx, tx, req.Id)
    if err != nil {
        return nil, fmt.Errorf("error fetching detailed dishes: %w", err)
    }

    detailedSets, err := or.fetchDetailedSets(ctx, tx, req.Id)
    if err != nil {
        return nil, fmt.Errorf("error fetching detailed sets: %w", err)
    }

    if err := tx.Commit(ctx); err != nil {
        or.logger.Error("Error committing transaction: " + err.Error())
        return nil, fmt.Errorf("error committing transaction: %w", err)
    }

    return &order.OrderDetailedListResponse{
        Data: []*order.OrderDetailedResponse{
            {
                Id:             req.Id,
                GuestId:        req.GuestId,
                UserId:         req.UserId,
                TableNumber:    req.TableNumber,
                OrderHandlerId: req.OrderHandlerId,
                Status:         req.Status,
                TotalPrice:     req.TotalPrice,
                DataSet:        detailedSets,
                DataDish:       detailedDishes,
                IsGuest:        req.IsGuest,
                Topping:        req.Topping,
                TrackingOrder:  req.TrackingOrder,
                TakeAway:       req.TakeAway,
                ChiliNumber:    req.ChiliNumber,
                TableToken:     req.TableToken,
                OrderName:      req.OrderName,
                CurrentVersion:        newVersion,
                ParentOrderId:  req.ParentOrderId,
            },
        },
        Pagination: &order.PaginationInfo{
            CurrentPage: 1,
            TotalPages:  1,
            TotalItems:  1,
            PageSize:    1,
        },
    }, nil
}


// Similarly, update the fetchDetailedDishes function
func (or *OrderRepository) fetchDetailedDishes(ctx context.Context, tx pgx.Tx, orderID int64) ([]*order.OrderDetailedDish, error) {
    // We'll use a CTE (Common Table Expression) to ensure we get the latest modifications
    query := `
        WITH latest_modifications AS (
            SELECT dish_id, quantity
            FROM dish_order_items doi
            WHERE doi.order_id = $1
            AND doi.modification_number = (
                SELECT MAX(modification_number)
                FROM dish_order_items
                WHERE order_id = $1 AND dish_id = doi.dish_id
            )
        )
        SELECT d.id, lm.quantity, d.name, d.price, d.description, d.image, d.status
        FROM latest_modifications lm
        JOIN dishes d ON lm.dish_id = d.id
    `
    
    // Add context timeout to prevent long-running queries
    queryCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
    defer cancel()
    
    rows, err := tx.Query(queryCtx, query, orderID)
    if err != nil {
        return nil, fmt.Errorf("error querying order dishes (order_id=%d): %w", orderID, err)
    }
    defer rows.Close()

    // Pre-allocate the slice with a reasonable capacity
    dishes := make([]*order.OrderDetailedDish, 0, 10)
    
    for rows.Next() {
        dish := &order.OrderDetailedDish{}
        err := rows.Scan(
            &dish.DishId,
            &dish.Quantity,
            &dish.Name,
            &dish.Price,
            &dish.Description,
            &dish.Image,
            &dish.Status,
        )
        if err != nil {
            // Include more context in error message
            return nil, fmt.Errorf("error scanning dish row (order_id=%d): %w", orderID, err)
        }
        dishes = append(dishes, dish)
    }

    // Check for errors from iteration
    if err = rows.Err(); err != nil {
        return nil, fmt.Errorf("error iterating dish rows (order_id=%d): %w", orderID, err)
    }

    // Validate we got at least some results
    if len(dishes) == 0 {
        or.logger.Info(fmt.Sprintf("No dishes found for order_id=%d", orderID))
    }

    return dishes, nil
}
// -------------------------------------------------- update order end  -----------------------

// -------------------------------------------------- update ordder adding set and dishes start  -----------------------
// func (or *OrderRepository) AddingSetsDishesOrder(ctx context.Context, req *order.UpdateOrderRequest) (*order.OrderDetailedListResponse, error) {
//     // Log entry into function with order ID and request details
//     or.logger.Info(fmt.Sprintf("[OrderRepository.AddingSetsDishesOrder] Starting order update process. OrderID: %d, UserID: %d, TableNumber: %d",
//         req.Id, req.UserId, req.TableNumber))
    
//     tx, err := or.db.Begin(ctx)
//     if err != nil {
//         or.logger.Error(fmt.Sprintf("[OrderRepository.AddingSetsDishesOrder] Transaction start failed: %s", err.Error()))
//         return nil, fmt.Errorf("error starting transaction: %w", err)
//     }
//     defer tx.Rollback(ctx)

//     // Version check logging
//     or.logger.Info(fmt.Sprintf("[OrderRepository.AddingSetsDishesOrder] Checking version for order %d. Requested version: %d",
//         req.Id, req.Version))

//     var currentVersion int32
//     var isGuest bool
//     err = tx.QueryRow(ctx, `
//         SELECT version, is_guest 
//         FROM orders 
//         WHERE id = $1`, req.Id).Scan(&currentVersion, &isGuest)
//     if err != nil {
//         if err == pgx.ErrNoRows {
//             or.logger.Error(fmt.Sprintf("[OrderRepository.AddingSetsDishesOrder] Order not found. ID: %d", req.Id))
//             return nil, fmt.Errorf("order not found with ID: %d", req.Id)
//         }
//         or.logger.Error(fmt.Sprintf("[OrderRepository.AddingSetsDishesOrder] Error fetching order details: %s", err.Error()))
//         return nil, fmt.Errorf("error fetching order: %w", err)
//     }

//     // Log version check results
//     or.logger.Info(fmt.Sprintf("[OrderRepository.AddingSetsDishesOrder] Current version: %d, Is Guest: %v", 
//         currentVersion, isGuest))

//     if req.Version != 0 && req.Version != currentVersion {
//         or.logger.Warning(fmt.Sprintf("[OrderRepository.AddingSetsDishesOrder] Version mismatch. Expected: %d, Got: %d",
//             currentVersion, req.Version))
//         return nil, fmt.Errorf("order version mismatch: expected %d, got %d", currentVersion, req.Version)
//     }

//     newVersion := currentVersion + 1
//     or.logger.Info(fmt.Sprintf("[OrderRepository.AddingSetsDishesOrder] Incrementing version to: %d", newVersion))

//     // Fetch version history
//     versionHistory, err := or.fetchVersionHistory(ctx, tx, req.Id)
//     if err != nil {
//         or.logger.Error(fmt.Sprintf("[OrderRepository.AddingSetsDishesOrder] Error fetching version history: %s", err.Error()))
//         return nil, fmt.Errorf("error fetching version history: %w", err)
//     }

//     // Calculate summary
//     totalSummary, err := or.calculateTotalSummary(ctx, tx, req.Id)
//     if err != nil {
//         or.logger.Error(fmt.Sprintf("[OrderRepository.AddingSetsDishesOrder] Error calculating total summary: %s", err.Error()))
//         return nil, fmt.Errorf("error calculating total summary: %w", err)
//     }

//     // Fetch detailed items
//     or.logger.Info(fmt.Sprintf("[OrderRepository.AddingSetsDishesOrder] Fetching detailed items for order %d", req.Id))
    
//     detailedDishes, err := or.fetchDetailedDishes(ctx, tx, req.Id)
//     if err != nil {
//         or.logger.Error(fmt.Sprintf("[OrderRepository.AddingSetsDishesOrder] Error fetching detailed dishes: %s", err.Error()))
//         return nil, fmt.Errorf("error fetching detailed dishes: %w", err)
//     }

//     detailedSets, err := or.fetchDetailedSets(ctx, tx, req.Id)
//     if err != nil {
//         or.logger.Error(fmt.Sprintf("[OrderRepository.AddingSetsDishesOrder] Error fetching detailed sets: %s", err.Error()))
//         return nil, fmt.Errorf("error fetching detailed sets: %w", err)
//     }

//     if err := tx.Commit(ctx); err != nil {
//         or.logger.Error(fmt.Sprintf("[OrderRepository.AddingSetsDishesOrder] Transaction commit failed: %s", err.Error()))
//         return nil, fmt.Errorf("error committing transaction: %w", err)
//     }

//     or.logger.Info(fmt.Sprintf("[OrderRepository.AddingSetsDishesOrder] 121212 Returning order details - CurrentVersion: %d, VersionHistory: %+v, TotalSummary: %+v", 
//     newVersion, 
//     versionHistory, 
//     totalSummary))
//     return &order.OrderDetailedListResponse{
//         Data: []*order.OrderDetailedResponse{
//             {
//                 Id:             req.Id,
//                 GuestId:        req.GuestId,
//                 UserId:         req.UserId,
//                 TableNumber:    req.TableNumber,
//                 OrderHandlerId: req.OrderHandlerId,
//                 Status:         req.Status,
//                 TotalPrice:     req.TotalPrice,
//                 DataSet:        detailedSets,
//                 DataDish:       detailedDishes,
//                 IsGuest:        req.IsGuest,
//                 Topping:        req.Topping,
//                 TrackingOrder:  req.TrackingOrder,
//                 TakeAway:       req.TakeAway,
//                 ChiliNumber:    req.ChiliNumber,
//                 TableToken:     req.TableToken,
//                 OrderName:      req.OrderName,
           
//                 ParentOrderId:  req.ParentOrderId,
//                 // New fields for version tracking
//                 CurrentVersion: newVersion,
//                 VersionHistory: versionHistory,
//                 TotalSummary:  totalSummary,
//             },
//         },
//         Pagination: &order.PaginationInfo{
//             CurrentPage: 1,
//             TotalPages:  1,
//             TotalItems:  1,
//             PageSize:    1,
//         },
//     }, nil
// }




// -------------------------------------------------- update ordder adding set and dishes end -----------------------


// Improved version of getVersionChanges to include price information
func (or *OrderRepository) getVersionChanges(ctx context.Context, tx pgx.Tx, orderID int64, version int32) ([]*order.OrderItemChange, error) {
    rows, err := tx.Query(ctx, `
        WITH version_changes AS (
            -- Get dish changes
            SELECT 
                'DISH' as item_type,
                d.id as item_id,
                d.name as item_name,
                di.quantity as quantity_changed,
                d.price
            FROM dish_order_items di
            JOIN dishes d ON di.dish_id = d.id
            WHERE di.order_id = $1 AND di.modification_number = $2
            
            UNION ALL
            
            -- Get set changes
            SELECT 
                'SET' as item_type,
                s.id as item_id,
                s.name as item_name,
                si.quantity as quantity_changed,
                s.price
            FROM set_order_items si
            JOIN sets s ON si.set_id = s.id
            WHERE si.order_id = $1 AND si.modification_number = $2
        )
        SELECT 
            item_type,
            item_id,
            item_name,
            quantity_changed,
            price
        FROM version_changes
        ORDER BY item_type, item_name`,
        orderID, version)
    if err != nil {
        return nil, fmt.Errorf("error querying version changes: %w", err)
    }
    defer rows.Close()

    var changes []*order.OrderItemChange
    for rows.Next() {
        var change order.OrderItemChange
        err := rows.Scan(
            &change.ItemType,
            &change.ItemId,
            &change.ItemName,
            &change.QuantityChanged,
            &change.Price,
        )
        if err != nil {
            return nil, fmt.Errorf("error scanning version change: %w", err)
        }
        changes = append(changes, &change)
    }

    return changes, nil
}

// new -------------


func (or *OrderRepository) fetchVersionHistory(ctx context.Context, tx pgx.Tx, orderID int64) ([]*order.OrderVersionSummary, error) {
    or.logger.Info(fmt.Sprintf("[OrderRepository.fetchVersionHistory] Fetching version history for order %d", orderID))
    
    // Combined query to get version summaries and changes in one go
    rows, err := tx.Query(ctx, `
        WITH version_summary AS (
            SELECT 
                m.modification_number,
                m.modification_type,
                m.modified_at,
                -- Count dishes in this version
                (SELECT COUNT(*) 
                 FROM dish_order_items di 
                 WHERE di.order_id = m.order_id 
                 AND di.modification_number = m.modification_number) as dishes_count,
                -- Count sets in this version
                (SELECT COUNT(*) 
                 FROM set_order_items si 
                 WHERE si.order_id = m.order_id 
                 AND si.modification_number = m.modification_number) as sets_count,
                -- Calculate version total price
                COALESCE(
                    (SELECT SUM(d.price * di.quantity)
                     FROM dish_order_items di
                     JOIN dishes d ON di.dish_id = d.id
                     WHERE di.order_id = m.order_id 
                     AND di.modification_number = m.modification_number), 0
                ) +
                COALESCE(
                    (SELECT SUM(s.price * si.quantity)
                     FROM set_order_items si
                     JOIN sets s ON si.set_id = s.id
                     WHERE si.order_id = m.order_id 
                     AND si.modification_number = m.modification_number), 0
                ) as version_total_price
            FROM order_modifications m
            WHERE m.order_id = $1
        ),
        version_changes AS (
            -- Get dish changes
            SELECT 
                modification_number,
                'DISH' as item_type,
                d.id as item_id,
                d.name as item_name,
                di.quantity as quantity_changed,
                d.price
            FROM dish_order_items di
            JOIN dishes d ON di.dish_id = d.id
            WHERE di.order_id = $1
            
            UNION ALL
            
            -- Get set changes
            SELECT 
                modification_number,
                'SET' as item_type,
                s.id as item_id,
                s.name as item_name,
                si.quantity as quantity_changed,
                s.price
            FROM set_order_items si
            JOIN sets s ON si.set_id = s.id
            WHERE si.order_id = $1
        )
        SELECT 
            vs.modification_number,
            vs.modification_type,
            vs.modified_at,
            vs.dishes_count,
            vs.sets_count,
            vs.version_total_price,
            vc.item_type,
            vc.item_id,
            vc.item_name,
            vc.quantity_changed,
            vc.price
        FROM version_summary vs
        LEFT JOIN version_changes vc ON vs.modification_number = vc.modification_number
        ORDER BY vs.modification_number ASC, vc.item_type, vc.item_name`,
        orderID)
    if err != nil {
        or.logger.Error(fmt.Sprintf("[OrderRepository.fetchVersionHistory] Query error: %s", err.Error()))
        return nil, fmt.Errorf("error querying version history: %w", err)
    }
    defer rows.Close()

    // Map to store version summaries
    summaryMap := make(map[int32]*order.OrderVersionSummary)
    
    // Scan rows and build version summaries with changes
    for rows.Next() {
        var (
            versionNum       int32
            modType         string
            modifiedAt      time.Time
            dishesCount     int32
            setsCount       int32
            versionTotal    int32
            // Change fields (can be null)
            itemType        sql.NullString
            itemID         sql.NullInt64
            itemName       sql.NullString
            quantityChanged sql.NullInt32
            price          sql.NullInt32
        )
        
        err := rows.Scan(
            &versionNum,
            &modType,
            &modifiedAt,
            &dishesCount,
            &setsCount,
            &versionTotal,
            &itemType,
            &itemID,
            &itemName,
            &quantityChanged,
            &price,
        )
        if err != nil {
            or.logger.Error(fmt.Sprintf("[OrderRepository.fetchVersionHistory] Scan error: %s", err.Error()))
            return nil, fmt.Errorf("error scanning version history: %w", err)
        }

        // Get or create version summary
        summary, exists := summaryMap[versionNum]
        if !exists {
            summary = &order.OrderVersionSummary{
                VersionNumber:     versionNum,
                ModificationType: modType,
                ModifiedAt:       timestamppb.New(modifiedAt),
                TotalDishesCount: dishesCount,
                TotalSetsCount:   setsCount,
                VersionTotalPrice: versionTotal,
                Changes:          make([]*order.OrderItemChange, 0),
            }
            summaryMap[versionNum] = summary
        }

        // Add change if present
        if itemType.Valid {
            change := &order.OrderItemChange{
                ItemType:       itemType.String,
                ItemId:        itemID.Int64,
                ItemName:      itemName.String,
                QuantityChanged: quantityChanged.Int32,
                Price:         price.Int32,
            }
            summary.Changes = append(summary.Changes, change)
        }
    }

    // Convert map to sorted slice
    var summaries []*order.OrderVersionSummary
    for _, summary := range summaryMap {
        summaries = append(summaries, summary)
    }
    sort.Slice(summaries, func(i, j int) bool {
        return summaries[i].VersionNumber < summaries[j].VersionNumber
    })

    or.logger.Info(fmt.Sprintf("[OrderRepository.fetchVersionHistory] Found %d versions", len(summaries)))
    return summaries, nil
}


// new -----

func (or *OrderRepository) calculateTotalSummary(ctx context.Context, tx pgx.Tx, orderID int64) (*order.OrderTotalSummary, error) {
    or.logger.Info(fmt.Sprintf("[OrderRepository.calculateTotalSummary] Starting calculation for order %d", orderID))
    
    var summary order.OrderTotalSummary
    
    // Complex query to calculate all summary metrics across all versions
    err := tx.QueryRow(ctx, `
        WITH order_stats AS (
            SELECT COUNT(DISTINCT modification_number) as total_versions
            FROM order_modifications
            WHERE order_id = $1
        ),
        dish_totals AS (
            SELECT 
                COUNT(DISTINCT dish_id) as unique_dishes,
                SUM(quantity) as total_dishes,
                SUM(quantity * price) as dish_total_price
            FROM dish_order_items di
            JOIN dishes d ON di.dish_id = d.id
            WHERE di.order_id = $1
        ),
        set_totals AS (
            SELECT 
                COUNT(DISTINCT set_id) as unique_sets,
                SUM(quantity) as total_sets,
                SUM(quantity * price) as set_total_price
            FROM set_order_items si
            JOIN sets s ON si.set_id = s.id
            WHERE si.order_id = $1
        )
        SELECT 
            os.total_versions,
            COALESCE(dt.total_dishes, 0),
            COALESCE(st.total_sets, 0),
            COALESCE(dt.dish_total_price, 0) + COALESCE(st.set_total_price, 0) as cumulative_total
        FROM order_stats os
        LEFT JOIN dish_totals dt ON true
        LEFT JOIN set_totals st ON true`,
        orderID).Scan(
            &summary.TotalVersions,
            &summary.TotalDishesOrdered,
            &summary.TotalSetsOrdered,
            &summary.CumulativeTotalPrice,
        )
    if err != nil {
        if err == pgx.ErrNoRows {
            or.logger.Warning(fmt.Sprintf("[OrderRepository.calculateTotalSummary] No data found for order %d", orderID))
            return &order.OrderTotalSummary{}, nil
        }
        or.logger.Error(fmt.Sprintf("[OrderRepository.calculateTotalSummary] Error calculating summary: %s", err.Error()))
        return nil, fmt.Errorf("error calculating total summary: %w", err)
    }

    // Query for most ordered items (combining both dishes and sets)
    rows, err := tx.Query(ctx, `
        WITH combined_items AS (
            -- Dish totals across all versions
            SELECT 
                'DISH' as item_type,
                d.id as item_id,
                d.name as item_name,
                SUM(di.quantity) as total_quantity
            FROM dish_order_items di
            JOIN dishes d ON di.dish_id = d.id
            WHERE di.order_id = $1
            GROUP BY d.id, d.name
            
            UNION ALL
            
            -- Set totals across all versions
            SELECT 
                'SET' as item_type,
                s.id as item_id,
                s.name as item_name,
                SUM(si.quantity) as total_quantity
            FROM set_order_items si
            JOIN sets s ON si.set_id = s.id
            WHERE si.order_id = $1
            GROUP BY s.id, s.name
        )
        SELECT 
            item_type,
            item_id,
            item_name,
            total_quantity
        FROM combined_items
        ORDER BY total_quantity DESC
        LIMIT 5`, // Limiting to top 5 most ordered items
        orderID)
    if err != nil {
        or.logger.Error(fmt.Sprintf("[OrderRepository.calculateTotalSummary] Error querying most ordered items: %s", err.Error()))
        return nil, fmt.Errorf("error querying most ordered items: %w", err)
    }
    defer rows.Close()

    // Scan most ordered items
    for rows.Next() {
        var item order.OrderItemCount
        if err := rows.Scan(&item.ItemType, &item.ItemId, &item.ItemName, &item.TotalQuantity); err != nil {
            or.logger.Error(fmt.Sprintf("[OrderRepository.calculateTotalSummary] Error scanning most ordered item: %s", err.Error()))
            return nil, fmt.Errorf("error scanning most ordered item: %w", err)
        }
        summary.MostOrderedItems = append(summary.MostOrderedItems, &item)
    }

    or.logger.Info(fmt.Sprintf("[OrderRepository.calculateTotalSummary] Successfully calculated summary for order %d", orderID))
    return &summary, nil
}



func (or *OrderRepository) fetchDetailedSets(ctx context.Context, tx pgx.Tx, orderID int64) ([]*order.OrderSetDetailed, error) {
    // Start logging
    or.logger.Info(fmt.Sprintf("[fetchDetailedSets] Starting to fetch sets for order ID: %d", orderID))

    query := `
        WITH set_dishes AS (
            -- First CTE: Aggregate dishes for each set into a JSONB array
            SELECT sd.set_id,
                   jsonb_agg(
                       jsonb_build_object(
                           'dish_id', d.id,  -- Changed 'id' to 'dish_id' to match the struct
                           'name', d.name,
                           'description', d.description,
                           'price', d.price,
                           'image', d.image,
                           'quantity', sd.quantity
                       ) ORDER BY d.id  -- Added ordering for consistency
                   ) as dishes
            FROM set_dishes sd
            JOIN dishes d ON sd.dish_id = d.id
            GROUP BY sd.set_id
        ),
        unique_set_orders AS (
            -- Second CTE: Get unique set orders to prevent duplicates
            SELECT DISTINCT ON (soi.set_id) 
                   soi.set_id,
                   soi.quantity,
                   soi.order_id  -- Added for logging
            FROM set_order_items soi
            WHERE soi.order_id = $1
            ORDER BY soi.set_id, soi.created_at DESC  -- Added ordering for consistent results
        )
        SELECT s.id, 
               s.name, 
               s.description, 
               s.user_id, 
               s.created_at, 
               s.updated_at,
               s.is_favourite, 
               s.like_by, 
               s.is_public, 
               s.image, 
               s.price, 
               uso.quantity,
               COALESCE(sd.dishes, '[]'::jsonb) as dishes
        FROM unique_set_orders uso
        JOIN sets s ON uso.set_id = s.id
        LEFT JOIN set_dishes sd ON s.id = sd.set_id
        ORDER BY s.id  -- Added ordering for consistent results
    `
    
    // Log query execution
    or.logger.Info(fmt.Sprintf("[fetchDetailedSets] Executing query for order ID: %d", orderID))
    
    rows, err := tx.Query(ctx, query, orderID)
    if err != nil {
        or.logger.Error(fmt.Sprintf("[fetchDetailedSets] Failed to execute query: %v", err))
        return nil, fmt.Errorf("error querying set details: %w", err)
    }
    defer rows.Close()

    var sets []*order.OrderSetDetailed
    
    setCount := 0
    for rows.Next() {
        setCount++
        or.logger.Info(fmt.Sprintf("[fetchDetailedSets] Processing set #%d", setCount))
        
        set := &order.OrderSetDetailed{}
        var createdAt, updatedAt time.Time
        var likeBy []int64
        var dishesJSON []byte
        
        err := rows.Scan(
            &set.Id,
            &set.Name,
            &set.Description,
            &set.UserId,
            &createdAt,
            &updatedAt,
            &set.IsFavourite,
            &likeBy,
            &set.IsPublic,
            &set.Image,
            &set.Price,
            &set.Quantity,
            &dishesJSON,
        )
        if err != nil {
            or.logger.Error(fmt.Sprintf("[fetchDetailedSets] Error scanning set row: %v", err))
            return nil, fmt.Errorf("error scanning set row: %w", err)
        }
        
        // Log set details before processing dishes
        or.logger.Info(fmt.Sprintf("[fetchDetailedSets] Successfully scanned set ID: %d, Name: %s, Quantity: %d", 
            set.Id, set.Name, set.Quantity))

        // Log raw dishes JSON for debugging

        
        var dishes []*order.OrderDetailedDish
        if err := json.Unmarshal(dishesJSON, &dishes); err != nil {
            or.logger.Error(fmt.Sprintf("[fetchDetailedSets] Error unmarshaling dishes for set %d: %v", set.Id, err))
            return nil, fmt.Errorf("error unmarshaling dishes: %w", err)
        }
        
        // Log dishes count and details
        or.logger.Info(fmt.Sprintf("[fetchDetailedSets] Set %d contains %d dishes", set.Id, len(dishes)))

        
        set.CreatedAt = timestamppb.New(createdAt)
        set.UpdatedAt = timestamppb.New(updatedAt)
        set.LikeBy = likeBy
        set.Dishes = dishes
        
        sets = append(sets, set)
    }

    if err = rows.Err(); err != nil {
        or.logger.Error(fmt.Sprintf("[fetchDetailedSets] Error during row iteration: %v", err))
        return nil, fmt.Errorf("error iterating set rows: %w", err)
    }

    or.logger.Info(fmt.Sprintf("[fetchDetailedSets] Successfully fetched %d sets for order %d", len(sets), orderID))

    // Log final result summary
    for _, set := range sets {
        or.logger.Info(fmt.Sprintf("[fetchDetailedSets] Final set summary - ID: %d, Name: %s, Dishes: %d, Quantity: %d", 
            set.Id, set.Name, len(set.Dishes), set.Quantity))
    }

    return sets, nil
}


// new -----------------------------------------------------
func (or *OrderRepository) AddingSetsDishesOrder(ctx context.Context, req *order.UpdateOrderRequest) (*order.OrderDetailedListResponse, error) {
    or.logger.Info(fmt.Sprintf("[OrderRepository.AddingSetsDishesOrder] Starting update. OrderID: %d", req.Id))
    
    tx, err := or.db.Begin(ctx)
    if err != nil {
        or.logger.Error(fmt.Sprintf("[OrderRepository.AddingSetsDishesOrder] Transaction start failed: %s", err.Error()))
        return nil, fmt.Errorf("error starting transaction: %w", err)
    }
    defer tx.Rollback(ctx)

    // Version check
    var currentVersion int32
    var isGuest bool
    err = tx.QueryRow(ctx, `
        SELECT version, is_guest 
        FROM orders 
        WHERE id = $1`, req.Id).Scan(&currentVersion, &isGuest)
    if err != nil {
        if err == pgx.ErrNoRows {
            or.logger.Error(fmt.Sprintf("[OrderRepository.AddingSetsDishesOrder] Order not found. ID: %d", req.Id))
            return nil, fmt.Errorf("order not found with ID: %d", req.Id)
        }
        or.logger.Error(fmt.Sprintf("[OrderRepository.AddingSetsDishesOrder] Error fetching order: %s", err.Error()))
        return nil, fmt.Errorf("error fetching order: %w", err)
    }

    if req.Version != 0 && req.Version != currentVersion {
        or.logger.Warning(fmt.Sprintf("[OrderRepository.AddingSetsDishesOrder] Version mismatch. OrderID: %d, Expected: %d, Got: %d",
            req.Id, currentVersion, req.Version))
        return nil, fmt.Errorf("order version mismatch: expected %d, got %d", currentVersion, req.Version)
    }

    newVersion := currentVersion + 1
    now := time.Now()

    // Add new dishes
    for _, dish := range req.DishItems {
        var exists bool
        err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM dishes WHERE id = $1)", dish.DishId).Scan(&exists)
        if err != nil || !exists {
            or.logger.Error(fmt.Sprintf("[OrderRepository.AddingSetsDishesOrder] Invalid dish ID: %d", dish.DishId))
            return nil, fmt.Errorf("invalid dish ID: %d", dish.DishId)
        }

        _, err = tx.Exec(ctx, `
            INSERT INTO dish_order_items (
                order_id, dish_id, quantity, created_at, updated_at,
                order_name, modification_type, modification_number
            ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
            req.Id, dish.DishId, dish.Quantity, now, now,
            req.OrderName, "ADDED", newVersion)
        if err != nil {
            or.logger.Error(fmt.Sprintf("[OrderRepository.AddingSetsDishesOrder] Dish insert failed: %s", err.Error()))
            return nil, fmt.Errorf("error adding dish: %w", err)
        }
    }

    // Add new sets
    for _, set := range req.SetItems {
        var exists bool
        err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM sets WHERE id = $1)", set.SetId).Scan(&exists)
        if err != nil || !exists {
            or.logger.Error(fmt.Sprintf("[OrderRepository.AddingSetsDishesOrder] Invalid set ID: %d", set.SetId))
            return nil, fmt.Errorf("invalid set ID: %d", set.SetId)
        }

        _, err = tx.Exec(ctx, `
            INSERT INTO set_order_items (
                order_id, set_id, quantity, created_at, updated_at,
                order_name, modification_type, modification_number
            ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
            req.Id, set.SetId, set.Quantity, now, now,
            req.OrderName, "ADDED", newVersion)
        if err != nil {
            or.logger.Error(fmt.Sprintf("[OrderRepository.AddingSetsDishesOrder] Set insert failed: %s", err.Error()))
            return nil, fmt.Errorf("error adding set: %w", err)
        }
    }

    // Update order version and timestamp
    _, err = tx.Exec(ctx, `
        UPDATE orders 
        SET version = $1, updated_at = $2 
        WHERE id = $3`,
        newVersion, now, req.Id)
    if err != nil {
        or.logger.Error(fmt.Sprintf("[OrderRepository.AddingSetsDishesOrder] Version update failed: %s", err.Error()))
        return nil, fmt.Errorf("error updating order version: %w", err)
    }

    // Record modification
    _, err = tx.Exec(ctx, `
        INSERT INTO order_modifications (
            order_id, modification_number, modification_type,
            modified_by_user_id, order_name
        ) VALUES ($1, $2, $3, $4, $5)`,
        req.Id, newVersion, "ADD_ITEMS", req.OrderHandlerId, req.OrderName)
    if err != nil {
        or.logger.Error(fmt.Sprintf("[OrderRepository.AddingSetsDishesOrder] Modification record failed: %s", err.Error()))
        return nil, fmt.Errorf("error recording modification: %w", err)
    }


    //     // Fetch version history
    versionHistory, err := or.fetchVersionHistory(ctx, tx, req.Id)
    if err != nil {
        or.logger.Error(fmt.Sprintf("[OrderRepository.AddingSetsDishesOrder] Error fetching version history: %s", err.Error()))
        return nil, fmt.Errorf("error fetching version history: %w", err)
    }
    // Calculate new totals
    totalSummary, err := or.calculateTotalSummary(ctx, tx, req.Id)
    if err != nil {
        or.logger.Error(fmt.Sprintf("[OrderRepository.AddingSetsDishesOrder] Total calculation failed: %s", err.Error()))
        return nil, fmt.Errorf("error calculating totals: %w", err)
    }

    // Fetch updated items
    detailedDishes, err := or.fetchDetailedDishes(ctx, tx, req.Id)
    if err != nil {
        or.logger.Error(fmt.Sprintf("[OrderRepository.AddingSetsDishesOrder] Dish fetch failed: %s", err.Error()))
        return nil, fmt.Errorf("error fetching dishes: %w", err)
    }

    detailedSets, err := or.fetchDetailedSets(ctx, tx, req.Id)
    if err != nil {
        or.logger.Error(fmt.Sprintf("[OrderRepository.AddingSetsDishesOrder] Set fetch failed: %s", err.Error()))
        return nil, fmt.Errorf("error fetching sets: %w", err)
    }

    if err := tx.Commit(ctx); err != nil {
        or.logger.Error(fmt.Sprintf("[OrderRepository.AddingSetsDishesOrder] Commit failed: %s", err.Error()))
        return nil, fmt.Errorf("transaction commit failed: %w", err)
    }

    or.logger.Info(fmt.Sprintf("[OrderRepository.AddingSetsDishesOrder] Update successful. OrderID: %d, NewVersion: %d", 
        req.Id, newVersion))

    return &order.OrderDetailedListResponse{
        Data: []*order.OrderDetailedResponse{
            {
                Id:             req.Id,
                GuestId:        req.GuestId,
                UserId:         req.UserId,
                TableNumber:    req.TableNumber,
                OrderHandlerId: req.OrderHandlerId,
                Status:         req.Status,
                TotalPrice:     req.TotalPrice,
                DataSet:        detailedSets,
                DataDish:       detailedDishes,
                IsGuest:        isGuest,
                Topping:        req.Topping,
                TrackingOrder:  req.TrackingOrder,
                TakeAway:       req.TakeAway,
                ChiliNumber:    req.ChiliNumber,
                TableToken:     req.TableToken,
                OrderName:      req.OrderName,
                ParentOrderId:  req.ParentOrderId,
                CurrentVersion: newVersion,
                VersionHistory: versionHistory,
                TotalSummary: totalSummary,
            },
        },
        Pagination: &order.PaginationInfo{
            CurrentPage: 1,
            TotalPages:  1,
            TotalItems:  int64(len(detailedDishes) + len(detailedSets)),
            PageSize:    100,
        },
    }, nil
}

// new function remove set or dishes -------------------------- start 


func (or *OrderRepository) RemovingSetsDishesOrder(ctx context.Context, req *order.UpdateOrderRequest) (*order.OrderDetailedListResponse, error) {
    or.logger.Info(fmt.Sprintf("[OrderRepository.RemovingSetsDishesOrder] Starting removal. OrderID: %d", req.Id))
    
    tx, err := or.db.Begin(ctx)
    if err != nil {
        or.logger.Error(fmt.Sprintf("[OrderRepository.RemovingSetsDishesOrder] Transaction start failed: %s", err.Error()))
        return nil, fmt.Errorf("error starting transaction: %w", err)
    }
    defer tx.Rollback(ctx)

    // Version check (same as adding)
    var currentVersion int32
    var isGuest bool
    err = tx.QueryRow(ctx, `
        SELECT version, is_guest 
        FROM orders 
        WHERE id = $1`, req.Id).Scan(&currentVersion, &isGuest)
    if err != nil {
        if err == pgx.ErrNoRows {
            or.logger.Error(fmt.Sprintf("[OrderRepository.RemovingSetsDishesOrder] Order not found. ID: %d", req.Id))
            return nil, fmt.Errorf("order not found with ID: %d", req.Id)
        }
        or.logger.Error(fmt.Sprintf("[OrderRepository.RemovingSetsDishesOrder] Error fetching order: %s", err.Error()))
        return nil, fmt.Errorf("error fetching order: %w", err)
    }

    if req.Version != 0 && req.Version != currentVersion {
        or.logger.Warning(fmt.Sprintf("[OrderRepository.RemovingSetsDishesOrder] Version mismatch. OrderID: %d, Expected: %d, Got: %d",
            req.Id, currentVersion, req.Version))
        return nil, fmt.Errorf("order version mismatch: expected %d, got %d", currentVersion, req.Version)
    }

    newVersion := currentVersion + 1
    now := time.Now()

    // Remove dishes
    for _, dish := range req.DishItems {
        var exists bool
        err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM dishes WHERE id = $1)", dish.DishId).Scan(&exists)
        if err != nil || !exists {
            or.logger.Error(fmt.Sprintf("[OrderRepository.RemovingSetsDishesOrder] Invalid dish ID: %d", dish.DishId))
            return nil, fmt.Errorf("invalid dish ID: %d", dish.DishId)
        }

        _, err = tx.Exec(ctx, `
            INSERT INTO dish_order_items (
                order_id, dish_id, quantity, created_at, updated_at,
                order_name, modification_type, modification_number
            ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
            req.Id, dish.DishId, dish.Quantity, now, now,
            req.OrderName, "REMOVED", newVersion)
        if err != nil {
            or.logger.Error(fmt.Sprintf("[OrderRepository.RemovingSetsDishesOrder] Dish removal record failed: %s", err.Error()))
            return nil, fmt.Errorf("error recording dish removal: %w", err)
        }
    }

    // Remove sets
    for _, set := range req.SetItems {
        var exists bool
        err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM sets WHERE id = $1)", set.SetId).Scan(&exists)
        if err != nil || !exists {
            or.logger.Error(fmt.Sprintf("[OrderRepository.RemovingSetsDishesOrder] Invalid set ID: %d", set.SetId))
            return nil, fmt.Errorf("invalid set ID: %d", set.SetId)
        }

        _, err = tx.Exec(ctx, `
            INSERT INTO set_order_items (
                order_id, set_id, quantity, created_at, updated_at,
                order_name, modification_type, modification_number
            ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
            req.Id, set.SetId, set.Quantity, now, now,
            req.OrderName, "REMOVED", newVersion)
        if err != nil {
            or.logger.Error(fmt.Sprintf("[OrderRepository.RemovingSetsDishesOrder] Set removal record failed: %s", err.Error()))
            return nil, fmt.Errorf("error recording set removal: %w", err)
        }
    }

    // Update order version and timestamp (same as adding)
    _, err = tx.Exec(ctx, `
        UPDATE orders 
        SET version = $1, updated_at = $2 
        WHERE id = $3`,
        newVersion, now, req.Id)
    if err != nil {
        or.logger.Error(fmt.Sprintf("[OrderRepository.RemovingSetsDishesOrder] Version update failed: %s", err.Error()))
        return nil, fmt.Errorf("error updating order version: %w", err)
    }

    // Record modification with REMOVE_ITEMS type
    _, err = tx.Exec(ctx, `
        INSERT INTO order_modifications (
            order_id, modification_number, modification_type,
            modified_by_user_id, order_name
        ) VALUES ($1, $2, $3, $4, $5)`,
        req.Id, newVersion, "REMOVE_ITEMS", req.OrderHandlerId, req.OrderName)
    if err != nil {
        or.logger.Error(fmt.Sprintf("[OrderRepository.RemovingSetsDishesOrder] Modification record failed: %s", err.Error()))
        return nil, fmt.Errorf("error recording modification: %w", err)
    }

    // The following parts remain identical to the adding function
    versionHistory, err := or.fetchVersionHistory(ctx, tx, req.Id)
    if err != nil {
        or.logger.Error(fmt.Sprintf("[OrderRepository.RemovingSetsDishesOrder] Error fetching version history: %s", err.Error()))
        return nil, fmt.Errorf("error fetching version history: %w", err)
    }

    totalSummary, err := or.calculateTotalSummary(ctx, tx, req.Id)
    if err != nil {
        or.logger.Error(fmt.Sprintf("[OrderRepository.RemovingSetsDishesOrder] Total calculation failed: %s", err.Error()))
        return nil, fmt.Errorf("error calculating totals: %w", err)
    }

    detailedDishes, err := or.fetchDetailedDishes(ctx, tx, req.Id)
    if err != nil {
        or.logger.Error(fmt.Sprintf("[OrderRepository.RemovingSetsDishesOrder] Dish fetch failed: %s", err.Error()))
        return nil, fmt.Errorf("error fetching dishes: %w", err)
    }

    detailedSets, err := or.fetchDetailedSets(ctx, tx, req.Id)
    if err != nil {
        or.logger.Error(fmt.Sprintf("[OrderRepository.RemovingSetsDishesOrder] Set fetch failed: %s", err.Error()))
        return nil, fmt.Errorf("error fetching sets: %w", err)
    }

    if err := tx.Commit(ctx); err != nil {
        or.logger.Error(fmt.Sprintf("[OrderRepository.RemovingSetsDishesOrder] Commit failed: %s", err.Error()))
        return nil, fmt.Errorf("transaction commit failed: %w", err)
    }

    or.logger.Info(fmt.Sprintf("[OrderRepository.RemovingSetsDishesOrder] Removal successful. OrderID: %d, NewVersion: %d", 
        req.Id, newVersion))

    return &order.OrderDetailedListResponse{
        Data: []*order.OrderDetailedResponse{
            {
                Id:             req.Id,
                GuestId:        req.GuestId,
                UserId:         req.UserId,
                TableNumber:    req.TableNumber,
                OrderHandlerId: req.OrderHandlerId,
                Status:         req.Status,
                TotalPrice:     req.TotalPrice,
                DataSet:        detailedSets,
                DataDish:       detailedDishes,
                IsGuest:        isGuest,
                Topping:        req.Topping,
                TrackingOrder:  req.TrackingOrder,
                TakeAway:       req.TakeAway,
                ChiliNumber:    req.ChiliNumber,
                TableToken:     req.TableToken,
                OrderName:      req.OrderName,
                ParentOrderId:  req.ParentOrderId,
                CurrentVersion: newVersion,
                VersionHistory: versionHistory,
                TotalSummary:   totalSummary,
            },
        },
        Pagination: &order.PaginationInfo{
            CurrentPage: 1,
            TotalPages:  1,
            TotalItems:  int64(len(detailedDishes) + len(detailedSets)),
            PageSize:    100,
        },
    }, nil
}





// new function remove set or dishes -------------------------- end

// new function dish delivery set or dishes -------------------------- start

func (or *OrderRepository) MarkDishesDelivered(ctx context.Context, req *order.CreateDishDeliveryRequest) (*order.OrderDetailedResponseWithDelivery, error) {
    const (
        deliveryStatus    = "DELIVERED"
        modificationType  = "DELIVER_ITEMS"
        validOrderStatus  = "IN_PROGRESS"
    )

    // Initial validation
 
    tx, err := or.db.Begin(ctx)
    if err != nil {
        or.logger.Error(fmt.Sprintf("[MarkDishesDelivered] Transaction start failed: %s", err.Error()))
        return nil, fmt.Errorf("error starting transaction: %w", err)
    }
    defer tx.Rollback(ctx)

    // Fetch current order state
    var (
        currentOrder      order.Order
        guestID           sql.NullInt64
        userID            sql.NullInt64
        tableNumber       sql.NullInt64
    )
    err = tx.QueryRow(ctx, `
        SELECT 
            id, version, is_guest, guest_id, user_id, table_number, status,
            total_price, order_handler_id, topping, tracking_order, take_away,
            chili_number, table_token, order_name, parent_order_id
        FROM orders 
        WHERE id = $1`, req.OrderId).Scan(
            &currentOrder.Id,
            &currentOrder.Version,
            &currentOrder.IsGuest,
            &guestID,
            &userID,
            &tableNumber,
            &currentOrder.Status,
            &currentOrder.TotalPrice,
            &currentOrder.OrderHandlerId,
            &currentOrder.Topping,
            &currentOrder.TrackingOrder,
            &currentOrder.TakeAway,
            &currentOrder.ChiliNumber,
            &currentOrder.TableToken,
            &currentOrder.OrderName,
            &currentOrder.ParentOrderId,
        )
    if err != nil {
        if err == pgx.ErrNoRows {
            or.logger.Error(fmt.Sprintf("[MarkDishesDelivered] Order not found. ID: %d", req.OrderId))
            return nil, fmt.Errorf("order not found")
        }
        or.logger.Error(fmt.Sprintf("[MarkDishesDelivered] Order fetch error: %s", err.Error()))
        return nil, fmt.Errorf("failed to retrieve order: %w", err)
    }
    currentOrder.GuestId = nullInt64ToProtoInt64(guestID)
    currentOrder.UserId = nullInt64ToProtoInt64(userID)
    currentOrder.TableNumber = nullInt64ToProtoInt64(tableNumber)
    
    // Validate order state
    if currentOrder.Status != validOrderStatus {
        or.logger.Warning(fmt.Sprintf(
            "[MarkDishesDelivered] Invalid order status. OrderID: %d, Status: %s",
            req.OrderId, currentOrder.Status))
        return nil, fmt.Errorf("order cannot be modified in current status: %s", currentOrder.Status)
    }

    // Strict version check


    newVersion := currentOrder.Version + 1
    now := time.Now().UTC()

    // Process dish deliveries
    for _, dish := range req.DishItems {
        // Validate quantity
        if dish.Quantity <= 0 {
            or.logger.Error(fmt.Sprintf(
                "[MarkDishesDelivered] Invalid quantity. DishID: %d, Qty: %d",
                dish.DishId, dish.Quantity))
            return nil, fmt.Errorf("invalid quantity for dish %d: must be positive", dish.DishId)
        }

        // Calculate net ordered quantity
        var netQuantity int64
        err = tx.QueryRow(ctx, `
            SELECT COALESCE(SUM(
                CASE modification_type
                    WHEN 'ADDED' THEN quantity
                    WHEN 'REMOVED' THEN -quantity
                    ELSE quantity
                END
            ), 0)
            FROM dish_order_items
            WHERE order_id = $1 AND dish_id = $2`,
            req.OrderId, dish.DishId).Scan(&netQuantity)
        if err != nil {
            or.logger.Error(fmt.Sprintf(
                "[MarkDishesDelivered] Net quantity error. DishID: %d: %s",
                dish.DishId, err.Error()))
            return nil, fmt.Errorf("failed to calculate ordered quantity: %w", err)
        }

        if netQuantity <= 0 {
            or.logger.Error(fmt.Sprintf(
                "[MarkDishesDelivered] Dish not in order. DishID: %d",
                dish.DishId))
            return nil, fmt.Errorf("dish %d not found in order", dish.DishId)
        }

        // Calculate existing deliveries
        var delivered int64
        err = tx.QueryRow(ctx, `
            SELECT COALESCE(SUM(quantity_delivered), 0)
            FROM dish_deliveries
            WHERE order_id = $1 AND dish_id = $2`,
            req.OrderId, dish.DishId).Scan(&delivered)
        if err != nil {
            or.logger.Error(fmt.Sprintf(
                "[MarkDishesDelivered] Delivery check error. DishID: %d: %s",
                dish.DishId, err.Error()))
            return nil, fmt.Errorf("failed to check existing deliveries: %w", err)
        }

        // Validate delivery quantity
        remaining := netQuantity - delivered
        if dish.Quantity > remaining {
            or.logger.Error(fmt.Sprintf(
                "[MarkDishesDelivered] Over-delivery. DishID: %d, Remaining: %d, Attempt: %d",
                dish.DishId, remaining, dish.Quantity))
            return nil, fmt.Errorf("cannot deliver %d of %d remaining for dish %d", 
                dish.Quantity, remaining, dish.DishId)
        }

        // Insert delivery record
        _, err = tx.Exec(ctx, `
            INSERT INTO dish_deliveries (
                order_id, order_name, guest_id, user_id, table_number,
                dish_id, quantity_delivered, delivery_status, delivered_at,
                delivered_by_user_id, modification_number, version
            ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`,
            req.OrderId,
            req.OrderName,
            getNullableID(currentOrder.IsGuest, guestID),
            getNullableID(!currentOrder.IsGuest, userID),
            getTableNumber(tableNumber),
            dish.DishId,
            dish.Quantity,
            deliveryStatus,
            now,
            req.UserId,
            newVersion,
            newVersion,
        )
        if err != nil {
            or.logger.Error(fmt.Sprintf(
                "[MarkDishesDelivered] Delivery insert failed. DishID: %d: %s",
                dish.DishId, err.Error()))
            return nil, fmt.Errorf("failed to record delivery: %w", err)
        }
    }

    // Update order version
    _, err = tx.Exec(ctx, `
        UPDATE orders 
        SET version = $1, updated_at = $2 
        WHERE id = $3`,
        newVersion, now, req.OrderId)
    if err != nil {
        or.logger.Error(fmt.Sprintf(
            "[MarkDishesDelivered] Version update failed: %s", 
            err.Error()))
        return nil, fmt.Errorf("failed to update order version: %w", err)
    }

    // Record modification
    _, err = tx.Exec(ctx, `
        INSERT INTO order_modifications (
            order_id, modification_number, modification_type,
            modified_by_user_id, order_name, version
        ) VALUES ($1, $2, $3, $4, $5, $6)`,
        req.OrderId, newVersion, modificationType, 
        req.UserId, req.OrderName, newVersion)
    if err != nil {
        or.logger.Error(fmt.Sprintf(
            "[MarkDishesDelivered] Modification record failed: %s", 
            err.Error()))
        return nil, fmt.Errorf("failed to record modification: %w", err)
    }

    // Fetch updated order details
    versionHistory, err := or.fetchVersionHistory(ctx, tx, req.OrderId)
    if err != nil {
        or.logger.Error(fmt.Sprintf("[MarkDishesDelivered] Version history fetch failed: %s", err.Error()))
        return nil, fmt.Errorf("error fetching version history: %w", err)
    }

    totalSummary, err := or.calculateTotalSummary(ctx, tx, req.OrderId)
    if err != nil {
        or.logger.Error(fmt.Sprintf("[MarkDishesDelivered] Total calculation failed: %s", err.Error()))
        return nil, fmt.Errorf("error calculating totals: %w", err)
    }

    detailedDishes, err := or.fetchDetailedDishes(ctx, tx, req.OrderId)
    if err != nil {
        or.logger.Error(fmt.Sprintf("[MarkDishesDelivered] Dish fetch failed: %s", err.Error()))
        return nil, fmt.Errorf("error fetching dishes: %w", err)
    }

    detailedSets, err := or.fetchDetailedSets(ctx, tx, req.OrderId)
    if err != nil {
        or.logger.Error(fmt.Sprintf("[MarkDishesDelivered] Set fetch failed: %s", err.Error()))
        return nil, fmt.Errorf("error fetching sets: %w", err)
    }
// for delivery start 

deliveryHistory, err := or.fetchDeliveryHistory(ctx,tx, req.OrderId)
if err != nil {
    or.logger.Error(fmt.Sprintf("[MarkDishesDelivered] Delivery history fetch failed: %s", err.Error()))
    return nil, fmt.Errorf("error fetching delivery history: %w", err)
}
orderDeliveryStatus, totalDelivered, lastDeliveryAt, err := or.calculateDeliveryStatus(ctx,tx, req.OrderId)

if err != nil {
    or.logger.Error(fmt.Sprintf("[MarkDishesDelivered] Delivery status calculation failed: %s", err.Error()))
    return nil, fmt.Errorf("error calculating delivery status: %w", err)
}
var pbLastDeliveryAt *timestamppb.Timestamp
if lastDeliveryAt != nil {
    pbLastDeliveryAt = timestamppb.New(*lastDeliveryAt)
}

// Commit the transaction after all operations
err = tx.Commit(ctx)
if err != nil {
    or.logger.Error(fmt.Sprintf("[MarkDishesDelivered] Transaction commit failed: %s", err.Error()))
    return nil, fmt.Errorf("error committing transaction: %w", err)
}
// for delivery end ---------------

    or.logger.Info(fmt.Sprintf("[MarkDishesDelivered] Delivery successful. OrderID: %d, NewVersion: %d", 
        req.OrderId, newVersion))

        return &order.OrderDetailedResponseWithDelivery{
            Id:                  currentOrder.Id,
            GuestId:             currentOrder.GuestId,
            UserId:              currentOrder.UserId,
            TableNumber:         currentOrder.TableNumber,
            OrderHandlerId:      currentOrder.OrderHandlerId,
            Status:              currentOrder.Status,
            TotalPrice:          currentOrder.TotalPrice,
            DataSet:             detailedSets,
            DataDish:            detailedDishes,
            IsGuest:             currentOrder.IsGuest,
            Topping:             currentOrder.Topping,
            TrackingOrder:       currentOrder.TrackingOrder,
            TakeAway:            currentOrder.TakeAway,
            ChiliNumber:         currentOrder.ChiliNumber,
            TableToken:          currentOrder.TableToken,
            OrderName:           currentOrder.OrderName,
            CurrentVersion:      newVersion,
            ParentOrderId:       currentOrder.ParentOrderId,
            VersionHistory:      versionHistory,
            TotalSummary:        totalSummary,
            DeliveryHistory:     deliveryHistory,
            CurrentDeliveryStatus: orderDeliveryStatus,
            TotalItemsDelivered: totalDelivered,
            LastDeliveryAt:      pbLastDeliveryAt,
        }, nil
    
}

// Helper functions
func getNullableID(shouldInclude bool, id sql.NullInt64) interface{} {
    if shouldInclude && id.Valid {
        return id.Int64
    }
    return nil
}

func getTableNumber(tableNumber sql.NullInt64) interface{} {
    if tableNumber.Valid {
        return tableNumber.Int64
    }
    return nil
}
// 1. Add these methods to your OrderRepository
func (or *OrderRepository) fetchDeliveryHistory(ctx context.Context, tx pgx.Tx, orderID int64) ([]*order.DishDelivery, error) {
    rows, err := tx.Query(ctx, `
        SELECT 
            dd.id, 
            dd.order_id, 
            dd.order_name, 
            dd.guest_id, 
            dd.user_id, 
            dd.table_number,
            dd.dish_id, 
            dd.quantity_delivered, 
            dd.delivery_status, 
            dd.delivered_at,
            dd.delivered_by_user_id, 
            dd.is_guest, 
            dd.created_at, 
            dd.updated_at,
            dd.modification_number
        FROM dish_deliveries dd
        WHERE dd.order_id = $1
        ORDER BY dd.delivered_at`, orderID)
    if err != nil {
        return nil, fmt.Errorf("error fetching delivery history: %w", err)
    }
    defer rows.Close()

    var deliveries []*order.DishDelivery
    for rows.Next() {
        var (
            dd                  order.DishDelivery
            dishID              int64
            quantityDelivered   int32
            deliveredAt         time.Time
            createdAt           time.Time
            updatedAt           time.Time
            modNumber           int32
        )

        err := rows.Scan(
            &dd.Id,
            &dd.OrderId,
            &dd.OrderName,
            &dd.GuestId,
            &dd.UserId,
            &dd.TableNumber,
            &dishID,
            &quantityDelivered,
            &dd.DeliveryStatus,
            &deliveredAt,
            &dd.DeliveredByUserId,
            &dd.IsGuest,
            &createdAt,
            &updatedAt,
            &modNumber,
        )
        if err != nil {
            return nil, fmt.Errorf("error scanning delivery row: %w", err)
        }

        // Convert timestamps to protobuf format
        dd.DeliveredAt = timestamppb.New(deliveredAt)
        dd.CreatedAt = timestamppb.New(createdAt)
        dd.UpdatedAt = timestamppb.New(updatedAt)

        // Build DishOrderItem with delivery context
        dishItem := &order.DishOrderItem{
            Id:                 dd.Id,  // Using delivery ID as proxy
            DishId:             dishID,
            Quantity:           int64(quantityDelivered),
            CreatedAt:          timestamppb.New(createdAt),
            UpdatedAt:          timestamppb.New(updatedAt),
            OrderName:          dd.OrderName,
            ModificationType:   "DELIVER_ITEMS",  // Static type for deliveries
            ModificationNumber: modNumber,
        }

        dd.DishItems = []*order.DishOrderItem{dishItem}
        deliveries = append(deliveries, &dd)
    }
    
    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("error after scanning rows: %w", err)
    }

    return deliveries, nil
}
func (or *OrderRepository) calculateDeliveryStatus(ctx context.Context, tx pgx.Tx, orderID int64) (order.DeliveryStatus, int32, *time.Time, error) {
    // Query to get net ordered and delivered quantities per dish
    netOrderedDeliveredQuery := `
        WITH net_ordered AS (
            SELECT 
                dish_id, 
                SUM(CASE modification_type 
                    WHEN 'ADDED' THEN quantity 
                    WHEN 'REMOVED' THEN -quantity 
                    ELSE quantity 
                END) AS net_ordered
            FROM dish_order_items
            WHERE order_id = $1
            GROUP BY dish_id
        ),
        total_delivered AS (
            SELECT 
                dish_id, 
                SUM(quantity_delivered) AS total_delivered
            FROM dish_deliveries
            WHERE order_id = $1
            GROUP BY dish_id
        )
        SELECT 
            n.dish_id,
            n.net_ordered,
            COALESCE(d.total_delivered, 0) AS total_delivered
        FROM net_ordered n
        LEFT JOIN total_delivered d ON n.dish_id = d.dish_id`

    rows, err := tx.Query(ctx, netOrderedDeliveredQuery, orderID)
    if err != nil {
        return order.DeliveryStatus_PENDING, 0, nil, fmt.Errorf("failed to query delivery data: %w", err)
    }
    defer rows.Close()

    var (
        totalDelivered    int32
        allFullyDelivered = true
        anyDelivered      = false
    )

    // Check each dish's delivery status
    for rows.Next() {
        var dishID, netOrdered, delivered int32
        if err := rows.Scan(&dishID, &netOrdered, &delivered); err != nil {
            return order.DeliveryStatus_PENDING, 0, nil, 
                fmt.Errorf("failed to scan dish delivery data: %w", err)
        }

        totalDelivered += delivered

        if delivered < netOrdered {
            allFullyDelivered = false
        }
        if delivered > 0 {
            anyDelivered = true
        }
    }

    if err := rows.Err(); err != nil {
        return order.DeliveryStatus_PENDING, 0, nil, 
            fmt.Errorf("error processing delivery data rows: %w", err)
    }

    // Determine overall delivery status using protobuf enum
    var deliveryStatus order.DeliveryStatus
    switch {
    case allFullyDelivered && totalDelivered > 0:
        deliveryStatus = order.DeliveryStatus_FULLY_DELIVERED
    case anyDelivered:
        deliveryStatus = order.DeliveryStatus_PARTIALLY_DELIVERED
    default:
        deliveryStatus = order.DeliveryStatus_PENDING  // Using PENDING as default state
    }

    // Get last delivery timestamp using transaction
    var lastDeliveryAt *time.Time
    err = tx.QueryRow(ctx, `
        SELECT MAX(delivered_at) 
        FROM dish_deliveries 
        WHERE order_id = $1`, orderID).Scan(&lastDeliveryAt)
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return deliveryStatus, totalDelivered, nil, nil
        }
        return order.DeliveryStatus_PENDING, 0, nil, 
            fmt.Errorf("failed to retrieve last delivery time: %w", err)
    }

    // Handle potential zero time value
    if lastDeliveryAt != nil && lastDeliveryAt.IsZero() {
        lastDeliveryAt = nil
    }

    return deliveryStatus, totalDelivered, lastDeliveryAt, nil
}
// new function dish delivery set or dishes -------------------------- end


func nullInt64ToProtoInt64(n sql.NullInt64) int64 {
    if n.Valid {
        return n.Int64
    }
    return 0 // Protobuf default value for missing int64
}