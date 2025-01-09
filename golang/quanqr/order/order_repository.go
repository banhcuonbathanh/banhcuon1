package order_grpc

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"math"

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

// Helper function with corrected transaction type
func (or *OrderRepository) fetchDetailedSets(ctx context.Context, tx pgx.Tx, orderID int64) ([]*order.OrderSetDetailed, error) {
    // First, fetch all sets with their basic information
    query := `
        WITH set_dishes AS (
            SELECT sd.set_id,
                   jsonb_agg(
                       jsonb_build_object(
                           'id', d.id,
                           'name', d.name,
                           'description', d.description,
                           'price', d.price,
                           'image', d.image,
                           'quantity', sd.quantity
                       )
                   ) as dishes
            FROM set_dishes sd
            JOIN dishes d ON sd.dish_id = d.id
            GROUP BY sd.set_id
        )
        SELECT s.id, s.name, s.description, s.user_id, s.created_at, s.updated_at,
               s.is_favourite, s.like_by, s.is_public, s.image, s.price, soi.quantity,
               COALESCE(sd.dishes, '[]'::jsonb) as dishes
        FROM set_order_items soi
        JOIN sets s ON soi.set_id = s.id
        LEFT JOIN set_dishes sd ON s.id = sd.set_id
        WHERE soi.order_id = $1
    `
    
    rows, err := tx.Query(ctx, query, orderID)
    if err != nil {
        return nil, fmt.Errorf("error querying set details: %w", err)
    }
    defer rows.Close()

    var sets []*order.OrderSetDetailed
    for rows.Next() {
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
            return nil, fmt.Errorf("error scanning set row: %w", err)
        }
        
        // Parse dishes JSON into the appropriate struct
        var dishes []*order.OrderDetailedDish
        if err := json.Unmarshal(dishesJSON, &dishes); err != nil {
            return nil, fmt.Errorf("error unmarshaling dishes: %w", err)
        }
        
        set.CreatedAt = timestamppb.New(createdAt)
        set.UpdatedAt = timestamppb.New(updatedAt)
        set.LikeBy = likeBy
        set.Dishes = dishes
        
        sets = append(sets, set)
    }

    if err = rows.Err(); err != nil {
        return nil, fmt.Errorf("error iterating set rows: %w", err)
    }

    return sets, nil
}

// Helper function to fetch dishes for a set using pgxpool transaction
func (or *OrderRepository) fetchSetDishes(ctx context.Context,  tx pgx.Tx, setID int64) ([]*order.OrderDetailedDish, error) {
    query := `
        SELECT d.id, sd.quantity, d.name, d.price, d.description, d.image, d.status
        FROM set_dishes sd
        JOIN dishes d ON sd.dish_id = d.id
        WHERE sd.set_id = $1
    `
    
    rows, err := tx.Query(ctx, query, setID)
    if err != nil {
        return nil, fmt.Errorf("error querying set dishes: %w", err)
    }
    defer rows.Close()

    var dishes []*order.OrderDetailedDish
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
            return nil, fmt.Errorf("error scanning dish row: %w", err)
        }
        dishes = append(dishes, dish)
    }

    if err = rows.Err(); err != nil {
        return nil, fmt.Errorf("error iterating dish rows: %w", err)
    }

    return dishes, nil
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
func (or *OrderRepository) AddingSetsDishesOrder(ctx context.Context, req *order.UpdateOrderRequest) (*order.OrderDetailedListResponse, error) {
    or.logger.Info(fmt.Sprintf("addingSetsDishesOrder order with ID repository: %d", req.Id))
    
    // Start a database transaction since we'll be making multiple related changes
    tx, err := or.db.Begin(ctx)
    if err != nil {
        or.logger.Error("Error starting transaction: " + err.Error())
        return nil, fmt.Errorf("error starting transaction: %w", err)
    }
    defer tx.Rollback(ctx) // Ensure rollback in case of errors

    // Check current version and guest status
    var currentVersion int32
    var isGuest bool
    err = tx.QueryRow(ctx, `
        SELECT version, is_guest 
        FROM orders 
        WHERE id = $1`, req.Id).Scan(&currentVersion, &isGuest)
    if err != nil {
        if err == pgx.ErrNoRows {
            return nil, fmt.Errorf("order not found with ID: %d", req.Id)
        }
        or.logger.Error(fmt.Sprintf("Error fetching order: %s", err.Error()))
        return nil, fmt.Errorf("error fetching order: %w", err)
    }

    // Validate version to prevent concurrent modifications
    if req.Version != 0 && req.Version != currentVersion {
        return nil, fmt.Errorf("order version mismatch: expected %d, got %d", currentVersion, req.Version)
    }

    newVersion := currentVersion + 1

    // Update the order's main record
    _, err = tx.Exec(ctx, `
        UPDATE orders 
        SET 
            version = $1,
            status = $2,
            total_price = $3,
            topping = $4,
            tracking_order = $5,
            take_away = $6,
            chili_number = $7,
            table_token = $8,
            order_name = $9,
            updated_at = CURRENT_TIMESTAMP
        WHERE id = $10`,
        newVersion,
        req.Status,
        req.TotalPrice,
        req.Topping,
        req.TrackingOrder,
        req.TakeAway,
        req.ChiliNumber,
        req.TableToken,
        req.OrderName,
        req.Id)
    if err != nil {
        return nil, fmt.Errorf("error updating order: %w", err)
    }

    // Record the modification in order_modifications table
    _, err = tx.Exec(ctx, `
        INSERT INTO order_modifications (
            order_id, 
            modification_number, 
            modification_type, 
            modified_at, 
            modified_by_user_id
        ) VALUES ($1, $2, $3, CURRENT_TIMESTAMP, $4)`,
        req.Id, 
        newVersion, 
        "UPDATE", // or could be "ADD_ITEMS" depending on your business logic
        req.UserId)
    if err != nil {
        return nil, fmt.Errorf("error recording modification: %w", err)
    }

    // Record dish items for this version
    for _, dish := range req.DishItems {
        _, err = tx.Exec(ctx, `
            INSERT INTO dish_order_items (
                order_id, 
                dish_id, 
                quantity, 
                modification_number
            ) VALUES ($1, $2, $3, $4)`,
            req.Id,
            dish.DishId,
            dish.Quantity,
            newVersion)
        if err != nil {
            return nil, fmt.Errorf("error recording dish item: %w", err)
        }
    }

    // Record set items for this version
    for _, set := range req.SetItems {
        _, err = tx.Exec(ctx, `
            INSERT INTO set_order_items (
                order_id, 
                set_id, 
                quantity, 
                modification_number
            ) VALUES ($1, $2, $3, $4)`,
            req.Id,
            set.SetId,
            set.Quantity,
            newVersion)
        if err != nil {
            return nil, fmt.Errorf("error recording set item: %w", err)
        }
    }

    // Fetch version history and summaries
    versionHistory, err := or.fetchVersionHistory(ctx, tx, req.Id)
    if err != nil {
        return nil, fmt.Errorf("error fetching version history: %w", err)
    }

    // Calculate total summary across all versions
    totalSummary, err := or.calculateTotalSummary(ctx, tx, req.Id)
    if err != nil {
        return nil, fmt.Errorf("error calculating total summary: %w", err)
    }

    // Fetch detailed items for the response
    detailedDishes, err := or.fetchDetailedDishes(ctx, tx, req.Id)
    if err != nil {
        return nil, fmt.Errorf("error fetching detailed dishes: %w", err)
    }

    detailedSets, err := or.fetchDetailedSets(ctx, tx, req.Id)
    if err != nil {
        return nil, fmt.Errorf("error fetching detailed sets: %w", err)
    }

    // Commit the transaction
    if err := tx.Commit(ctx); err != nil {
        or.logger.Error("Error committing transaction: " + err.Error())
        return nil, fmt.Errorf("error committing transaction: %w", err)
    }

    // Construct and return the response
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
                CurrentVersion: newVersion,
                ParentOrderId:  req.ParentOrderId,
                VersionHistory: versionHistory,
                TotalSummary:  totalSummary,
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

// Helper function to calculate total summary
func (or *OrderRepository) calculateTotalSummary(ctx context.Context, tx pgx.Tx, orderID int64) (*order.OrderTotalSummary, error) {
    var summary order.OrderTotalSummary
    
    err := tx.QueryRow(ctx, `
        WITH order_stats AS (
            SELECT 
                COUNT(DISTINCT modification_number) as total_versions,
                SUM(CASE WHEN item_type = 'DISH' THEN quantity ELSE 0 END) as total_dishes,
                SUM(CASE WHEN item_type = 'SET' THEN quantity ELSE 0 END) as total_sets,
                SUM(price * quantity) as total_price
            FROM (
                SELECT 
                    di.modification_number,
                    'DISH' as item_type,
                    di.quantity,
                    d.price
                FROM dish_order_items di
                JOIN dishes d ON di.dish_id = d.id
                WHERE di.order_id = $1
                UNION ALL
                SELECT 
                    si.modification_number,
                    'SET' as item_type,
                    si.quantity,
                    s.price
                FROM set_order_items si
                JOIN sets s ON si.set_id = s.id
                WHERE si.order_id = $1
            ) all_items
        )
        SELECT 
            total_versions,
            total_dishes,
            total_sets,
            total_price
        FROM order_stats`,
        orderID).Scan(
            &summary.TotalVersions,
            &summary.TotalDishesOrdered,
            &summary.TotalSetsOrdered,
            &summary.CumulativeTotalPrice,
        )
    if err != nil {
        return nil, fmt.Errorf("error calculating total summary: %w", err)
    }

    // Fetch most ordered items
    summary.MostOrderedItems, err = or.fetchMostOrderedItems(ctx, tx, orderID)
    if err != nil {
        return nil, fmt.Errorf("error fetching most ordered items: %w", err)
    }

    return &summary, nil
}

// Helper function to fetch version changes


// fetchMostOrderedItems retrieves the most frequently ordered items (both dishes and sets)
// across all versions of an order, combining quantities from different modifications
func (or *OrderRepository) fetchMostOrderedItems(ctx context.Context, tx pgx.Tx, orderID int64) ([]*order.OrderItemCount, error) {
    // Query both dishes and sets, combining their quantities across all modifications
    rows, err := tx.Query(ctx, `
        WITH combined_items AS (
            -- First, get all dish orders with their quantities
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
            
            -- Then get all set orders with their quantities
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
        -- Select top ordered items, ordered by quantity
        SELECT 
            item_type,
            item_id,
            item_name,
            total_quantity
        FROM combined_items
        ORDER BY total_quantity DESC
        LIMIT 5  -- Limit to top 5 most ordered items; adjust as needed
    `, orderID)
    if err != nil {
        return nil, fmt.Errorf("error querying most ordered items: %w", err)
    }
    defer rows.Close()

    var items []*order.OrderItemCount
    for rows.Next() {
        var item order.OrderItemCount
        err := rows.Scan(
            &item.ItemType,
            &item.ItemId,
            &item.ItemName,
            &item.TotalQuantity,
        )
        if err != nil {
            return nil, fmt.Errorf("error scanning most ordered item: %w", err)
        }
        items = append(items, &item)
    }

    return items, nil
}
// -------------------------------------------------- update ordder adding set and dishes end -----------------------


// fetchVersionHistory fetches the complete history of changes for an order
func (or *OrderRepository) fetchVersionHistory(ctx context.Context, tx pgx.Tx, orderID int64) ([]*order.OrderVersionSummary, error) {
    var summaries []*order.OrderVersionSummary
    
    rows, err := tx.Query(ctx, `
        WITH version_items AS (
            SELECT 
                m.modification_number,
                m.modification_type,
                m.modified_at,
                COUNT(DISTINCT d.id) as dishes_count,
                COUNT(DISTINCT s.id) as sets_count,
                SUM(CASE 
                    WHEN d.id IS NOT NULL THEN d.price * di.quantity
                    WHEN s.id IS NOT NULL THEN s.price * si.quantity
                    ELSE 0
                END) as version_total_price
            FROM order_modifications m
            LEFT JOIN dish_order_items di ON m.order_id = di.order_id AND m.modification_number = di.modification_number
            LEFT JOIN set_order_items si ON m.order_id = si.order_id AND m.modification_number = si.modification_number
            LEFT JOIN dishes d ON di.dish_id = d.id
            LEFT JOIN sets s ON si.set_id = s.id
            WHERE m.order_id = $1
            GROUP BY m.modification_number, m.modification_type, m.modified_at
        )
        SELECT 
            modification_number,
            modification_type,
            modified_at,
            dishes_count,
            sets_count,
            version_total_price
        FROM version_items
        ORDER BY modification_number ASC`,
        orderID)
    if err != nil {
        return nil, fmt.Errorf("error querying version history: %w", err)
    }
    defer rows.Close()
    
    // Create a slice to store our version data
    type versionInfo struct {
        summary    *order.OrderVersionSummary  // Note the pointer here
        modifiedAt time.Time
    }
    var versionsToProcess []versionInfo
    
    // Collect all rows
    for rows.Next() {
        // Create new instances for each iteration
        versionData := versionInfo{
            summary: &order.OrderVersionSummary{}, // Initialize a new pointer
        }
        
        err := rows.Scan(
            &versionData.summary.VersionNumber,
            &versionData.summary.ModificationType,
            &versionData.modifiedAt,
            &versionData.summary.TotalDishesCount,
            &versionData.summary.TotalSetsCount,
            &versionData.summary.VersionTotalPrice,
        )
        if err != nil {
            return nil, fmt.Errorf("error scanning version history: %w", err)
        }
        versionsToProcess = append(versionsToProcess, versionData)
    }
    
    // Process each version
    for _, versionData := range versionsToProcess {
        // Set the timestamp using the modified time
        versionData.summary.ModifiedAt = timestamppb.New(versionData.modifiedAt)
        
        // Get the changes for this version
        changes, err := or.getVersionChanges(ctx, tx, orderID, versionData.summary.VersionNumber)
        if err != nil {
            return nil, fmt.Errorf("error getting version changes: %w", err)
        }
        versionData.summary.Changes = changes
        
        // Append the pointer to the summary
        summaries = append(summaries, versionData.summary)
    }

    return summaries, nil
}

// getVersionChanges retrieves the specific changes made in a particular version of the order
func (or *OrderRepository) getVersionChanges(ctx context.Context, tx pgx.Tx, orderID int64, version int32) ([]*order.OrderItemChange, error) {
    // Query to get all changes (both dishes and sets) for a specific version
    rows, err := tx.Query(ctx, `
        WITH version_changes AS (
            -- Get dish changes
            SELECT 
                'DISH' as item_type,
                d.id as item_id,
                d.name as item_name,
                di.quantity as quantity_changed,
                d.price as price
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
                s.price as price
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