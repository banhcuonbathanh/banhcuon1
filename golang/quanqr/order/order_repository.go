package order_grpc

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strings"

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
            COALESCE(o.order_name, '') as order_name,
            COALESCE(o.current_version, 0) as current_version,
            o.created_at,
            o.updated_at
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

    var detailedOrders []*order.OrderDetailedResponseWithDelivery
    for rows.Next() {
        var o order.OrderDetailedResponseWithDelivery
        
        // Create nullable variables for fields that can be NULL
        var (
            guestId         sql.NullInt64
            userId          sql.NullInt64
            tableNumber     sql.NullInt64
            orderHandlerId  sql.NullInt64
            totalPrice      sql.NullInt32
            status          sql.NullString
            topping         sql.NullString
            trackingOrder   sql.NullString
            chiliNumber     sql.NullInt64
            orderName       sql.NullString
            currentVersion  sql.NullInt32
            createdAt       sql.NullTime
            updatedAt       sql.NullTime
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
            &currentVersion,
            &createdAt,
            &updatedAt,
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
        if currentVersion.Valid {
            o.CurrentVersion = currentVersion.Int32
        }
        if createdAt.Valid {
            o.CreatedAt = timestamppb.New(createdAt.Time)
        }
        if updatedAt.Valid {
            o.UpdatedAt = timestamppb.New(updatedAt.Time)
        }
        
        // Initialize delivery-related fields
        o.DeliveryHistory = []*order.DishDelivery{}
        o.TotalItemsDelivered = 0
        
        // Fetch order version history
        versionQuery := `
            SELECT 
                version_number,
                modification_type,
                modified_at
            FROM order_versions
            WHERE order_id = $1
            ORDER BY version_number
        `
        versionRows, err := or.db.Query(ctx, versionQuery, o.Id)
        if err != nil {
            or.logger.Error("Error fetching version history: " + err.Error())
            return nil, fmt.Errorf("error fetching version history: %w", err)
        }
        defer versionRows.Close()

        o.VersionHistory = []*order.OrderVersionSummary{}
        for versionRows.Next() {
            var version order.OrderVersionSummary
            var modifiedAt time.Time
            
            err := versionRows.Scan(
                &version.VersionNumber,
                &version.ModificationType,
                &modifiedAt,
            )
            if err != nil {
                or.logger.Error("Error scanning version history: " + err.Error())
                return nil, fmt.Errorf("error scanning version history: %w", err)
            }
            
            version.ModifiedAt = timestamppb.New(modifiedAt)
            
            // Fetch dishes for this version
            versionDishQuery := `
                SELECT 
                    d.id,
                    vd.quantity,
                    d.name,
                    d.price,
                    d.description,
                    d.image,
                    d.status
                FROM version_dishes vd
                JOIN dishes d ON vd.dish_id = d.id
                WHERE vd.order_id = $1 AND vd.version_number = $2
            `
            versionDishRows, err := or.db.Query(ctx, versionDishQuery, o.Id, version.VersionNumber)
            if err != nil {
                or.logger.Error("Error fetching version dish details: " + err.Error())
                return nil, fmt.Errorf("error fetching version dish details: %w", err)
            }
            defer versionDishRows.Close()

            var versionDishes []*order.OrderDetailedDish
            for versionDishRows.Next() {
                var dish order.OrderDetailedDish
                err := versionDishRows.Scan(
                    &dish.DishId,
                    &dish.Quantity,
                    &dish.Name,
                    &dish.Price,
                    &dish.Description,
                    &dish.Image,
                    &dish.Status,
                )
                if err != nil {
                    or.logger.Error("Error scanning version dish detail: " + err.Error())
                    return nil, fmt.Errorf("error scanning version dish detail: %w", err)
                }
                versionDishes = append(versionDishes, &dish)
            }
            version.DishesOrdered = versionDishes
            
            // Fetch sets for this version
            versionSetQuery := `
                SELECT 
                    s.id,
                    s.name,
                    s.description,
                    s.user_id,
                    s.is_favourite,
                    s.is_public,
                    s.image,
                    s.price,
                    vs.quantity,
                    s.created_at,
                    s.updated_at
                FROM version_sets vs
                JOIN sets s ON vs.set_id = s.id
                WHERE vs.order_id = $1 AND vs.version_number = $2
            `
            versionSetRows, err := or.db.Query(ctx, versionSetQuery, o.Id, version.VersionNumber)
            if err != nil {
                or.logger.Error("Error fetching version set details: " + err.Error())
                return nil, fmt.Errorf("error fetching version set details: %w", err)
            }
            defer versionSetRows.Close()

            var versionSets []*order.OrderSetDetailed
            for versionSetRows.Next() {
                var set order.OrderSetDetailed
                var setCreatedAt, setUpdatedAt time.Time
                var userID sql.NullInt32
                
                err := versionSetRows.Scan(
                    &set.Id,
                    &set.Name,
                    &set.Description,
                    &userID,
                    &set.IsFavourite,
                    &set.IsPublic,
                    &set.Image,
                    &set.Price,
                    &set.Quantity,
                    &setCreatedAt,
                    &setUpdatedAt,
                )
                if err != nil {
                    or.logger.Error("Error scanning version set detail: " + err.Error())
                    return nil, fmt.Errorf("error scanning version set detail: %w", err)
                }

                if userID.Valid {
                    set.UserId = userID.Int32
                }
                
                set.CreatedAt = timestamppb.New(setCreatedAt)
                set.UpdatedAt = timestamppb.New(setUpdatedAt)
                
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
                versionSets = append(versionSets, &set)
            }
            version.SetOrdered = versionSets
            
            o.VersionHistory = append(o.VersionHistory, &version)
        }
        
        // Fetch delivery history
        deliveryQuery := `
            SELECT 
                id,
                order_id,
                order_name,
                guest_id,
                user_id,
                table_number,
                quantity_delivered,
                delivery_status,
                delivered_at,
                delivered_by_user_id,
                created_at,
                updated_at,
                dish_id,
                is_guest,
                modification_number
            FROM dish_deliveries
            WHERE order_id = $1
            ORDER BY delivered_at DESC
        `
        deliveryRows, err := or.db.Query(ctx, deliveryQuery, o.Id)
        if err != nil {
            or.logger.Error("Error fetching delivery history: " + err.Error())
            return nil, fmt.Errorf("error fetching delivery history: %w", err)
        }
        defer deliveryRows.Close()

        totalDelivered := int32(0)
        var lastDeliveryAt time.Time
        
        for deliveryRows.Next() {
            var delivery order.DishDelivery
            var (
                deliveredAt, createdAt, updatedAt sql.NullTime
                guestId, userId, tableNumber, deliveredById, dishId sql.NullInt64
                orderName sql.NullString
                quantityDelivered sql.NullInt32
                deliveryStatus sql.NullString
                modificationNumber sql.NullInt32
            )
            
            err := deliveryRows.Scan(
                &delivery.Id,
                &delivery.OrderId,
                &orderName,
                &guestId,
                &userId,
                &tableNumber,
                &quantityDelivered,
                &deliveryStatus,
                &deliveredAt,
                &deliveredById,
                &createdAt,
                &updatedAt,
                &dishId,
                &delivery.IsGuest,
                &modificationNumber,
            )
            if err != nil {
                or.logger.Error("Error scanning delivery history: " + err.Error())
                return nil, fmt.Errorf("error scanning delivery history: %w", err)
            }
            
            if orderName.Valid {
                delivery.OrderName = orderName.String
            }
            if guestId.Valid {
                delivery.GuestId = guestId.Int64
            }
            if userId.Valid {
                delivery.UserId = userId.Int64
            }
            if tableNumber.Valid {
                delivery.TableNumber = tableNumber.Int64
            }
            if quantityDelivered.Valid {
                delivery.QuantityDelivered = quantityDelivered.Int32
                totalDelivered += quantityDelivered.Int32
            }
            if deliveryStatus.Valid {
                delivery.DeliveryStatus = deliveryStatus.String
            }
            if deliveredAt.Valid {
                delivery.DeliveredAt = timestamppb.New(deliveredAt.Time)
                lastDeliveryAt = deliveredAt.Time
            }
            if deliveredById.Valid {
                delivery.DeliveredByUserId = deliveredById.Int64
            }
            if createdAt.Valid {
                delivery.CreatedAt = timestamppb.New(createdAt.Time)
            }
            if updatedAt.Valid {
                delivery.UpdatedAt = timestamppb.New(updatedAt.Time)
            }
            if dishId.Valid {
                delivery.DishId = dishId.Int64
            }
            if modificationNumber.Valid {
                delivery.ModificationNumber = modificationNumber.Int32
            }
            
            o.DeliveryHistory = append(o.DeliveryHistory, &delivery)
        }
        
        // Update delivery summary fields
        o.TotalItemsDelivered = totalDelivered
        if !lastDeliveryAt.IsZero() {
            o.LastDeliveryAt = timestamppb.New(lastDeliveryAt)
        }
        
        // Determine current delivery status based on history
        if len(o.DeliveryHistory) == 0 {
            o.CurrentDeliveryStatus = order.DeliveryStatus_PENDING
        } else {
            // Logic to determine overall delivery status
            // This is a simplified example - adjust based on your business rules
            allComplete := true
            for _, delivery := range o.DeliveryHistory {
                if delivery.DeliveryStatus != "DELIVERED" {
                    allComplete = false
                    break
                }
            }
            
            if allComplete {
                o.CurrentDeliveryStatus = order.DeliveryStatus_FULLY_DELIVERED
            } else {
                o.CurrentDeliveryStatus = order.DeliveryStatus_PARTIALLY_DELIVERED
            }
        }

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

    if err := tx.Commit(ctx); err != nil {
        or.logger.Error("Error committing transaction: " + err.Error())
        return nil, fmt.Errorf("error committing transaction: %w", err)
    }

    // Create the response with OrderDetailedResponseWithDelivery instead of OrderDetailedResponse
    return &order.OrderDetailedListResponse{
        Data: []*order.OrderDetailedResponseWithDelivery{
            {
                Id:             req.Id,
                GuestId:        req.GuestId,
                UserId:         req.UserId,
                TableNumber:    req.TableNumber,
                OrderHandlerId: req.OrderHandlerId,
                Status:         req.Status,
                TotalPrice:     req.TotalPrice,
                IsGuest:        req.IsGuest,
                Topping:        req.Topping,
                TrackingOrder:  req.TrackingOrder,
                TakeAway:       req.TakeAway,
                ChiliNumber:    req.ChiliNumber,
                TableToken:     req.TableToken,
                OrderName:      req.OrderName,
                CurrentVersion: newVersion,
                // Parent order ID is no longer in the struct based on your schema
                // ParentOrderId:  req.ParentOrderId,
                CreatedAt:      timestamppb.New(createdAt),
                UpdatedAt:      timestamppb.New(updatedAt),
                // Initialize empty slices for the new fields
                VersionHistory:       []*order.OrderVersionSummary{},
                DeliveryHistory:      []*order.DishDelivery{},
                CurrentDeliveryStatus: order.DeliveryStatus_PENDING,
                TotalItemsDelivered:   0,
                // LastDeliveryAt is omitted since there's no delivery yet
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




// Improved version of getVersionChanges to include price information
// func (or *OrderRepository) getVersionChanges(ctx context.Context, tx pgx.Tx, orderID int64, version int32) ([]*order.OrderItemChange, error) {
//     rows, err := tx.Query(ctx, `
//         WITH version_changes AS (
//             -- Get dish changes
//             SELECT 
//                 'DISH' as item_type,
//                 d.id as item_id,
//                 d.name as item_name,
//                 di.quantity as quantity_changed,
//                 d.price
//             FROM dish_order_items di
//             JOIN dishes d ON di.dish_id = d.id
//             WHERE di.order_id = $1 AND di.modification_number = $2
            
//             UNION ALL
            
//             -- Get set changes
//             SELECT 
//                 'SET' as item_type,
//                 s.id as item_id,
//                 s.name as item_name,
//                 si.quantity as quantity_changed,
//                 s.price
//             FROM set_order_items si
//             JOIN sets s ON si.set_id = s.id
//             WHERE si.order_id = $1 AND si.modification_number = $2
//         )
//         SELECT 
//             item_type,
//             item_id,
//             item_name,
//             quantity_changed,
//             price
//         FROM version_changes
//         ORDER BY item_type, item_name`,
//         orderID, version)
//     if err != nil {
//         return nil, fmt.Errorf("error querying version changes: %w", err)
//     }
//     defer rows.Close()

//     // var changes []*order.OrderItemChange
//     // for rows.Next() {
//     //     var change order.OrderItemChange
//     //     err := rows.Scan(
//     //         &change.ItemType,
//     //         &change.ItemId,
//     //         &change.ItemName,
//     //         &change.QuantityChanged,
//     //         &change.Price,
//     //     )
//     //     if err != nil {
//     //         return nil, fmt.Errorf("error scanning version change: %w", err)
//     //     }
//     //     changes = append(changes, &change)
//     // }

//     return , nil
// }

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
                // TotalDishesCount: dishesCount,
                // TotalSetsCount:   setsCount,
                // VersionTotalPrice: versionTotal,
                // Changes:          make([]*order.OrderItemChange, 0),
            }
            summaryMap[versionNum] = summary
        }

        // Add change if present
        // if itemType.Valid {
        //     change := &order.OrderItemChange{
        //         ItemType:       itemType.String,
        //         ItemId:        itemID.Int64,
        //         ItemName:      itemName.String,
        //         QuantityChanged: quantityChanged.Int32,
        //         Price:         price.Int32,
        //     }
        //     summary.Changes = append(summary.Changes, change)
        // }
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
        LIMIT 1`, // Limiting to top 5 most ordered items
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

    // Fetch version history
    versionHistory, err := or.fetchVersionHistory(ctx, tx, req.Id)
    if err != nil {
        or.logger.Error(fmt.Sprintf("[OrderRepository.AddingSetsDishesOrder] Error fetching version history: %s", err.Error()))
        return nil, fmt.Errorf("error fetching version history: %w", err)
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

    // Fetch delivery history - similar to FetchOrdersByCriteria
    deliveryHistory, lastDeliveryAt, totalItemsDelivered, deliveryStatus, err := or.getOrderDeliveryHistory(ctx, req.Id)
    if err != nil {
        or.logger.Error(fmt.Sprintf("[OrderRepository.AddingSetsDishesOrder] Error fetching delivery history: %s", err.Error()))
        // Continue even if delivery history fails
    }

    if err := tx.Commit(ctx); err != nil {
        or.logger.Error(fmt.Sprintf("[OrderRepository.AddingSetsDishesOrder] Commit failed: %s", err.Error()))
        return nil, fmt.Errorf("transaction commit failed: %w", err)
    }

    or.logger.Info(fmt.Sprintf("[OrderRepository.AddingSetsDishesOrder] Update successful. OrderID: %d, NewVersion: %d", 
        req.Id, newVersion))

    // Create timestamps
    nowProto := timestamppb.New(now)

    return &order.OrderDetailedListResponse{
        Data: []*order.OrderDetailedResponseWithDelivery{
            {
                Id:             req.Id,
                GuestId:        req.GuestId,
                UserId:         req.UserId,
                TableNumber:    req.TableNumber,
                OrderHandlerId: req.OrderHandlerId,
                Status:         req.Status,
                TotalPrice:     req.TotalPrice,
                IsGuest:        isGuest,
                Topping:        req.Topping,
                TrackingOrder:  req.TrackingOrder,
                TakeAway:       req.TakeAway,
                ChiliNumber:    req.ChiliNumber,
                TableToken:     req.TableToken,
                OrderName:      req.OrderName,
       
                CurrentVersion: newVersion,
                VersionHistory: versionHistory,
             
                // Add delivery-related fields
                DeliveryHistory:       deliveryHistory,
                LastDeliveryAt:        lastDeliveryAt,
                TotalItemsDelivered:   totalItemsDelivered,
                CurrentDeliveryStatus: deliveryStatus,
                // Add timestamp fields
                CreatedAt: nowProto,
                UpdatedAt: nowProto,
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

    // Fetch delivery history
    deliveryHistory, lastDeliveryAt, totalItemsDelivered, deliveryStatus, err := or.getOrderDeliveryHistory(ctx, req.Id)
    if err != nil {
        or.logger.Error(fmt.Sprintf("[OrderRepository.RemovingSetsDishesOrder] Error fetching delivery history: %s", err.Error()))
        // Continue even if delivery history fails
    }

    if err := tx.Commit(ctx); err != nil {
        or.logger.Error(fmt.Sprintf("[OrderRepository.RemovingSetsDishesOrder] Commit failed: %s", err.Error()))
        return nil, fmt.Errorf("transaction commit failed: %w", err)
    }

    or.logger.Info(fmt.Sprintf("[OrderRepository.RemovingSetsDishesOrder] Removal successful. OrderID: %d, NewVersion: %d", 
        req.Id, newVersion))

    // Create timestamps
    nowProto := timestamppb.New(now)

    return &order.OrderDetailedListResponse{
        Data: []*order.OrderDetailedResponseWithDelivery{
            {
                Id:             req.Id,
                GuestId:        req.GuestId,
                UserId:         req.UserId,
                TableNumber:    req.TableNumber,
                OrderHandlerId: req.OrderHandlerId,
                Status:         req.Status,
                TotalPrice:     req.TotalPrice,
                IsGuest:        isGuest,
                Topping:        req.Topping,
                TrackingOrder:  req.TrackingOrder,
                TakeAway:       req.TakeAway,
                ChiliNumber:    req.ChiliNumber,
                TableToken:     req.TableToken,
                OrderName:      req.OrderName,
        
                CurrentVersion: newVersion,
                VersionHistory: versionHistory,
             
                // Add delivery-related fields
                DeliveryHistory:       deliveryHistory,
                LastDeliveryAt:        lastDeliveryAt,
                TotalItemsDelivered:   totalItemsDelivered,
                CurrentDeliveryStatus: deliveryStatus,
                // Add timestamp fields
                CreatedAt: nowProto,
                UpdatedAt: nowProto,
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
    )
    
    // var validOrderStatuses = map[string]bool{
    //     "IN_PROGRESS": true,
    //     "PENDING":     true,
    // }

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
        userID           sql.NullInt64
        tableNumber      sql.NullInt64
        parentOrderID    sql.NullInt64  // Add this to handle NULL parent_order_id
    )
    err = tx.QueryRow(ctx, `
    SELECT 
        id, version, is_guest, guest_id, user_id, table_number, status,
        total_price, order_handler_id, topping, tracking_order, take_away,
        chili_number, table_token, order_name, COALESCE(parent_order_id, 0)
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
        &parentOrderID,
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
    currentOrder.ParentOrderId = nullInt64ToProtoInt64(parentOrderID)
    
    // Validate order state
    // if !validOrderStatuses[currentOrder.Status] {
    //     or.logger.Warning(fmt.Sprintf(
    //         "[MarkDishesDelivered] Invalid order status. OrderID: %d, Status: %s",
    //         req.OrderId, currentOrder.Status))
    //     return nil, fmt.Errorf("order cannot be modified in current status: %s", currentOrder.Status)
    // }

    // Strict version check


    newVersion := currentOrder.Version + 1
    now := time.Now().UTC()

    // Process dish deliveries


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


    if err != nil {
        or.logger.Error(fmt.Sprintf("[MarkDishesDelivered] Total calculation failed: %s", err.Error()))
        return nil, fmt.Errorf("error calculating totals: %w", err)
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

fmt.Printf("\n=== Delivery History golang/quanqr/order/order_repository.go for Order ID: %d ===\n", req.OrderId)

for i, delivery := range deliveryHistory {
    fmt.Printf("\nDelivery #%d:\n", i+1)
    fmt.Printf("  Modification Number: %d\n", delivery.ModificationNumber)
    fmt.Printf("  Quantity Delivered: %d\n", delivery.QuantityDelivered)

    fmt.Printf("  ---------------------------\n")
}

        return &order.OrderDetailedResponseWithDelivery{
            Id:                  currentOrder.Id,
            GuestId:             currentOrder.GuestId,
            UserId:              currentOrder.UserId,
            TableNumber:         currentOrder.TableNumber,
            OrderHandlerId:      currentOrder.OrderHandlerId,
            Status:              currentOrder.Status,
            TotalPrice:          currentOrder.TotalPrice,
    
            IsGuest:             currentOrder.IsGuest,
            Topping:             currentOrder.Topping,
            TrackingOrder:       currentOrder.TrackingOrder,
            TakeAway:            currentOrder.TakeAway,
            ChiliNumber:         currentOrder.ChiliNumber,
            TableToken:          currentOrder.TableToken,
            OrderName:           currentOrder.OrderName,
            CurrentVersion:      newVersion,
  
            VersionHistory:      versionHistory,

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

// craete order start 

// create order end 




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

    deliveries := make([]*order.DishDelivery, 0)
    for rows.Next() {
        var (
            dd                  order.DishDelivery
            guestID            sql.NullInt64
            userID             sql.NullInt64
            tableNumber        sql.NullInt64
            dishID             int64
            quantityDelivered  int32
            deliveredAt        time.Time
            createdAt         time.Time
            updatedAt         time.Time
            modNumber         int32
        )
        err := rows.Scan(
            &dd.Id,
            &dd.OrderId,
            &dd.OrderName,
            &guestID,
            &userID,
            &tableNumber,
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

        // Handle nullable fields
        dd.GuestId = nullInt64ToProtoInt64(guestID)
        dd.UserId = nullInt64ToProtoInt64(userID)
        dd.TableNumber = nullInt64ToProtoInt64(tableNumber)

        // Convert timestamps
        dd.DeliveredAt = timestamppb.New(deliveredAt)
        dd.CreatedAt = timestamppb.New(createdAt)
        dd.UpdatedAt = timestamppb.New(updatedAt)

        // Assign modification number directly to DishDelivery
        dd.ModificationNumber = modNumber
        
        // Assign quantity delivered directly to DishDelivery
        dd.QuantityDelivered = quantityDelivered

        // Build DishOrderItem with delivery context


        deliveries = append(deliveries, &dd)
    }
    
    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("error after scanning rows: %w", err)
    }

    return deliveries, nil
}


// new for AddingDishesToOrder start

func (or *OrderRepository) AddingDishesToOrder(ctx context.Context, req *order.CreateDishOrderItemWithOrderID) (*order.DishOrderItem, error) {
    or.logger.Info(fmt.Sprintf("golang/quanqr/order/order_repository.go Adding dishes to order ID %d: dish_id=%d, quantity=%d", req.OrderId, req.DishId, req.Quantity))

    tx, err := or.db.Begin(ctx)
    if err != nil {
        or.logger.Error("Error starting transaction: " + err.Error())
        return nil, fmt.Errorf("error starting transaction: %w", err)
    }
    defer tx.Rollback(ctx)

    // First, verify the order exists and get its current version
    var currentVersion int32
    var isGuest bool
    var orderHandlerId int64
    err = tx.QueryRow(ctx, `
        SELECT version, is_guest, order_handler_id 
        FROM orders 
        WHERE id = $1`,
        req.OrderId,
    ).Scan(&currentVersion, &isGuest, &orderHandlerId)
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return nil, fmt.Errorf("order not found: %d", req.OrderId)
        }
        return nil, fmt.Errorf("error fetching order details: %w", err)
    }

    // Verify the dish exists
    var exists bool
    err = tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM dishes WHERE id = $1)", req.DishId).Scan(&exists)
    if err != nil {
        return nil, fmt.Errorf("error verifying dish existence: %w", err)
    }
    if !exists {
        return nil, fmt.Errorf("dish with id %d does not exist", req.DishId)
    }

    now := time.Now()
    nextVersion := currentVersion + 1

    // Update the order's version
    _, err = tx.Exec(ctx, `
        UPDATE orders 
        SET version = $1, 
            updated_at = $2 
        WHERE id = $3`,
        nextVersion,
        now,
        req.OrderId,
    )
    if err != nil {
        return nil, fmt.Errorf("error updating order version: %w", err)
    }

    // Create a new order modification record
    _, err = tx.Exec(ctx, `
        INSERT INTO order_modifications (
            order_id, modification_number, modification_type, 
            modified_by_user_id, order_name, modified_at
        )
        VALUES ($1, $2, $3, $4, $5, $6)`,
        req.OrderId,
        nextVersion,
        "ADD",
        orderHandlerId,
        req.OrderName,
        now,
    )
    if err != nil {
        return nil, fmt.Errorf("error creating modification record: %w", err)
    }

    // Insert the new dish order item
    var newItemId int64
    err = tx.QueryRow(ctx, `
        INSERT INTO dish_order_items (
            order_id, dish_id, quantity, created_at, updated_at,
            order_name, modification_type, modification_number
        )
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
        RETURNING id`,
        req.OrderId,
        req.DishId,
        req.Quantity,
        now,
        now,
        req.OrderName,
        "ADD",
        nextVersion,
    ).Scan(&newItemId)
    if err != nil {
        return nil, fmt.Errorf("error inserting dish order item: %w", err)
    }

    if err := tx.Commit(ctx); err != nil {
        return nil, fmt.Errorf("error committing transaction: %w", err)
    }

    // Return the created dish order item
    return &order.DishOrderItem{
        Id:                newItemId,
        DishId:           req.DishId,
        Quantity:         req.Quantity,
        CreatedAt:        timestamppb.New(now),
        UpdatedAt:        timestamppb.New(now),
        OrderName:        req.OrderName,
        ModificationType: "ADD",
        ModificationNumber: nextVersion,
    }, nil
}

// new for AddingDishesToOrder end 

// new for AddingSetToOrder start

func (or *OrderRepository) AddingSetToOrder(ctx context.Context, req *order.CreateSetOrderItemWithOrderID) (*order.ResponseSetOrderItemWithOrderID, error) {
    or.logger.Info(fmt.Sprintf(" golang/quanqr/order/order_repository.go Adding set to order ID %d: set_id=%d, quantity=%d", req.OrderId, req.SetId, req.Quantity))

    tx, err := or.db.Begin(ctx)
    if err != nil {
        or.logger.Error("Error starting transaction: " + err.Error())
        return nil, fmt.Errorf("error starting transaction: %w", err)
    }
    defer tx.Rollback(ctx)

    // First, verify the order exists and get its current version
    var currentVersion int32
    var isGuest bool
    var orderHandlerId int64
    err = tx.QueryRow(ctx, `
        SELECT version, is_guest, order_handler_id 
        FROM orders 
        WHERE id = $1`,
        req.OrderId,
    ).Scan(&currentVersion, &isGuest, &orderHandlerId)
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return nil, fmt.Errorf("order not found: %d", req.OrderId)
        }
        return nil, fmt.Errorf("error fetching order details: %w", err)
    }

    // Verify the set exists and get its dishes
    var exists bool
    err = tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM sets WHERE id = $1)", req.SetId).Scan(&exists)
    if err != nil {
        return nil, fmt.Errorf("error verifying set existence: %w", err)
    }
    if !exists {
        return nil, fmt.Errorf("set with id %d does not exist", req.SetId)
    }

    now := time.Now()
    nextVersion := currentVersion + 1

    // Update the order's version
    _, err = tx.Exec(ctx, `
        UPDATE orders 
        SET version = $1, 
            updated_at = $2 
        WHERE id = $3`,
        nextVersion,
        now,
        req.OrderId,
    )
    if err != nil {
        return nil, fmt.Errorf("error updating order version: %w", err)
    }

    // Create a new order modification record
    _, err = tx.Exec(ctx, `
        INSERT INTO order_modifications (
            order_id, modification_number, modification_type, 
            modified_by_user_id, order_name, modified_at
        )
        VALUES ($1, $2, $3, $4, $5, $6)`,
        req.OrderId,
        nextVersion,
        "ADD",
        orderHandlerId,
        req.OrderName,
        now,
    )
    if err != nil {
        return nil, fmt.Errorf("error creating modification record: %w", err)
    }

    // Insert the new set order item
    var newSetItemId int64
    err = tx.QueryRow(ctx, `
        INSERT INTO set_order_items (
            order_id, set_id, quantity, created_at, updated_at,
            order_name, modification_type, modification_number
        )
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
        RETURNING id`,
        req.OrderId,
        req.SetId,
        req.Quantity,
        now,
        now,
        req.OrderName,
        "ADD",
        nextVersion,
    ).Scan(&newSetItemId)
    if err != nil {
        return nil, fmt.Errorf("error inserting set order item: %w", err)
    }

    // Get the set's dishes
    rows, err := tx.Query(ctx, `
        SELECT 
            d.id,
            d.name,
            d.price,
            d.description,
            d.image,
            d.status,
            sd.quantity
        FROM set_dishes sd
        JOIN dishes d ON sd.dish_id = d.id
        WHERE sd.set_id = $1`,
        req.SetId,
    )
    if err != nil {
        return nil, fmt.Errorf("error fetching set dishes: %w", err)
    }
    defer rows.Close()

    var dishes []*order.OrderDetailedDish
    for rows.Next() {
        var dish order.OrderDetailedDish
        var dishQuantity int64
        err := rows.Scan(
            &dish.DishId,
            &dish.Name,
            &dish.Price,
            &dish.Description,
            &dish.Image,
            &dish.Status,
            &dishQuantity,
        )
        if err != nil {
            return nil, fmt.Errorf("error scanning dish row: %w", err)
        }
        // Multiply the dish quantity by the set quantity ordered
        dish.Quantity = dishQuantity * req.Quantity
        dishes = append(dishes, &dish)
    }

    if err := tx.Commit(ctx); err != nil {
        return nil, fmt.Errorf("error committing transaction: %w", err)
    }

    // Prepare the response
    setItem := &order.SetOrderItem{
        Id:                newSetItemId,
        SetId:            req.SetId,
        Quantity:         req.Quantity,
        CreatedAt:        timestamppb.New(now),
        UpdatedAt:        timestamppb.New(now),
        OrderName:        req.OrderName,
        ModificationType: "ADD",
        ModificationNumber: nextVersion,
    }

    return &order.ResponseSetOrderItemWithOrderID{
        Set:    setItem,
        Dishes: dishes,
    }, nil
}

// new for AddingSetToOrder end


// create order start 
func (or *OrderRepository) CreateOrder(ctx context.Context, req *order.CreateOrderRequest) (*order.OrderDetailedResponseWithDelivery, error) {
    or.logger.Info(fmt.Sprintf("golang/quanqr/order/order_repository.go Creating new order: %+v", req))
    
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

    var orderId int64
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
    ).Scan(&orderId, &createdAt, &updatedAt)

    if err != nil {
        or.logger.Error("Error creating order: " + err.Error())
        return nil, fmt.Errorf("error creating order: %w", err)
    }

    // Create initial order modification record
    modificationQuery := `
        INSERT INTO order_modifications (
            order_id, modification_number, modification_type, 
            modified_by_user_id, order_name, modified_at
        )
        VALUES ($1, $2, $3, $4, $5, $6)
    `
    _, err = tx.Exec(ctx, modificationQuery,
        orderId,
        1, // First modification
        "ADD",
        req.OrderHandlerId,
        req.OrderName,
        now, // Adding modified_at timestamp
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
            orderId, dish.DishId, dish.Quantity, now, now, 
            req.OrderName, "INITIAL", 1)
        if err != nil {
            or.logger.Error(fmt.Sprintf("Error inserting order dish: %s", err.Error()))
            return nil, fmt.Errorf("error inserting order dish: %w", err)
        }

        // We no longer create dish_deliveries records here - they will be created when actually delivering dishes
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
            orderId, set.SetId, set.Quantity, now, now, 
            req.OrderName, "INITIAL", 1)
        if err != nil {
            or.logger.Error(fmt.Sprintf("Error inserting order set: %s", err.Error()))
            return nil, fmt.Errorf("error inserting order set: %w", err)
        }
    }

    // Get version history
    versionHistory, err := or.getOrderVersionHistory(ctx, tx, orderId)
    if err != nil {
        or.logger.Error("Error getting version history: " + err.Error())
        return nil, fmt.Errorf("error getting version history: %w", err)
    }

    if err := tx.Commit(ctx); err != nil {
        or.logger.Error("Error committing transaction: " + err.Error())
        return nil, fmt.Errorf("error committing transaction: %w", err)
    }

    // Populate detailed response with delivery information
    response := &order.OrderDetailedResponseWithDelivery{
        Id:                 orderId,
        GuestId:            req.GuestId,
        UserId:             req.UserId,
        IsGuest:            req.IsGuest,
        TableNumber:        req.TableNumber,
        OrderHandlerId:     req.OrderHandlerId,
        Status:             req.Status,
        CreatedAt:          timestamppb.New(createdAt),
        UpdatedAt:          timestamppb.New(updatedAt),
        TotalPrice:         req.TotalPrice,
        Topping:            req.Topping,
        TrackingOrder:      trackingOrder,
        TakeAway:           req.TakeAway,
        ChiliNumber:        req.ChiliNumber,
        TableToken:         req.TableToken,
        OrderName:          req.OrderName,
        CurrentVersion:     1, // Initial version
        VersionHistory:     versionHistory,
        DeliveryHistory:    []*order.DishDelivery{}, // Empty for new orders
        CurrentDeliveryStatus: order.DeliveryStatus_PENDING, // New orders haven't started delivery
        TotalItemsDelivered: 0, // No items delivered for new orders
        LastDeliveryAt:     nil, // No delivery yet
    }
    or.logger.Info(fmt.Sprintf("golang/quanqr/order/order_repository.go 1212121 Created order successfully: %+v", response))
    return response, nil
}
// Updated helper function to get order version history with dishes and sets
func (or *OrderRepository) getOrderVersionHistory(ctx context.Context, tx pgx.Tx, orderId int64) ([]*order.OrderVersionSummary, error) {
    // First, get the basic modification info
    modQuery := `
    SELECT 
        modification_number, 
        modification_type, 
        modified_by_user_id, 
        modified_at
    FROM 
        order_modifications
    WHERE 
        order_id = $1
    ORDER BY 
        modification_number ASC
    `
    
    modRows, err := tx.Query(ctx, modQuery, orderId)
    if err != nil {
        return nil, fmt.Errorf("error querying order modifications: %w", err)
    }
    defer modRows.Close()
    
    // Collect all modification data first to avoid nested queries with same transaction
    var modificationsData []struct {
        versionNumber     int32
        modificationType  string
        modifiedByUserId  sql.NullInt64
        modifiedAt        time.Time
    }
    
    for modRows.Next() {
        var modData struct {
            versionNumber     int32
            modificationType  string
            modifiedByUserId  sql.NullInt64
            modifiedAt        time.Time
        }
        
        if err := modRows.Scan(
            &modData.versionNumber,
            &modData.modificationType,
            &modData.modifiedByUserId,
            &modData.modifiedAt,
        ); err != nil {
            return nil, fmt.Errorf("error scanning order modification row: %w", err)
        }
        
        modificationsData = append(modificationsData, modData)
    }
    
    if err := modRows.Err(); err != nil {
        return nil, fmt.Errorf("error iterating order modifications: %w", err)
    }
    
    // Now process each version - important: modRows is now closed
    var versions []*order.OrderVersionSummary
    
    for _, modData := range modificationsData {
        version := &order.OrderVersionSummary{
            VersionNumber:    modData.versionNumber,
            ModificationType: modData.modificationType,
            ModifiedAt:       timestamppb.New(modData.modifiedAt),
        }
        
        // For each version, get the associated dishes
        dishQuery := `
        SELECT 
            doi.dish_id,
            d.name,
            doi.quantity,
            d.price,
            d.description,
            d.image
        FROM 
            dish_order_items doi
        JOIN 
            dishes d ON doi.dish_id = d.id
        WHERE 
            doi.order_id = $1 AND
            doi.modification_number = $2
        `
        
        dishRows, err := tx.Query(ctx, dishQuery, orderId, version.VersionNumber)
        if err != nil {
            return nil, fmt.Errorf("error querying version dishes: %w", err)
        }
        
        var dishes []*order.OrderDetailedDish
        for dishRows.Next() {
            var dish order.OrderDetailedDish
            var dishId int64
            var quantity int64
            
            if err := dishRows.Scan(
                &dishId,
                &dish.Name,
                &quantity,
                &dish.Price,
                &dish.Description,
                &dish.Image,
            ); err != nil {
                dishRows.Close()
                return nil, fmt.Errorf("error scanning dish row: %w", err)
            }
            
            dish.DishId = dishId
            dish.Quantity = quantity
            dish.Price = dish.Price * int32(quantity)
            dish.Status = ""
            
            dishes = append(dishes, &dish)
        }
        
        // Important: Close rows before starting a new query
        dishRows.Close()
        
        if err := dishRows.Err(); err != nil {
            return nil, fmt.Errorf("error iterating dish rows: %w", err)
        }
        
        // For each version, get the associated sets
        setQuery := `
        SELECT 
            soi.set_id,
            s.name,
            soi.quantity,
            s.price,
            s.description,
            s.image,
            s.is_public,
            s.user_id,
            s.created_at,
            s.updated_at
        FROM 
            set_order_items soi
        JOIN 
            sets s ON soi.set_id = s.id
        WHERE 
            soi.order_id = $1 AND
            soi.modification_number = $2
        `
        
        setRows, err := tx.Query(ctx, setQuery, orderId, version.VersionNumber)
        if err != nil {
            return nil, fmt.Errorf("error querying version sets: %w", err)
        }
        
        // Get all set data first before querying for set dishes
        var setsData []struct {
            set       *order.OrderSetDetailed
            setId     int64
        }
        
        for setRows.Next() {
            var set order.OrderSetDetailed
            var setId int64
            var quantity int64
            var createdAt, updatedAt time.Time
            
            if err := setRows.Scan(
                &setId,
                &set.Name,
                &quantity,
                &set.Price,
                &set.Description,
                &set.Image,
                &set.IsPublic,
                &set.UserId,
                &createdAt,
                &updatedAt,
            ); err != nil {
                setRows.Close()
                return nil, fmt.Errorf("error scanning set row: %w", err)
            }
            
            set.Id = setId
            set.Quantity = quantity
            set.Price = set.Price * int32(quantity)
            set.CreatedAt = timestamppb.New(createdAt)
            set.UpdatedAt = timestamppb.New(updatedAt)
            
            setsData = append(setsData, struct {
                set   *order.OrderSetDetailed
                setId int64
            }{
                set:   &set,
                setId: setId,
            })
        }
        
        // Important: Close the set rows before querying for set dishes
        setRows.Close()
        
        if err := setRows.Err(); err != nil {
            return nil, fmt.Errorf("error iterating set rows: %w", err)
        }
        
        // Now that setRows is closed, we can process sets and their dishes
        var sets []*order.OrderSetDetailed
        
        for _, setData := range setsData {
            // Fetch dishes in this set
            setDishesQuery := `
            SELECT 
                d.id,
                d.name,
                d.price,
                d.description,
                d.image,
                sd.quantity 
            FROM 
                set_dishes sd
            JOIN 
                dishes d ON sd.dish_id = d.id
            WHERE 
                sd.set_id = $1
            `
            
            setDishRows, err := tx.Query(ctx, setDishesQuery, setData.setId)
            if err != nil {
                return nil, fmt.Errorf("error querying set dishes: %w", err)
            }
            
            var setDishes []*order.OrderDetailedDish
            for setDishRows.Next() {
                var setDish order.OrderDetailedDish
                var dishId int64
                var dishQuantity int64
                
                if err := setDishRows.Scan(
                    &dishId,
                    &setDish.Name,
                    &setDish.Price,
                    &setDish.Description,
                    &setDish.Image,
                    &dishQuantity,
                ); err != nil {
                    setDishRows.Close()
                    return nil, fmt.Errorf("error scanning set dish row: %w", err)
                }
                
                setDish.DishId = dishId
                setDish.Quantity = dishQuantity
                setDish.Status = ""
                
                setDishes = append(setDishes, &setDish)
            }
            
            // Important: Close rows before proceeding
            setDishRows.Close()
            
            if err := setDishRows.Err(); err != nil {
                return nil, fmt.Errorf("error iterating set dish rows: %w", err)
            }
            
            // Assign the dishes to the set
            setData.set.Dishes = setDishes
            
            sets = append(sets, setData.set)
        }
        
        // Add dishes and sets to the version summary
        version.DishesOrdered = dishes
        version.SetOrdered = sets
        
        versions = append(versions, version)
    }
    
    return versions, nil
}
// create order end


// start of fetch order with criterial 

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
                o.created_at,
                o.updated_at,
                COALESCE(o.topping, '') as topping,
                COALESCE(o.tracking_order, '') as tracking_order,
                COALESCE(o.take_away, false) as take_away,
                COALESCE(o.chili_number, 0) as chili_number,
                o.table_token,
                COALESCE(o.order_name, '') as order_name,
                COALESCE(o.version, 1) as current_version,
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

    var detailedOrders []*order.OrderDetailedResponseWithDelivery
    var totalItems int64
    for rows.Next() {
        var o order.OrderDetailedResponseWithDelivery
        var (
            guestId        sql.NullInt64
            userId         sql.NullInt64
            tableNumber    sql.NullInt64
            orderHandlerId sql.NullInt64
            totalPrice     sql.NullInt32
            status         sql.NullString
            createdAt      time.Time
            updatedAt      time.Time
            topping        sql.NullString
            trackingOrder  sql.NullString
            chiliNumber    sql.NullInt64
            orderName      sql.NullString
            currentVersion sql.NullInt32
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
            &createdAt,
            &updatedAt,
            &topping,
            &trackingOrder,
            &o.TakeAway,
            &chiliNumber,
            &o.TableToken,
            &orderName,
            &currentVersion,
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
        if currentVersion.Valid {
            o.CurrentVersion = currentVersion.Int32
        } else {
            o.CurrentVersion = 1 // Default to version 1 if not specified
        }

        // Set timestamps
        o.CreatedAt = timestamppb.New(createdAt)
        o.UpdatedAt = timestamppb.New(updatedAt)

        // Set default delivery status
        o.CurrentDeliveryStatus = order.DeliveryStatus_PENDING
        o.TotalItemsDelivered = 0

        // Fetch version history for this order using a new transaction
        // Since getOrderVersionHistory expects a pgx.Tx but we have a pgxpool.Pool
        tx, err := or.db.Begin(ctx)
        if err != nil {
            or.logger.Error("Error starting transaction for version history: " + err.Error())
            // Continue without version history
        } else {
            versionHistory, err := or.getOrderVersionHistory(ctx, tx, o.Id)
            if err != nil {
                or.logger.Error("Error fetching version history: " + err.Error())
                tx.Rollback(ctx) // Rollback on error
            } else {
                o.VersionHistory = versionHistory
                tx.Commit(ctx) // Commit if successful
            }
        }

        // Fetch delivery history for this order
        deliveryHistory, lastDeliveryAt, totalItemsDelivered, deliveryStatus, err := or.getOrderDeliveryHistory(ctx, o.Id)
        if err != nil {
            or.logger.Error("Error fetching delivery history: " + err.Error())
            // Continue processing other orders even if delivery history fails for one
        } else {
            o.DeliveryHistory = deliveryHistory
            o.LastDeliveryAt = lastDeliveryAt
            o.TotalItemsDelivered = totalItemsDelivered
            o.CurrentDeliveryStatus = deliveryStatus
        }

        detailedOrders = append(detailedOrders, &o)
    }

    totalPages := int32(math.Ceil(float64(totalItems) / float64(req.PageSize)))

    // Create the response with the correct type
    response := &order.OrderDetailedListResponse{
        Data: detailedOrders,
        Pagination: &order.PaginationInfo{
            CurrentPage: req.Page,
            TotalPages:  totalPages,
            TotalItems:  totalItems,
            PageSize:    req.PageSize,
        },
    }

    return response, nil
}
func (or *OrderRepository) getOrderDeliveryHistory(ctx context.Context, orderId int64) ([]*order.DishDelivery, *timestamppb.Timestamp, int32, order.DeliveryStatus, error) {
    // Query to get delivery history
    query := `
        SELECT 
            id,
            order_id,
            order_name,
            guest_id,
            user_id,
            table_number,
            quantity_delivered,
            delivery_status,
            delivered_at,
            delivered_by_user_id,
            created_at,
            updated_at,
            dish_id,
            is_guest,
            modification_number
        FROM 
            dish_deliveries
        WHERE 
            order_id = $1
        ORDER BY 
            delivered_at ASC
    `
    
    rows, err := or.db.Query(ctx, query, orderId)
    if err != nil {
        return nil, nil, 0, order.DeliveryStatus_PENDING, fmt.Errorf("error querying delivery history: %w", err)
    }
    defer rows.Close()
    
    var deliveries []*order.DishDelivery
    var lastDeliveryTime time.Time
    var totalDelivered int32
    var lastStatus string
    
    for rows.Next() {
        var delivery order.DishDelivery
        var deliveredAt, createdAt, updatedAt time.Time
        var deliveredByUserId sql.NullInt64
        var guestId sql.NullInt64
        var userId sql.NullInt64
        var tableNumber sql.NullInt64
        var isGuest sql.NullBool
        
        if err := rows.Scan(
            &delivery.Id,
            &delivery.OrderId,
            &delivery.OrderName,
            &guestId,
            &userId,
            &tableNumber,
            &delivery.QuantityDelivered,
            &delivery.DeliveryStatus,
            &deliveredAt,
            &deliveredByUserId,
            &createdAt,
            &updatedAt,
            &delivery.DishId,
            &isGuest,
            &delivery.ModificationNumber,
        ); err != nil {
            return nil, nil, 0, order.DeliveryStatus_PENDING, fmt.Errorf("error scanning delivery row: %w", err)
        }
        
        // Handle NULL values
        if deliveredByUserId.Valid {
            delivery.DeliveredByUserId = deliveredByUserId.Int64
        }
        
        if guestId.Valid {
            delivery.GuestId = guestId.Int64
        }
        
        if userId.Valid {
            delivery.UserId = userId.Int64
        }
        
        if tableNumber.Valid {
            delivery.TableNumber = tableNumber.Int64
        }
        
        if isGuest.Valid {
            delivery.IsGuest = isGuest.Bool
        }
        
        delivery.DeliveredAt = timestamppb.New(deliveredAt)
        delivery.CreatedAt = timestamppb.New(createdAt)
        delivery.UpdatedAt = timestamppb.New(updatedAt)
        
        // Track total delivered
        totalDelivered += delivery.QuantityDelivered
        
        // Track the latest delivery time and status
        if deliveredAt.After(lastDeliveryTime) {
            lastDeliveryTime = deliveredAt
            lastStatus = delivery.DeliveryStatus
        }
        
        deliveries = append(deliveries, &delivery)
    }
    
    if err := rows.Err(); err != nil {
        return nil, nil, 0, order.DeliveryStatus_PENDING, fmt.Errorf("error iterating delivery rows: %w", err)
    }
    
    // Determine overall delivery status
    var overallStatus order.DeliveryStatus
    switch strings.ToUpper(lastStatus) {
    case "PENDING":
        overallStatus = order.DeliveryStatus_PENDING
    case "IN_PROGRESS":
        overallStatus = order.DeliveryStatus_PARTIALLY_DELIVERED
    case "DELIVERED":
        overallStatus = order.DeliveryStatus_FULLY_DELIVERED
    case "CANCELLED":
        overallStatus = order.DeliveryStatus_CANCELLED
    default:
        overallStatus = order.DeliveryStatus_PENDING
    }
    
    // Create timestamp for last delivery
    var lastDeliveryTimestamp *timestamppb.Timestamp
    if !lastDeliveryTime.IsZero() {
        lastDeliveryTimestamp = timestamppb.New(lastDeliveryTime)
    }
    
    return deliveries, lastDeliveryTimestamp, totalDelivered, overallStatus, nil
}

// add delivery to order start

func (or *OrderRepository) AddingDishesDeliveryToOrder(ctx context.Context, req *order.CreateDishDeliveryRequest) (*order.DishDelivery, error) {
    or.logger.Info("Adding dish delivery to order repository golang/quanqr/order/order_repository.go")

    // SQL query to insert a new dish delivery
    query := `
        INSERT INTO dish_deliveries (
            order_id,
            order_name,
            guest_id, 
            user_id,
            table_number,
            dish_id,
            quantity_delivered,
            delivery_status,
            delivered_at,
            delivered_by_user_id,
            created_at,
            updated_at,
            is_guest,
            modification_number
        ) VALUES (
            $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13,
            (SELECT COALESCE(MAX(modification_number), 0) + 1 FROM dish_deliveries WHERE order_id = $1)
        ) RETURNING 
            id, 
            order_id,
            order_name,
            guest_id, 
            user_id,
            table_number, 
            dish_id,
            quantity_delivered,
            delivery_status,
            delivered_at,
            delivered_by_user_id,
            created_at,
            updated_at,
            is_guest,
            modification_number
    `

    // Current time for created_at/updated_at if not provided
    now := time.Now()
    createdAt := now
    updatedAt := now
    deliveredAt := now

    // Use provided timestamps if available
    if req.CreatedAt != nil {
        createdAt = req.CreatedAt.AsTime()
    }
    if req.UpdatedAt != nil {
        updatedAt = req.UpdatedAt.AsTime()
    }
    if req.DeliveredAt != nil {
        deliveredAt = req.DeliveredAt.AsTime()
    }

    // Execute the query
    row := or.db.QueryRow(
        ctx,
        query,
        req.OrderId,
        req.OrderName,
        sql.NullInt64{Int64: req.GuestId, Valid: req.GuestId != 0},
        sql.NullInt64{Int64: req.UserId, Valid: req.UserId != 0},
        sql.NullInt64{Int64: req.TableNumber, Valid: req.TableNumber != 0},
        req.DishId,
        req.QuantityDelivered,
        req.DeliveryStatus,
        deliveredAt,
        sql.NullInt64{Int64: req.DeliveredByUserId, Valid: req.DeliveredByUserId != 0},
        createdAt,
        updatedAt,
        req.IsGuest,
    )

    // Prepare response
    var delivery order.DishDelivery
    var guestId, userId, tableNumber, deliveredByUserId sql.NullInt64

    // Scan the result into the response
    err := row.Scan(
        &delivery.Id,
        &delivery.OrderId,
        &delivery.OrderName,
        &guestId,
        &userId,
        &tableNumber,
        &delivery.DishId,
        &delivery.QuantityDelivered,
        &delivery.DeliveryStatus,
        &deliveredAt,
        &deliveredByUserId,
        &createdAt,
        &updatedAt,
        &delivery.IsGuest,
        &delivery.ModificationNumber,
    )
    if err != nil {
        or.logger.Error("Error inserting dish delivery: " + err.Error())
        return nil, fmt.Errorf("error inserting dish delivery: %w", err)
    }

    // Handle NULL values
    if guestId.Valid {
        delivery.GuestId = guestId.Int64
    }
    if userId.Valid {
        delivery.UserId = userId.Int64
    }
    if tableNumber.Valid {
        delivery.TableNumber = tableNumber.Int64
    }
    if deliveredByUserId.Valid {
        delivery.DeliveredByUserId = deliveredByUserId.Int64
    }

    // Set timestamps in protobuf format
    delivery.DeliveredAt = timestamppb.New(deliveredAt)
    delivery.CreatedAt = timestamppb.New(createdAt)
    delivery.UpdatedAt = timestamppb.New(updatedAt)

    // Update order status based on delivery (optional, depending on business logic)
    if err := or.updateOrderDeliveryStatus(ctx, req.OrderId); err != nil {
        or.logger.Error("Error updating order delivery status: " + err.Error())
        // Continue despite the error, as the delivery was created successfully
    }

    return &delivery, nil
}

// Helper function to update the order's delivery status
func (or *OrderRepository) updateOrderDeliveryStatus(ctx context.Context, orderId int64) error {
    // This function would implement business logic to update the order's status
    // based on the deliveries. For example, checking if all items are delivered,
    // and updating the order status accordingly.
    // 
    // Implementing full logic would require knowledge of additional database
    // schema details and business rules.
    
    // Placeholder implementation that could be expanded
    query := `
        WITH order_items AS (
            SELECT SUM(quantity) as total_items
            FROM order_items
            WHERE order_id = $1
        ),
        delivered_items AS (
            SELECT SUM(quantity_delivered) as total_delivered
            FROM dish_deliveries
            WHERE order_id = $1
        )
        UPDATE orders
        SET status = CASE
            WHEN d.total_delivered >= o.total_items THEN 'Delivered'
            WHEN d.total_delivered > 0 THEN 'Partially Delivered'
            ELSE status
        END
        FROM order_items o, delivered_items d
        WHERE id = $1
    `
    
    _, err := or.db.Exec(ctx, query, orderId)
    if err != nil {
        return fmt.Errorf("error updating order status: %w", err)
    }
    
    return nil
}
// add delivery to to order end